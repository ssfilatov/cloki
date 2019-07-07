package querier

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/ssfilatov/cloki/core/clickhouse"
	"github.com/ssfilatov/cloki/core/utils"
	"github.com/ssfilatov/cloki/lokiproto"
	"strings"
	"time"
)

const rowLimit = 256

type Querier struct {
	cfg *utils.Config
	ch  *clickhouse.Clickhouse
}

func New(cfg *utils.Config, ch *clickhouse.Clickhouse) *Querier {
	return &Querier{
		cfg,
		ch,
	}
}

func (q *Querier) Query(ctx context.Context, req *lokiproto.QueryRequest) (*lokiproto.QueryResponse, error) {
	logQL := req.Query
	//parser is limited to EQ only, so it's actually not logQL
	filters, err := utils.ParseExpr(logQL)
	if err != nil {
		return nil, err
	}
	//remove unavailable filters
	for label := range *filters {
		if !q.LabelValid(label) {
			delete(*filters, label)
		}
	}
	var asc bool
	if req.Direction == lokiproto.Direction_FORWARD {
		asc = true
	} else {
		asc = false
	}
	logEntry := &clickhouse.LogEntryQuery{
		Filters: filters,
		Start:   req.Start.Format("2006-01-02 15:04:05"),
		End:     req.End.Format("2006-01-02 15:04:05"),
		Asc:     asc,
		Regex:   req.Regex,
		Limit:   req.Limit,
	}
	entryList, err := q.ch.GetLogEntries(logEntry)
	entries := make([]lokiproto.Entry, 0, len(entryList))
	for _, entryString := range entryList {
		tokens := strings.Split(entryString, "\t")
		if len(tokens) != 2 {
			log.Errorf("error spliting entry: %s", entryString)
			continue
		}
		layout := "2006-01-02 15:04:05"
		timestamp, err := time.Parse(layout, tokens[0])
		if err != nil {
			log.Errorf("error parsing time: %s", err)
			continue
		}
		entries = append(entries, lokiproto.Entry{
			Timestamp: timestamp,
			Line:      tokens[1],
		})
	}
	//single stream is used
	stream := &lokiproto.Stream{
		Labels:  string(req.Query),
		Entries: entries,
	}
	result := &lokiproto.QueryResponse{
		Streams: make([]*lokiproto.Stream, 0, 1),
	}
	result.Streams = append(result.Streams, stream)
	return result, err
}

func (q *Querier) LabelValid(label string) bool {
	for _, labelAvailable := range *q.cfg.LabelList {
		if label == labelAvailable {
			return true
		}
	}
	return false
}

func (q *Querier) Label(ctx context.Context, req *lokiproto.LabelRequest) (*lokiproto.LabelResponse, error) {
	if req.Values == false {
		return &lokiproto.LabelResponse{
			Values: *q.cfg.LabelList,
		}, nil
	} else {
		if !q.LabelValid(req.Name) {
			return &lokiproto.LabelResponse{
				Values: []string{},
			}, nil
		} else {
			rows, err := q.ch.GetColumnValues(req.Name, rowLimit)
			if err != nil {
				return nil, err
			}
			return &lokiproto.LabelResponse{
				Values: rows,
			}, nil
		}
	}
}
