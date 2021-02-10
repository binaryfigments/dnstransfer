package axfr

import (
	"net"
	"strings"

	"github.com/miekg/dns"
)

type Data struct {
	Records []dns.RR `json:"records,omitempty"`
}

// Get function, main function of this AXFR lib.
func Get(hostname string, nameserver string) (*Data, error) {
	data := new(Data)
	domain := strings.ToLower(hostname)
	/*
		domain, err := publicsuffix.EffectiveTLDPlusOne(hostname)
		if err != nil {
			return data, err
		}
	*/
	msg := new(dns.Msg)
	msg.SetAxfr(dns.Fqdn(domain))

	transfer := new(dns.Transfer)
	answerChan, err := transfer.In(msg, net.JoinHostPort(nameserver, "53"))
	if err != nil {
		return data, err
	}

	for t := range answerChan {
		if t.Error != nil {
			return data, t.Error
		}
		data.Records = t.RR
	}

	return data, err
}
