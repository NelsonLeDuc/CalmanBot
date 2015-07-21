package utility

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func ParseJSON(bytes []byte, path string) string {

	var stuff interface{}

	json.Unmarshal(bytes, &stuff)
	elements := strings.Split(path, ".")

	for _, el := range elements {

		converted := ConvertedComponent(el, stuff)
		num, err := strconv.ParseInt(converted, 10, 64)

		if err == nil {
			if num < 0 {
				stuff = ""
				break
			}

			arr := stuff.([]interface{})
			stuff = arr[num]
		} else {
			switch t := stuff.(type) {
			case map[string]interface{}:
				stuff = t[converted]
			default:
				stuff = ""
				break
			}

		}
	}

	switch t := stuff.(type) {
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	default:
		return stuff.(string)
	}
}

func ConvertedComponent(s string, stuff interface{}) string {

	if s == "{_randomInt_}" {
		switch t := stuff.(type) {
		case []interface{}:
			length := len(t)
			var num int
			if length > 0 {
				rand.Seed(time.Now().UnixNano())
				num = rand.Intn(length)
			} else {
				num = -1
			}
			return strconv.Itoa(num)
		default:
			return s
		}
	}

	return s
}
