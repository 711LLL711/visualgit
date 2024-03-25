package main

import (
	"flag"
)

func main() {
	var path string
	var email string
	flag.StringVar(&path, "add", "", "Path to scan")
	flag.StringVar(&email, "email", "youremail@qq.com", "Email of commit")
	flag.Parse()
	if path != "" {
		Scan(path)
		return
	}
	Pic(email)
}
