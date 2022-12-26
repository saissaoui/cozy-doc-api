package services

import (
	"cozy-doc-api/connectors"
	"cozy-doc-api/models"
	"cozy-doc-api/utils"
	"encoding/json"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	batchSize              = 30
	maxParalellBulkProcess = 8
)

//go:generate mockery --name=DocsService --output ./../mocks --case=underscore
type DocsService interface {
	InsertDocs(req *models.DocumentRequest) error
}

type DocsServiceImpl struct {
	CouchDbClient connectors.CouchDbClient
}

func InitService(config utils.AppConfig) (*DocsServiceImpl, error) {
	client, err := connectors.NewClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create database connector")
	}

	err = client.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to ping database ")
	}
	return &DocsServiceImpl{
		CouchDbClient: client,
	}, nil
}

func worker(wg *sync.WaitGroup, database string, service DocsServiceImpl, wrappersChan chan models.BulkDocsWrapper, errChan chan error, resultsChan chan int) {
	for w := range wrappersChan {
		err := service.CouchDbClient.BulkInsertDocuments(w, database)
		if err != nil {
			errChan <- err
		} else {
			resultsChan <- len(w.Documents)
		}
	}
	wg.Done()
}

func (service DocsServiceImpl) createWorkerPool(database string, wrappersChan chan models.BulkDocsWrapper, errChan chan error, resultsChan chan int) {
	var wg sync.WaitGroup
	for i := 0; i < maxParalellBulkProcess; i++ {
		wg.Add(1)
		go worker(&wg, database, service, wrappersChan, errChan, resultsChan)
	}
	wg.Wait()
	close(resultsChan)
}

func logResults(done chan bool, resultsChan chan int) {
	total := 0
	for result := range resultsChan {
		total += result
		utils.Logger.Info("docs inserted ", zap.Int("count", result), zap.Int("total", total))
	}
	done <- true
}

func logErrors(errChan chan error) {
	for err := range errChan {
		utils.Logger.Error("error from workers", zap.Error(err))
	}
}

// InsertDocs creates a database if it does'nt exist and inserts an array of docs
func (service DocsServiceImpl) InsertDocs(req *models.DocumentRequest) error {

	var wrappersChan = make(chan models.BulkDocsWrapper, 8)
	var resultsChan = make(chan int)
	var errChan = make(chan error)
	var done = make(chan bool)

	err := service.CouchDbClient.CreateDatabase(req.Database)
	if err != nil {
		return errors.Wrap(err, "InsertDocs error : CreateDatabase went wrong ")
	}
	if len(req.Documents) < batchSize {
		err := bulkInsertAll(req, service)
		if err != nil {
			return err
		}
	} else {
		go bulkInsertParallelized(req, wrappersChan)
		go logResults(done, resultsChan)
		go logErrors(errChan)
		service.createWorkerPool(req.Database, wrappersChan, errChan, resultsChan)
		<-done
	}

	return nil
}
func bulkInsertParallelized(req *models.DocumentRequest, wrappersChan chan models.BulkDocsWrapper) {
	docsToInsert := make([]*json.RawMessage, 0)
	for _, doc := range req.Documents {
		docsToInsert = append(docsToInsert, doc)
		if len(docsToInsert) == batchSize {
			docsWrapper := models.BulkDocsWrapper{
				Documents: docsToInsert,
			}
			wrappersChan <- docsWrapper
			docsToInsert = make([]*json.RawMessage, 0)
		}
	}
	if len(docsToInsert) > 0 {
		docsWrapper := models.BulkDocsWrapper{
			Documents: docsToInsert,
		}
		wrappersChan <- docsWrapper
	}
	close(wrappersChan)
}

func bulkInsertAll(req *models.DocumentRequest, service DocsServiceImpl) error {
	docsWrapper := models.BulkDocsWrapper{
		Documents: req.Documents,
	}
	err := service.CouchDbClient.BulkInsertDocuments(docsWrapper, req.Database)
	if err != nil {
		return errors.Wrap(err, "InsertDocs error : InsertDocument went wrong ")
	}
	return nil
}
