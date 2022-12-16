package services

import (
	"cozy-doc-api/models"
	"cozy-doc-api/utils"
)

//go:generate mockery --name=DocsService --output ./../mocks --case=underscore
type DocsService interface {
	BulkInsertDocs(req *models.DocumentRequest) error
}

type DocsServiceImpl struct {
}

func (service DocsServiceImpl) BulkInsertDocs(req *models.DocumentRequest) error {
	utils.Logger.Info("Hello")
	return nil
}

func InitService(config utils.AppConfig) *DocsServiceImpl {

	return &DocsServiceImpl{}

}
