package gordap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getResourcePrefix(resourceType BootstrapType) (string, error) {
	switch resourceType {
	case IPV4:
		return "ipv4.json", nil
	case IPV6:
		return "ipv6.json", nil
	case ASN:
		return "asn.json", nil
	case DNS:
		return "dns.json", nil
	case TAGS:
		return "object-tags.json", nil
	default:
		return "", errors.New("Bootstrap():  resource type not supported")
	}
}

//https://data.iana.org/rdap/
func bootstrap(resourceType BootstrapType) (*BootstrapRecord, error) {

	bootstrapResourseSuffix, err := getResourcePrefix(resourceType)

	resp, err := http.Get("https://data.iana.org/rdap/" + bootstrapResourseSuffix)
	if err != nil {
		return nil, fmt.Errorf("Bootstrap() failed. HTTP error: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Bootstrap() failed. read error: %w", err)
	}

	var records BootstrapRecord

	if err := json.Unmarshal(body, &records); err != nil {
		return nil, fmt.Errorf("Bootstrap() failed. Unmarshal error: %w", err)
	}

	return &records, nil
}

//BootstrapIP get bootstrap information for IP address
func BootstrapIP(ipaddress string) (Refs, error) {
	bsRecord, _ := bootstrap(IPV4)
	return (*bsRecord).GetEndpoints(ipaddress), nil
}
