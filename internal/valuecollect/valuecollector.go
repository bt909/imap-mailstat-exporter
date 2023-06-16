// Package valuecollect connects to imap, collect values and put them in channels for prometheus exporter functions
package valuecollect

import (
	"crypto/tls"
	"fmt"
	"imap-mailstat-exporter/internal/configread"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/prometheus/client_golang/prometheus"
)

var Configfile string

type imapStatsCollector struct {
	entireMails *prometheus.Desc
	unseenMails *prometheus.Desc
}

// provide metric "layout"
func NewImapStatsCollector() *imapStatsCollector {
	return &imapStatsCollector{
		entireMails: prometheus.NewDesc(
			prometheus.BuildFQName("imap", "entire", "mails"),
			"The total number of mails in folder",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
		unseenMails: prometheus.NewDesc(
			prometheus.BuildFQName("imap", "unseen", "mails"),
			"The total number of unseen mails in folder",
			[]string{"mailboxname", "mailboxfoldername"}, nil,
		),
	}

}

// count entire mails and return values and "cleaned" names for using as metric labels (replace characters not allowed in labels)
func countEntire(c *client.Client, mailbox *imap.MailboxStatus, mailboxfolder string) (mailboxfoldername string, mailboxname string, messages uint32) {
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
		log.Fatal(err)
	}
	messages = uint32(len(ids))
	return metricname, namespacename, messages
}

// put metrics description in description channel
func (valuecollector *imapStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- valuecollector.entireMails
	ch <- valuecollector.unseenMails
}

// collect values and put them in metrics channel
func (valuecollector *imapStatsCollector) Collect(ch chan<- prometheus.Metric) {
	config := configread.GetConfig(Configfile)
	sliceLength := len(config.Accounts)
	var wg sync.WaitGroup
	wg.Add(sliceLength)

	for account := range config.Accounts {
		go func(account int) {
			defer wg.Done()
			fmt.Println("Fetch metrics for", config.Accounts[account].Mailaddress, "using server", config.Accounts[account].Serveraddress, "at", time.Now().Format(time.RFC850))

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
					log.Fatalf("failed to dial IMAP server: %v", err)
				}
			} else {
				c, err = client.DialTLS(serverconnection.String(), nil)
				if err != nil {
					log.Fatalf("failed to dial IMAP server: %v", err)
				}
			}

			defer c.Close()

			if err := c.Login(config.Accounts[account].Username, config.Accounts[account].Password); err != nil {
				log.Fatalf("failed to login: %v", err)
			}

			selectedInbox, err := c.Select("INBOX", true)
			if err != nil {
				log.Fatalf("failed to select INBOX: %v", err)
			}

			metricSeenInbox, namespaceSeenInBox, countEntireInbox := countEntire(c, selectedInbox, config.Accounts[account].Name)
			var metricnameSeenInbox []string
			metricnameSeenInbox = append(metricnameSeenInbox, metricSeenInbox, namespaceSeenInBox)
			ch <- prometheus.MustNewConstMetric(valuecollector.entireMails, prometheus.GaugeValue, float64(countEntireInbox), metricnameSeenInbox...)

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
					log.Fatalf("failed to select %s: %v", mboxfolder.String(), err)
				}

				metricSeen, namespaceSeen, countEntire := countEntire(c, selected, config.Accounts[account].Name)
				var metricnameSeen []string
				metricnameSeen = append(metricnameSeen, metricSeen, namespaceSeen)
				ch <- prometheus.MustNewConstMetric(valuecollector.entireMails, prometheus.GaugeValue, float64(countEntire), metricnameSeen...)

				metricUnseen, namespaceUnseen, countUnseen := countUnseen(c, selected, config.Accounts[account].Name)
				var metricnameUnseen []string
				metricnameUnseen = append(metricnameUnseen, metricUnseen, namespaceUnseen)
				ch <- prometheus.MustNewConstMetric(valuecollector.unseenMails, prometheus.GaugeValue, float64(countUnseen), metricnameUnseen...)
			}

			if err := c.Logout(); err != nil {
				log.Fatalf("failed to logout: %v", err)
			}
		}(account)
	}
	wg.Wait()

}
