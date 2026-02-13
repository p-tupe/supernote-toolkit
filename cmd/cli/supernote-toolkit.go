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
	txt := flag.Bool("txt", false, "Extract TXT")
	recurse := flag.Bool("recurse", true, "Recurse directories")
	force := flag.Bool("force", false, "Force convert all .note files")
	device := flag.String("device", "", "Chose a specific device (A5X2 | A6X2)")

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
		if *txt {
			options = append(options, i.ConvertTXT)
		}

		allNotes, err = i.FilterFreshNotes(allNotes, *output, options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	var d *i.Device
	if *device != "" {
		switch *device {
		case "A5X2":
			d = i.A5X2
		case "A6X2":
			d = i.A6X2
		default:
			log.Fatalln("Invalid device (must be one of: A5X2, A6X2):", *device)
		}
	}

	var wg sync.WaitGroup
	for _, note := range allNotes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			notebook, err := i.NewNotebook(note, d)
			if err != nil {
				log.Println(err)
				return
			}

			op := filepath.Join(*output, note.Parents)

			if *png {
				notebook.ToPNG(op)
			}

			if *pdf {
				notebook.ToPDF(op)
			}

			if *txt {
				notebook.ToTXT(op)
			}
		}()
	}
	wg.Wait()
}
