package yagolib

import (
	"fmt"
	"os"
)

func Exit(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

/*
func StdErrPrint(a ...interface{}) {
	fmt.Fprint(os.Stderr, a)
}

func StdErrPrintln(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a)
}

func StdErrPrintf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a)
}
*/
