package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("API Key set:", os.Getenv("ANTHROPIC_API_KEY") != "")
	fmt.Println("API Key length:", len(os.Getenv("ANTHROPIC_API_KEY")))
}
