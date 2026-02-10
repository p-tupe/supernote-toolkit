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
	// txt := flag.Bool("txt", false, "Extract TXT")
	recurse := flag.Bool("recurse", true, "Recurse directories")
	force := flag.Bool("force", false, "Force convert all .note files")

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

	if !*force {
		options := []string{}
		if *png {
			options = append(options, i.ConvertPNG)
		}
		if *pdf {
			options = append(options, i.ConvertPDF)
		}

		allNotes, err = i.FilterFreshNotes(allNotes, *output, options)
		if err != nil {
			log.Fatalln(err)
		}
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

			wg.Done()
		}()
	}
	wg.Wait()
}
