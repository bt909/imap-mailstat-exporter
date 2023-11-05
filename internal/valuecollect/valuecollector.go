// Package valuecollect connects to imap, collect values and put them in channels for prometheus exporter functions
package valuecollect

import (
	"crypto/tls"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/emersion/go-imap"
	quota "github.com/emersion/go-imap-quota"
	"github.com/emersion/go-imap/client"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/bt909/imap-mailstat-exporter/internal/configread"
)

type imapStatsCollector struct {
	up                   *prometheus.Desc
	allMails             *prometheus.Desc
	unseenMails          *prometheus.Desc
	mailboxQuotaUsed     *prometheus.Desc
	mailboxQuotaAvail    *prometheus.Desc
	levelQuotaUsed       *prometheus.Desc
	levelQuotaAvail      *prometheus.Desc
	storageQuotaUsed     *prometheus.Desc
	storageQuotaAvail    *prometheus.Desc
	messageQuotaUsed     *prometheus.Desc
	messageQuotaAvail    *prometheus.Desc
	oldestUnseen         *prometheus.Desc
	fetchDuration        *prometheus.Desc
	allMailsOld          *prometheus.Desc
	unseenMailsOld       *prometheus.Desc
	mailboxQuotaUsedOld  *prometheus.Desc
	mailboxQuotaAvailOld *prometheus.Desc
	levelQuotaUsedOld    *prometheus.Desc
	levelQuotaAvailOld   *prometheus.Desc
	storageQuotaUsedOld  *prometheus.Desc
	storageQuotaAvailOld *prometheus.Desc
	messageQuotaUsedOld  *prometheus.Desc
	messageQuotaAvailOld *prometheus.Desc
	oldestUnseenOld      *prometheus.Desc
	info                 *prometheus.Desc
	configfile           string
	logger               log.Logger
	oldestunseenfeature  bool
	migrationmode        bool
	version              string
}

