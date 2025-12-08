package service

import (
	"djtracker/internal/utils"
	"fmt"
	"github.com/hcl/audioduration"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileTrack struct {
	Path string
	Name string
	Ext  string
}

// MapExtType mappe l'extension du fichier vers un entier
// représentant la valeur attendu par la librairie audioduration.Duration()
func (f *FileTrack) MapExtType() int {
	if f == nil {
		return -1
	}

	switch strings.ToLower(f.Ext) {
	case ".flac":
		return 0
	case ".mp3":
		return 2
	default:
		return -1
	}
}

// FindFile cherche le fichier name dans le dossier root
func FindFile(root, name string) (*FileTrack, error) {
	var found FileTrack
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(info.Name())
		trimName := strings.TrimSuffix(info.Name(), ext)

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(trimName), strings.ToLower(name)) {
			found = FileTrack{
				Path: path,
				Name: trimName,
				Ext:  ext,
			}
			return fs.SkipAll
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if found.Name == "" {
		return nil, fmt.Errorf("file %s not found", name)
	}

	return &found, nil
}

func (s *Service) findTrackDuration(fileTrack *FileTrack) (time.Duration, error) {
	file, err := os.Open(fileTrack.Path)
	if err != nil {
		return -1, err
	}
	defer utils.SafeDeferClose(file, s.log)

	duration, err := audioduration.Duration(file, fileTrack.MapExtType())
	if err != nil {
		return -1, err
	}

	return time.Duration(duration * float64(time.Second)), nil
}
