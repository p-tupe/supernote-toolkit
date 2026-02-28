package internal

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type NoteFile struct {
	// path of file on disk
	Path string
	// parent directories, if any
	Parents string
}

// GetNotesFromDir takes in a string path `dir` and returns a list of `.note` files
// in the format of [NoteFile] struct.
//
// if recurse is true, it includes the files in subfolders, and attaches the
// recursed subfolders as `Parents` in the NoteFile struct.
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

// FilterFreshNotes takes in a list of [NoteFile]s `allNotes` and a string path
// of the destination `output`. It checks if, for a give .note file, any
// corresponding files in the output exist and are newer than that file; removing
// that file from the returned list.
func FilterFreshNotes(allNotes []NoteFile, output string, convertTo []string) ([]NoteFile, error) {
	freshNotes := make([]NoteFile, 0, len(allNotes))

	pdfConvert := slices.Contains(convertTo, ConvertPDF)
	pngConvert := slices.Contains(convertTo, ConvertPNG)
	txtConvert := slices.Contains(convertTo, ExtractTXT)

	for _, note := range allNotes {
		noteFile, err := os.Stat(note.Path)
		if err != nil {
			return nil, err
		}

		if pdfConvert {
			name := strings.TrimSuffix(noteFile.Name(), filepath.Ext(noteFile.Name()))
			pdfFilePath := filepath.Join(output, note.Parents, name+".pdf")

			pdfFile, err := os.Stat(pdfFilePath)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					freshNotes = append(freshNotes, note)
					continue
				}
				return nil, err
			}

			if pdfFile.ModTime().Before(noteFile.ModTime()) {
				freshNotes = append(freshNotes, note)
				continue
			}
		}

		if pngConvert {
			name := strings.TrimSuffix(noteFile.Name(), filepath.Ext(noteFile.Name()))
			pngDirPath := filepath.Join(output, note.Parents, name)

			pngDir, err := os.Stat(pngDirPath)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					freshNotes = append(freshNotes, note)
					continue
				}
				return nil, err
			}

			if pngDir.ModTime().Before(noteFile.ModTime()) {
				freshNotes = append(freshNotes, note)
				continue
			}
		}

		if txtConvert {
			name := strings.TrimSuffix(noteFile.Name(), filepath.Ext(noteFile.Name()))
			txtFilePath := filepath.Join(output, note.Parents, name+".txt")

			txtFile, err := os.Stat(txtFilePath)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					freshNotes = append(freshNotes, note)
					continue
				}
				return nil, err
			}

			if txtFile.ModTime().Before(noteFile.ModTime()) {
				freshNotes = append(freshNotes, note)
				continue
			}

		}
	}

	return freshNotes, nil
}
