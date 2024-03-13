package utils

import (
	"fmt"
	"strings"

	"autenticacion-ms/cmd/config/errors"
	"autenticacion-ms/cmd/entity"
)

// ValidateDocumentType validate slice document types in struct model.IDDocuments of course values allowed in countries without exceptions
func ValidateDocumentType(documentTypesAllowed string, countryExceptions string, documentType string, country string) error {
	arrDocumentTypes := strings.Split(documentTypesAllowed, ",")
	arrCountryExceptions := strings.Split(countryExceptions, ",")
	if len(arrCountryExceptions) == 0 || !Contains(arrCountryExceptions, country) {
		if !Contains(arrDocumentTypes, documentType) {
			return errors.BadRequest([]entity.Detail{{Message: fmt.Sprintf("The Document Type '%s' isn't allowed for the country '%s'", documentType, country)}}, errors.CONSUMER)
		}
	}
	return nil
}

// contains checks if a string is present in a slice
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
