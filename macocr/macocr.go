package macocr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/chamzzzzzz/gocr"
)

type Creator struct {
}

func (c *Creator) GetType() string {
	return "macocr"
}

func (c *Creator) Create(option *gocr.Option) (gocr.Recognizer, error) {
	p := &Recognizer{Option: *option}
	if option.Spec != nil {
		if err := json.Unmarshal(option.Spec, &p.Spec); err != nil {
			return nil, err
		}
	}
	return p, nil
}

type Spec struct {
	BinPath string
}

type Recognizer struct {
	gocr.Option
	Spec
	path string
}

func (p *Recognizer) GetType() string {
	return p.Type
}

func (p *Recognizer) GetID() string {
	return p.ID
}

func (p *Recognizer) GetOption() gocr.Option {
	return p.Option
}

func (p *Recognizer) Recognize(ctx context.Context, document *gocr.Document) (*gocr.Result, error) {
	if p.path == "" {
		path, err := p.LookupBinPath()
		if err != nil {
			return nil, err
		}
		p.path = path
	}
	cmd := exec.Command(p.path, document.Path)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	result := &gocr.Result{}
	output := strings.TrimSuffix(out.String(), "\n")
	for i, line := range strings.Split(output, "\n") {
		if i == 0 {
			fields := strings.SplitN(line, " ", 2)
			if len(fields) == 2 {
				resolution := strings.Split(fields[0], "x")
				if len(resolution) == 2 {
					result.Size.Width, _ = strconv.Atoi(resolution[0])
					result.Size.Height, _ = strconv.Atoi(resolution[1])
				}
			}
		} else {
			fields := strings.SplitN(line, " ", 3)
			if len(fields) == 3 {
				observation := &gocr.Observation{}
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

func (p *Recognizer) LookupBinPath() (path string, err error) {
	file := "mac-ocr-cli"
	if p.BinPath != "" {
		file = p.BinPath
	}
	path, err = exec.LookPath(file)
	if err != nil {
		if file == "mac-ocr-cli" {
			err = fmt.Errorf("mac-ocr-cli not installed. brew install chamzzzzzz/tap/mac-ocr-cli")
		} else {
			err = fmt.Errorf("%s not found", file)
		}
		return
	}
	return
}
