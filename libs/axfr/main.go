package axfr

import (
	"net"
	"strings"

	"github.com/miekg/dns"
	"golang.org/x/net/publicsuffix"
)

// Data struct is the main struct
type Data struct {
	Domain       string     `json:"domain"`
	Records      []*records `json:"records,omitempty"`
	Error        bool       `json:"error,omitempty"`
	ErrorMessage string     `json:"errormessage,omitempty"`
}

type records struct {
	Nameserver string   `json:"host,omitempty"`
	Records    []string `json:"records,omitempty"`
	Transfer   []dns.RR `json:"transfer,omitempty"`
}

// Get function, main function of this module.
func Get(hostname string, nameserver string) *Data {
	results := new(Data)
	domain := strings.ToLower(hostname)
	domain, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		results.Error = true
		results.ErrorMessage = err.Error()
		return results
	}
	results.Domain = domain

	servers, err := net.LookupNS(domain)
	if err != nil {
		results.Error = true
		results.ErrorMessage = err.Error()
		return results
	}

	for _, server := range servers {
		r := new(records)
		r.Nameserver = server.Host

		msg := new(dns.Msg)
		msg.SetAxfr(dns.Fqdn(domain))

		transfer := new(dns.Transfer)
		answerChan, err := transfer.In(msg, net.JoinHostPort(server.Host, "53"))
		if err != nil {
			continue
		}

		for a := range answerChan {
			if a.Error != nil {
				break
			}
			r.Transfer = a.RR

			for _, rr := range a.RR {
				r.Records = append(r.Records, rr.String())
			}
		}
		results.Records = append(results.Records, r)
	}
	return results
}
