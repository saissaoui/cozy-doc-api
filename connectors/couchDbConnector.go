package connectors

import (
	"bytes"
	"cozy-doc-api/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//go:generate mockery --name=CouchDbClient --output ./../mocks --case=underscore
type CouchDbClient interface {
	Ping() error
	CreateDatabase(name string) error
	InsertDocument(doc json.RawMessage, database string) error
}
type CouchDbClientImpl struct {
	httpClient     *http.Client
	basicAuthToken string
	url            string
}

type CouchHttpReponseBody struct {
	Ok     bool   `json:"ok"`
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

// NewClient initialises a new couch db client
func NewClient(cfg utils.AppConfig) (CouchDbClient, error) {
	user := cfg.CouchDbUser
	if user == "" {
		return nil, errors.New("unable to retrieve Couch Db User")
	}

	password := cfg.CouchDbPassword
	if password == "" {
		return nil, errors.New("unable to retrieve Couch Db Password")
	}

	host := cfg.CouchDbHost
	if host == "" {
		return nil, errors.New("unable to retrieve Couch Db Host")
	}

	port := cfg.CouchDbPort
	if port == 0 {
		return nil, errors.New("unable to retrieve Couch Db Port")
	}

	basicAuthToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.CouchDbUser, cfg.CouchDbPassword)))

	return &CouchDbClientImpl{
		httpClient: &http.Client{
			Timeout: 50 * time.Second,
		},
		basicAuthToken: basicAuthToken,
		url:            fmt.Sprintf("%s:%d", host, port),
	}, nil
}

// Ping tests if the connection to the database is up
func (c CouchDbClientImpl) Ping() error {
	request, err := http.NewRequest(http.MethodGet, c.url, nil)
	if err != nil {
		return errors.Wrap(err, "Unable to create ping request")
	}
	request.Header.Add("accept", "application/json")

	httpResponse, err := c.httpClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "Unable to request database")

	}
	if httpResponse.StatusCode != 200 {
		return fmt.Errorf("Ping database went wrong, status %d", httpResponse.StatusCode)
	}

	return nil
}

// CreateDatabase creates new database with given name, if the database exists the function logs a warning
func (c CouchDbClientImpl) CreateDatabase(name string) error {
	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", c.url, name), nil)
	if err != nil {
		return errors.Wrap(err, "Unable to create database creation request")
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("Authorization", "Basic "+c.basicAuthToken)

	httpResponse, error := c.httpClient.Do(request)
	if error != nil {
		return error
	}

	defer httpResponse.Body.Close()

	switch httpResponse.StatusCode {
	case 201, 202:
		utils.Logger.Info("Database creation success", zap.String("database", "name"))
		return nil
	case 412:
		utils.Logger.Warn("Database already exists", zap.String("database", "name"))
		return nil
	default:
		resp := new(CouchHttpReponseBody)
		var errorCause string
		body, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			errorCause = "unknown error, unable to read response body"

		} else {
			err := json.Unmarshal(body, &resp)
			if err != nil {
				errorCause = "Unknown error, unable unmarshal body"
			}
			errorCause = resp.Reason
		}
		return fmt.Errorf("database creation failed : status %d , error : %s", httpResponse.StatusCode, errorCause)
	}
}

func (c CouchDbClientImpl) InsertDocument(doc json.RawMessage, database string) error {
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", c.url, database), bytes.NewReader(doc))
	if err != nil {
		return errors.Wrap(err, "Unable to create insert document request")
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Basic "+c.basicAuthToken)

	httpResponse, err := c.httpClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "Unable to insert document")

	}

	if httpResponse.StatusCode != 201 && httpResponse.StatusCode != 202 {
		return fmt.Errorf("insert document went wrong with status : %d", httpResponse.StatusCode)
	}
	return nil
}
