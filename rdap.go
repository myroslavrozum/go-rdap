package rdap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

func query(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("query() failed. HTTP error: %w", err)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

//https://rdap.arin.net/registry/ip/172.217.7.17
func rdap(ipaddr string) (net.IP, error) {
	bsr, _ := BootstrapIP(ipaddr)

	for _, endpoint := range bsr.HTTPS {
		endpoint += "ip/" + ipaddr
		body, _ := query(endpoint)

		var ra RDAPAnswer
		json.Unmarshal(body, &ra)
		for _, e := range ra.Entities {
			if e.VcardArray != nil {
				for _, vc := range e.VcardArray {
					for k, v := range vc {
						fmt.Printf("%-10v:%-10v\n", k, v)
					}
				}
				fmt.Println("=======================")
			}
		}
	}
	return net.ParseIP(ipaddr), nil
}
