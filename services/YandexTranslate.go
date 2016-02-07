package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"../config"
)

const (
	URL_ROOT       = "https://translate.yandex.net/api/v1.5/tr.json"
	LANGS_PATH     = "getLangs"
	TRANSLATE_PATH = "translate"
	DETECT_PATH    = "detect"
)

type YTranslate struct {
	apiKey string
}

type ResponseDetect struct {
	Code    int
	Message string
	Lang    string
}
type ResponseLanguages struct {
	Code    int
	Message string
	Dirs    []string
	Langs   map[string]string
}

type ResponseTranslate struct {
	Code     int
	Message  string
	Lang     string
	Text     []string
	Detected map[string]string
}

func NewYTranslate() *YTranslate {
	return &YTranslate{apiKey: config.YTranslate_API}
}

func (tr *YTranslate) Detect(text string) string {
	builtParams := url.Values{"key": {tr.apiKey}, "text": {text}}
	resp, err := http.PostForm(absUrl(DETECT_PATH), builtParams)
	if err != nil {
		fmt.Println("services/YTranslate/Detect(): troubles")
		return config.YTranslate_Fallback
	}
	defer resp.Body.Close()

	var response ResponseDetect
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("services/YTranslate/Detect(): (%v) %v", response.Code, response.Message)
		return config.YTranslate_Fallback
	}

	if response.Code != 200 {
		fmt.Println("services/YTranslate/Detect(): (%v) %v", response.Code, response.Message)
		return config.YTranslate_Fallback
	}

	return response.Lang
}

func (tr *YTranslate) GetLangs(ui string) (*ResponseLanguages, error) {
	resp, err := http.PostForm(absUrl(LANGS_PATH), url.Values{"key": {tr.apiKey}, "ui": {ui}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response ResponseLanguages
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("(%v) %v", response.Code, response.Message)
	}

	return &response, nil
}

func (tr *YTranslate) Translate(lang, text string) (*ResponseTranslate, error) {
	builtParams := url.Values{"key": {tr.apiKey}, "lang": {lang}, "text": {text}, "options": {"1"}}
	resp, err := http.PostForm(absUrl(TRANSLATE_PATH), builtParams)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response ResponseTranslate
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Code != 200 {
		return nil, fmt.Errorf("(%v) %v", response.Code, response.Message)
	}

	return &response, nil
}

func absUrl(route string) string {
	return URL_ROOT + "/" + route
}
