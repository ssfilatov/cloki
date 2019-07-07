package clickhouse

import (
	"fmt"
	"github.com/ssfilatov/cloki/core/utils"
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

func NewClickhouse(config *utils.Config) (*Clickhouse, error) {

	ch := &Clickhouse{
		fmt.Sprintf("%s/?user=%s&password=%s&database=%s&query=",
			config.Clickhouse.URL,
			config.Clickhouse.User,
			url.QueryEscape(config.Clickhouse.Password),
			config.Clickhouse.Database),
		config.Clickhouse.LogTableName,
		config.Clickhouse.TimestampColumn,
	}
	err := ch.Ping()
	if err != nil {
		return nil, fmt.Errorf("error trying to connect to clickhouse: %s", err)
	}
	return ch, nil
}

type Clickhouse struct {
	uriString       string
	table           string
	timestampColumn string
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
	defer resp.Body.Close()
	byteResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	stringResp := string(byteResp)
	if stringResp == "" {
		return make([]string, 0), nil
	}
	lines := strings.Split(stringResp, "\n")
	lines = lines[:len(lines)-1]
	return lines, nil
}

type LogEntryQuery struct {
	Filters *map[string]string
	Start   string
	End     string
	Asc     bool
	Regex   string
	Limit   uint32
}

func (ch *Clickhouse) GetLogEntries(q *LogEntryQuery) ([]string, error) {
	queryString := fmt.Sprintf("SELECT %s, message FROM %s", ch.timestampColumn, ch.table)
	queryString = fmt.Sprintf(`%s WHERE %s BETWEEN toDateTime('%s') AND toDateTime('%s')`,
		queryString, ch.timestampColumn, q.Start, q.End)
	for column, value := range *q.Filters {
		queryString = fmt.Sprintf(`%s AND %s='%s'`, queryString, column, value)
	}
	if q.Asc {
		queryString = fmt.Sprintf("%s ORDER BY %s ASC", queryString, ch.timestampColumn)
	} else {
		queryString = fmt.Sprintf("%s ORDER BY %s", queryString, ch.timestampColumn)
	}
	if q.Regex != "" {
		queryString = fmt.Sprintf("%s LIKE %s", queryString, q.Regex)
	}
	queryString = url.PathEscape(fmt.Sprintf("%s LIMIT %d", queryString, q.Limit))
	resp, err := http.Get(ch.uriString + queryString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	byteResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	stringResp := string(byteResp)
	if stringResp == "" {
		return make([]string, 0), nil
	}
	lines := strings.Split(stringResp, "\n")
	lines = lines[:len(lines)-1]
	return lines, nil
}
