package services

import (
	"cozy-doc-api/utils"
	"math/rand"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (s *Suite) TestpsertDocs_HappyPath() {
	req := utils.FakeDocRequest("test", utils.FakeDocs(rand.Intn(5)+1))
	s.store.On("CreateDatabase", req.Database).Return(nil)
	s.store.On("BulkInsertDocuments", mock.Anything, req.Database).Return(nil)

	err := s.docsService.InsertDocs(req)
	s.Require().NoError(err)
}

func (s *Suite) TestpsertDocs_DatabaseCreationError() {
	req := utils.FakeDocRequest("test", utils.FakeDocs(rand.Intn(5)+1))
	e := errors.New("InsertDocs error : CreateDatabase went wrong")
	s.store.On("CreateDatabase", req.Database).Return(e)
	s.store.On("BulkInsertDocuments", mock.Anything, req.Database).Return(nil)

	err := s.docsService.InsertDocs(req)
	require.Error(s.T(), err)
	assert.Equal(s.T(), e, errors.Cause(err))
}

func (s *Suite) TestpsertDocs_InsertDocError() {
	req := utils.FakeDocRequest("test", utils.FakeDocs(rand.Intn(5)+1))
	e := errors.New("InsertDocs error : InsertDocument went wrong")
	s.store.On("CreateDatabase", req.Database).Return(nil)
	s.store.On("BulkInsertDocuments", mock.Anything, req.Database).Return(e)

	err := s.docsService.InsertDocs(req)
	require.Error(s.T(), err)
	assert.Equal(s.T(), e, errors.Cause(err))
}
