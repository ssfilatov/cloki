package chloki

import (
	"fmt"
	"github.com/prometheus/log"
	"github.com/ssfilatov/clickhouse-loki-adapter/core/clickhouse"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

const Forward = "forward"
const Backward = "backward"
const DefaultDirection = Backward

type LogQuery struct {
	query     string
	limit     int32
	start     string
	end       string
	direction string
	regexp    string
}

func NewGetLogsQueryFromCtx(ctx *fasthttp.RequestCtx) (*LogQuery, error) {
	var logQuery LogQuery
	logQuery.query = string(ctx.QueryArgs().Peek("query"))
	byteLimit := ctx.QueryArgs().Peek("limit")
	logQuery.limit = readInt32(byteLimit)
	logQuery.start = string(ctx.QueryArgs().Peek("start"))
	logQuery.end = string(ctx.QueryArgs().Peek("end"))
	direction := string(ctx.QueryArgs().Peek("direction"))
	if direction == "" {
		direction = DefaultDirection
	} else if direction != Forward && direction != Backward {
		return nil, fmt.Errorf("wrong direction value: %s", direction)
	}
	logQuery.direction = direction
	logQuery.regexp = string(ctx.QueryArgs().Peek("regexp"))
	return &logQuery, nil
}

func queryHandler(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	logQuery, err := NewGetLogsQueryFromCtx(ctx)
	if err != nil {
		log.Errorf("Error parsing query: %s", err)
		ctx.Error("Error parsing query", fasthttp.StatusBadRequest)
	}
	ch, err := clickhouse.NewClickhouse(config)
	if err != nil {
		log.Errorf("Error connecting to clickhouse: %s", err)
		ctx.Error("Error connecting to clickhouse", fasthttp.StatusInternalServerError)
	}
	_, err = ch.GetLogEntries(logQuery.start, logQuery.end)
	if err != nil {
		log.Errorf("Error quering clickhouse: %s", err)
		ctx.Error("Error quering clickhouse", fasthttp.StatusInternalServerError)
	}

}
