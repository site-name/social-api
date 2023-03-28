package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	entries, err := os.ReadDir("./api/schemas")
	if err != nil {
		log.Fatalln(err)
	}

	var buf bytes.Buffer

	for _, entry := range entries {
		filePath := path.Join("./api/schemas", entry.Name())
		if !strings.HasSuffix(filePath, ".graphqls") {
			continue
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalln(err)
		}

		_, err = buf.Write(data)
		if err != nil {
			log.Fatalln(err)
		}
		buf.WriteByte('\n')
	}

	file, err := os.Create("out.graphql")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = io.Copy(file, &buf)
	if err != nil {
		log.Fatalln(err)
	}
}
