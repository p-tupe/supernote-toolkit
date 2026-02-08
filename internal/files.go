package internal

import (
	"os"
	"path/filepath"
)

type NoteFile struct {
	// path of file on disk
	Path string
	// parent directories, if any
	Parents string
}

func GetNotesFromDir(dir string, recurse bool, parents string) ([]NoteFile, error) {
	allNotes := make([]NoteFile, 0)

	allFiles, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range allFiles {
		if f.IsDir() && recurse {
			moreNotes, err := GetNotesFromDir(filepath.Join(dir, f.Name()), recurse, filepath.Join(parents, f.Name()))
			if err != nil {
				return nil, err
			}

			allNotes = append(allNotes, moreNotes...)
		}

		if f.Type().IsRegular() {
			ext := filepath.Ext(f.Name())
			if ext == ".note" {
				path := filepath.Join(dir, f.Name())
				allNotes = append(allNotes, NoteFile{path, parents})
			}
		}
	}

	return allNotes, nil
}
