package rdap

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"unicode/utf8"
)

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

type remark struct {
	Title       string   `json:"title"`
	Description []string `json:"description"`
}

//Entity ...
type Entity struct {
	Handle          string        `json:"handle"`
	VcardArrayRaw   []interface{} `json:"vcardArray"`
	EntitiesRaw     []interface{} `json:"entities"`
	Port43          string        `json:"port43"`
	Status          []string      `json:"status"`
	Remarks         []remark      `json:"remarks"`
	ObjectClassName string        `json:"objectClassName"`
	VcardArray      []vCard
	Entities        []Entity
}

func (e *Entity) UnmarshalJSON(data []byte) error {
	var tmpEntity map[string]interface{}

	if err := json.Unmarshal(data, &tmpEntity); err != nil {
		log.Printf("%s (%T): %v\n", string(data), data, err.Error())
		return err
	}

	if _, exists := tmpEntity[`handle`]; exists {
		(*e).Handle = tmpEntity[`handle`].(string)
	}
	if _, exists := tmpEntity[`vcardArray`]; exists {
		(*e).VcardArrayRaw = tmpEntity[`vcardArray`].([]interface{})
	}
	if _, exists := tmpEntity[`entities`]; exists {
		(*e).EntitiesRaw = tmpEntity[`entities`].([]interface{})
	}
	if _, exists := tmpEntity[`port43`]; exists {
		(*e).Port43 = tmpEntity[`port43`].(string)
	}
	if _, exists := tmpEntity[`status`]; exists {
		arr := make([]string, len(tmpEntity[`status`].([]interface{})))
		for _, v := range tmpEntity[`status`].([]interface{}) {
			arr = append(arr, v.(string))
		}
		(*e).Status = arr
	}
	if _, exists := tmpEntity[`remarks`]; exists {
		(*e).Remarks = tmpEntity[`remarks`].([]remark)
	}
	if _, exists := tmpEntity[`objectClassName`]; exists {
		(*e).ObjectClassName = tmpEntity[`objectClassName`].(string)
	}

	if len((*e).VcardArrayRaw) > 0 {
		(*e).processRawVcard()
	}

	for _, entTmp := range (*e).EntitiesRaw {
		log.Println(entTmp)
	}
	return nil
}

func (e *Entity) processRawVcard() {
	if len((*e).VcardArrayRaw) < 2 {
		log.Println((*e).VcardArrayRaw)
		return
	}

	if (*e).VcardArrayRaw[0] != `vcard` {
		return
	}

	vc := make(vCard, len((*e).VcardArrayRaw))

	for _, entry := range (*e).VcardArrayRaw[1].([]interface{}) {
		k, v := processVcardEntry(entry.([]interface{}))
		vc[k] = v
	}
	(*e).VcardArray = append((*e).VcardArray, vc)
}

func fixUtf(r rune) rune {
	if r == utf8.RuneError {
		return -1
	}
	return r
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

	val := strings.Map(fixUtf, strings.Join(value, ", "))

	return key, val
}

//RDAPAnswer ...
type RDAPAnswer struct {
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
