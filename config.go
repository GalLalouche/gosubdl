package main

import (
	t "gosubdl/types"
	"path/filepath"
)
type Config struct {
  File string `arg:""`
  Mode t.MediaType `short:"m"`
}

func (c Config) FileName() string {
  return filepath.Base(c.File)
}
