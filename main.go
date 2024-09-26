package main

import (
	"errors"
	"fmt"
	"gosubdl/common"
	r "gosubdl/requests"
	t "gosubdl/types"
	"regexp"
	"strconv"

	"github.com/alecthomas/kong"
  "github.com/joho/godotenv"
)

func main() {
  _ = godotenv.Load() // Ignoring failed loads, since this might run from outside the source dir.
	var config Config
	kong.Parse(&config)
	switch config.Mode {
	case t.Movie:
		downloadMovieSubtitles(config.FileName())
	case t.Tv:
		downloadTvSubtitles(config.FileName())
	}
}

func downloadMovieSubtitles(fileName string) {
	downloadSubtitles(fileName, r.GetMovieSdIds, fetchMovieSubtitles)
}

func downloadTvSubtitles(fileName string) {
	season, err := extractSeason(fileName)
	if err != nil {
		panic(err)
	}
	fetchTvSubtitles := func(id t.SdId) (t.Subtitles, error) {
		subs, err := r.GetTvSeasonSubtitles(id, season)
		if err != nil {
			return t.Subtitles{}, err
		}
		return chooseSubtitles(subs)
	}
	downloadSubtitles(fileName, r.GetTvSdIds, fetchTvSubtitles)
}

func downloadSubtitles(
	fileName string,
	sdIdFetcher func(string) ([]t.NameAndSdId, error),
	subtitlesFetcher func(t.SdId) (t.Subtitles, error),
) {
	sdId, err := fetchSdId(fileName, sdIdFetcher)
	if err != nil {
		panic(err)
	}
	sub, err := subtitlesFetcher(sdId.Id)
	if err != nil {
		panic(err)
	}
	if err := r.DownloadSubtitles(sub); err != nil {
		panic(err)
	}
}

func extractSeason(file string) (int, error) {
	re := regexp.MustCompile(`(?i)S(\d\d)E\d\d`)
	return strconv.Atoi(re.FindStringSubmatch(file)[1])
}

func fetchSdId(fileName string, fetcher func(string) ([]t.NameAndSdId, error)) (t.NameAndSdId, error) {
	fmt.Printf("Fetching SD IDs for %s\n", fileName)
	sdIds, err := fetcher(fileName)
	if err != nil {
		return t.NameAndSdId{}, err
	}
	if len(sdIds) == 0 {
		return t.NameAndSdId{}, errors.New("Could not fetch any IDs")
	}
	fmt.Printf("Fetched %d SD IDs, please input a number matching the correct name\n", len(sdIds))
	printList(sdIds)
	i := readDigit(len(sdIds))
	fmt.Println("")
	return sdIds[i], nil
}

func fetchMovieSubtitles(id t.SdId) (t.Subtitles, error) {
	subs, err := r.GetMovieSubtitles(id)
	if err != nil {
		return t.Subtitles{}, err
	}
	return chooseSubtitles(subs)
}

func chooseSubtitles(subs []t.Subtitles) (t.Subtitles, error) {
	if len(subs) == 0 {
		return t.Subtitles{}, errors.New("Could not fetch any subtitles")
	}
	fmt.Printf("Fetched %d subtitles, please input a number matching the correct name\n", len(subs))
	printList(subs)
	i := readNum(len(subs))
	return subs[i], nil
}

func printList[T any](list []T) {
	for i, sub := range list {
		fmt.Printf("%d %+v\n", i, sub)
	}
}

func readNum(number int) int {
	var i int
	for _, err := fmt.Scanf("%d", &i); err != nil && i > 0 && i < number; {
		fmt.Printf("Invalid input, please type a number between 0 and %d\n", number)
	}
	return i
}

func readDigit(number int) int {
	common.AllowReadingSingleChar()
	c := common.ReadChar()
	i := int(c - '0')
	if i < 0 || i > 9 {
		fmt.Printf("Invalid input %v, please enter a number between 0 and %d\n", c, number)
		c = common.ReadChar()
		i = int(c - '0')
	}
	return i
}
