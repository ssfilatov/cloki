package querier

import (
	"context"
	cloki "github.com/ssfilatov/cloki/core"
	"github.com/ssfilatov/cloki/core/clickhouse"
	"github.com/ssfilatov/cloki/lokiproto"
	"time"
)

const rowLimit = 256

type Querier struct {
	cfg *cloki.Config
	ch  *clickhouse.Clickhouse
}

func New(cfg *cloki.Config, ch *clickhouse.Clickhouse) *Querier {
	return &Querier{
		cfg,
		ch,
	}
}

func (q *Querier) Query(ctx context.Context, req *lokiproto.QueryRequest) (*lokiproto.QueryResponse, error) {

	resp, _, err := q.ReadBatch(req.Limit)
	return resp, err
}

func (q *Querier) ReadBatch(size uint32) (*lokiproto.QueryResponse, uint32, error) {
	sampleLabels := "test label"
	sampleEntry := lokiproto.Entry{
		Timestamp: time.Now(),
		Line:      "sample line",
	}
	streams := map[string]*lokiproto.Stream{}
	respSize := uint32(0)
	for ; respSize < size; respSize++ {
		labels, entry := sampleLabels, sampleEntry
		stream, ok := streams[labels]
		if !ok {
			stream = &lokiproto.Stream{
				Labels: labels,
			}
			streams[labels] = stream
		}
		stream.Entries = append(stream.Entries, entry)
	}

	result := lokiproto.QueryResponse{
		Streams: make([]*lokiproto.Stream, 0, len(streams)),
	}
	for _, stream := range streams {
		result.Streams = append(result.Streams, stream)
	}
	return &result, respSize, nil
}

func (q *Querier) Label(ctx context.Context, req *lokiproto.LabelRequest) (*lokiproto.LabelResponse, error) {
	if req.Values == false {
		return &lokiproto.LabelResponse{
			Values: *q.cfg.LabelList,
		}, nil
	} else {
		nameInLabels := false
		for _, label := range *q.cfg.LabelList {
			if req.Name == label {
				nameInLabels = true
				break
			}
		}
		if nameInLabels == false {
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
