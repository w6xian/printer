//go:build windows
// +build windows

package xprint

import (
	"embed"
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/w6xian/printer"
)

//go:embed PDFtoPrinter.exe
var pdfPrinter embed.FS

func DefaultPrinter() (string, error) {
	return printer.Default()
}

func AllPrinters() ([]string, error) {
	return printer.All()
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
	// https://raw.githubusercontent.com/emendelson/pdftoprinter/refs/heads/main/PDFtoPrinter.exe

	// todo pdfPrinter 生成 PDFtoPrinter.exe
	// pdfPrinter 是 PDFtoPrinter.exe 的字节流 生成临时文件
	tempPath, err := os.MkdirTemp("", "xprint")
	if err != nil {
		return err
	}
	f, err := pdfPrinter.ReadFile("PDFtoPrinter.exe")
	if err != nil {
		return err
	}
	name := filepath.Join(tempPath, "PDFtoPrinter.exe")
	if err := os.WriteFile(name, f, 0755); err != nil {
		return err
	}
	// 清除临时文件夹
	defer os.RemoveAll(tempPath)
	return runCmd(name, pdfPath, opts.Printer)
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
