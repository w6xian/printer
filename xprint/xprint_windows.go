//go:build windows
// +build windows

package xprint

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/w6xian/printer"
)

func DefaultPrinter() (string, error) {
	return printer.Default()
}

func printRawPlatform(data []byte, opts Options) error {
	p, err := printer.Open(opts.Printer)
	if err != nil {
		return err
	}
	defer p.Close()

	if err := p.StartRawDocument(opts.JobName); err != nil {
		return err
	}
	defer p.EndDocument()

	if err := p.StartPage(); err != nil {
		return err
	}
	if len(data) > 0 {
		if _, err := p.Write(data); err != nil {
			return err
		}
	}
	if err := p.EndPage(); err != nil {
		return err
	}
	return nil
}

func printPDFPlatform(pdfPath string, opts Options) error {
	sumatra, err := sumatraPath(opts)
	if err != nil {
		return err
	}
	args := []string{"-print-to", opts.Printer, "-silent"}
	if opts.Copies > 1 {
		args = append(args, "-print-settings", fmt.Sprintf("copies=%d", opts.Copies))
	}
	args = append(args, pdfPath)
	return runCmd(sumatra, args...)
}

func chromePath(opts Options) (string, error) {
	if p := strings.TrimSpace(opts.Chrome); p != "" {
		return p, nil
	}
	if p := strings.TrimSpace(os.Getenv("XPRINT_CHROME")); p != "" {
		return p, nil
	}
	for _, exe := range []string{"msedge.exe", "chrome.exe"} {
		if p, err := exec.LookPath(exe); err == nil {
			return p, nil
		}
	}
	return "", ErrNoChrome
}

func sumatraPath(opts Options) (string, error) {
	if p := strings.TrimSpace(opts.Sumatra); p != "" {
		return p, nil
	}
	if p := strings.TrimSpace(os.Getenv("XPRINT_SUMATRA")); p != "" {
		return p, nil
	}
	if p, err := exec.LookPath("SumatraPDF.exe"); err == nil {
		return p, nil
	}
	for _, p := range []string{
		filepath.Join(os.Getenv("ProgramFiles"), "SumatraPDF", "SumatraPDF.exe"),
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "SumatraPDF", "SumatraPDF.exe"),
	} {
		if strings.TrimSpace(p) == "" {
			continue
		}
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", ErrNoSumatra
}
