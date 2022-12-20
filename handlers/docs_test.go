package handlers

import (
	"bytes"
	"cozy-doc-api/mocks"
	"cozy-doc-api/models"
	"cozy-doc-api/services"
	"cozy-doc-api/utils"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	service services.DocsService
	w       *httptest.ResponseRecorder
	ctx     *gin.Context
}

// SetupTest setup suite
func (s *Suite) SetupTest() {
	s.w = httptest.NewRecorder()
	s.fakeContext()
}

// TeardownTest teardown suite
func (s *Suite) TeardownTest() {
	s.w.Flush()
}

func (s *Suite) fakeContext() {
	s.ctx, _ = gin.CreateTestContext(s.w)
	var err error
	s.ctx.Request, err = http.NewRequest("GET", fake.DomainName(), nil)
	s.Require().NoError(err)
}

func TestHTTPSuite(t *testing.T) {
	s := new(Suite)
	suite.Run(t, s)
}
func (s Suite) TestInsertDocs_HappyPath() {
	docsService := new(mocks.DocsService)
	s.service = docsService
	req := utils.FakeDocRequest("test", utils.FakeDocs(10))
	b, err := json.Marshal(req)
	s.Require().NoError(err)
	s.ctx.Request, err = http.NewRequest("POST", fake.DomainName(), bytes.NewBuffer(b))
	s.Require().NoError(err)
	s.ctx.Request.Header.Add("content-type", "application/json")

	docsService.On("InsertDocs", req).Return(nil)
	InsertDocs(s.service)(s.ctx)
	s.Assert().Equal(http.StatusOK, s.w.Code)
}

func (s *Suite) TestGetFizzBuzz_Invalid() {
	invalidErrTest := func(tt *testing.T, name string, body interface{}, svc services.DocsService) {
		tt.Run(name, func(t *testing.T) {
			s.fakeContext()
			b, err := json.Marshal(body)
			s.Require().NoError(err)
			s.ctx.Request, err = http.NewRequest("POST", fake.DomainName(), bytes.NewBuffer(b))
			require.NoError(t, err)
			s.ctx.Request.Header.Add("content-type", "application/json")
			InsertDocs(svc)(s.ctx)
			assert.Equal(t, http.StatusBadRequest, s.w.Code)
			mock.AssertExpectationsForObjects(t, svc)
		})
	}
	docsService := new(mocks.DocsService)
	s.service = docsService

	invalidErrTest(s.T(), "InvalidBody", "invalid", s.service)

	invalidErrTest(s.T(), "InvalidRequest", &models.DocumentRequest{}, s.service)

}

func (s *Suite) TestGetFizzBuzz_Internal() {
	docsService := new(mocks.DocsService)
	s.service = docsService

	req := utils.FakeDocRequest("test", utils.FakeDocs(10))
	b, err := json.Marshal(req)
	s.Require().NoError(err)
	s.ctx.Request, err = http.NewRequest("POST", fake.DomainName(), bytes.NewBuffer(b))
	s.Require().NoError(err)
	s.ctx.Request.Header.Add("content-type", "application/json")

	docsService.On("InsertDocs", mock.Anything).Return(errors.New("internal"))

	InsertDocs(s.service)(s.ctx)
	s.Assert().Equal(http.StatusInternalServerError, s.w.Code)
	mock.AssertExpectationsForObjects(s.T(), docsService)
}
