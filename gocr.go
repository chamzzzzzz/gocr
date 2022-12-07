package gocr

import (
	"github.com/chamzzzzzz/gocr/driver"
	_ "github.com/chamzzzzzz/gocr/driver/macOCR"
)

type (
	Size        = driver.Size
	Point       = driver.Point
	BoudingBox  = driver.BoudingBox
	Observation = driver.Observation
	Image       = driver.Image
	Result      = driver.Result
)

type OCR struct {
	recognizer driver.Recognizer
}

func (ocr *OCR) Recognize(file string) (*Result, error) {
	return ocr.recognizer.Recognize(file)
}

func (ocr *OCR) Driver() driver.Driver {
	return ocr.recognizer.Driver()
}

func Open(driverName, paramName string) (*OCR, error) {
	recognizer, err := driver.Open(driverName, paramName)
	if err != nil {
		return nil, err
	}
	return &OCR{recognizer}, nil
}

func Drivers() []string {
	return driver.Drivers()
}
