package cmd

import (
	"bufio"
	"os"
	"strings"

	"github.com/binaryfigments/dnstransfer/libs/axfr"
	"github.com/binaryfigments/dnstransfer/libs/ns"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "This command checks a list of domain names for zone transfers.",
	Long:  `This command checks a list of domain names for zone transfers.`,
	Run: func(cmd *cobra.Command, args []string) {

		// get nameserver from flags
		nameserverFlag, _ := cmd.Flags().GetString("nameserver")
		fileFlag, _ := cmd.Flags().GetString("file")

		log.WithFields(log.Fields{
			"file":       fileFlag,
			"nameserver": nameserverFlag,
		}).Info("Starting dnstransfer command")

		// open file
		file, err := os.Open(fileFlag)
		if err != nil {
			log.WithFields(log.Fields{
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

		for _, domain := range domains {
			// Get nameservers
			domainnameservers, err := ns.Get(domain, nameserverFlag)
			if err != nil {
				// TODO: build in a retry (in ns function above)
				log.WithFields(log.Fields{
					"error":  err,
					"domain": domain,
				}).Warn("Get Nameservers")
			}

			for _, domainnameserver := range domainnameservers {
				log.WithFields(log.Fields{
					"domain":     domain,
					"nameserver": domainnameserver,
				}).Info("Start zone transfer check")

				// DNS RCODEs: http://www.iana.org/assignments/dns-parameters/dns-parameters.xhtml#dns-parameters-6

				axfrdata, err := axfr.Get(domain, domainnameserver)
				if err != nil {
					// Check if bad xfr code
					if strings.Contains(err.Error(), "bad xfr rcode") {
						log.WithFields(log.Fields{
							"error":      err,
							"domain":     domain,
							"nameserver": domainnameserver,
						}).Info("Zone transfer failed")
					}
					// strings.Contains(err, "tcp"):
					if strings.Contains(err.Error(), "red tcp") {
						log.WithFields(log.Fields{
							"error":      err,
							"domain":     domain,
							"nameserver": domainnameserver,
						}).Info("TCP timeout on port 53")
					}
				}
				if len(axfrdata.Records) > 0 {
					log.WithFields(log.Fields{
						"domain":     domain,
						"nameserver": domainnameserver,
					}).Error("Zone can be transfered!")
				}
			}

		}
		// fin
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	// TODO: Logging
	// log.SetFormatter(&log.JSONFormatter{})
	// log.SetLevel(log.WarnLevel)
	log.SetOutput(os.Stdout)
	checkCmd.PersistentFlags().String("file", "domains.txt", "File with domain names.")
}
