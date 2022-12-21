package utils

import (
	"cozy-doc-api/models"
	"encoding/json"
	"fmt"
	"math/rand"
)

func FakeDocRequest(database string, docs []*json.RawMessage) *models.DocumentRequest {
	return &models.DocumentRequest{
		Database:  database,
		Documents: docs,
	}
}

func FakeDocs(size int) (docs []*json.RawMessage) {
	docs = make([]*json.RawMessage, 0)
	for index := 0; index < rand.Intn(5)+1; index++ {
		doc := fmt.Sprintf("{\"_id\":%[1]d,\"docVal\": \"doc_%[1]d\"}", index)
		bytes, err := json.Marshal(doc)
		if err != nil {
			panic(err)
		}
		r := json.RawMessage(bytes)
		docs = append(docs, &r)
	}
	return
}
