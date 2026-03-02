//go:build !windows
// +build !windows

package xprint

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func DefaultPrinter() (string, error) {
	out, err := combinedOutput("lpstat", "-d")
	if err != nil {
		return "", err
	}
	s := strings.TrimSpace(out)
	const prefix = "system default destination:"
	if strings.HasPrefix(strings.ToLower(s), prefix) {
		if i := strings.Index(s, ":"); i >= 0 {
			return strings.TrimSpace(s[i+1:]), nil
		}
		return strings.TrimSpace(s[len(prefix):]), nil
	}
	return "", nil
}

func printRawPlatform(data []byte, opts Options) error {
	dir := strings.TrimSpace(opts.TempDir)
	if dir == "" {
		dir = os.TempDir()
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	p := filepath.Join(dir, fmt.Sprintf("xprint-%s.bin", randToken(10)))
	if err := os.WriteFile(p, data, 0o600); err != nil {
		return err
	}
	if !opts.KeepTemp {
		defer os.Remove(p)
	}
	args := []string{"-d", opts.Printer, "-o", "raw", "-t", opts.JobName, p}
	return runCmd("lp", args...)
}

func printPDFPlatform(pdfPath string, opts Options) error {
	args := []string{"-d", opts.Printer, "-n", fmt.Sprintf("%d", opts.Copies), "-t", opts.JobName, pdfPath}
	return runCmd("lp", args...)
}

func chromePath(opts Options) (string, error) {
	if p := strings.TrimSpace(opts.Chrome); p != "" {
		return p, nil
	}
	if p := strings.TrimSpace(os.Getenv("XPRINT_CHROME")); p != "" {
		return p, nil
	}

	if runtime.GOOS == "darwin" {
		for _, p := range []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
		} {
			if _, err := os.Stat(p); err == nil {
				return p, nil
			}
		}
	}

	for _, exe := range []string{
		"google-chrome",
		"google-chrome-stable",
		"chromium",
		"chromium-browser",
		"microsoft-edge",
	} {
		if p, err := exec.LookPath(exe); err == nil {
			return p, nil
		}
	}
	return "", ErrNoChrome
}

func combinedOutput(exe string, args ...string) (string, error) {
	cmd := exec.Command(exe, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		s := strings.TrimSpace(string(out))
		if s == "" {
			return "", err
		}
		return "", fmt.Errorf("%w: %s", err, s)
	}
	return string(out), nil
}
