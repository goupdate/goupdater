package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	ea := flag.Bool("ea", false, "enable updates?")
	upload := flag.Bool("upload", false, "upload this?")

	flag.Parse()

	exitnow := make(chan bool)

	if *ea {
		Update(*upload, exitnow)
		if *upload {
			fmt.Println("uploaded to server")
			return
		}
	}

	ticker := time.NewTicker(time.Second)

	//do some stuff,
	for {
		select {
		case <-ticker.C:
			fmt.Println("im working...")

		case <-exitnow:
			fmt.Println("i was updated, restart")
		}
	}
}
