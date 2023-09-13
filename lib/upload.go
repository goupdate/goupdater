package lib

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

/*
upload new version

// params:
//	key
//  project
//  branch
//  size
//  crc
*/
func (u *ClientInfo) Upload(branch string) error {
	client := &fasthttp.Client{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Minute * 5,
	}

	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)

	result := "start"
	go func() {
		if u.Verbose {
			fmt.Printf("upload new version to %s - %s (%s) - %s\n", u.Server, u.Project, u.Branch, result)
		}
	}()

	me, err := getMyInfo()
	if err != nil {
		result = "get my info failed: " + err.Error()
		return err
	}

	args.Add("key", u.Key)
	args.Add("project", u.Project)
	args.Add("branch", u.Branch)
	args.Add("size", fmt.Sprintf("%d", me.size))
	args.Add("crc", fmt.Sprintf("%d", me.crc))

	body, crcb, err := getMyBody()
	if err != nil {
		result = "get my body failed: " + err.Error()
		return err
	}
	if body.Len() == 0 {
		result = "get my body failed: len is 0"
		return fmt.Errorf("len of my body is 0")
	}
	if crcb == 0 {
		result = "get my body failed: crc is 0"
		return fmt.Errorf("crc of by body is 0")
	}

	args.AddBytesV("file", body.Bytes())
	args.Add("crcb", fmt.Sprintf("%d", crcb)) //crc32 of input file

	status, buf, err := client.Post(nil, "https://"+u.Server+"//upload", args)
	if err != nil {
		result = "failed: " + err.Error()
		return err
	}
	if status != http.StatusOK {
		result = "failed: incorrect status"
		return fmt.Errorf("incorrect status: %d", status)
	}

	if string(buf) != "OK" {
		result = "failed: incorrect answer"
		return fmt.Errorf("incorrect server's body answer: %s", string(buf))
	}

	result = "OK"
	return nil
}
