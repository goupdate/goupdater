package web

import (
	"encoding/json"
	"net/http"

	"goupdater/server/internal/config"
	"goupdater/server/internal/log"
	"goupdater/server/internal/storage"

	"github.com/valyala/fasthttp"
)

// checks if new file exists to download
// method: GET
// parameters:
//   key
//   project
//   branch
// answer:
//   200 - found + info JSON
//   404 - not found
//   403 - access denied

func GetInfo(ctx *fasthttp.RequestCtx) {
	project, branch, allowed := IsAllowed(ctx)
	if !allowed {
		ctx.Error("Access denied", http.StatusForbidden)
		return
	}

	result := "OK"
	defer log.Action(ctx.RemoteIP().String(), "info", project, branch, result)

	item, err := storage.Storage.Find(project, branch)
	if err != nil {
		result = err.Error()
		log.Error_log.Printf("storage-find: %s", err.Error())
		ctx.Error("Internal error, check error.log", http.StatusInternalServerError)
		return
	}

	if item.Size == 0 {
		result = "not found"
		ctx.Error("Not found", http.StatusNotFound)
		return
	}

	// do not tell where we store it
	item.Path = ""

	ctx.SetContentType("application/json")

	e := json.NewEncoder(ctx)
	e.Encode(&item)
}

// returns file if exists for given key,projectname,projectbranch
func Download(ctx *fasthttp.RequestCtx) {
	project, branch, allowed := IsAllowed(ctx)
	if !allowed {
		ctx.Error("Access denied", 403)
		return
	}

	result := "OK"
	defer log.Action(ctx.RemoteIP().String(), "download", project, branch, result)

	item, err := storage.Storage.Find(project, branch)
	if err != nil {
		result = err.Error()
		log.Error_log.Printf("storage-find: %s", err.Error())
		ctx.Error("Internal error, check error.log", http.StatusInternalServerError)
		return
	}
	if item.Size == 0 {
		result = "not found"
		ctx.Error("Not found", http.StatusNotFound)
		return
	}

	fasthttp.ServeFileBytesUncompressed(ctx, []byte(item.Path))
}

func IsAllowed(ctx *fasthttp.RequestCtx) (project, branch string, allowed bool) {
	key := string(ctx.FormValue("key"))
	project = string(ctx.FormValue("project"))
	branch = string(ctx.FormValue("branch"))

	allowed = false

	for _, v := range config.Config.Projects {
		for _, b := range v.Branches {
			if v.Key == key && v.Name == project && b == branch {
				allowed = true
				break
			}
		}
	}
	return
}
