package log

import (
	"fmt"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/MasterDimmy/zipologger"
	"github.com/valyala/fasthttp"
)

var (
	Access_log   = zipologger.NewLogger("./logs/access.log", 5, 5, 5, false)
	Error_log    = zipologger.NewLogger("./logs/error.log", 5, 5, 5, false).SetAlsoToStdout(true)
	action_log   = zipologger.NewLogger("./logs/action.log", 5, 5, 5, false)
	Requests_log = zipologger.NewLogger("./logs/requests.log", 2, 2, 2, false)
)

var DumpRequests bool

func LogCtxRequest(h func(_ *fasthttp.RequestCtx)) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		ipv4 := ctx.RemoteIP().To4().String()

		h(ctx) //execute handler, then log it

		Access_log.Printf("[%s] %s - %4d|%4s|%.2f ms| %s | %d bytes", ipv4, string(ctx.Response.Header.ContentType()), ctx.Response.StatusCode(), string(ctx.Method()), float64(time.Since(start).Nanoseconds())/1000000.0, string(ctx.Request.URI().Path()), len(ctx.Response.Body()))

		if DumpRequests {
			dumpAnswer(ctx)
		}
	})
}

type CtxLogger struct {
	Logger *zipologger.Logger
}

func (l CtxLogger) Printf(format string, w1 ...interface{}) {
	if len(w1) == 0 {
		l.Logger.Printf(format, nil)
		return
	}
	if len(w1) == 1 {
		l.Logger.Printf(format, w1[0])
		return
	}
	l.Logger.Printf(format, w1[0], w1[1:]...)
}

func Action(ip, action, project, branch, result string) {
	action_log.Printf("%s : %s : %s (%s) - %s", ip, action, project, branch, result)
}

var url_replacer = strings.NewReplacer("/", "_", ".", "")

func dumpAnswer(ctx *fasthttp.RequestCtx) {
	//not found for /api/sdfsdfdsf - do not log this
	if ctx.Response.StatusCode() == 404 {
		return
	}

	ipv4 := ctx.RemoteIP().String()

	req_url := string(ctx.Request.URI().Path())
	upd_url, ok := ctx.UserValue("URL").(string)
	if ok && upd_url != "" {
		req_url = upd_url
	}

	ret_log := func() *zipologger.Logger {
		return zipologger.GetLoggerBySuffix(url_replacer.Replace(req_url), "./logs/answers_url_", 3, 10, 60, false)
	}

	body := ""
	enc := string(ctx.Response.Header.Peek("Content-Encoding"))
	if enc == "gzip" {
		bodyb, err := ctx.Response.BodyGunzip()
		body = string(bodyb)
		if err != nil {
			ret_log().Printf("BodyGunzip: %s", err.Error())
			return
		}
	} else {
		body = string(ctx.Response.Body())
	}

	if len(body) > 2*1024*1024 {
		body = body[:2*1024*1024]
	}

	body = strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, string(body))

	params := ""
	method := string(ctx.Method())
	if method == "POST" {
		params = string(ctx.PostBody())
		params, _ = url.QueryUnescape(params)
		if len(params) > 10000 {
			params = params[:10000]
		}
	} else if method == "GET" {
		params = string(ctx.QueryArgs().String())
	}

	ret := fmt.Sprintf("[%s] %s - %4d|%4s| %s | %s | %d bytes =>\n%s",
		ipv4, string(ctx.Response.Header.ContentType()),
		ctx.Response.StatusCode(), method,
		ctx.Request.RequestURI(),
		params,
		len(ctx.Response.Body()), body)

	if ctx.Response.StatusCode() != 404 { //для not found файлы не создаем
		//every log adds date-time prefix to message ("2006/01/02 15:04:05 ")
		//cut it from answers_log return, to avoid printing date-time twice (pim 25.11.2022)
		if len(ret) > 20 {
			ret = ret[20:]
		}
		ret_log().Print(ret)
	}
}
