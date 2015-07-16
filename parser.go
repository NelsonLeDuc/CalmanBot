package main

import (
	"io"
    "encoding/json"
    "strings"
    "strconv"
    "math/rand"
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
    
    return stuff.(string)
}

func ParseMessageJSON(reader io.Reader) Message {
    message := new(Message)
    json.NewDecoder(reader).Decode(message)
    
    return *message
}

func ConvertedComponent(s string, stuff interface{}) string {
    
    if s == "{_randomInt_}" {
        switch t := stuff.(type) {
        case []interface{}:
            rand.Seed(time.Now().UnixNano())
            num := rand.Intn(len(t))
            return strconv.Itoa(num)
        default:
            return s
        }
    }
    
    return s
}