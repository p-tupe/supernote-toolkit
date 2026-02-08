package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"sync"

	i "github.com/p-tupe/supernote-toolkit/internal"
)

func main() {
	input := flag.String("input", "", "Folder of .note files")
	output := flag.String("output", "./output", "Folder for converted files")
	pdf := flag.Bool("pdf", true, "Convert to PDF")
	png := flag.Bool("png", false, "Convert to PNG")
	txt := flag.Bool("txt", false, "Extract to TXT")
	recurse := flag.Bool("recurse", false, "Recurse directories")
	_ = flag.Bool("force", false, "Convert every .note file (don't skip those already converted)")

	flag.Parse()

	if *input == "" {
		flag.Usage()
		fmt.Printf("\nExample: ./supernote-toolkit -input /path/to/notes\n")
		return
	}

	allNotes, err := i.GetNotesFromDir(*input, *recurse, "")
	if err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup
	for _, note := range allNotes {
		wg.Add(1)
		go func() {
			notebook, err := i.NewNotebook(note)

			if err != nil {
				log.Fatalln(err)
			}

			op := filepath.Join(*output, note.Parents)

			if *png {
				notebook.ToPNG(op)
			}

			if *pdf {
				notebook.ToPDF(op)
			}

			if *txt {
				// TODO
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
