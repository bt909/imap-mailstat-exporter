package valuecollector

import (
	"fmt"
	"imap-mailstat-exporter/pkg/configread"
	"log"
	"strings"
	"sync"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/prometheus/client_golang/prometheus"
)

type imapStatsCollector struct {
	seenMails   *prometheus.Desc
	unseenMails *prometheus.Desc
}

func NewImapStatsCollector() *imapStatsCollector {
	return &imapStatsCollector{
		seenMails: prometheus.NewDesc(
			prometheus.BuildFQName("imap", "total", "mails"),
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

func countSeen(c *client.Client, mailbox *imap.MailboxStatus, mailboxfolder string) (mailboxfoldername string, mailboxname string, messages uint32) {
	mailboxfolder = strings.ReplaceAll(mailboxfolder, " ", "_")
	mailboxname = strings.ReplaceAll(mailbox.Name, ".", "_")
	messages = mailbox.Messages
	return mailboxfolder, mailboxname, messages
}

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

func (valuecollector *imapStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- valuecollector.seenMails
	ch <- valuecollector.unseenMails
}

func (valuecollector *imapStatsCollector) Collect(ch chan<- prometheus.Metric) {
	config := configread.GetConfig()
	sliceLength := len(config.Accounts)
	var wg sync.WaitGroup
	wg.Add(sliceLength)

	for account := range config.Accounts {
		go func(account int) {
			defer wg.Done()
			fmt.Println("Fetch metrics for", config.Accounts[account].Mailaddress, "using server", config.Accounts[account].Serveraddress)

			var serverconnection strings.Builder
			serverconnection.WriteString(config.Accounts[account].Serveraddress)
			serverconnection.WriteString(":")
			serverconnection.WriteString(fmt.Sprint(config.Accounts[account].Serverport))
			c, err := client.DialTLS(serverconnection.String(), nil)
			if err != nil {
				log.Fatalf("failed to dial IMAP server: %v", err)
			}
			defer c.Close()

			if err := c.Login(config.Accounts[account].Mailaddress, config.Accounts[account].Password); err != nil {
				log.Fatalf("failed to login: %v", err)
			}

			selectedInbox, err := c.Select("INBOX", true)
			if err != nil {
				log.Fatalf("failed to select INBOX: %v", err)
			}

			metricSeenInbox, namespaceSeenInBox, countSeenInbox := countSeen(c, selectedInbox, config.Accounts[account].Name)
			var metricnameSeenInbox []string
			metricnameSeenInbox = append(metricnameSeenInbox, metricSeenInbox, namespaceSeenInBox)
			ch <- prometheus.MustNewConstMetric(valuecollector.seenMails, prometheus.GaugeValue, float64(countSeenInbox), metricnameSeenInbox...)

			metricUnseenInbox, namespaceUnseenInbox, countUnseenInbox := countUnseen(c, selectedInbox, config.Accounts[account].Name)
			var metricnameUnseenInbox []string
			metricnameUnseenInbox = append(metricnameUnseenInbox, metricUnseenInbox, namespaceUnseenInbox)
			ch <- prometheus.MustNewConstMetric(valuecollector.unseenMails, prometheus.GaugeValue, float64(countUnseenInbox), metricnameUnseenInbox...)

			/* countSeen(c, selectedInbox, config.Accounts[account].Name)
			countUnseen(c, selectedInbox, config.Accounts[account].Name) */

			for _, f := range config.Accounts[account].Folders {

				var mboxfolder strings.Builder
				mboxfolder.WriteString("INBOX.")
				mboxfolder.WriteString(f)
				selected, err := c.Select(mboxfolder.String(), true)
				if err != nil {
					log.Fatalf("failed to select %s: %v", mboxfolder.String(), err)
				}

				metricSeen, namespaceSeen, countSeen := countSeen(c, selected, config.Accounts[account].Name)
				var metricnameSeen []string
				metricnameSeen = append(metricnameSeen, metricSeen, namespaceSeen)
				ch <- prometheus.MustNewConstMetric(valuecollector.seenMails, prometheus.GaugeValue, float64(countSeen), metricnameSeen...)

				metricUnseen, namespaceUnseen, countUnseen := countUnseen(c, selected, config.Accounts[account].Name)
				var metricnameUnseen []string
				metricnameUnseen = append(metricnameUnseen, metricUnseen, namespaceUnseen)
				ch <- prometheus.MustNewConstMetric(valuecollector.unseenMails, prometheus.GaugeValue, float64(countUnseen), metricnameUnseen...)

				//countUnseen(c, selected, config.Accounts[account].Name)
			}

			if err := c.Logout(); err != nil {
				log.Fatalf("failed to logout: %v", err)
			}
		}(account)
	}
	wg.Wait()

}
