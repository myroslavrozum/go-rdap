package gordap

//Event `event` data from RDAP answer
type Event struct {
	EventAction string   `json:"eventAction"`
	EventDate   string   `json:"eventDate"`
	Roles       []string `json:"roles"`
	Links       []Link   `json:"links"`
	Events      []Event  `json:"events"`
}
