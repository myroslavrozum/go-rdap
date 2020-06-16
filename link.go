package gordap

//Link `link` data from RDAP answer
type Link struct {
	Value string `json:"value"`
	Rel   string `json:"rel"`
	Type  string `json:"type"`
	Href  string `json:"href"`
}
