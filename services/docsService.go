package services

import (
	"cozy-doc-api/connectors"
	"cozy-doc-api/models"
	"cozy-doc-api/utils"

	"github.com/pkg/errors"
)

//go:generate mockery --name=DocsService --output ./../mocks --case=underscore
type DocsService interface {
	InsertDocs(req *models.DocumentRequest) error
}

type DocsServiceImpl struct {
	CouchDbClient connectors.CouchDbClient
}

// InsertDocs creates a database if it does'nt exist and inserts an array of docs
func (service DocsServiceImpl) InsertDocs(req *models.DocumentRequest) error {

	err := service.CouchDbClient.CreateDatabase(req.Database)
	if err != nil {
		return errors.Wrap(err, "InsertDocs error : CreateDatabase went wrong ")
	}

	for _, doc := range req.Documents {
		err = service.CouchDbClient.InsertDocument(doc, req.Database)
		if err != nil {
			return errors.Wrap(err, "InsertDocs error : InsertDocument went wrong ")
		}
	}

	return nil
}

func InitService(config utils.AppConfig) (*DocsServiceImpl, error) {
	client, err := connectors.NewClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create database connector")
	}

	if client != nil {
		client.Ping()
	}

	return &DocsServiceImpl{
		CouchDbClient: client,
	}, nil
}
