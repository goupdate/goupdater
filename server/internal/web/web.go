package web

import (
	"time"

	"goupdater/server/internal/log"

	"github.com/MasterDimmy/zipologger"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func CreateAndListen(dumpRequests bool) *fasthttp.Server {
	router := router.New()
	router.POST("/info", GetInfo)
	router.POST("/download", Download)
	router.POST("/upload", UploadNewFile)

	log.DumpRequests = dumpRequests

	srv := &fasthttp.Server{
		Handler:            log.LogCtxRequest(router.Handler),
		Concurrency:        100,
		DisableKeepalive:   false,
		MaxRequestBodySize: 100 * 1024 * 1024,
		ReadTimeout:        300 * time.Second,
		TCPKeepalive:       true,
		MaxConnsPerIP:      10,
		MaxRequestsPerConn: 20,
		GetOnly:            false,

		LogAllErrors: true,

		NoDefaultServerHeader: true,
		NoDefaultContentType:  true,
		Logger: log.CtxLogger{
			Logger: zipologger.NewLogger("./logs/all.log", 3, 3, 3, false),
		},
		Name: "goupdater-server",
	}

	return srv
}
