package gordap

//Answer Answer from RDAP server
type Answer struct {
	RdapConformance []string `json:"rdapConformance"`
	Notices         []Notice `json:"notices"`
	Handle          string   `json:"handle"`
	StartAddress    string   `json:"startAddress"`
	EndAddress      string   `json:"endAddress"`
	IPVersion       string   `json:"ipVersion"`
	Name            string   `json:"name"`
	Type            string   `json:"type"`
	ParentHandle    string   `json:"parentHandle"`
	Events          []Event  `json:"events"`
	Links           []Link   `json:"links"`
	Entities        []Entity `json:"entities"`
	Port43          string   `json:"port43"`
	Status          []string `json:"status"`
	ObjectClassName string
	//	Cidr0Cidrs                 interface{} `json:"cidr0_cidrs"`
	//	ArinOriginas0Originautnums interface{} `json:"arin_originas0_originautnums"`
}
