package utils

import (
	"fmt"
	URL "net/url"
	"strings"
)

func GetUrlOperation(url string) string {
	if urlFragments := strings.Split(url, ":"); len(urlFragments) > 1 {
		return urlFragments[1]
	}
	return ""
}

func AppendUrlOperation(url string, operation string) string {
	newUrl := fmt.Sprintf("%v:%v", url, operation)
	return newUrl
}

func UrlPathFormat(url string, values ...string) string {
	urlFragments := strings.Split(url, "/")
	indexValue := 0
	for i := 0; i < len(urlFragments); i++ {
		if b := strings.Contains(urlFragments[i], "{") && strings.Contains(urlFragments[i], "}"); b {
			url = strings.ReplaceAll(url, urlFragments[i], values[indexValue])

			indexValue = indexValue + 1
		}
	}
	return url
}

func URLParse(protocol, subdomain, domain, path string) (*URL.URL, error) {
	var urlFormat string
	if subdomain != "" {
		urlFormat = fmt.Sprintf("%v://%v.%v%v", protocol, subdomain, domain, path)
	} else {
		urlFormat = fmt.Sprintf("%v://%v%v", protocol, domain, path)
	}
	url, err := URL.Parse(urlFormat)
	if err != nil {
		textErr := fmt.Sprintf("Error building path %s", path)
		return nil, fmt.Errorf(textErr)
	}

	return url, nil
}

func AddQueryParam(url *URL.URL, key string, value string) {
	q := url.Query()
	q.Set(key, value)
	url.RawQuery = q.Encode()
}
