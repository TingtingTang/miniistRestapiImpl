package main

import (
	"fmt"
	"os"
)

// Return a unique identification with UUID format
// From Russ Cox's post
// Refer to: https://groups.google.com/forum/#!msg/golang-nuts/d0nF_k4dSx4/rPGgfXv6QCoJ
//
func Uuidv4() string {
	f, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	// uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4],
	uuid := fmt.Sprintf("%x%x%x%x%x", b[0:4],
		b[4:6],
		b[6:8],
		b[8:10],
		b[10:])

	return uuid
}
