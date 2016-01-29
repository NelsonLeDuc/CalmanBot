package utility

import (
	"net/url"
	"strings"
)

var processors []processor

func ProcessedString(input string) string {
	for _, p := range processors {
		if p.CanProcess(input) {
			input = p.Process(input)
			break
		}
	}

	return input
}

type processor interface {
	CanProcess(string) bool
	Process(string) string
}

func init() {
	processors = []processor{
		imgurProcessor{},
	}
}

//imgur
type imgurProcessor struct {
}

func (p imgurProcessor) CanProcess(str string) bool {
	URL, _ := url.Parse(str)

	if (URL.Scheme == "http" || URL.Scheme == "https") &&
		URL.Host == "imgur.com" &&
		strings.HasSuffix(URL.Path, ".gif") {
		return true
	}

	return false
}

func (p imgurProcessor) Process(str string) string {
	URL, _ := url.Parse(str)
	URL.Path += "v"

	return URL.String()
}
