package web

import (
	"errors"
	"fmt"
	"hash/crc32"
	"net/http"
	"strconv"

	"goupdater/server/internal/log"
	"goupdater/server/internal/storage"

	"github.com/valyala/fasthttp"
)

// uploads new file for project-branch to storage
// method: POST
// params:
//	key
//  project
//  branch
//  size -- original
//  crc -- original crc32
//  file -- new body
//  crcb -- new body crc32
func UploadNewFile(ctx *fasthttp.RequestCtx) {
	project, branch, allowed := IsAllowed(ctx)
	if !allowed {
		ctx.Error("Access denied", 403)
		return
	}

	result := "OK"
	defer log.Action(ctx.RemoteIP().String(), "upload", project, branch, result)

	serveErr := func(err error) {
		log.Error_log.Printf("storage-store: %s", err.Error())
		ctx.Error("Internal error, check error.log", http.StatusInternalServerError)
	}

	file := ctx.FormValue("file")
	size := string(ctx.FormValue("size"))
	size_, _ := strconv.Atoi(size)

	if len(file) == 0 {
		serveErr(errors.New("size of new body is zero"))
		return
	}

	crcb_here := crc32.ChecksumIEEE(file)

	crcb := string(ctx.FormValue("crcb"))
	crcb_, _ := strconv.Atoi(crcb)
	if crcb_ == 0 || crcb_here != uint32(crcb_) {
		serveErr(fmt.Errorf("crc of file body is not same: %d vs %d", crcb_, crcb_here))
	}

	var crc int64
	crc_ := string(ctx.FormValue("crc"))
	if len(crc_) > 0 {
		crc__, err := strconv.Atoi(crc_)
		if err == nil {
			crc = int64(crc__)
		}
	}
	if crc == 0 {
		serveErr(errors.New("no crc is given"))
		return
	}

	err := storage.Storage.Store(project, branch, crc, file, int64(size_))
	if err != nil {
		serveErr(err)
		return
	}

	ctx.SetStatusCode(200)
	ctx.WriteString("OK")
}
