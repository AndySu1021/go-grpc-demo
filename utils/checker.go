package utils

import (
	"fmt"
	"os"
)

func CheckErr(err error) {
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		if err != nil {
			return
		}
	}
}
