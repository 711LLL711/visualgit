package main

import (
	"flag"
)

func main() {
	var path string
	var email string
	flag.StringVar(&path, "add", "", "Path to scan")
	flag.StringVar(&email, "email", "youremail@gmail.com", "Email to send the report")
	flag.Parse()
	if path != "" {
		Scan(path)
		return
	}
	Pic()
}
