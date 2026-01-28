package cli

import (
	"log"
	"os"

	i "github.com/p-tupe/supernote-toolkit/internal"
)

func Execute() {
	input := os.Args[1]

	file, err := os.Open(input)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	_, err = i.NewNotebook(file)
	if err != nil {
		log.Fatalln(err)
	}
}
