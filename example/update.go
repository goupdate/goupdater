package main

import (
	"fmt"
	"time"

	"github.com/goupdate/goupdater"
)

var BRANCH = "main"
var PROJECT = "some-project"

func Update(uploadthis bool, exitnow chan bool) {
	var update = goupdater.New("127.0.0.1:1980",
		"key to upload and download",
		PROJECT,
		BRANCH)

	update.Verbose(true)

	if uploadthis {
		err := update.Upload(BRANCH)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
		} else {
			fmt.Println("Upload success!")
		}
		return
	}

	go func() {
		for {
			// check if there any new version of this file available?
			ok, err := update.Check()
			if err != nil {
				fmt.Printf("update check: %s", err.Error())
				return
			}

			if ok {
				fmt.Println("new version found! download and upgrade!")
				err := update.DownloadAndReplaceMe()
				if err == nil {
					fmt.Println("i was upgraded")
					exitnow <- true
					return
				} else {
					fmt.Printf("error during upgrade: %s\n", err.Error())
				}
			} else {
				fmt.Println("no new version found")
			}
			time.Sleep(time.Minute)
		}
	}()
}
