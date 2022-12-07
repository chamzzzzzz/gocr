package gocr

import (
	"os"
	"testing"
)

func TestDrivers(t *testing.T) {
	t.Log(Drivers())
}

func TestRecognize(t *testing.T) {
	dn := os.Getenv("GOCR_TEST_OCR_DN")
	pn := os.Getenv("GOCR_TEST_OCR_PN")
	file := os.Getenv("GOCR_TEST_OCR_FILE")
	ocr, err := Open(dn, pn)
	if err != nil {
		t.Error(err)
	} else {
		result, err := ocr.Recognize(file)
		if err != nil {
			t.Error(err)
		} else {
			t.Log(result.Image.File)
			t.Log(result.Image.Size)
			for _, observation := range result.Observations {
				t.Log(observation)
			}
		}
	}
}
