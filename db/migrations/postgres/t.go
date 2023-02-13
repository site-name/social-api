package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatalln(err)
	}

	// var m = map[string]struct{}{}

	for _, entry := range entries {
		if strings.Contains(entry.Name(), "_update_") {
			data, err := os.ReadFile(entry.Name())
			if err != nil {
				log.Fatalln(err)
			}

			factor := strings.Split(entry.Name(), "_")[2]

			for _, create := range entries {
				if strings.HasSuffix(create.Name(), "_create_"+factor) {
					createData, err := os.ReadFile(create.Name())
					if err != nil {
						log.Fatalln(err)
					}

					createData = append(createData, data...)

					err = os.WriteFile(create.Name(), createData, 0644)
					if err != nil {
						log.Fatalln(err)
					}

					err = os.Remove(entry.Name())
					if err != nil {
						log.Fatalln(err)
					}
				}
			}
		}
	}

}
