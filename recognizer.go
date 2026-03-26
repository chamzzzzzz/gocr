package gocr

import (
	"context"
	"encoding/json"
	"fmt"
)

type Option struct {
	Type      string
	Id        string
	AppId     string
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
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}

type Point struct {
	X int `json:"x,omitempty"`
	Y int `json:"y,omitempty"`
}

type BoundingBox struct {
	Origin Point `json:"origin,omitempty"`
	Size   Size  `json:"size,omitempty"`
}

type Observation struct {
	Confidence  int          `json:"confidence,omitempty"`
	Text        string       `json:"text,omitempty"`
	BoundingBox *BoundingBox `json:"bounding_box,omitempty"`
}

type Image struct {
	File string `json:"file,omitempty"`
	*Size
}

type Result struct {
	Code         string         `json:"code,omitempty"`
	Message      string         `json:"message,omitempty"`
	Image        *Image         `json:"image,omitempty"`
	Observations []*Observation `json:"observations,omitempty"`
}

type Recognizer interface {
	GetType() string
	GetId() string
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
	r.recognizers[recognizer.GetId()] = recognizer
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
