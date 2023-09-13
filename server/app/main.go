package main

import (
	"flag"
	"os"

	"goupdater/server/internal/config"
	"goupdater/server/internal/log"
	"goupdater/server/internal/storage"
	"goupdater/server/internal/web"
	"goupdater/server/other/tlsgen"

	"github.com/MasterDimmy/go-ctrlc"
	"github.com/MasterDimmy/zipologger"
	"github.com/valyala/fasthttp"
)

func main() {
	defer zipologger.Wait()

	dumpRequests := flag.Bool("dumpreq", false, "do dump requests into log?")
	flag.Parse()

	cc := &ctrlc.CtrlC{}

	err := config.Read()
	if err != nil {
		log.Error_log.Printf("cant read config: %s", err.Error())
		return
	}

	err = storage.Create()
	if err != nil {
		log.Error_log.Printf("cant create storage: %s", err.Error())
		return
	}

	defer cc.DeferThisToWaitCtrlC()

	var srv *fasthttp.Server
	go cc.InterceptKill(true, func() {
		srv.Shutdown()
	})

	srv = web.CreateAndListen(*dumpRequests)

	_, err = os.Open("tls.cert")
	if err != nil {
		tlsgen.Generate()
		log.Error_log.Print("tls.cert not exists, recreating")
	}

	_, err = os.Open("tls.key")
	if err != nil {
		tlsgen.Generate()
		log.Error_log.Print("tls.key not exists, recreating")
	}

	err = srv.ListenAndServeTLS(config.Config.Listen, "tls.cert", "tls.key")
	if err != nil {
		log.Error_log.Printf("cant serve %s , error: %s", config.Config.Listen, err.Error())
		return
	}
}
