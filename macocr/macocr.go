package macocr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

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

func (p *Recognizer) GetId() string {
	return p.Id
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
	if err := json.Unmarshal(out.Bytes(), result); err != nil {
		return nil, err
	}
	if result.Code != "0" {
		return result, fmt.Errorf("code: %s, message: %s", result.Code, result.Message)
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
