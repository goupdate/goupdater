package lib

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

type serverFileInfo struct {
	Uploaded string

	Crc32 int64 //of stored body, calculated on client
	Size  int64 //in bytes
}

// true if remote version is newer than present
func (u *ClientInfo) Check() (bool, error) {
	item, err := u.info()
	if err != nil {
		return false, err
	}
	//no new version available
	if item == nil {
		return false, nil
	}

	if u.Verbose {
		fmt.Printf("got remote info: %+v\n", item)
	}

	me, err := getMyInfo()
	if err != nil {
		return false, err
	}

	newer := item.Size != me.size || item.Crc32 != me.crc

	if u.Verbose {
		if newer {
			fmt.Printf("me [%+v] differs from remote [%+v]\n", *me, item)
		} else {
			fmt.Printf("i'm same as remote: [%+v]\n", item)
		}
	}

	return newer, nil
}

func (u *ClientInfo) info() (*serverFileInfo, error) {
	client := &fasthttp.Client{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		ReadTimeout: time.Second * 5,
	}

	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)

	args.Add("key", u.Key)
	args.Add("project", u.Project)
	args.Add("branch", u.Branch)

	status, body, err := client.Post(nil, "https://"+u.Server+"//info", args)
	if err != nil {
		return nil, err
	}
	if status == http.StatusForbidden {
		return nil, errors.New("Access Forbidden. Check Params")
	}

	if status == http.StatusNotFound {
		return nil, nil //no any version available yet
	}

	if status != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("incorrect status: %d", status))
	}

	var item = &serverFileInfo{}
	err = json.Unmarshal(body, item)
	return item, err
}
