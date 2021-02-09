package handler

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"testing/fstest"
	"time"
)

type SpaTestSuite struct {
	suite.Suite
}

type FailWriter struct {
	statusCode int
	header     http.Header
}

func (f *FailWriter) Header() http.Header {
	if f.header == nil {
		f.header = make(map[string][]string)
	}
	return f.header
}

func (f *FailWriter) WriteHeader(statusCode int) {
	f.statusCode = statusCode
}

func (f *FailWriter) Write(_ []byte) (n int, err error) {
	return 0, errors.New("this is failing")
}

func (suite *SpaTestSuite) TestServeHttp_Index() {
	mockIndex := `<html>Hello World</html>`
	mockFs := fstest.MapFS{}
	handler := SpaHandler{
		Dist:  mockFs,
		Index: mockIndex,
	}
	assert.HTTPBodyContains(suite.T(), handler.ServeHTTP, http.MethodGet, "/", nil, mockIndex, "Expected %s", mockIndex)
}

func (suite *SpaTestSuite) TestServeHttp_IndexOnNotExistingRoute() {
	mockIndex := `<html>Hello World</html>`
	mockFs := fstest.MapFS{}
	handler := SpaHandler{
		Dist:  mockFs,
		Index: mockIndex,
	}
	assert.HTTPBodyContains(suite.T(), handler.ServeHTTP, http.MethodGet, "/does/not/exist", nil, mockIndex, "Expected %s", mockIndex)
}

func (suite *SpaTestSuite) TestServeHttp_Assets() {
	mockIndex := `<html>Hello World</html>`
	mockFs := fstest.MapFS{}
	styleCss := fstest.MapFile{
		Data:    []byte("body{color: blue;}"),
		Mode:    os.ModePerm,
		ModTime: time.Now(),
		Sys:     nil,
	}
	mockFs["dist/assets/style.css"] = &styleCss
	handler := SpaHandler{
		Dist:  mockFs,
		Index: mockIndex,
	}
	assert.HTTPBodyContains(suite.T(), handler.ServeHTTP, http.MethodGet, "/assets/style.css", nil, "body{color: blue;}")
}

func (suite *SpaTestSuite) TestServeHttp_WriterFail() {
	mockIndex := `<html>Hello World</html>`
	mockFs := fstest.MapFS{}
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	failingWriter := &FailWriter{}
	handler := SpaHandler{
		Dist:  mockFs,
		Index: mockIndex,
	}
	handler.ServeHTTP(failingWriter, request)
	assert.Equal(suite.T(), http.StatusInternalServerError, failingWriter.statusCode, "Expected Internal Server Error")
}

func TestSpaHandler(t *testing.T) {
	suite.Run(t, new(SpaTestSuite))
}
