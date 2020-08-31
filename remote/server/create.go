package server

import (
	"os"
	"path/filepath"

	"github.com/make-os/lobe/remote/repo"
	"gopkg.in/src-d/go-git.v4"
)

// InitRepository creates a bare git repository
func (sv *Server) InitRepository(name string) error {
	return repo.InitRepository(name, sv.rootDir, sv.gitBinPath)
}

// HasRepository returns true if a valid repository exist
// for the given name
func (sv *Server) HasRepository(name string) bool {

	path := filepath.Join(sv.rootDir, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	if _, err := git.PlainOpen(path); err != nil {
		return false
	}

	return true
}
