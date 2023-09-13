package lib

import (
	"bytes"
	"compress/gzip"
	"errors"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
)

type myInfo struct {
	size int64
	crc  int64 //crc32
}

/*
	Collect info about myself
*/

func getMyInfo() (*myInfo, error) {
	cmd := os.Args[0]
	cmd = filepath.Clean(cmd)

	f, err := os.OpenFile(cmd, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := crc32.NewIEEE()

	n, err := io.Copy(h, f)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, errors.New("zero bytes for crc sum")
	}

	return &myInfo{
		size: n,
		crc:  int64(h.Sum32()),
	}, nil
}

// body, crc32 of it
func getMyBody() (*bytes.Buffer, int64, error) {
	cmd := os.Args[0]
	cmd = filepath.Clean(cmd)

	f, err := os.OpenFile(cmd, os.O_RDONLY, 0644)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()

	var buf []byte
	b := bytes.NewBuffer(buf)
	g := gzip.NewWriter(b)
	n, err := io.CopyBuffer(g, f, nil)
	if err != nil {
		return nil, 0, err
	}
	err = g.Close()
	if err != nil {
		return nil, 0, err
	}

	bCopy := b.Bytes()

	h := crc32.NewIEEE()
	n, err = io.Copy(h, bytes.NewReader(bCopy))
	if err != nil {
		return nil, 0, err
	}
	if n == 0 {
		return nil, 0, errors.New("zero bytes for crc sum")
	}

	return b, int64(h.Sum32()), err
}
