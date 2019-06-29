package chloki

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

func queryHandler(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	query := ctx.QueryArgs().Peek("query")
	limit := ctx.QueryArgs().Peek("limit")
	start := ctx.QueryArgs().Peek("start")
	end := ctx.QueryArgs().Peek("end")
	direction := ctx.QueryArgs().Peek("direction")
	regexp := ctx.QueryArgs().Peek("regexp")

}

func getLabelsHandler(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	fmt.Fprint(ctx, "Welcome!\n")
}

func getLabelValuesHandler(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {
	fmt.Fprintf(ctx, "hello, %s!\n", ps.ByName("name"))
}