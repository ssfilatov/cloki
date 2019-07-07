package cloki

import (
	"github.com/ssfilatov/cloki/core/clickhouse"
	"github.com/ssfilatov/cloki/core/querier"
	"github.com/ssfilatov/cloki/core/utils"
)

type CLoki struct {
	Query *querier.Querier
	Cfg   *utils.Config
}

func NewCLoki(cfg *utils.Config) (*CLoki, error) {
	ch, err := clickhouse.NewClickhouse(cfg)
	if err != nil {
		return nil, err
	}
	q := querier.New(cfg, ch)
	return &CLoki{
		q,
		cfg,
	}, nil
}
