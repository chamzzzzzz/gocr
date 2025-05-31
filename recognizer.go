package gocr

import (
	"context"
	"encoding/json"
	"fmt"
)

type Option struct {
	Type      string
	ID        string
	AppID     string
	AppKey    string
	AppSecret string
	AppURL    string
	Spec      json.RawMessage
}

type Document struct {
	Name string
	Path string
	Data []byte
}

type Size struct {
	Width  int
	Height int
}

type Point struct {
	X int
	Y int
}

type BoudingBox struct {
	Origin Point
	Size   Size
}

type Observation struct {
	Confidence int
	Text       string
	BoudingBox BoudingBox
}

type Result struct {
	Size         Size
	Observations []*Observation
}

type Recognizer interface {
	GetType() string
	GetID() string
	GetOption() Option
	Recognize(ctx context.Context, document *Document) (*Result, error)
}

type Creator interface {
	GetType() string
	Create(option *Option) (Recognizer, error)
}

type Workspace struct {
	creators    map[string]Creator
	recognizers map[string]Recognizer
}

func NewWorkspace() *Workspace {
	return &Workspace{
		creators:    make(map[string]Creator),
		recognizers: make(map[string]Recognizer),
	}
}

func (r *Workspace) RegisterCreator(creator Creator) error {
	r.creators[creator.GetType()] = creator
	return nil
}

func (r *Workspace) AddRecognizer(recognizer Recognizer) error {
	r.recognizers[recognizer.GetID()] = recognizer
	return nil
}

func (r *Workspace) CreateRecognizer(option *Option) (Recognizer, error) {
	creator, ok := r.creators[option.Type]
	if !ok {
		return nil, fmt.Errorf("unknown recognizer type: %s", option.Type)
	}
	recognizer, err := creator.Create(option)
	if err != nil {
		return nil, err
	}
	return recognizer, nil
}

func (r *Workspace) GetRecognizer(id string) Recognizer {
	return r.recognizers[id]
}
