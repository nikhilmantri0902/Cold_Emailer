package api

import (
	"encoding/json"
	"log"
)

// MarshalMetadata safely marshals a map to a JSON string
func MarshalMetadata(m map[string]interface{}) string {
	b, err := json.Marshal(m)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// LogEnrichmentSummary logs the enrichment summary
func LogEnrichmentSummary(companiesAdded, companiesSkipped, contactsAdded, contactsSkipped, errors int) {
	log.Printf("Enrichment complete: companies added=%d, skipped=%d, contacts added=%d, skipped=%d, errors=%d", companiesAdded, companiesSkipped, contactsAdded, contactsSkipped, errors)
}
