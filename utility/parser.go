package utility

import (
	"encoding/json"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type RandomProvider func(int) int

var LinearProvider = func(n int) int {
	rand.Seed(time.Now().UnixNano())
	sqrt := 1.0 - math.Sqrt(1.0-rand.Float64())
	selection := sqrt * float64(n)
	return int(selection)
}

var UniformProvider = func(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func ParseJSON(bytes []byte, path string, rp RandomProvider) string {

	var stuff interface{}

	json.Unmarshal(bytes, &stuff)
	elements := strings.Split(path, ".")

	for _, el := range elements {

		converted := ConvertedComponent(el, stuff, rp)
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

func ConvertedComponent(s string, stuff interface{}, rp RandomProvider) string {

	if s == "{_randomInt_}" {
		switch t := stuff.(type) {
		case []interface{}:
			length := len(t)
			var num int
			if length > 0 {
				num = rp(length)
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
