package rdap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
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

		var ra Answer
		json.Unmarshal(body, &ra)
		for _, e := range ra.Entities {
			if len(e.VcardArrayRaw) > 0 {
				processVcardArray(e.VcardArrayRaw)
			}
		}
	}
	return net.ParseIP(ipaddr), nil
}

func processVcardEntry(vcEntry []interface{}) (string, string) {
	key := vcEntry[0].(string)
	value := make([]string, 0)

	for _, vcLine := range vcEntry[1:] {
		switch vcLine.(type) {
		case string:
			kv := vcLine.(string)
			if kv == `text` {
				continue
			}
			kv = strings.Replace(kv, "\n", ", ", -1)
			value = append(value, kv)

		case []string:
			kv := strings.Join(vcLine.([]string), " ")
			value = append(value, kv)

		case map[string]interface{}:
			kv := vcLine.(map[string]interface{})
			if len(kv) == 0 {
				continue
			}
			for _, v := range kv {
				switch v.(type) {
				case string:
					v := v.(string)
					v = strings.Replace(v, "\n", ", ", -1)
					value = append(value, v)
				case []string:
					v := strings.Join(v.([]string), " ")
					value = append(value, v)
				default:
					value = append(value, fmt.Sprint(v))
				}
			}
		default:
			value = append(value, fmt.Sprint(vcLine))

		}
	}

	return key, strings.Join(value, ", ")
}

func processVcardArray(vcArray []interface{}) {
	if len(vcArray) < 2 {
		log.Println(vcArray)
		return
	}

	if vcArray[0] != `vcard` {
		return
	}

	for _, entry := range vcArray[1].([]interface{}) {
		k, v := processVcardEntry(entry.([]interface{}))
		fmt.Printf("%-10s%-10s\n", k, v)
	}
	fmt.Println("==========")
}
