package clickhouse

import (
	"database/sql"
	"fmt"
	_ "github.com/mailru/go-clickhouse"
	"github.com/ssfilatov/clickhouse-loki-adapter/core/chloki"
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

func NewClickhouse(config *chloki.Config) (*Clickhouse, error) {
	connectionString := fmt.Sprintf("http://%s/%s", config.Clickhouse.Server, config.Clickhouse.Database)
	connect, err := sql.Open("clickhouse", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to clickhouse: %s", err)
	}
	ch := &Clickhouse{
		connect,
		config.Clickhouse.LogTableName,
	}
	return ch, nil
}

type Clickhouse struct {
	conn  *sql.DB
	table string
}

func (c *Clickhouse) GetLogEntries(start, end string) ([]*LogEntry, error) {
	rows, err := c.conn.Query(`
		SELECT 
			timestamp,
			project_id,
			instance_id,
			msec,
			tag,
		    level,
			message
		FROM
			%s`)
	if err != nil {
		return nil, err
	}
	var logEntries []*LogEntry
	for rows.Next() {
		var entry LogEntry
		if err := rows.Scan(
			&entry.Timestamp,
			&entry.ProjectID,
			&entry.InstanceID,
			&entry.Msec,
			&entry.Tag,
			&entry.Level,
			&entry.Message,
		); err != nil {
			return nil, err
		}
		logEntries = append(logEntries, &entry)
	}
	return logEntries, nil
}
