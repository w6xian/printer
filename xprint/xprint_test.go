package xprint

import "testing"

func TestRunCmdSuccess(t *testing.T) {
	if err := runCmd("go", "env", "GOROOT"); err != nil {
		t.Fatalf("runCmd failed: %v", err)
	}
}

func TestValidateInputs(t *testing.T) {
	if err := PrintPDF("", Options{Printer: "p"}); err == nil {
		t.Fatalf("expected error for empty pdfPath")
	}
	if err := PrintURL("", Options{Printer: "p"}); err == nil {
		t.Fatalf("expected error for empty url")
	}
}
