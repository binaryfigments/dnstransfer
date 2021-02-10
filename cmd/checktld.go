package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/binaryfigments/dnstransfer/libs/axfr"
	"github.com/binaryfigments/dnstransfer/libs/ns"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checktldCmd)
	log.SetOutput(os.Stderr)
	checktldCmd.PersistentFlags().String("file", "tlds-alpha-by-domain.txt", "File with domain names.")
}

// checkCmd represents the check command
var checktldCmd = &cobra.Command{
	Use:   "checktld",
	Short: "This command checks a list TLD's for zone transfers",
	Long:  `This command checks a list TLD's for zone transfers`,
	Run: func(cmd *cobra.Command, args []string) {
		// log.Out = os.Stderr <------------------------------- see init for this

		// get nameserver from flags
		nameserverFlag, _ := cmd.Flags().GetString("nameserver")
		fileFlag, _ := cmd.Flags().GetString("file")

		log.WithFields(logrus.Fields{
			"file":       fileFlag,
			"check":      "checktld",
			"nameserver": nameserverFlag,
		}).Info("Starting check")

		// Set the output format
		jo, _ := cmd.Flags().GetBool("json")
		if jo == true {
			log.SetFormatter(&logrus.JSONFormatter{})
		}

		// Set the log level
		do, _ := cmd.Flags().GetBool("debug")
		if do == true {
			log.SetLevel(logrus.InfoLevel)
		} else {
			log.SetLevel(logrus.ErrorLevel)
		}

		// You could set this to any `io.Writer` such as a file
		lf, _ := cmd.Flags().GetString("logfile")
		if lf != "" {
			logfile, err := os.OpenFile(lf, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err == nil {
				log.Out = logfile
			} else {
				log.Error("Failed to log to file, using default stderr")
			}
		}

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
			domain = strings.ToLower(domain + ".")
			wg.Add(1)
			go transferTldWorker(domain, nameserverFlag, &wg)
		}
		wg.Wait()

		log.SetLevel(logrus.InfoLevel)
		log.WithFields(logrus.Fields{
			"file":       fileFlag,
			"check":      "checktld",
			"nameserver": nameserverFlag,
		}).Info("Check finished")
	},
}

func transferTldWorker(domain string, nameserverFlag string, wg *sync.WaitGroup) {
	defer wg.Done()
	// start here
	// Get nameservers
	var domainnameservers []string
	err := retry(2, 2*time.Second, func() (err error) {
		domainnameservers, err = ns.Get(domain, nameserverFlag)
		return
	})
	// domainnameservers, err := ns.Get(domain, nameserverFlag)
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
			fmt.Printf("%v,%v,%v\n", domain, domainnameserver, true)
			log.WithFields(logrus.Fields{
				"error":        err,
				"domain":       domain,
				"nameserver":   domainnameserver,
				"transferable": true,
			}).Warn("Zone can be transfered!")
		}
	}
}
