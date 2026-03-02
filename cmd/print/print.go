// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// print command prints text documents to selected printer.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/w6xian/printer"
	"github.com/w6xian/printer/xprint"
)

var (
	copies    = flag.Int("n", 1, "number of copies to print")
	printerId = flag.String("p", findDefaultPrinter(), "printer name or printer index from printer list")
	doList    = flag.Bool("l", false, "list printers")
)

func findDefaultPrinter() string {
	p, err := xprint.DefaultPrinter()
	if err != nil {
		return ""
	}
	return p
}

func listPrinters() error {
	defaultPrinter, _ := xprint.DefaultPrinter()
	if runtime.GOOS == "windows" {
		printers, err := printer.ReadNames()
		if err != nil {
			return err
		}
		for i, p := range printers {
			s := " "
			if p == defaultPrinter {
				s = "*"
			}
			fmt.Printf(" %s %d. %s\n", s, i, p)
		}
		return nil
	}

	out, err := exec.Command("lpstat", "-p").CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	lines := strings.Split(string(out), "\n")
	printers := make([]string, 0, len(lines))
	for _, line := range lines {
		f := strings.Fields(line)
		if len(f) >= 2 && strings.ToLower(f[0]) == "printer" {
			printers = append(printers, f[1])
		}
	}
	for i, p := range printers {
		s := " "
		if p == defaultPrinter {
			s = "*"
		}
		fmt.Printf(" %s %d. %s\n", s, i, p)
	}
	return nil
}

func selectPrinter() (string, error) {
	n, err := strconv.Atoi(*printerId)
	if err != nil {
		// must be a printer name
		return *printerId, nil
	}
	var printers []string
	if runtime.GOOS == "windows" {
		printers, err = printer.ReadNames()
		if err != nil {
			return "", err
		}
	} else {
		out, err := exec.Command("lpstat", "-p").CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			f := strings.Fields(line)
			if len(f) >= 2 && strings.ToLower(f[0]) == "printer" {
				printers = append(printers, f[1])
			}
		}
	}
	if n < 0 {
		return "", fmt.Errorf("printer index (%d) cannot be negative", n)
	}
	if n >= len(printers) {
		return "", fmt.Errorf("printer index (%d) is too large, there are only %d printers", n, len(printers))
	}
	return printers[n], nil
}

func printTextFile(printerName, documentName, path string) error {
	output, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(strings.ReplaceAll(string(output), "\r\n", "\n"), "\n")
	var buf bytes.Buffer
	for _, line := range lines {
		if line == "" {
			buf.WriteString("\r\n")
			continue
		}
		buf.WriteString(line)
		buf.WriteString("\r\n")
	}
	return xprint.PrintRaw(buf.Bytes(), xprint.Options{
		Printer: printerName,
		Copies:  *copies,
		JobName: documentName,
	})
}

func printDocument(path string) error {
	if *copies < 0 {
		return fmt.Errorf("number of copies to print (%d) cannot be negative", *copies)
	}

	printerName, err := selectPrinter()
	if err != nil {
		return err
	}

	documentName := filepath.Base(path)
	if strings.HasPrefix(strings.ToLower(path), "http://") || strings.HasPrefix(strings.ToLower(path), "https://") {
		return xprint.PrintURL(path, xprint.Options{
			Printer: printerName,
			Copies:  *copies,
			JobName: documentName,
		})
	}
	if strings.HasSuffix(strings.ToLower(path), ".pdf") {
		return xprint.PrintPDF(path, xprint.Options{
			Printer: printerName,
			Copies:  *copies,
			JobName: documentName,
		})
	}
	return printTextFile(printerName, documentName, path)
}

func usage() {
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "usage: print [-n=<copies>] [-p=<printer>] <file-path|pdf|url>\n")
	fmt.Fprintf(os.Stderr, "       or\n")
	fmt.Fprintf(os.Stderr, "       print -l\n")
	fmt.Fprintln(os.Stderr)
	flag.PrintDefaults()
	os.Exit(1)
}

func exit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if *doList {
		exit(listPrinters())
	}
	switch len(flag.Args()) {
	case 0:
		fmt.Fprintf(os.Stderr, "no document path to print provided\n")
	case 1:
		exit(printDocument(flag.Arg(0)))
	default:
		fmt.Fprintf(os.Stderr, "too many parameters provided\n")
	}
	usage()
}
