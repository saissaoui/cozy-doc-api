package models

import "encoding/json"

type BulkDocsWrapper struct {
	Documents []*json.RawMessage `json:"docs"`
}

type CouchHttpReponseBody struct {
	Ok     bool   `json:"ok"`
	Error  string `json:"error"`
	Reason string `json:"reason"`
}
