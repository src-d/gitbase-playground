package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/src-d/enry.v1"

	"github.com/src-d/gitbase-web/server/handler"
	"github.com/src-d/gitbase-web/server/serializer"
	"github.com/src-d/gitbase-web/server/service"
)

type UASTGetLanguagesSuite struct {
	suite.Suite
	handler http.Handler
}

func TestUASTGetLanguagesSuite(t *testing.T) {
	if !isIntegration() {
		t.Skip("use the env var GITBASEPG_INTEGRATION_TESTS=true to run this test")
	}

	q := new(UASTGetLanguagesSuite)
	r := chi.NewRouter()
	r.Use(lg.RequestLogger(logrus.New()))
	r.Post("/detect-lang", handler.APIHandlerFunc(handler.DetectLanguage()))
	r.Get("/get-languages", handler.APIHandlerFunc(handler.GetLanguages(bblfshServerURL())))

	q.handler = r

	suite.Run(t, q)
}

func UnmarshalGetLanguagesResponse(b []byte) []service.Language {
	var resBody struct {
		Data []service.Language `json:"data"`
	}
	json.Unmarshal(b, &resBody)
	return resBody.Data
}

func (suite *UASTGetLanguagesSuite) TestSameEnryLanguage() {
	req, _ := http.NewRequest("GET", "/get-languages", strings.NewReader(""))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Require().Equal(http.StatusOK, res.Code, res.Body.String())

	escapeForJSON := func(s string) string {
		return strings.Replace(strings.Replace(s, "\"", "\\\"", -1),
			"\n", "\\n", -1)
	}

	for _, lang := range UnmarshalGetLanguagesResponse(res.Body.Bytes()) {
		langName := lang.Name
		suite.T().Run(langName, func(t *testing.T) {
			require := require.New(t)

			content, filename := suite.getContentAndFilename(langName)
			jsonRequest := fmt.Sprintf(`{ "content": "%s", "filename": "%s" }`,
				escapeForJSON(content), filename)
			req, _ := http.NewRequest("POST", "/detect-lang", strings.NewReader(jsonRequest))

			res = httptest.NewRecorder()
			suite.handler.ServeHTTP(res, req)

			require.Equal(http.StatusOK, res.Code, res.Body.String())

			detectedLang, detectedLangType := handler.UnmarshalDetectLangResponse(res.Body.Bytes())

			if langName == "Bash" {
				require.NotEqual(langName, detectedLang)
				t.Skip("TEST FAILURE IS A KNOWN ISSUE")
			}

			require.Equal(langName, detectedLang)
			require.Equal(enry.Programming, detectedLangType)
		})
	}
}

func (suite *UASTGetLanguagesSuite) getContentAndFilename(lang string) (string, string) {
	suite.T().Helper()

	switch lang {
	case "Bash":
		return "echo 'Hello World!'", "hello.sh"
	case "C#":
		return `
class HelloWorldProgram
{
    public static void Main()
    {
        System.Console.WriteLine("Hello, world!");
    }
}
`, "hello.cs"
	case "C++":
		return `
#include <iostream>

int main()
{
  std::cout << "Hello World!" << std::endl;
  return 0;
}
`, "hello.cpp"
	case "Go":
		return `
package main

import "fmt"

func main() {
    fmt.Println("Hello World!")
}
`, "hello.go"
	case "Java":
		return `
public class HelloWorld {

    public static void main(String[] args) {
        System.out.println("Hello World!");
    }

}
`, "hello.java"
	case "JavaScript":
		return "console.log('Hello World!')", "hello.js"
	case "PHP":
		return `
<?php
  echo "Hello World!";
?>
`, "hello.php"
	case "Python":
		return "print('Hello World!')", "hello.py"
	case "Ruby":
		return "puts 'Hello world!'", "hello.rb"
	}

	return "", ""
}

type UASTParseSuite struct {
	suite.Suite
	handler http.Handler
}

func TestUASTParseSuite(t *testing.T) {
	q := new(UASTParseSuite)
	q.handler = lg.RequestLogger(logrus.New())(handler.APIHandlerFunc(handler.Parse(bblfshServerURL())))

	if !isIntegration() {
		t.Skip("use the env var GITBASEPG_INTEGRATION_TESTS=true to run this test")
	}

	suite.Run(t, q)
}

func (suite *UASTParseSuite) TestSuccess() {
	jsonRequest := `{ "content": "console.log('test')", "language": "javascript" }`
	req, _ := http.NewRequest("POST", "/parse", strings.NewReader(jsonRequest))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Require().Equal(http.StatusOK, res.Code, res.Body.String())

	var resBody serializer.Response
	err := json.Unmarshal(res.Body.Bytes(), &resBody)
	suite.Nil(err)

	suite.Equal(res.Code, resBody.Status)
	suite.NotEmpty(resBody.Data)
}

func (suite *UASTParseSuite) TestError() {
	jsonRequest := `{ "content": "function(} ][", "language": "javascript" }`
	req, _ := http.NewRequest("POST", "/parse", strings.NewReader(jsonRequest))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusBadRequest, res.Code)
}

type UASTFilterSuite struct {
	suite.Suite
	handler http.Handler
}

func TestUASTFilterSuite(t *testing.T) {
	q := new(UASTFilterSuite)
	q.handler = lg.RequestLogger(logrus.New())(handler.APIHandlerFunc(handler.Filter()))

	suite.Run(t, q)
}

