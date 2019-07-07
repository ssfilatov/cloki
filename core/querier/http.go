package querier

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/ssfilatov/cloki/core/utils"
	"github.com/ssfilatov/cloki/lokiproto"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultQueryLimit = 100
	defaulSince       = 1 * time.Hour
)

func (q *Querier) QueryHandler(w http.ResponseWriter, r *http.Request) {
	request, err := httpRequestToQueryRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Debugf("request: %v", request)

	result, err := q.Query(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// nolint
func directionParam(values url.Values, name string, def lokiproto.Direction) (lokiproto.Direction, error) {
	value := values.Get(name)
	if value == "" {
		return def, nil
	}

	d, ok := lokiproto.Direction_value[strings.ToUpper(value)]
	if !ok {
		return lokiproto.Direction_FORWARD, fmt.Errorf("invalid direction '%s'", value)
	}
	return lokiproto.Direction(d), nil
}

func httpRequestToQueryRequest(httpRequest *http.Request) (*lokiproto.QueryRequest, error) {
	params := httpRequest.URL.Query()
	now := time.Now()
	queryRequest := lokiproto.QueryRequest{
		Regex: params.Get("regexp"),
		Query: params.Get("query"),
	}

	limit, err := utils.IntParam(params, "limit", defaultQueryLimit)
	if err != nil {
		return nil, err
	}
	queryRequest.Limit = uint32(limit)

	queryRequest.Start, err = utils.UnixNanoTimeParam(params, "start", now.Add(-defaulSince))
	if err != nil {
		return nil, err
	}

	queryRequest.End, err = utils.UnixNanoTimeParam(params, "end", now)
	if err != nil {
		return nil, err
	}

	queryRequest.Direction, err = directionParam(params, "direction", lokiproto.Direction_BACKWARD)
	if err != nil {
		return nil, err
	}

	return &queryRequest, nil
}

func (q *Querier) LabelHandler(w http.ResponseWriter, r *http.Request) {
	name, ok := mux.Vars(r)["name"]
	req := &lokiproto.LabelRequest{
		Values: ok,
		Name:   name,
	}
	resp, err := q.Label(r.Context(), req)
	if err != nil {
		log.Errorf("error getting labels: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Errorf("error encoding labels resp: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