// provide metric "layout"
func NewImapStatsCollector(configfile string, logger log.Logger, oldestunseenfeature bool, migrationmode bool, version string) *imapStatsCollector {
	return &imapStatsCollector{
		up:   prometheus.NewDesc("mailstat_up", "Was talking to all accounts imap successfully.", nil, nil),
		info: prometheus.NewDesc("mailstat_info", "Info metric for imap-mailstat-exporter.", []string{"version"}, nil),
		allMails: prometheus.NewDesc(
			prometheus.BuildFQName("mailstat", "mails", "all"),
			"The total number of mails in folder.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		unseenMails: prometheus.NewDesc(
			prometheus.BuildFQName("mailstat", "mails", "unseen"),
			"The total number of unseen mails in folder.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		mailboxQuotaUsed: prometheus.NewDesc(
			prometheus.BuildFQName("mailstat", "mailbox_quota", "used"),
			"How many mailboxes are used.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		mailboxQuotaAvail: prometheus.NewDesc(
			prometheus.BuildFQName("mailstat", "mailbox_quota", "avail"),
			"How many mailboxes are available according your quota.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		levelQuotaUsed: prometheus.NewDesc(
			prometheus.BuildFQName("mailstat", "level_quota", "used"),
			"How many levels are used.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		levelQuotaAvail: prometheus.NewDesc(
			prometheus.BuildFQName("mailstat", "level_quota", "avail"),
			"How many levels are available according your quota.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		storageQuotaUsed: prometheus.NewDesc(
			prometheus.BuildFQName("mailstat", "storage_quota", "used_bytes"),
			"How many storage is used.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		storageQuotaAvail: prometheus.NewDesc(
			prometheus.BuildFQName("mailstat", "storage_quota", "avail_bytes"),
			"How many storage is available according your quota.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		messageQuotaUsed: prometheus.NewDesc(
			prometheus.BuildFQName("mailstat", "message_quota", "used"),
			"How many messages are used.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		messageQuotaAvail: prometheus.NewDesc(
			prometheus.BuildFQName("mailstat", "message_quota", "avail"),
			"How many messages available according your quota.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		oldestUnseen: prometheus.NewDesc(
			prometheus.BuildFQName("mailstat", "mails", "oldest_unseen_timestamp"),
			"Timestamp in unix format of oldest unseen mail.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		fetchDuration: prometheus.NewDesc(
			prometheus.BuildFQName("mailstat", "fetch", "duration_seconds"),
			"Duration for fetching the metrics for the given account",
			[]string{"mailboxname"}, nil,
		),
		allMailsOld: prometheus.NewDesc(
			prometheus.BuildFQName("imap_mailstat", "mails_all", "quantity"),
			"The total number of mails in folder.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		unseenMailsOld: prometheus.NewDesc(
			prometheus.BuildFQName("imap_mailstat", "mails_unseen", "quantity"),
			"The total number of unseen mails in folder.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		mailboxQuotaUsedOld: prometheus.NewDesc(
			prometheus.BuildFQName("imap_mailstat", "mails_mailboxquotaused", "quantity"),
			"How many mailboxes are used.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		mailboxQuotaAvailOld: prometheus.NewDesc(
			prometheus.BuildFQName("imap_mailstat", "mails_mailboxquotaavail", "quantity"),
			"How many mailboxes are available according your quota.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		levelQuotaUsedOld: prometheus.NewDesc(
			prometheus.BuildFQName("imap_mailstat", "mails_levelquotaused", "quantity"),
			"How many levels are used.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		levelQuotaAvailOld: prometheus.NewDesc(
			prometheus.BuildFQName("imap_mailstat", "mails_levelquotaavail", "quantity"),
			"How many levels are available according your quota.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		storageQuotaUsedOld: prometheus.NewDesc(
			prometheus.BuildFQName("imap_mailstat", "mails_storagequotaused", "kilobytes"),
			"How many storage is used.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		storageQuotaAvailOld: prometheus.NewDesc(
			prometheus.BuildFQName("imap_mailstat", "mails_storagequotaavail", "kilobytes"),
			"How many storage is available according your quota.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		messageQuotaUsedOld: prometheus.NewDesc(
			prometheus.BuildFQName("imap_mailstat", "mails_messagequotaused", "quantity"),
			"How many messages are used.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		messageQuotaAvailOld: prometheus.NewDesc(
			prometheus.BuildFQName("imap_mailstat", "mails_messagequotaavail", "quantity"),
			"How many messages available according your quota.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		oldestUnseenOld: prometheus.NewDesc(
			prometheus.BuildFQName("imap_mailstat", "mails_oldestunseen", "timestamp"),
			"Timestamp in unix format of oldest unseen mail.",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		configfile:          configfile,
		logger:              logger,
		oldestunseenfeature: oldestunseenfeature,
		migrationmode:       migrationmode,
		version:             version,
	}

}

// count all mails and return values and "cleaned" names for using as metric labels (replace characters not allowed in labels)
func countAllmails(c *client.Client, mailbox *imap.MailboxStatus, mailboxfolder string) (mailboxfoldername string, mailboxname string, messages uint32) {
	mailboxfolder = strings.ReplaceAll(mailboxfolder, " ", "_")
	mailboxname = strings.ReplaceAll(mailbox.Name, ".", "_")
	messages = mailbox.Messages
	return mailboxfolder, mailboxname, messages
}

// count unseen mails and return values and "cleaned" names for using as metric labels (replace characters not allowed in labels)
func countUnseen(c *client.Client, mailbox *imap.MailboxStatus, logger log.Logger, mailboxname string, oldestunseenfeature bool) (metricname string, namespacename string, messages uint32, oldestunseen int64) {
	metricname = strings.ReplaceAll(mailboxname, " ", "_")
	namespacename = strings.ReplaceAll(mailbox.Name, ".", "_")
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}
	ids, err := c.Search(criteria)
	if err != nil {
		level.Error(logger).Log("msg", "Error in searching mails", "mailboxname", mailboxname)
		return
	}
	// if feature flag is enabled and there are unseen mail ids, we try to get the date from the envelope header and convert to unix timestamp
	if len(ids) > 0 && oldestunseenfeature {
		seqset := new(imap.SeqSet)
		seqset.AddNum(ids...)
		mails := make(chan *imap.Message, len(ids))
		err := c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, mails)
		if err != nil {
			level.Error(logger).Log("msg", "Error in getting dates for unseen mails", "mailboxname", mailboxname)
			return
		}
		var dates []int64
		for msg := range mails {
			dates = append(dates, msg.Envelope.Date.Unix())
			// this sort is to ensure first value is oldest
			sort.Slice(dates, func(i, j int) bool {
				return dates[i] < dates[j]
			})
			oldestunseen = dates[0]
		}
	}
	messages = uint32(len(ids))
	return metricname, namespacename, messages, oldestunseen
}

