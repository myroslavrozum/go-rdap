package gordap

import (
	"net"
	"strconv"
	"strings"
)

//BootstrapType type of data requested
// https://data.iana.org/rdap/
type BootstrapType int

const (
	//IPV4 Bootstrap data on IPv4 networks
	IPV4 = iota
	//IPV6 Bootstrap data on IPv6 networks
	IPV6
	//ASN Bootstrap data on Autonomous System
	ASN
	//DNS Bootstrap data on DNS names
	DNS
	//TAGS Bootstrap data on Object Tag
	TAGS
)

//Service ...
type Service [][]string

//Refs ...
type Refs struct {
	HTTP  []string
	HTTPS []string
}

//BootstrapRecord ...
type BootstrapRecord struct {
	Services    []Service `json:"services"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	Publication string    `json:"publication"`
}

//GetServices returns services array from bootstrap record
func (br *BootstrapRecord) GetServices() []Service {
	return (*br).Services
}

//GetEndpoints returns RDAP endpoints for IP address
func (br *BootstrapRecord) GetEndpoints(ip string) Refs {
	ipaddr := net.ParseIP(ip)

	maskSize := 8
	ipv4Mask := net.CIDRMask(maskSize, 32)
	networkAddress := ipaddr.Mask(ipv4Mask).String()
	var ref Refs

	//log.Println(ipaddr.DefaultMask().String())
	for _, service := range br.GetServices() {
		networks, refs := service[0], service[1]
		for _, network := range networks {
			entry := networkAddress + `/` + strconv.Itoa(maskSize)
			if network == entry {
				for _, r := range refs {
					if strings.HasPrefix(r, `http`) && !strings.HasPrefix(r, `https`) {
						ref.HTTP = append(ref.HTTP, r)
					}
					if strings.HasPrefix(r, `https`) {
						ref.HTTPS = append(ref.HTTP, r)
					}
				}
				return ref
			}
		}
	}
	return ref
}
