package utils

import (
	"os"
	"testing"
)

func TestPrintToPdf(t *testing.T) {
	filename := "file.pdf"
	err := ColumnPrintToPDF("Pvz6E94NYDg2JjQemzVL3rAkWQjnwp", filename, nil)

	if err != nil {
		t.Fatal("PrintToPDF test is failure", err)
	} else {
		_ = os.Remove(filename)
	}
}
