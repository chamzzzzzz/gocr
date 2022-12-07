package driver

import (
	"fmt"
	"sort"
	"sync"
)

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

type Image struct {
	File string
	Size Size
}

type Result struct {
	Image        Image
	Observations []*Observation
}

type Recognizer interface {
	Recognize(file string) (*Result, error)
	Driver() Driver
}

type Driver interface {
	Open(paramName string) (Recognizer, error)
}

var (
	drivers = make(map[string]Driver)
	mu      sync.RWMutex
)

func Register(name string, driver Driver) {
	mu.Lock()
	defer mu.Unlock()
	if driver == nil {
		panic("gocr/driver: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("gocr/driver: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func Drivers() []string {
	mu.RLock()
	defer mu.RUnlock()
	list := make([]string, 0, len(drivers))
	for name := range drivers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

func Open(name, paramName string) (Recognizer, error) {
	mu.RLock()
	driver, ok := drivers[name]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown driver %q (forgotten import?)", name)
	}
	return driver.Open(paramName)
}
