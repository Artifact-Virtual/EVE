package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello from Go!")
	fmt.Println("Working directory:", os.Getenv("PWD"))
	fmt.Println("API Key set:", os.Getenv("ANTHROPIC_API_KEY") != "")
}
