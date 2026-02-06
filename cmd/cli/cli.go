package cli

import (
	"log"
	"os"

	// "runtime/pprof"

	i "github.com/p-tupe/supernote-toolkit/internal"
)

func Execute() {
	input := os.Args[1]

	if input == "" {
		log.Fatalln("Please add a .note file as an argument")
	}

	// file, err := os.Create("cpu.pprof")
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// err = pprof.StartCPUProfile(file)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer pprof.StopCPUProfile()

	notebook, err := i.NewNotebook(input)
	if err != nil {
		log.Fatalln(err)
	}

	notebook.ToPNG()
}
