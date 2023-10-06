package main

import (
	"flag"
	"fmt"

	//"goupdater"
	"github.com/goupdate/goupdater"
)

var RELEASE = "true" //set this to false via go build if dont want to run goupdater

func main() {
	fmt.Println("version 1.6")
	fmt.Printf("release: %s\n", RELEASE)

	if RELEASE=="true" {
		fmt.Println("starting updater")

		// upload to remote server?
		uploadToBranch := flag.Bool("goupload", false, "--goupload=true to upload this app to server to <branch>")
		currentBranch := flag.String("branch", "", "set remote branch to search updates")
		flag.Parse()
		flag.Usage()

		if *currentBranch == "" {
			fmt.Println("please set working branch")
			return
		}

		var update = goupdater.New("127.0.0.1:1980",
			"i can upload and check for updates with this key",
			"TestProject",
			*currentBranch)

		update.Verbose(true)

		if *uploadToBranch && *currentBranch != "" {
			err := update.Upload(*currentBranch)
			if err != nil {
				fmt.Printf("error: %s\n", err.Error())
			} else {
				fmt.Println("Upload success!")
			}
			return
		}

		// check if there any new version of this file available?
		ok, err := update.Check()
		if err != nil {
			fmt.Println("update check: " + err.Error())
			return
		}

		if ok {
			fmt.Println("new version found! download and upgrade!")
			err := update.DownloadAndReplaceMe()
			if err == nil {
				fmt.Println("i was upgraded")
				return
			} else {
				fmt.Printf("error during upgrade: %s\n", err.Error())
			}
		} else {
			fmt.Println("no new version found")
		}

	} else {
		fmt.Println("updater is OFF")
	}

	return
}
