package masker

import (
	"net/http"
	"strings"

	"autenticacion-ms/cmd/utils"
)

// MaskerHeaders is for obfuscated or masked headers
func MaskerHeaders(headers http.Header, headersToMasker ...string) http.Header {
	var headersObfuscated = headers.Clone()
	for _, key := range headersToMasker {
		value := headersObfuscated.Get(key)
		if value != "" {
			headersObfuscated.Set(key, Password(value))
		}
	}
	return headersObfuscated
}

// MaskerHeaders is for obfuscated or masked headers
func MaskerHeadersV2(headers http.Header, headersToMasker ...string) map[string][]string {

	//initialize and map header original in headersMap
	var headerResult map[string][]string
	apigeeHeaders := utils.MakeNewHeadersToCopy()
	request := &http.Request{}
	request.Header = headers
	headersMap := apigeeHeaders.GetHeadersInMap(request)

	headersObfuscated := headers.Clone()
	for _, key := range headersToMasker {
		value := headersObfuscated.Get(key)
		if value != "" {
			headersObfuscated.Set(key, Password(value))
		}
	}
	//replace http.Header to map[string][]string
	headerResult = headersObfuscated
	//validate same header without sensibility and replace with original
	for headObs, valueObs := range headerResult {
		for headMap := range headersMap {
			if strings.EqualFold(headObs, headMap) {
				delete(headerResult, headObs)
				headerResult[headMap] = []string{valueObs[0]}
			}
		}
	}
	//return object map[string][]string
	return headerResult
}
