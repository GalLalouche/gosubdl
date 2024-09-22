package types

import (
	"encoding/json"
	"errors"
	"fmt"
)

type results struct {
	Results   []Media     `json:"results"`
	Subtitles []Subtitles `json:"subtitles"`
}

func (t *results) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	var realResults struct {
		Status    bool        `json:"status"`
		Results   []Media     `json:"results"`
		Subtitles []Subtitles `json:"subtitles"`
		Error     *string     `json:"error"`
	}
	if err := json.Unmarshal(data, &realResults); err != nil {
		return err
	}
	if !realResults.Status {
		return errors.New(fmt.Sprintf("GET failed: '%s'", *realResults.Error))
	}
	t.Results = realResults.Results
	t.Subtitles = realResults.Subtitles
	return nil
}