// returns quota related values and "cleaned" names for using as metric labels (replace characters not allowed in labels)
func getMailboxUsed(qc *quota.Client, mailbox *imap.MailboxStatus, logger log.Logger, mailboxname string) (metricname string, namespacename string, mailboxUsed map[string]uint32, mailboxAvail map[string]uint32) {
	metricname = strings.ReplaceAll(mailboxname, " ", "_")
	namespacename = strings.ReplaceAll(mailbox.Name, ".", "_")

	mailboxUsed = make(map[string]uint32)
	mailboxAvail = make(map[string]uint32)

	// Retrieve quotas for INBOX
	quotas, err := qc.GetQuotaRoot("INBOX")
	if err != nil {
		level.Error(logger).Log("msg", "Error in getting quota for INBOX", "mailboxname", mailboxname)
		return
	}

	// put quota values in return values (index 0 is used, index 1 is available, convert from kilobytes to bytes as common prometheus unit)
	for _, quota := range quotas {
		for name, usage := range quota.Resources {
			mailboxUsed[name+"_used"] = usage[0] * 1000
			mailboxAvail[name+"_avail"] = usage[1] * 1000
		}
	}

	return metricname, namespacename, mailboxUsed, mailboxAvail
}

// put metrics description in description channel
func (valuecollector *imapStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- valuecollector.up
	ch <- valuecollector.info
	ch <- valuecollector.allMails
	ch <- valuecollector.unseenMails
	ch <- valuecollector.mailboxQuotaUsed
	ch <- valuecollector.mailboxQuotaAvail
	ch <- valuecollector.levelQuotaUsed
	ch <- valuecollector.levelQuotaAvail
	ch <- valuecollector.storageQuotaUsed
	ch <- valuecollector.storageQuotaAvail
	ch <- valuecollector.messageQuotaUsed
	ch <- valuecollector.messageQuotaAvail
	ch <- valuecollector.oldestUnseen
	ch <- valuecollector.fetchDuration
	ch <- valuecollector.allMailsOld
	ch <- valuecollector.unseenMailsOld
	ch <- valuecollector.mailboxQuotaUsedOld
	ch <- valuecollector.mailboxQuotaAvailOld
	ch <- valuecollector.levelQuotaUsedOld
	ch <- valuecollector.levelQuotaAvailOld
	ch <- valuecollector.storageQuotaUsedOld
	ch <- valuecollector.storageQuotaAvailOld
	ch <- valuecollector.messageQuotaUsedOld
	ch <- valuecollector.messageQuotaAvailOld
	ch <- valuecollector.oldestUnseenOld
}

