package util

import (
	"encoding/json"
	"fmt"
	// "os"
)

func Dd(data ...any) {
	for _, d := range data {
		b, err := json.MarshalIndent(d, "", "  ")
		if err != nil {
			fmt.Printf("Error dumping: %v\n", err)
			continue
		}
		fmt.Println(string(b))
	}
	// os.Exit(1)
}