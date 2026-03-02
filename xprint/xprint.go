package xprint

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	ErrNoPrinter = errors.New("xprint: no printer specified and no default printer")
	ErrNoChrome  = errors.New("xprint: chrome/chromium not found")
	ErrNoSumatra = errors.New("xprint: sumatrapdf not found")
)

type Options struct {
	Printer  string
	Copies   int
	JobName  string
	Chrome   string
	Sumatra  string
	TempDir  string
	KeepTemp bool
}

func (o Options) withDefaults() Options {
	if o.Copies <= 0 {
		o.Copies = 1
	}
	if strings.TrimSpace(o.JobName) == "" {
		o.JobName = "print"
	}
	return o
}

func PrintRaw(data []byte, opts Options) error {
	opts = opts.withDefaults()
	if strings.TrimSpace(opts.Printer) == "" {
		p, err := DefaultPrinter()
		if err != nil {
			return err
		}
		if strings.TrimSpace(p) == "" {
			return ErrNoPrinter
		}
		opts.Printer = p
	}
	for i := 0; i < opts.Copies; i++ {
		if err := printRawPlatform(data, opts); err != nil {
			return err
		}
	}
	return nil
}

func PrintPDF(pdfPath string, opts Options) error {
	opts = opts.withDefaults()
	if strings.TrimSpace(pdfPath) == "" {
		return errors.New("xprint: empty pdfPath")
	}
	if strings.TrimSpace(opts.Printer) == "" {
		p, err := DefaultPrinter()
		if err != nil {
			return err
		}
		if strings.TrimSpace(p) == "" {
			return ErrNoPrinter
		}
		opts.Printer = p
	}
	return printPDFPlatform(pdfPath, opts)
}

func PrintURL(url string, opts Options) error {
	opts = opts.withDefaults()
	if strings.TrimSpace(url) == "" {
		return errors.New("xprint: empty url")
	}
	if strings.TrimSpace(opts.Printer) == "" {
		p, err := DefaultPrinter()
		if err != nil {
			return err
		}
		if strings.TrimSpace(p) == "" {
			return ErrNoPrinter
		}
		opts.Printer = p
	}

	pdfPath, cleanup, err := urlToTempPDF(url, opts)
	if err != nil {
		return err
	}
	if cleanup != nil && !opts.KeepTemp {
		defer cleanup()
	}
	return PrintPDF(pdfPath, opts)
}

func urlToTempPDF(url string, opts Options) (string, func(), error) {
	chrome, err := chromePath(opts)
	if err != nil {
		return "", nil, err
	}

	dir := strings.TrimSpace(opts.TempDir)
	if dir == "" {
		dir = os.TempDir()
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", nil, err
	}
	pdfPath := filepath.Join(dir, fmt.Sprintf("xprint-%s.pdf", randToken(10)))

	args := []string{
		"--headless",
		"--disable-gpu",
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-extensions",
		"--print-to-pdf=" + pdfPath,
		url,
	}
	if err := runCmd(chrome, args...); err != nil {
		return "", nil, err
	}
	return pdfPath, func() { _ = os.Remove(pdfPath) }, nil
}

func runCmd(exe string, args ...string) error {
	cmd := exec.Command(exe, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err == nil {
		return nil
	}
	s := strings.TrimSpace(out.String())
	if s == "" {
		return err
	}
	return fmt.Errorf("%w: %s", err, s)
}

func randToken(n int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	rb := make([]byte, n)
	if _, err := rand.Read(rb); err != nil {
		seed := uint64(os.Getpid())*0x9e3779b97f4a7c15 + uint64(n)*0x100000001b3
		for i := 0; i < n; i++ {
			seed ^= seed << 13
			seed ^= seed >> 7
			seed ^= seed << 17
			b[i] = alphabet[int(seed%uint64(len(alphabet)))]
		}
		return string(b)
	}
	for i := 0; i < n; i++ {
		b[i] = alphabet[int(rb[i])%len(alphabet)]
	}
	return string(b)
}