// collect values and put them in metrics channel
func (valuecollector *imapStatsCollector) Collect(ch chan<- prometheus.Metric) {
	config := configread.GetConfig(valuecollector.configfile)
	// set up variable to 1 and if any account has some errors it is set to 0 and used as up metric
	up := 1
	ch <- prometheus.MustNewConstMetric(valuecollector.info, prometheus.GaugeValue, 1, valuecollector.version)
	sliceLength := len(config.Accounts)
	var wg sync.WaitGroup
	wg.Add(sliceLength)

	for account := range config.Accounts {
		start := time.Now()
		go func(account int) {
			defer wg.Done()
			level.Info(valuecollector.logger).Log("msg", "Start metrics fetch", "address", config.Accounts[account].Mailaddress, "server", config.Accounts[account].Serveraddress)

			var serverconnection strings.Builder
			var c *client.Client
			var err error

			serverconnection.WriteString(config.Accounts[account].Serveraddress)
			serverconnection.WriteString(":")
			serverconnection.WriteString(fmt.Sprint(config.Accounts[account].Serverport))
			if config.Accounts[account].Starttls {
				c, err = client.Dial(serverconnection.String())
				if err != nil {
					level.Error(valuecollector.logger).Log("msg", "Failed to dial IMAP server", "address", config.Accounts[account].Mailaddress, "server", config.Accounts[account].Serveraddress)
					up = 0
					return
				}
				tlsConfig := &tls.Config{ServerName: config.Accounts[account].Serveraddress}
				if err := c.StartTLS(tlsConfig); err != nil {
					level.Error(valuecollector.logger).Log("msg", "Failed to start TLS secured connection via StartTLS", "address", config.Accounts[account].Mailaddress, "server", config.Accounts[account].Serveraddress)
					up = 0
					return
				}
			} else {
				c, err = client.DialTLS(serverconnection.String(), nil)
				level.Info(valuecollector.logger).Log("msg", "Connection setup", "duration", time.Since(start), "address", config.Accounts[account].Mailaddress)
				if err != nil {
					level.Error(valuecollector.logger).Log("msg", "Failed to dial server via TLS", "address", config.Accounts[account].Mailaddress, "server", config.Accounts[account].Serveraddress)
					up = 0
					return
				}
			}

			defer c.Close()

			startLogin := time.Now()
			if err := c.Login(config.Accounts[account].Username, config.Accounts[account].Password); err != nil {
				level.Error(valuecollector.logger).Log("msg", "Failed to login", "address", config.Accounts[account].Mailaddress, "server", config.Accounts[account].Serveraddress)
				up = 0
				return
			}
			level.Info(valuecollector.logger).Log("msg", "IMAP login", "duration", time.Since(startLogin), "address", config.Accounts[account].Mailaddress, "server", config.Accounts[account].Serveraddress)

			defer c.Logout()

			selectedInbox, err := c.Select("INBOX", true)
			if err != nil {
				level.Error(valuecollector.logger).Log("msg", "Failed to select", "folder", "Inbox", "address", config.Accounts[account].Mailaddress)
				up = 0
				return
			}

			metricSeenInbox, namespaceSeenInBox, countAllmailsInbox := countAllmails(c, selectedInbox, config.Accounts[account].Name)
			var metricnameSeenInbox []string
			metricnameSeenInbox = append(metricnameSeenInbox, metricSeenInbox, namespaceSeenInBox)
			ch <- prometheus.MustNewConstMetric(valuecollector.allMails, prometheus.GaugeValue, float64(countAllmailsInbox), metricnameSeenInbox...)
			if valuecollector.migrationmode {
				ch <- prometheus.MustNewConstMetric(valuecollector.allMailsOld, prometheus.GaugeValue, float64(countAllmailsInbox), metricnameSeenInbox...)
			}

			metricUnseenInbox, namespaceUnseenInbox, countUnseenInbox, _ := countUnseen(c, selectedInbox, valuecollector.logger, config.Accounts[account].Name, valuecollector.oldestunseenfeature)
			var metricnameUnseenInbox []string
			metricnameUnseenInbox = append(metricnameUnseenInbox, metricUnseenInbox, namespaceUnseenInbox)
			ch <- prometheus.MustNewConstMetric(valuecollector.unseenMails, prometheus.GaugeValue, float64(countUnseenInbox), metricnameUnseenInbox...)
			if valuecollector.migrationmode {
				ch <- prometheus.MustNewConstMetric(valuecollector.unseenMailsOld, prometheus.GaugeValue, float64(countUnseenInbox), metricnameUnseenInbox...)
			}

			metricOldestUnseenInbox, namespaceOldestUnseenInbox, _, timestampOldestUnseenInbox := countUnseen(c, selectedInbox, valuecollector.logger, config.Accounts[account].Name, valuecollector.oldestunseenfeature)
			var metricnameOldestUnseenInbox []string
			if timestampOldestUnseenInbox > 0 {
				metricnameOldestUnseenInbox = append(metricnameOldestUnseenInbox, metricOldestUnseenInbox, namespaceOldestUnseenInbox)
				ch <- prometheus.MustNewConstMetric(valuecollector.oldestUnseen, prometheus.GaugeValue, float64(timestampOldestUnseenInbox), metricnameOldestUnseenInbox...)
				if valuecollector.migrationmode {
					ch <- prometheus.MustNewConstMetric(valuecollector.oldestUnseenOld, prometheus.GaugeValue, float64(timestampOldestUnseenInbox), metricnameOldestUnseenInbox...)
				}
			}

			qc := quota.NewClient(c)

			// Check for server support and set metrics only for accounts with metrics available
			if quotaSupport, _ := qc.SupportQuota(); quotaSupport {
				level.Info(valuecollector.logger).Log("msg", "Fetching quota related metrics", "address", config.Accounts[account].Mailaddress)

				metricMailboxQuotaUsed, namespaceMailboxQuotaUsed, countMailboxQuotaUsed, _ := getMailboxUsed(qc, selectedInbox, valuecollector.logger, config.Accounts[account].Name)
				var metricnameMailboxQuotaUsed []string
				metricnameMailboxQuotaUsed = append(metricnameMailboxQuotaUsed, metricMailboxQuotaUsed, namespaceMailboxQuotaUsed)
				ch <- prometheus.MustNewConstMetric(valuecollector.mailboxQuotaUsed, prometheus.GaugeValue, float64(countMailboxQuotaUsed["MAILBOX_used"]), metricnameMailboxQuotaUsed...)
				if valuecollector.migrationmode {
					ch <- prometheus.MustNewConstMetric(valuecollector.mailboxQuotaUsedOld, prometheus.GaugeValue, float64(countMailboxQuotaUsed["MAILBOX_used"]), metricnameMailboxQuotaUsed...)
				}

				metricMailboxQuotaAvail, namespaceMailboxQuotaAvail, _, countMailboxQuotaAvail := getMailboxUsed(qc, selectedInbox, valuecollector.logger, config.Accounts[account].Name)
				var metricnameMailboxQuotaAvail []string
				metricnameMailboxQuotaAvail = append(metricnameMailboxQuotaAvail, metricMailboxQuotaAvail, namespaceMailboxQuotaAvail)
				ch <- prometheus.MustNewConstMetric(valuecollector.mailboxQuotaAvail, prometheus.GaugeValue, float64(countMailboxQuotaAvail["MAILBOX_avail"]), metricnameMailboxQuotaAvail...)
				if valuecollector.migrationmode {
					ch <- prometheus.MustNewConstMetric(valuecollector.mailboxQuotaAvailOld, prometheus.GaugeValue, float64(countMailboxQuotaAvail["MAILBOX_avail"]), metricnameMailboxQuotaAvail...)
				}

				metricLevelQuotaUsed, namespaceLevelQuotaUsed, countLevelQuotaUsed, _ := getMailboxUsed(qc, selectedInbox, valuecollector.logger, config.Accounts[account].Name)
				var metricnameLevelQuotaUsed []string
				metricnameLevelQuotaUsed = append(metricnameLevelQuotaUsed, metricLevelQuotaUsed, namespaceLevelQuotaUsed)
				ch <- prometheus.MustNewConstMetric(valuecollector.levelQuotaUsed, prometheus.GaugeValue, float64(countLevelQuotaUsed["LEVEL_used"]), metricnameLevelQuotaUsed...)
				if valuecollector.migrationmode {
					ch <- prometheus.MustNewConstMetric(valuecollector.levelQuotaUsedOld, prometheus.GaugeValue, float64(countLevelQuotaUsed["LEVEL_used"]), metricnameLevelQuotaUsed...)
				}

				metricLevelQuotaAvail, namespaceLevelQuotaAvail, _, countLevelQuotaAvail := getMailboxUsed(qc, selectedInbox, valuecollector.logger, config.Accounts[account].Name)
				var metricnameLevelQuotaAvail []string
				metricnameLevelQuotaAvail = append(metricnameLevelQuotaAvail, metricLevelQuotaAvail, namespaceLevelQuotaAvail)
				ch <- prometheus.MustNewConstMetric(valuecollector.levelQuotaAvail, prometheus.GaugeValue, float64(countLevelQuotaAvail["LEVEL_avail"]), metricnameLevelQuotaAvail...)
				if valuecollector.migrationmode {
					ch <- prometheus.MustNewConstMetric(valuecollector.levelQuotaAvailOld, prometheus.GaugeValue, float64(countLevelQuotaAvail["LEVEL_avail"]), metricnameLevelQuotaAvail...)
				}

				metricStorageQuotaUsed, namespaceStorageQuotaUsed, countStorageQuotaUsed, _ := getMailboxUsed(qc, selectedInbox, valuecollector.logger, config.Accounts[account].Name)
				var metricnameStorageQuotaUsed []string
				metricnameStorageQuotaUsed = append(metricnameStorageQuotaUsed, metricStorageQuotaUsed, namespaceStorageQuotaUsed)
				ch <- prometheus.MustNewConstMetric(valuecollector.storageQuotaUsed, prometheus.GaugeValue, float64(countStorageQuotaUsed["STORAGE_used"]), metricnameStorageQuotaUsed...)
				if valuecollector.migrationmode {
					ch <- prometheus.MustNewConstMetric(valuecollector.storageQuotaUsedOld, prometheus.GaugeValue, float64(countStorageQuotaUsed["STORAGE_used"]/1000), metricnameStorageQuotaUsed...)
				}

				metricStorageQuotaAvail, namespaceStorageQuotaAvail, _, countStorageQuotaAvail := getMailboxUsed(qc, selectedInbox, valuecollector.logger, config.Accounts[account].Name)
				var metricnameStorageQuotaAvail []string
				metricnameStorageQuotaAvail = append(metricnameStorageQuotaAvail, metricStorageQuotaAvail, namespaceStorageQuotaAvail)
				ch <- prometheus.MustNewConstMetric(valuecollector.storageQuotaAvail, prometheus.GaugeValue, float64(countStorageQuotaAvail["STORAGE_avail"]), metricnameStorageQuotaAvail...)
				if valuecollector.migrationmode {
					ch <- prometheus.MustNewConstMetric(valuecollector.storageQuotaAvailOld, prometheus.GaugeValue, float64(countStorageQuotaAvail["STORAGE_avail"]/1000), metricnameStorageQuotaAvail...)
				}

				metricMessageQuotaUsed, namespaceMessageQuotaUsed, countMessageQuotaUsed, _ := getMailboxUsed(qc, selectedInbox, valuecollector.logger, config.Accounts[account].Name)
				var metricnameMessageQuotaUsed []string
				metricnameMessageQuotaUsed = append(metricnameMessageQuotaUsed, metricMessageQuotaUsed, namespaceMessageQuotaUsed)
				ch <- prometheus.MustNewConstMetric(valuecollector.messageQuotaUsed, prometheus.GaugeValue, float64(countMessageQuotaUsed["MESSAGE_used"]), metricnameMessageQuotaUsed...)
				if valuecollector.migrationmode {
					ch <- prometheus.MustNewConstMetric(valuecollector.messageQuotaUsedOld, prometheus.GaugeValue, float64(countMessageQuotaUsed["MESSAGE_used"]), metricnameMessageQuotaUsed...)
				}

				metricMessageQuotaAvail, namespaceMessageQuotaAvail, _, countMessageQuotaAvail := getMailboxUsed(qc, selectedInbox, valuecollector.logger, config.Accounts[account].Name)
				var metricnameMessageQuotaAvail []string
				metricnameMessageQuotaAvail = append(metricnameMessageQuotaAvail, metricMessageQuotaAvail, namespaceMessageQuotaAvail)
				ch <- prometheus.MustNewConstMetric(valuecollector.messageQuotaAvail, prometheus.GaugeValue, float64(countMessageQuotaAvail["MESSAGE_avail"]), metricnameMessageQuotaAvail...)
				if valuecollector.migrationmode {
					ch <- prometheus.MustNewConstMetric(valuecollector.messageQuotaAvailOld, prometheus.GaugeValue, float64(countMessageQuotaAvail["MESSAGE_avail"]), metricnameMessageQuotaAvail...)
				}
			}

			for _, f := range config.Accounts[account].Folders {

				var mboxfolder strings.Builder
				mboxfolder.WriteString("INBOX.")
				mboxfolder.WriteString(f)
				selected, err := c.Select(mboxfolder.String(), true)
				if err != nil {
					level.Error(valuecollector.logger).Log("msg", "Failed to select", "address", config.Accounts[account].Mailaddress, "folder", mboxfolder)
					up = 0
					return
				}

				metricSeen, namespaceSeen, countAllmails := countAllmails(c, selected, config.Accounts[account].Name)
				var metricnameSeen []string
				metricnameSeen = append(metricnameSeen, metricSeen, namespaceSeen)
				ch <- prometheus.MustNewConstMetric(valuecollector.allMails, prometheus.GaugeValue, float64(countAllmails), metricnameSeen...)
				if valuecollector.migrationmode {
					ch <- prometheus.MustNewConstMetric(valuecollector.allMailsOld, prometheus.GaugeValue, float64(countAllmails), metricnameSeen...)
				}

				metricUnseen, namespaceUnseen, countUnseenMails, _ := countUnseen(c, selected, valuecollector.logger, config.Accounts[account].Name, valuecollector.oldestunseenfeature)
				var metricnameUnseen []string
				metricnameUnseen = append(metricnameUnseen, metricUnseen, namespaceUnseen)
				ch <- prometheus.MustNewConstMetric(valuecollector.unseenMails, prometheus.GaugeValue, float64(countUnseenMails), metricnameUnseen...)
				if valuecollector.migrationmode {
					ch <- prometheus.MustNewConstMetric(valuecollector.unseenMailsOld, prometheus.GaugeValue, float64(countUnseenMails), metricnameUnseen...)
				}

				metricOldestUnseen, namespaceOldestUnseen, _, timestampOldestUnseen := countUnseen(c, selected, valuecollector.logger, config.Accounts[account].Name, valuecollector.oldestunseenfeature)
				if timestampOldestUnseen > 0 {
					var metricnameOldestUnseen []string
					metricnameOldestUnseen = append(metricnameOldestUnseen, metricOldestUnseen, namespaceOldestUnseen)
					ch <- prometheus.MustNewConstMetric(valuecollector.oldestUnseen, prometheus.GaugeValue, float64(timestampOldestUnseen), metricnameOldestUnseen...)
					if valuecollector.migrationmode {
						ch <- prometheus.MustNewConstMetric(valuecollector.oldestUnseenOld, prometheus.GaugeValue, float64(timestampOldestUnseen), metricnameOldestUnseen...)
					}
				}
			}
			ch <- prometheus.MustNewConstMetric(valuecollector.fetchDuration, prometheus.GaugeValue, float64(time.Since(start).Seconds()), strings.ReplaceAll(strings.ReplaceAll(config.Accounts[account].Name, ".", "_"), " ", "_"))
			level.Info(valuecollector.logger).Log("msg", "Metric fetch", "duration", time.Since(start), "address", config.Accounts[account].Mailaddress)
		}(account)
	}
	wg.Wait()
	ch <- prometheus.MustNewConstMetric(valuecollector.up, prometheus.GaugeValue, float64(up))
}
