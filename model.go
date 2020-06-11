package rdap

//Link ...
type Link struct {
	Value string `json:"value"`
	Rel   string `json:"rel"`
	Type  string `json:"type"`
	Href  string `json:"href"`
}

//Notice ...
type Notice struct {
	Title       string   `json:"title"`
	Description []string `json:"description"`
	Links       []Link   `json:"links"`
}

//Event ...
type Event struct {
	EventAction string   `json:"eventAction"`
	EventDate   string   `json:"eventDate"`
	Roles       []string `json:"roles"`
	Links       []Link   `json:"links"`
	Events      []Event  `json:"events"`
}

//https://tools.ietf.org/html/rfc6350
//https://mariadesouza.com/2017/09/07/custom-unmarshal-json-in-golang/
type vCard map[string]string

type vCard_z struct {
	name       string
	attributes map[string]string
	text       string
}

type remark struct {
	Title       string   `json:"title"`
	Description []string `json:"description"`
}

//Entity ...
type Entity struct {
	Handle          string        `json:"handle"`
	VcardArrayRaw   []interface{} `json:"vcardArray"`
	VcardArray      []vCard
	Entities        []Entity `json:"entities"`
	Port43          string   `json:"port43"`
	Status          []string `json:"status"`
	ObjectClassName string
	Remarks         []remark `json:"remarks"`
}

//Answer ...
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
