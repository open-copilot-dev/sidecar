package domain

import (
	"os"
	"path/filepath"
)

var BaseDir string

func init() {
	BaseDir = filepath.Dir(os.Args[0])
}
