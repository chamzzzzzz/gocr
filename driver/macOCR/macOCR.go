package macOCR

import (
	"bytes"
	"fmt"
	"github.com/chamzzzzzz/gocr/driver"
	"os/exec"
	"strconv"
	"strings"
)

const (
	DriverName = "macOCR"
	ParamName  = "mac-ocr-cli"
)

type Recognizer struct {
	Path string
}

func (recognizer *Recognizer) Recognize(file string) (*driver.Result, error) {
	cmd := exec.Command(recognizer.Path, file)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	result := &driver.Result{}
	output := strings.TrimSuffix(out.String(), "\n")
	for i, line := range strings.Split(output, "\n") {
		if i == 0 {
			fields := strings.SplitN(line, " ", 2)
			if len(fields) == 2 {
				result.Image.File = fields[1]
				resolution := strings.Split(fields[0], "x")
				if len(resolution) == 2 {
					result.Image.Size.Width, _ = strconv.Atoi(resolution[0])
					result.Image.Size.Height, _ = strconv.Atoi(resolution[1])
				}
			}
		} else {
			fields := strings.SplitN(line, " ", 3)
			if len(fields) == 3 {
				observation := &driver.Observation{}
				if normalizeConfidence, err := strconv.ParseFloat(fields[0], 64); err == nil {
					observation.Confidence = int(normalizeConfidence * 100)
				}
				observation.Text = fields[2]
				boudingBox := strings.Split(strings.Trim(fields[1], "[]"), ",")
				if len(boudingBox) == 4 {
					observation.BoudingBox.Origin.X, _ = strconv.Atoi(boudingBox[0])
					observation.BoudingBox.Origin.Y, _ = strconv.Atoi(boudingBox[1])
					observation.BoudingBox.Size.Width, _ = strconv.Atoi(boudingBox[2])
					observation.BoudingBox.Size.Height, _ = strconv.Atoi(boudingBox[3])
				}
				result.Observations = append(result.Observations, observation)
			}
		}
	}
	return result, nil
}

func (recognizer *Recognizer) Driver() driver.Driver {
	return &Driver{}
}

type Driver struct {
}

func (driver *Driver) Open(paramName string) (driver.Recognizer, error) {
	if paramName == "" {
		paramName = ParamName
	}
	path, err := exec.LookPath(paramName)
	if err != nil {
		if paramName == ParamName {
			return nil, fmt.Errorf("mac-ocr-cli not installed. brew install chamzzzzzz/tap/mac-ocr-cli")
		} else {
			return nil, fmt.Errorf("%s not found.", paramName)
		}
	}
	return &Recognizer{path}, nil
}

func init() {
	driver.Register(DriverName, &Driver{})
}
