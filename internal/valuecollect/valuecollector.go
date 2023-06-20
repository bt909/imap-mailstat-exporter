// Package valuecollect connects to imap, collect values and put them in channels for prometheus exporter functions
package valuecollect

import (
	"crypto/tls"
	"fmt"
	"imap-mailstat-exporter/internal/configread"
	"imap-mailstat-exporter/utils"
	"strings"
	"sync"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var Configfile string
var Loglevel string

type imapStatsCollector struct {
	allMails    *prometheus.Desc
	unseenMails *prometheus.Desc
}

// provide metric "layout"
func NewImapStatsCollector() *imapStatsCollector {
	return &imapStatsCollector{
		allMails: prometheus.NewDesc(
			prometheus.BuildFQName("imap_mailstat", "mails_all", "quantity"),
			"The total number of mails in folder",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		unseenMails: prometheus.NewDesc(
			prometheus.BuildFQName("imap_mailstat", "mails_unseen", "quantity"),
			"The total number of unseen mails in folder",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
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
func countUnseen(c *client.Client, mailbox *imap.MailboxStatus, mailboxname string) (metricname string, namespacename string, messages uint32) {
	metricname = strings.ReplaceAll(mailboxname, " ", "_")
	namespacename = strings.ReplaceAll(mailbox.Name, ".", "_")
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}
	ids, err := c.Search(criteria)
	if err != nil {
		utils.Logger.Error("Error in searching unseen mails", zap.String("mailboxname", fmt.Sprint(mailboxname)), zap.Error(err))
		return
	}
	messages = uint32(len(ids))
	return metricname, namespacename, messages
}

// put metrics description in description channel
func (valuecollector *imapStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- valuecollector.allMails
	ch <- valuecollector.unseenMails
}

// collect values and put them in metrics channel
func (valuecollector *imapStatsCollector) Collect(ch chan<- prometheus.Metric) {
	utils.InitializeLogger(Loglevel)
	config := configread.GetConfig(Configfile)
	sliceLength := len(config.Accounts)
	var wg sync.WaitGroup
	wg.Add(sliceLength)

	for account := range config.Accounts {
		start := time.Now()
		go func(account int) {
			defer wg.Done()
			utils.Logger.Info("Start metrics fetch", zap.String("address", config.Accounts[account].Mailaddress), zap.String("server", config.Accounts[account].Serveraddress))

			var serverconnection strings.Builder
			var c *client.Client
			var err error

			serverconnection.WriteString(config.Accounts[account].Serveraddress)
			serverconnection.WriteString(":")
			serverconnection.WriteString(fmt.Sprint(config.Accounts[account].Serverport))
			if config.Accounts[account].Starttls == true {
				c, err = client.Dial(serverconnection.String())
				tlsConfig := &tls.Config{ServerName: config.Accounts[account].Serveraddress}
				if err := c.StartTLS(tlsConfig); err != nil {
					utils.Logger.Error("failed to dial IMAP server", zap.String("server", fmt.Sprint(config.Accounts[account].Serveraddress)), zap.String("address", fmt.Sprint(config.Accounts[account].Mailaddress)), zap.Error(err))
					return
				}
			} else {
				c, err = client.DialTLS(serverconnection.String(), nil)
				utils.Logger.Info("Connection setup", zap.String("duration", fmt.Sprint(time.Since(start))), zap.String("address", fmt.Sprint(config.Accounts[account].Mailaddress)))
				if err != nil {
					utils.Logger.Error("failed to dial IMAP server", zap.String("server", fmt.Sprint(config.Accounts[account].Serveraddress)), zap.String("address", fmt.Sprint(config.Accounts[account].Mailaddress)), zap.Error(err))
					return
				}
			}

			defer c.Close()

			startLogin := time.Now()
			if err := c.Login(config.Accounts[account].Username, config.Accounts[account].Password); err != nil {
				utils.Logger.Error("failed to login", zap.String("address", fmt.Sprint(config.Accounts[account].Mailaddress)), zap.Error(err))
				return
			}
			utils.Logger.Info("IMAP Login", zap.String("duration:", fmt.Sprint(time.Since(startLogin))), zap.String("address", config.Accounts[account].Mailaddress))

			selectedInbox, err := c.Select("INBOX", true)
			if err != nil {
				utils.Logger.Error("failed to select", zap.String("folder", "Inbox"), zap.Error(err))
				return
			}

			metricSeenInbox, namespaceSeenInBox, countAllmailsInbox := countAllmails(c, selectedInbox, config.Accounts[account].Name)
			var metricnameSeenInbox []string
			metricnameSeenInbox = append(metricnameSeenInbox, metricSeenInbox, namespaceSeenInBox)
			ch <- prometheus.MustNewConstMetric(valuecollector.allMails, prometheus.GaugeValue, float64(countAllmailsInbox), metricnameSeenInbox...)

			metricUnseenInbox, namespaceUnseenInbox, countUnseenInbox := countUnseen(c, selectedInbox, config.Accounts[account].Name)
			var metricnameUnseenInbox []string
			metricnameUnseenInbox = append(metricnameUnseenInbox, metricUnseenInbox, namespaceUnseenInbox)
			ch <- prometheus.MustNewConstMetric(valuecollector.unseenMails, prometheus.GaugeValue, float64(countUnseenInbox), metricnameUnseenInbox...)

			for _, f := range config.Accounts[account].Folders {

				var mboxfolder strings.Builder
				mboxfolder.WriteString("INBOX.")
				mboxfolder.WriteString(f)
				selected, err := c.Select(mboxfolder.String(), true)
				if err != nil {
					utils.Logger.Error("failed to select", zap.String("address", fmt.Sprint(config.Accounts[account].Mailaddress)), zap.String("folder", mboxfolder.String()), zap.Error(err))
					return
				}

				metricSeen, namespaceSeen, countAllmails := countAllmails(c, selected, config.Accounts[account].Name)
				var metricnameSeen []string
				metricnameSeen = append(metricnameSeen, metricSeen, namespaceSeen)
				ch <- prometheus.MustNewConstMetric(valuecollector.allMails, prometheus.GaugeValue, float64(countAllmails), metricnameSeen...)

				metricUnseen, namespaceUnseen, countUnseen := countUnseen(c, selected, config.Accounts[account].Name)
				var metricnameUnseen []string
				metricnameUnseen = append(metricnameUnseen, metricUnseen, namespaceUnseen)
				ch <- prometheus.MustNewConstMetric(valuecollector.unseenMails, prometheus.GaugeValue, float64(countUnseen), metricnameUnseen...)
			}

			if err := c.Logout(); err != nil {
				utils.Logger.Error("failed to logout", zap.String("address", fmt.Sprint(config.Accounts[account].Mailaddress)), zap.Error(err))
				return
			}
			utils.Logger.Info("Metric fetch", zap.String("duration:", fmt.Sprint(time.Since(start))), zap.String("address", config.Accounts[account].Mailaddress))
		}(account)
	}
	wg.Wait()
}
