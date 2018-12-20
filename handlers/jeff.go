package handlers

import (
    "errors"
    "fmt"
    "golang.org/x/net/html"
	"net/http"
	"net/url"
	"encoding/json"
)

func getBody(doc *html.Node) (*html.Node, error) {
    var b *html.Node
    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "generated-text" {
					b = n
				}
			}
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
    f(doc)
    if b != nil {
        return b, nil
    }
    return nil, errors.New("Missing <generated-text> in the node tree")
}

func main() {

}

func HandleJeffFetch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("n")
	if query == "" {
		return
	}

	data := url.Values{}
	data.Set("numberOfParagraphs", query)

	resp, err := http.PostForm("http://jeffsum.com/", data)
	if err != nil {
		fmt.Println(err)
		return 
	}
	defer resp.Body.Close()

    doc, _ := html.Parse(resp.Body)
    bn, err := getBody(doc)
    if err != nil {
        return
	}

	output := ""
	start := true

	for c := bn.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "p" {
			if !start {
				output += "\n\n"
			} else {
				start = false
			}
			for t := c.FirstChild; t != nil; t = t.NextSibling {
				output += t.Data
			}
		}
	}

	jsonMap := make(map[string]string)
	jsonMap["output"] = output

	json, _ := json.Marshal(jsonMap)
	w.Write(json)
}
