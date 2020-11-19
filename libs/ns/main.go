package ns

import "github.com/miekg/dns"

/*
 * Get nameservers
 * TODO: Rewrite
 */

// Get function that get the nameservers of a domain
func Get(domain string, nameserver string) ([]string, error) {
	var answer []string
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeNS)
	m.MsgHdr.RecursionDesired = true
	// m.SetEdns0(4096, true)
	c := new(dns.Client)
	in, _, err := c.Exchange(m, nameserver+":53")
	if err != nil {
		return answer, err
	}
	for _, ain := range in.Answer {
		if a, ok := ain.(*dns.NS); ok {
			answer = append(answer, a.Ns)
		}
	}
	if len(answer) < 1 {
		return answer, err
	}
	return answer, nil
}
