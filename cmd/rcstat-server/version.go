package main

import (
	"fmt"
)

const (
	Version = "0.0.1"
)

func PrintVersion() {
	fmt.Printf("rcstat server version: %s\n", Version)
}
