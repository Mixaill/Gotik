package services

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type YMusicSearchRes struct {
	Albums struct {
		Items []interface{} `json:"items"`
	} `json:"albums"`
	Artists struct {
		Items []interface{} `json:"items"`
	} `json:"artists"`
	Counts struct {
		Albums  int `json:"albums"`
		Artists int `json:"artists"`
		Tracks  int `json:"tracks"`
		Videos  int `json:"videos"`
	} `json:"counts"`
	Misspell struct{} `json:"misspell"`
	Pager    struct {
		Page    int `json:"page"`
		PerPage int `json:"perPage"`
		Total   int `json:"total"`
	} `json:"pager"`
	Text   string `json:"text"`
	Tracks struct {
		Items []struct {
			Album struct {
				Artists             []interface{} `json:"artists"`
				Cover               int           `json:"cover"`
				CoverURI            string        `json:"coverUri"`
				Genre               string        `json:"genre"`
				ID                  int           `json:"id"`
				OriginalReleaseYear int           `json:"originalReleaseYear"`
				Recent              bool          `json:"recent"`
				StorageDir          string        `json:"storageDir"`
				Title               string        `json:"title"`
				TrackCount          int           `json:"trackCount"`
				VeryImportant       bool          `json:"veryImportant"`
				Year                int           `json:"year"`
			} `json:"album"`
			Albums []struct {
				Artists             []interface{} `json:"artists"`
				Cover               int           `json:"cover"`
				CoverURI            string        `json:"coverUri"`
				Genre               string        `json:"genre"`
				ID                  int           `json:"id"`
				OriginalReleaseYear int           `json:"originalReleaseYear"`
				Recent              bool          `json:"recent"`
				StorageDir          string        `json:"storageDir"`
				Title               string        `json:"title"`
				TrackCount          int           `json:"trackCount"`
				VeryImportant       bool          `json:"veryImportant"`
				Year                int           `json:"year"`
			} `json:"albums"`
			Artists []struct {
				Composer bool `json:"composer"`
				Cover    struct {
					Prefix string `json:"prefix"`
					Type   string `json:"type"`
					URI    string `json:"uri"`
				} `json:"cover"`
				Decomposed []interface{} `json:"decomposed"`
				ID         int           `json:"id"`
				Name       string        `json:"name"`
				Various    bool          `json:"various"`
			} `json:"artists"`
			Available      bool     `json:"available"`
			DurationMillis int      `json:"durationMillis"`
			DurationMs     int      `json:"durationMs"`
			Explicit       bool     `json:"explicit"`
			ID             int      `json:"id"`
			Regions        []string `json:"regions"`
			StorageDir     string   `json:"storageDir"`
			Title          string   `json:"title"`
		} `json:"items"`
		PerPage int `json:"perPage"`
		Total   int `json:"total"`
	} `json:"tracks"`
	Videos struct {
		Items []interface{} `json:"items"`
	} `json:"videos"`
}

type YMusicFilenameRes struct {
	XMLName  xml.Name `xml:"track"`
	Filename string   `xml:"filename,attr"`
	Length   string   `xml:"track-length,attr"`
}

type YMusicDownloadInfoRes struct {
	XMLName xml.Name `xml:"download-info"`
	Host    string   `xml:"host"`
	Path    string   `xml:"path"`
	Ts      string   `xml:"ts"`
	Region  int      `xml:"region"`
	S       string   `xml:"s"`
}

type YMusic struct{}

func (ym *YMusic) makeRequest(str string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", str, nil)
	if err != nil {
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:39.0) Gecko/20100101 Firefox/39.0")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	return body
}

func (ym *YMusic) getSearchResult(str string) YMusicSearchRes {
	m := YMusicSearchRes{}

	request := "https://music.yandex.by/handlers/search.jsx" +
		"?text=" + url.QueryEscape(str) +
		"&type=tracks" +
		"&lang=ru" +
		"&external-domain=music.yandex.by" +
		"&overembed=false"

	response := ym.makeRequest(request)
	if response != nil {
		json.Unmarshal(response, &m)
	}
	return m
}

func (ym *YMusic) getFileName(sr YMusicSearchRes) YMusicFilenameRes {
	m := YMusicFilenameRes{}
	if sr.Tracks.Total == 0 {
		return m
	}

	request := "http://storage.music.yandex.ru/get/" + sr.Tracks.Items[0].StorageDir + "/2.xml"
	response := ym.makeRequest(request)

	if response != nil {
		xml.Unmarshal(response, &m)
	}
	return m
}

func (ym *YMusic) getFileDInfo(sr YMusicSearchRes, fn YMusicFilenameRes) YMusicDownloadInfoRes {
	m := YMusicDownloadInfoRes{}
	if sr.Tracks.Total == 0 || fn.Filename == "" {
		return m
	}
	request := "http://storage.music.yandex.ru/download-info/" + sr.Tracks.Items[0].StorageDir + "/" + fn.Filename
	response := ym.makeRequest(request)

	if response != nil {
		xml.Unmarshal(response, &m)
	}
	return m
}

func (ym *YMusic) GetTrack(text string) (io.ReadCloser, string) {
	var title string
	search := ym.getSearchResult(text)
	filename := ym.getFileName(search)
	dinfo := ym.getFileDInfo(search, filename)

	if dinfo.Host != "" && dinfo.Path != "" {
		title = search.Tracks.Items[0].Artists[0].Name + " - " + search.Tracks.Items[0].Title
		skey := md5.New()
		io.WriteString(skey, "XGRlBW9FXlekgbPrRHuSiA"+dinfo.Path[1:]+dinfo.S)
		request := "http://" + dinfo.Host + "/get-mp3/" + hex.EncodeToString(skey.Sum(nil)) + "/"
		request += dinfo.Ts + dinfo.Path + "?track-id=" + strconv.Itoa(search.Tracks.Items[0].ID) + "&region=" + strconv.Itoa(dinfo.Region) + "&from=service-search"

		fmt.Println(request)
		res, err := http.Get(request)
		if err != nil {
			fmt.Println(err)
		}

		return res.Body, title
	}
	return nil, title
}
