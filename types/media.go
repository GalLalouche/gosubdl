package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type SdId int

func (id *SdId) String() string {
	return strconv.Itoa(int(*id))
}

type MediaType int

const (
	Tv MediaType = iota
	Movie
)

func (t *MediaType) UnmarshalText(text []byte) error {
	res, err := mediaTypeFromString(string(text))
	if err == nil {
		t = &res
	}
	return err
}
func (t MediaType) String() string {
	switch t {
	case Tv:
		return "tv"
	case Movie:
		return "movie"
	default:
		panic("Unreachable")
	}
}

func mediaTypeFromString(s string) (MediaType, error) {
	switch s {
	case "tv":
		return Tv, nil
	case "movie":
		return Movie, nil
	default:
		return 0, errors.New(fmt.Sprintf("Invalid media type: '%s'", s))
	}
}

type Media struct {
	SdId         SdId
	Typ          MediaType
	Name         string
	ImdbId       *SdId
	TmdbId       *int
	FirstAirDate *time.Time
	ReleaseDate  *time.Time
	Year         int
}

func (m Media) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Media(SdId: %d, ", m.SdId))
	builder.WriteString(fmt.Sprintf("Type: %s, ", m.Typ))
	builder.WriteString(fmt.Sprintf("Name: %q, ", m.Name))

	if m.ImdbId != nil {
		builder.WriteString(fmt.Sprintf("ImdbId: %d, ", *m.ImdbId))
	} else {
		builder.WriteString("ImdbId: nil, ")
	}

	if m.TmdbId != nil {
		builder.WriteString(fmt.Sprintf("TmdbId: %d, ", *m.TmdbId))
	} else {
		builder.WriteString("TmdbId: nil, ")
	}

	if m.FirstAirDate != nil {
		builder.WriteString(fmt.Sprintf("FirstAirDate: %s, ", m.FirstAirDate.Format("2006-01-02")))
	} else {
		builder.WriteString("FirstAirDate: nil, ")
	}

	if m.ReleaseDate != nil {
		builder.WriteString(fmt.Sprintf("ReleaseDate: %s, ", m.ReleaseDate.Format("2006-01-02")))
	} else {
		builder.WriteString("ReleaseDate: nil, ")
	}

	builder.WriteString(fmt.Sprintf("Year: %d)", m.Year))

	return builder.String()
}

func (t *Media) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	var realMedia struct {
		Sdid         int        `json:"sd_id"`
		Typ          string     `json:"type"`
		Name         string     `json:"name"`
		ImdbId       *string    `json:"imdb_id"`
		TmdbId       *int       `json:"tmdb_id"`
		FirstAirDate *time.Time `json:"first_air_date"`
		ReleaseDate  *time.Time `json:"release_date"`
		Year         int        `json:"year"`
	}
	if err := json.Unmarshal(data, &realMedia); err != nil {
		return err
	}
	typ, err := mediaTypeFromString(realMedia.Typ)
	if err != nil {
		return err
	}

	var sdId *SdId
	if realMedia.ImdbId != nil && *realMedia.ImdbId != "" {
		// prefix is "tt_"
		asInt, err := strconv.Atoi((*realMedia.ImdbId)[2:])
		if err != nil {
			return err
		}
		sdId = new(SdId)
		*sdId = SdId(asInt)
	}
	*t = Media{
		SdId:         SdId(realMedia.Sdid),
		Typ:          typ,
		Name:         realMedia.Name,
		ImdbId:       sdId,
		TmdbId:       realMedia.TmdbId,
		FirstAirDate: realMedia.FirstAirDate,
		ReleaseDate:  realMedia.ReleaseDate,
		Year:         realMedia.Year,
	}
	return nil
}

type NameAndSdId struct {
	Id   SdId
	Typ  MediaType
	Year int
	Name string
}

func GetSdIdsFromMainResult(data []byte) ([]NameAndSdId, error) {
	if string(data) == "null" || string(data) == `""` {
		return nil, errors.New(fmt.Sprintf("Invalid empty or null JSON : '%s'", string(data)))
	}
	r := new(results)
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	var results []NameAndSdId
	for _, element := range r.Results {
		results = append(results, NameAndSdId{
			Id:   element.SdId,
			Typ:  element.Typ,
			Year: element.Year,
			Name: element.Name,
		})
	}
	return results, nil
}
