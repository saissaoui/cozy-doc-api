package services

import (
	"cozy-doc-api/mocks"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	docsService DocsService
	store       *mocks.CouchDbClient
}

// Run the suite
func TestServicesSuite(t *testing.T) {
	s := new(Suite)
	suite.Run(t, s)
}

// SetupTest setup suite
func (s *Suite) SetupTest() {
	s.store = new(mocks.CouchDbClient)
	s.docsService = DocsServiceImpl{
		CouchDbClient: s.store,
	}
}
