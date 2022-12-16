package utils

import (
	"cozy-doc-api/models"
	"encoding/json"
	"fmt"
)

func FakeDocRequest(database string, docs []json.RawMessage) *models.DocumentRequest {
	return &models.DocumentRequest{
		Database:  database,
		Documents: docs,
	}
}

func FakeDocs(size int) (docs []json.RawMessage) {
	docs = make([]json.RawMessage, 10)
	for index := 0; index < size; index++ {
		doc := fmt.Sprintf("{\"docID\":%[1]d,\"docVal\": \"doc_%[1]d\"}", index)
		bytes, err := json.Marshal(doc)
		if err != nil {
			panic(err)
		}
		docs[index] = bytes
	}
	return
}
