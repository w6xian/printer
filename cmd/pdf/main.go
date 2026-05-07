package main

import (
	"fmt"
	"log"

	"github.com/w6xian/printer/xprint"
)

func main() {
	dn, err := xprint.DefaultPrinter()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dn)
	err = xprint.PrintPDF("sample.pdf", xprint.Options{
		Printer: dn,
		Copies:  1,
		JobName: "sample.pdf",
	})
	if err != nil {
		log.Fatal(err)
	}
}
