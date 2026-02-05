package cli

import (
	"log"
	"os"

	i "github.com/p-tupe/supernote-toolkit/internal"
)

func Execute() {
	input := os.Args[1]

	if input == "" {
		log.Fatalln("Please add a .note file as an argument")
	}

	file, err := os.Open(input)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	notebook, err := i.NewNotebook(file)
	if err != nil {
		log.Fatalln(err)
	}

	notebook.ToPNG()
}
