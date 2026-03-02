//go:build !windows
// +build !windows

package printer

import (
	"errors"
	"time"
)

var errUnsupported = errors.New("printer: unsupported on this platform")

type Printer struct{}

type DriverInfo struct {
	Name        string
	Environment string
	DriverPath  string
	Attributes  uint32
}

type JobInfo struct {
	JobID           uint32
	UserMachineName string
	UserName        string
	DocumentName    string
	DataType        string
	Status          string
	StatusCode      uint32
	Priority        uint32
	Position        uint32
	TotalPages      uint32
	PagesPrinted    uint32
	Submitted       time.Time
}

func Default() (string, error) { return "", errUnsupported }

func ReadNames() ([]string, error) { return nil, errUnsupported }

func Open(name string) (*Printer, error) { return nil, errUnsupported }

func (p *Printer) Jobs() ([]JobInfo, error) { return nil, errUnsupported }

func (p *Printer) DriverInfo() (*DriverInfo, error) { return nil, errUnsupported }

func (p *Printer) StartDocument(name, datatype string) error { return errUnsupported }

func (p *Printer) StartRawDocument(name string) error { return errUnsupported }

func (p *Printer) Write(b []byte) (int, error) { return 0, errUnsupported }

func (p *Printer) Read(b []byte) (int, error) { return 0, errUnsupported }

func (p *Printer) EndDocument() error { return errUnsupported }

func (p *Printer) StartPage() error { return errUnsupported }

func (p *Printer) EndPage() error { return errUnsupported }

func (p *Printer) Close() error { return nil }
