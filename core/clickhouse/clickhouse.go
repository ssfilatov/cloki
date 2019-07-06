package clickhouse

import (
	"fmt"
	cloki "github.com/ssfilatov/cloki/core"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type LogEntry struct {
	Timestamp  string
	ProjectID  string
	InstanceID string
	Msec       int32
	Tag        string
	Level      string
	Message    string
}

func NewClickhouse(config *cloki.Config) (*Clickhouse, error) {

	ch := &Clickhouse{
		fmt.Sprintf("%s/?user=%s&password=%s&database=%s&query=",
			config.Clickhouse.URL,
			config.Clickhouse.User,
			url.QueryEscape(config.Clickhouse.Password),
			config.Clickhouse.Database),
		config.Clickhouse.LogTableName,
	}
	err := ch.Ping()
	if err != nil {
		return nil, fmt.Errorf("error trying to connect to clickhouse: %s", err)
	}
	return ch, nil
}

type Clickhouse struct {
	uriString string
	table     string
}

func (ch *Clickhouse) Ping() error {
	_, err := http.Get(ch.uriString + url.PathEscape("SELECT 1"))
	if err != nil {
		return err
	}
	return nil
}

func (ch *Clickhouse) GetColumnValues(name string, limit uint32) ([]string, error) {
	query := url.PathEscape(fmt.Sprintf("SELECT DISTINCT %s FROM %s LIMIT %d", name, ch.table, limit))
	resp, err := http.Get(ch.uriString + query)
	if err != nil {
		return nil, err
	}
	byteResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(byteResp), "\n"), nil
}

//func (ch *Clickhouse) GetLogEntries(name string, limit uint32) ([]string, error) {
//	query := url.PathEscape(fmt.Sprintf("SELECT message FROM %s LIMIT %d", ch.table, limit))
//	resp, err :=  http.Get(ch.uriString + query)
//	if err != nil {
//		return nil, err
//	}
//	byteResp, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//	return strings.Split(string(byteResp), "\n"), nil
//}
