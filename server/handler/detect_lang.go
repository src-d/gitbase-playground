package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gopkg.in/src-d/enry.v1"

	"github.com/src-d/gitbase-web/server/serializer"
)

type DetectLangRequest struct {
	Content  string `json:"content"`
	Filename string `json:"filename"`
}

// DetectLanguage returns a function that detects language by filename and content
func DetectLanguage() RequestProcessFunc {
	return func(r *http.Request) (res *serializer.Response, err error) {
		var req DetectLangRequest
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &req)
		if err != nil {
			return nil, serializer.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		lang := enry.GetLanguage(req.Filename, []byte(req.Content))
		langType := enry.GetLanguageType(lang)
		return serializer.NewDetectLangResponse(lang, langType), nil
	}
}
