package types

import (
	"encoding/json"
	"net/url"
	"strings"
)

type Subtitles struct {
	ReleaseName string `json:"release_name"`
	Url         string `json:"url"`
	Season      *int    `json:"season"`
	Episode     *int    `json:"episode"`
	Hi          bool   `json:"hi"`
}

const dlEndPoint = "https://dl.subdl.com/"

func (s *Subtitles) FullUrl() *url.URL {
	res, err := url.Parse(dlEndPoint + s.Url)
  if err != nil {
    panic(err)
  }
	return res
}

func SubtitlesFromResponse(bytes []byte) ([]Subtitles, error) {
	var res results
	if err := json.Unmarshal(bytes, &res); err != nil {
		return nil, err
	}
	return res.Subtitles, nil
}

func (s *Subtitles) FileName() string {
	u, err := url.Parse(s.Url)
	if err != nil {
		panic(err)
	}
	pos := strings.LastIndex(u.Path, "/")
	if pos == -1 {
		panic("couldn't find a period to indicate a file extension")
	}
	return u.Path[pos+1 : len(u.Path)]
}
