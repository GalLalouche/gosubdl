package requests

import (
	"fmt"
	t "gosubdl/types"
	"io"
	"net/http"
	"net/url"
	"os"
)

const endpoint = "https://api.subdl.com/api/v1/subtitles"

func getHelper[T any](
	r request,
  parseBody func(bytes []byte) ([]T, error),
) (*url.URL, <-chan []T, <-chan error) {
	rChan := make(chan []T)
	eChan := make(chan error)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)

	if err != nil {
		go func() {
			eChan <- err
		}()
		return nil, rChan, eChan
	}
	q := req.URL.Query()
	r.update(&q)
	req.URL.RawQuery = q.Encode()
	go func() {
		res, err := http.Get(req.URL.String())
		if err != nil {
			eChan <- err
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			eChan <- err
		}
		results, err := parseBody(body)
		if err != nil {
			eChan <- err
		} else {
			rChan <- results
		}
	}()
	return req.URL, rChan, eChan
}

func getSdIds(r request) (*url.URL, <-chan []t.NameAndSdId, <-chan error) {
  return getHelper(r, t.GetSdIdsFromMainResult)
}

func GetMovieSdIds(fileName string) ([]t.NameAndSdId, *url.URL, error) {
	url, rChan, eChan := getSdIds(releases(fileName, t.Movie))
	select {
	case result := <-rChan:
		return result, url, nil
	case err := <-eChan:
		return nil, url, err
	}
}

func GetTvSdIds(showName string) ([]t.NameAndSdId, *url.URL, error) {
	url, rChan, eChan := getSdIds(releases(showName, t.Tv))
	select {
	case result := <-rChan:
		return result, url, nil
	case err := <-eChan:
		return nil, url, err
	}
}

func getMovieSubtitles(id t.SdId) (*url.URL, <-chan []t.Subtitles, <-chan error) {
  return getHelper(movieSubtitles(id), t.SubtitlesFromResponse)
}

func getTvSeasonSubtitles(id t.SdId, season int) (*url.URL, <-chan []t.Subtitles, <-chan error) {
  return getHelper(tvSeasonSubtitles(id, season), t.SubtitlesFromResponse)
}

func GetMovieSubtitles(id t.SdId) ([]t.Subtitles, *url.URL, error) {
	url, rChan, eChan := getMovieSubtitles(id)
	select {
	case r := <-rChan:
		return r, url, nil
	case err := <-eChan:
		return nil, url, err
	}
}

func GetTvSeasonSubtitles(id t.SdId, season int) ([]t.Subtitles, *url.URL, error) {
	url, rChan, eChan := getTvSeasonSubtitles(id, season)
	select {
	case r := <-rChan:
		return r, url, nil
	case err := <-eChan:
		return nil, url, err
	}
}

func DownloadSubtitles(s t.Subtitles) error {
	fullUrl := s.FullUrl().String()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	target := fmt.Sprintf("%s/Downloads/%s", homeDir, s.FileName())

	fmt.Printf("Downloading subtitles %s to %s\n", fullUrl, target)
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(fullUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	written, err := io.Copy(out, resp.Body)
	fmt.Printf("Written %d bytes\n", written)
	return err
}
