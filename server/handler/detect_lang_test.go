package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pressly/lg"
	"github.com/stretchr/testify/suite"
	enry "gopkg.in/src-d/enry.v1"
)

type DetectLangSuite struct {
	HandlerUnitSuite
}

func (suite *DetectLangSuite) SetupTest() {
	h := APIHandlerFunc(DetectLanguage())
	suite.handler = lg.RequestLogger(suite.logger)(h)
}

func (suite *DetectLangSuite) TearDownTest() {}

// Tests
// -----------------------------------------------------------------------------

func TestDetectLangSuite(t *testing.T) {
	s := new(DetectLangSuite)

	suite.Run(t, s)
}

func (suite *DetectLangSuite) TestOnlyContent() {
	body := `{"content": "#!/usr/bin/env node\nconsole.log('Node')"}`
	req, _ := http.NewRequest("POST", "/detect-lang", strings.NewReader(body))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusOK, res.Code)

	lang, langType := UnmarshalDetectLangResponse(res.Body.Bytes())
	suite.Equal("JavaScript", lang)
	suite.Equal(enry.Programming, langType)
}

func (suite *DetectLangSuite) TestOnlyFilename() {
	body := `{"filename": "index.js"}`
	req, _ := http.NewRequest("POST", "/detect-lang", strings.NewReader(body))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusOK, res.Code)

	lang, langType := UnmarshalDetectLangResponse(res.Body.Bytes())
	suite.Equal("JavaScript", lang)
	suite.Equal(enry.Programming, langType)
}

func (suite *DetectLangSuite) TestDetect() {
	body := `{"filename": "foo.m", "content": "x_0=linspace(0,100,101);"}`
	req, _ := http.NewRequest("POST", "/detect-lang", strings.NewReader(body))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusOK, res.Code)

	lang, langType := UnmarshalDetectLangResponse(res.Body.Bytes())
	suite.Equal("MATLAB", lang)
	suite.Equal(enry.Programming, langType)
}

func (suite *DetectLangSuite) TestUnknownContent() {
	body := `{"content": "commit message"}`
	req, _ := http.NewRequest("POST", "/detect-lang", strings.NewReader(body))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusOK, res.Code)

	lang, langType := UnmarshalDetectLangResponse(res.Body.Bytes())
	suite.Equal("", lang)
	suite.Equal(enry.Unknown, langType)
}

func UnmarshalDetectLangResponse(b []byte) (string, enry.Type) {
	var resBody struct {
		Data struct {
			Language string `json:"language"`
			Type     int    `json:"type"`
		} `json:"data"`
	}
	json.Unmarshal(b, &resBody)
	return resBody.Data.Language, enry.Type(resBody.Data.Type)
}