func (suite *UASTFilterSuite) TestSuccess() {
	jsonRequest := `{ "protobufs": "` + uastProtoMsgBase64List + `", "filter": "//*" }`
	req, _ := http.NewRequest("POST", "/filter", strings.NewReader(jsonRequest))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Require().Equal(http.StatusOK, res.Code, res.Body.String())

	var resBody serializer.Response
	err := json.Unmarshal(res.Body.Bytes(), &resBody)
	suite.Nil(err)

	suite.Equal(res.Code, resBody.Status)
	suite.NotEmpty(resBody.Data)
}

func (suite *UASTFilterSuite) TestProtobufError() {
	jsonRequest := `{ "protobufs": "not-proto", "filter": "[" }`
	req, _ := http.NewRequest("POST", "/filter", strings.NewReader(jsonRequest))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *UASTFilterSuite) TestFilterError() {
	jsonRequest := `{ "protobufs": "` + uastProtoMsgBase64List + `", "filter": "[" }`
	req, _ := http.NewRequest("POST", "/filter", strings.NewReader(jsonRequest))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusBadRequest, res.Code)
}

type UASTModeSuite struct {
	suite.Suite
	handler http.Handler
}

func TestUASTModeSuite(t *testing.T) {
	q := new(UASTModeSuite)
	q.handler = lg.RequestLogger(logrus.New())(handler.APIHandlerFunc(handler.Parse(bblfshServerURL())))

	if !isIntegration() {
		t.Skip("use the env var GITBASEPG_INTEGRATION_TESTS=true to run this test")
	}

	suite.Run(t, q)
}

func (suite *UASTModeSuite) TestSuccess() {
	testCases := []string{
		`{ "content": "console.log('test')", "language": "javascript", "mode": "" }`,
		`{ "content": "console.log('test')", "language": "javascript", "mode": "native" }`,
		`{ "content": "console.log('test')", "language": "javascript", "mode": "annotated" }`,
		`{ "content": "console.log('test')", "language": "javascript", "mode": "semantic" }`,
	}

	for _, tc := range testCases {
		suite.T().Run(tc, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/parse", strings.NewReader(tc))

			res := httptest.NewRecorder()
			suite.handler.ServeHTTP(res, req)

			suite.Require().Equal(http.StatusOK, res.Code, res.Body.String())

			var resBody serializer.Response
			err := json.Unmarshal(res.Body.Bytes(), &resBody)
			suite.Nil(err)

			suite.Equal(res.Code, resBody.Status)
			suite.NotEmpty(resBody.Data)
		})
	}
}

func (suite *UASTModeSuite) TestWrongMode() {
	jsonRequest := `{ "content": "console.log('test')", "language": "javascript", "mode": "foo" }`
	req, _ := http.NewRequest("POST", "/parse", strings.NewReader(jsonRequest))

	res := httptest.NewRecorder()
	suite.handler.ServeHTTP(res, req)

	suite.Equal(http.StatusBadRequest, res.Code)
}

// JSON: [<UAST(console.log("test"))>]
// Easy to obtain in the frontend with SELECT UAST('console.log("test")', 'JavaScript') AS uast
// Gitbase v0.18.0-beta.1, Bblfsh v2.9.2-drivers
const uastProtoMsgBase64List = "AGJncgEAAAAECFcQAQNCAQIOOgUDBAUGB0IFCBYXGBkGEgRAcG9zBxIFQHJvbGUHEgVAdHlwZQoSCGNvbW1lbnRzCRIHcHJvZ3JhbQo6AwUJCkIDCwwUBRIDZW5kBxIFc3RhcnQQEg51YXN0OlBvc2l0aW9ucww6BAUNDg9CBBAREhMFEgNjb2wGEgRsaW5lCBIGb2Zmc2V0DxINdWFzdDpQb3NpdGlvbgIgFAIgAQIgEwhCBBASEhVQDAIgAANCARcGEgRGaWxlABA6BgMEBRobHEIGCB0fIBhWBhIEYm9keQwSCmRpcmVjdGl2ZXMMEgpzb3VyY2VUeXBlA0IBHggSBk1vZHVsZQkSB1Byb2dyYW0DQgEhDDoEAwQFIkIECCMlJgwSCmV4cHJlc3Npb24DQgEkCxIJU3RhdGVtZW50FRITRXhwcmVzc2lvblN0YXRlbWVudA46BQMEBScoQgUIKSwtPAsSCWFyZ3VtZW50cwgSBmNhbGxlZQRCAiorDBIKRXhwcmVzc2lvbgYSBENhbGwQEg5DYWxsRXhwcmVzc2lvbgNCAS4OOgUDBAUvMEIFMTc5OjsIEgZGb3JtYXQHEgVWYWx1ZQdCAwsyNFAICEIEEBMSM1AMAiASCEIEEDUSNlAMAiANAiAMBEICKzgKEghBcmd1bWVudA0SC3Vhc3Q6U3RyaW5nAhIABhIEdGVzdBA6BgMEBT0+P0IGQENHSElRChIIY29tcHV0ZWQIEgZvYmplY3QKEghwcm9wZXJ0eQdCAwtBFFAICEIEEDYSQlAMAiALB0IFRCpFK0YLEglRdWFsaWZpZWQMEgpJZGVudGlmaWVyCBIGQ2FsbGVlEhIQTWVtYmVyRXhwcmVzc2lvbgIwAAo6AwMFSkIDS09QBhIETmFtZQdCAwtMFFAICEIEEE0STlAMAiAIAiAHERIPdWFzdDpJZGVudGlmaWVyCRIHY29uc29sZQdCA1JPVVBJB0IDC0FTUAgIQgQQVBJNUAwCIAkFEgNsb2cIEgZtb2R1bGU="
