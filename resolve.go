package resolve

import (
	"fmt"
	"github.com/google/uuid"
	"net"
)

// Resolve returns IP addresses, excluding any defaults the are returned by a DNS server.
// If no error is returned the slice shound be of non-zero length.
func Resolve(host string) ([]net.IP, error) {
	// Note some DNS servers always return a default IP address, even for domain names
	// that dont exist. We cater for this by identifying the default ip address first.
	// make a domain name that cannot exist
	// e.g. https://community.bt.com/t5/Archive-Staging/DNS-server-returns-IP-addresses-even-for-domain-names-which/td-p/163037
	badHost := uuid.New().String()
	badIps, err := net.LookupIP(badHost)
	if err != nil {
		return nil, err
	}

	bad := map[string]bool{}
	for _, ip := range badIps {
		bad[ip.String()] = true
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	var good []net.IP
	for _, ip := range ips {
		if _, ok := bad[ip.String()]; !ok {
			good = append(good, ip)
		}
	}

	if len(good) == 0 {
		return nil, fmt.Errorf("no valid ip addresses")
	}

	return good, nil
}
