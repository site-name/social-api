package docextractor

import (
	"fmt"
	"os"
	"testing"
)

func TestPdfExtract(t *testing.T) {
	fmt.Println(os.Getwd())

	pdf := &pdfExtractor{}
	f, err := os.Open("file.pdf")
	if err != nil {
		t.Fatalf("error opening file: %v\n", err)
	}
	text, err := pdf.Extract(
		"file.pdf",
		f,
	)
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}

	fmt.Println(text)
}
