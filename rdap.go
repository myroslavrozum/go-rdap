package gordap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

//Rdap2Vcard https://rdap.arin.net/registry/ip/172.217.7.17
func Rdap2Vcard(ipaddr string) (net.IP, error) {
	bsr, _ := BootstrapIP(ipaddr)

	for _, endpoint := range bsr.HTTPS {
		endpoint += "ip/" + ipaddr
		body, _ := query(endpoint)

		var ra Answer
		json.Unmarshal(body, &ra)
		for _, e := range ra.Entities {
			printVcards(&e)
		}
	}
	return net.ParseIP(ipaddr), nil
}

func printVcards(e *Entity) {
	log.Println(e.Handle)
	if e.VcardArray != nil {
		for _, vc := range e.VcardArray {
			for k, v := range vc {
				log.Printf("%-10v:  %-10v\n", k, v)
			}
		}
		log.Println("=======================")
	}
	for _, tmpE := range e.Entities {
		printVcards(&tmpE)
	}
}
