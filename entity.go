package gordap

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"unicode/utf8"
)

type remark struct {
	Title       string   `json:"title"`
	Description []string `json:"description"`
}

//https://tools.ietf.org/html/rfc6350
//https://mariadesouza.com/2017/09/07/custom-unmarshal-json-in-golang/
type vCard map[string]string

//Entity from RFC7482
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

func processUnmarshaledEntity(input *map[string]interface{}) *Entity {
	var tmpEntity Entity

	//Unmarshalling Entity field-by-field
	//handle
	if _, exists := (*input)[`handle`]; exists {
		tmpEntity.Handle = (*input)[`handle`].(string)
	}

	//vcardArray
	if _, exists := (*input)[`vcardArray`]; exists {
		tmpEntity.VcardArrayRaw = (*input)[`vcardArray`].([]interface{})
	}

	//entities - recursive
	if _, exists := (*input)[`entities`]; exists {
		for _, i := range (*input)[`entities`].([]interface{}) {
			i := i.(map[string]interface{})
			tmpEntity.Entities = append(tmpEntity.Entities, *processUnmarshaledEntity(&i))
		}
		tmpEntity.EntitiesRaw = (*input)[`entities`].([]interface{})
	}

	//port43
	if _, exists := (*input)[`port43`]; exists {
		tmpEntity.Port43 = (*input)[`port43`].(string)
	}

	//status
	if _, exists := (*input)[`status`]; exists {
		tEs := (*input)[`status`].([]interface{})
		arr := make([]string, len(tEs))
		for i := 0; i < len(tEs); i++ {
			arr[i] = tEs[i].(string)
		}
		tmpEntity.Status = arr
	}

	//remarks
	if _, exists := (*input)[`remarks`]; exists {
		for _, i := range (*input)[`remarks`].([]interface{}) {
			i := i.(map[string]interface{})
			var r remark
			if _, exists := i[`Description`]; exists {
				r.Description = i[`Description`].([]string)
			}
			if _, exists := i[`title`]; exists {
				r.Title = i[`title`].(string)
			}
			tmpEntity.Remarks = append(tmpEntity.Remarks, r)
		}
	}

	//objectClassName
	if _, exists := (*input)[`objectClassName`]; exists {
		tmpEntity.ObjectClassName = (*input)[`objectClassName`].(string)
	}

	//processRawVcard. Converting raw vCard to dictinary an assigining it ti VcardArray
	if len(tmpEntity.VcardArrayRaw) > 0 {
		tmpEntity.processRawVcard()
	}

	return &tmpEntity
}

//Clone clones `src` to `this` by value
func (e *Entity) Clone(src *Entity) {
	(*e).Handle = (*src).Handle
	(*e).VcardArrayRaw = (*src).VcardArrayRaw
	(*e).EntitiesRaw = (*src).EntitiesRaw
	(*e).Entities = (*src).Entities
	(*e).Port43 = (*src).Port43
	(*e).Status = (*src).Status
	(*e).Remarks = (*src).Remarks
	(*e).ObjectClassName = (*src).ObjectClassName
	(*e).VcardArray = (*src).VcardArray
}

//UnmarshalJSON Custom Unmarshal Entity and processing VCardRaw into VcardArray
func (e *Entity) UnmarshalJSON(data []byte) error {
	var tmpEntity map[string]interface{}

	if err := json.Unmarshal(data, &tmpEntity); err != nil {
		return err
	}
	ent := processUnmarshaledEntity(&tmpEntity)
	(*e).Clone(ent)
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
		vcEntryKey, vcEntryValue := processVcardEntry(entry.([]interface{}))
		if _, exists := vc[vcEntryKey]; exists {
			vc[vcEntryKey] += ", " + vcEntryValue
		} else {
			vc[vcEntryKey] = vcEntryValue
		}
	}
	(*e).VcardArray = append((*e).VcardArray, vc)
}
func processVcardEntry(vcEntry []interface{}) (string, string) {
	key := vcEntry[0].(string)
	propMap := vcEntry[1].(map[string]interface{})
	//propType := vcEntry[2].(string)
	propValue := vcEntry[3].(interface{})

	var valueBuilder string
	var tag string

	valueBuilder = interfaceToString(propValue)

	propMapLength := len(propMap)
	if propMapLength > 0 {
		tags := make([]string, propMapLength)
		i := 0
		for _, propMapValue := range propMap {
			tags[i] = interfaceToString(propMapValue)
		}
		tag = strings.Join(tags, ", ")
		if tag != "" {
			valueBuilder = tag + " " + valueBuilder
		}
	}

	retval := valueBuilder
	retval = strings.Replace(retval, "\n", ", ", -1)
	retval = strings.TrimSpace(retval)
	retval = strings.Map(fixUtf, retval)

	return key, retval
}

func fixUtf(r rune) rune {
	if r == utf8.RuneError {
		return -1
	}
	return r
}

func interfaceToString(input interface{}) string {
	switch input.(type) {
	case string:
		return input.(string)
	case []string:
		return strings.Join(input.([]string), ",")
	case []interface{}:
		input := input.([]interface{})
		tmp := make([]string, len(input))
		for _, value := range input {
			tmp = append(tmp, fmt.Sprint(value))
		}
		retval := strings.Join(tmp, " ")
		retval = strings.TrimSpace(retval)

		return retval
	default:
		return fmt.Sprint(input)
	}
}
