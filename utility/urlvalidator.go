package utility

import (
	"net/http"
	"net/url"
)

func ValidateURL(u string, isImage bool) bool {

	if IsValidHTTPURLString(u) {

		resp, err := http.Get(u)
		defer resp.Body.Close()

		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if isImage {
				return ValidateImage(resp.Body)
			}

			return true
		} else {
			return false
		}
	}

	return true
}

func IsValidHTTPURLString(s string) bool {
	URL, _ := url.Parse(s)
	return (URL.Scheme == "http" || URL.Scheme == "https")
}
