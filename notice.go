package gordap

//Notice `notice` data from RDAP answer
type Notice struct {
	Title       string   `json:"title"`
	Description []string `json:"description"`
	Links       []Link   `json:"links"`
}
