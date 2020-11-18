package cmd

import (
	"bufio"
	"os"
	"strings"
	"sync"

	"github.com/binaryfigments/dnstransfer/libs/axfr"
	"github.com/binaryfigments/dnstransfer/libs/ns"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "This command checks a list of domain names for zone transfers.",
	Long:  `This command checks a list of domain names for zone transfers.`,
	Run: func(cmd *cobra.Command, args []string) {

		log.Out = os.Stdout

		// You could set this to any `io.Writer` such as a file
		logfile, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.Out = logfile
		} else {
			log.Info("Failed to log to file, using default stderr")
		}

		// get nameserver from flags
		nameserverFlag, _ := cmd.Flags().GetString("nameserver")
		fileFlag, _ := cmd.Flags().GetString("file")

		// open file
		file, err := os.Open(fileFlag)
		if err != nil {
			log.WithFields(logrus.Fields{
				"file": fileFlag,
			}).Fatal(err)
		}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		var domains []string
		for scanner.Scan() {
			domains = append(domains, scanner.Text())
		}
		file.Close()

		var wg sync.WaitGroup
		for _, domain := range domains {
			wg.Add(1)
			go transferWorker(domain, nameserverFlag, &wg)
		}
		wg.Wait()
	},
}

func transferWorker(domain string, nameserverFlag string, wg *sync.WaitGroup) {
	defer wg.Done()
	// start here
	// Get nameservers
	domainnameservers, err := ns.Get(domain, nameserverFlag)
	if err != nil {
		// TODO: build in a retry (in ns function above)
		log.WithFields(logrus.Fields{
			"error":        err,
			"domain":       domain,
			"nameserver":   nameserverFlag,
			"transferable": false,
		}).Warn("Get Nameservers")
	}

	for _, domainnameserver := range domainnameservers {
		// DNS RCODEs: http://www.iana.org/assignments/dns-parameters/dns-parameters.xhtml#dns-parameters-6
		axfrdata, err := axfr.Get(domain, domainnameserver)
		if err != nil {
			// Check if bad xfr code
			if strings.Contains(err.Error(), "bad xfr rcode") {
				log.WithFields(logrus.Fields{
					"error":        err,
					"domain":       domain,
					"nameserver":   domainnameserver,
					"transferable": false,
				}).Info("Zone transfer failed")
			}
			// strings.Contains(err, "tcp"):
			if strings.Contains(err.Error(), "red tcp") {
				log.WithFields(logrus.Fields{
					"error":        err,
					"domain":       domain,
					"nameserver":   domainnameserver,
					"transferable": false,
				}).Info("TCP timeout on port 53")
			}
		}
		if len(axfrdata.Records) > 0 {
			log.WithFields(logrus.Fields{
				"error":        err,
				"domain":       domain,
				"nameserver":   domainnameserver,
				"transferable": true,
			}).Error("Zone can be transfered!")
		}
	}
}

func init() {
	rootCmd.AddCommand(checkCmd)
	// TODO: Logging
	// log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.ErrorLevel)
	log.SetOutput(os.Stdout)
	checkCmd.PersistentFlags().String("file", "domains.txt", "File with domain names.")
}
