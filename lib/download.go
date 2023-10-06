package lib

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"errors"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/valyala/fasthttp"
)

// newsz, newcrc - should be
func (u *ClientInfo) DownloadAndReplaceMe(timeout time.Duration) error {
	item, err := u.info()
	if err != nil {
		return err
	}

	client := &fasthttp.Client{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		ReadTimeout:  timeout,
		WriteTimeout: time.Minute * 5,
	}

	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)

	args.Add("key", u.Key)
	args.Add("project", u.Project)
	args.Add("branch", u.Branch)

	var dst []byte
	status, buf, err := client.Post(dst, "https://"+u.Server+"//download", args)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("incorrect status: %d", status)
	}

	if len(buf) < 1000 {
		return fmt.Errorf("incorrect server's body answer: len: %d : %s", len(buf), string(buf))
	}

	b := bytes.NewReader(buf)
	g, err := gzip.NewReader(b)
	if err != nil {
		return err
	}

	cmd := os.Args[0]
	cmd = filepath.Clean(cmd)

	newbuf, err := ioutil.ReadAll(g)
	if err != nil {
		return err
	}

	crc := crc32.ChecksumIEEE(newbuf)

	if crc != uint32(item.Crc32) {
		return errors.New("crc of downloaded file is not as it should be")
	}

	if int64(len(newbuf)) != item.Size {
		return errors.New("downloaded file size is not as it should be")
	}

	err = ioutil.WriteFile(cmd+".new", newbuf, 0755)
	if err != nil {
		return err
	}

	return replaceMe(cmd+".new", cmd)
}
