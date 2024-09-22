package requests

import (
	"fmt"
	t "gosubdl/types"
	"io"
	"net/http"
	"os"
)

const endpoint = "https://api.subdl.com/api/v1/subtitles"

func getSdIds(r request) (<-chan []t.NameAndSdId, <-chan error) {
	rChan := make(chan []t.NameAndSdId)
	eChan := make(chan error)
	go func() {
		nr, err := http.NewRequest(http.MethodGet, endpoint, nil)

		if err != nil {
			eChan <- err
		}
		q := nr.URL.Query()
		r.update(&q)
		nr.URL.RawQuery = q.Encode()

		r, err := http.Get(nr.URL.String())
		if err != nil {
			eChan <- err
		}
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			eChan <- err
		}
		ids, err := t.GetSdIdsFromMainResult(body)
		if err != nil {
			eChan <- err
		}
		rChan <- ids
	}()
	return rChan, eChan
}

func GetMovieSdIds(fileName string) ([]t.NameAndSdId, error) {
	rChan, eChan := getSdIds(releases(fileName, t.Movie))
	select {
	case result := <-rChan:
		return result, nil
	case err := <-eChan:
		return nil, err
	}
}

func GetTvSdIds(showName string) ([]t.NameAndSdId, error) {
	rChan, eChan := getSdIds(releases(showName, t.Tv))
	select {
	case result := <-rChan:
		return result, nil
	case err := <-eChan:
		return nil, err
	}
}

func getMovieSubtitles(id t.SdId) (<-chan []t.Subtitles, <-chan error) {
	rChan := make(chan []t.Subtitles)
	eChan := make(chan error)
	go func() {
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)

		if err != nil {
			eChan <- err
		}
		q := req.URL.Query()
		r := movieSubtitles(id)
		r.update(&q)
		req.URL.RawQuery = q.Encode()

		fmt.Printf("req.URL: %s\n", req.URL)

		res, err := http.Get(req.URL.String())
		if err != nil {
			eChan <- err
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			eChan <- err
		}
		subs, err := t.SubtitlesFromResponse(body)
		if err != nil {
			eChan <- err
		} else {
			rChan <- subs
		}
	}()
	return rChan, eChan
}

func getTvSeasonSubtitles(id t.SdId, season int) (<-chan []t.Subtitles, <-chan error) {
	rChan := make(chan []t.Subtitles)
	eChan := make(chan error)
	go func() {
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)

		if err != nil {
			eChan <- err
		}
		// TODO reduce duplication (if possible...)
		q := req.URL.Query()
		r := tvSeasonSubtitles(id, season)
		r.update(&q)
		req.URL.RawQuery = q.Encode()

		res, err := http.Get(req.URL.String())
		if err != nil {
			eChan <- err
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			eChan <- err
		}
		subs, err := t.SubtitlesFromResponse(body)
		if err != nil {
			eChan <- err
			return
		}

		var ret []t.Subtitles
		for _, sub := range subs {
			isFullseason := sub.Episode == nil
			if isFullseason {
				ret = append(ret, sub)
			}
		}
		rChan <- ret
	}()
	return rChan, eChan
}

func GetMovieSubtitles(id t.SdId) ([]t.Subtitles, error) {
	rChan, eChan := getMovieSubtitles(id)
	select {
	case r := <-rChan:
		return r, nil
	case err := <-eChan:
		return nil, err
	}
}

func GetTvSeasonSubtitles(id t.SdId, season int) ([]t.Subtitles, error) {
	rChan, eChan := getTvSeasonSubtitles(id, season)
	select {
	case r := <-rChan:
		return r, nil
	case err := <-eChan:
		return nil, err
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
