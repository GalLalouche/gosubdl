package requests

import (
	"gosubdl/common"
	t "gosubdl/types"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type request struct {
	FilmName      *string      `json:"film_name"`      // Text search by film name.
	FileName      *string      `json:"file_name"`      // Search by file name.
	SdId          *string      `json:"sd_id"`          // Search by SubDL ID.
	ImdbId        *string      `json:"imdb_id"`        // Search by IMDb ID.
	TmdbId        *string      `json:"tmdb_id"`        // Search by TMDB ID.
	SeasonNumber  *int         `json:"season_number"`  // Specific season number for TV shows.
	EpisodeNumber *int         `json:"episode_number"` // Specific episode number for TV shows.
	Type          *t.MediaType `json:"type"`           // Type of the content, either `movie` or `tv`.
	Year          *int         `json:"year"`           // Release year of the movie or TV show.
	Languages     []string     `json:"languages"`      // Comma-separated language codes for subtitle languages.
	SubsPerPage   *int         `json:"subs_per_page"`  // limit of subtitles will see in the results default is 10, (max can be 30)
	Comment       *bool        `json:"comment"`        // send comment=1 to get author comment on subtitle
	Releases      *bool        `json:"releases"`       // send releases=1 to get releases list on subtitle
	Hi            *bool        `json:"hi"`             // send hi=1 to only get hearing impaired subtitles
	// This usually doesn't work, as full season subtitles are actually marked as such
	FullSeason *bool `json:"full_season"` // send full_season=1 to get all full season subtitles
}

func makeDefault() request {
	res := request{
		Languages: []string{"EN"},
	}
	res.Hi = common.Ptr(false)
	return res
}
func releases(fileName string, typ t.MediaType) request {
	res := makeDefault()
	res.Type = common.Ptr(typ)
	res.FileName = common.Ptr(fileName)
	res.SubsPerPage = common.Ptr(0)
	return res
}
func movieSubtitles(sdId t.SdId) request {
	res := makeDefault()
	res.SdId = common.Ptr(sdId.String())
	res.Type = common.Ptr(t.Movie)
	res.SubsPerPage = common.Ptr(30)
	return res
}
func tvSeasonSubtitles(sdId t.SdId, season int) request {
	res := makeDefault()
	res.SdId = common.Ptr(sdId.String())
	res.Type = common.Ptr(t.Tv)
	res.SeasonNumber = common.Ptr(season)
	res.SubsPerPage = common.Ptr(30)
	return res
}
func seasonSubtitles(sdId t.SdId, season int) request {
	res := makeDefault()
	res.SdId = common.Ptr(sdId.String())
	res.Type = common.Ptr(t.Tv)
	res.SubsPerPage = common.Ptr(30)
	res.SeasonNumber = common.Ptr(season)
	return res
}

func (r *request) update(q *url.Values) {
	boolToParam := func(b bool) string {
		if b {
			return "1"
		}
		return "0"
	}
	q.Add("api_key", apiKey())
	if r.FilmName != nil {
		q.Add("film_name", *r.FilmName)
	}
	if r.FileName != nil {
		q.Add("file_name", *r.FileName)
	}
	if r.SdId != nil {
		q.Add("sd_id", string(*r.SdId))
	}
	if r.ImdbId != nil {
		q.Add("imdb_id", string(*r.ImdbId))
	}
	if r.TmdbId != nil {
		q.Add("tmdb_id", string(*r.TmdbId))
	}
	if r.SeasonNumber != nil {
		q.Add("season_number", strconv.Itoa(*r.SeasonNumber))
	}
	if r.EpisodeNumber != nil {
		q.Add("episode_number", strconv.Itoa(*r.EpisodeNumber))
	}
	if r.Type != nil {
		q.Add("type", r.Type.String())
	}
	if r.Year != nil {
		q.Add("year", strconv.Itoa(*r.Year))
	}
	if r.Languages != nil {
		q.Add("languages", strings.Join(r.Languages, ","))
	}
	if r.SubsPerPage != nil {
		q.Add("subs_per_page", strconv.Itoa(*r.SubsPerPage))
	}
	if r.Comment != nil {
		q.Add("comment", boolToParam(*r.Comment))
	}
	if r.Releases != nil {
		q.Add("releases", boolToParam(*r.Releases))
	}
	if r.Hi != nil {
		q.Add("hi", boolToParam(*r.Hi))
	}
	if r.FullSeason != nil {
		q.Add("full_season", boolToParam(*r.FullSeason))
	}
}

var apiKeyDuringBuild string

func apiKey() string {
	res, exists := os.LookupEnv("SUBDL_API_KEY")
	if exists {
		return res
	}
	if apiKeyDuringBuild != "" {
		return apiKeyDuringBuild
	}
	panic("Missing SUBDL_API_KEY environment variable")
}
