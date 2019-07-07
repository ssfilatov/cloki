# Loki interface for clickhouse

### Example config

```yaml
debug: True
server:
    http_listen_host: localhost
    http_listen_port: 3100
label_list:
    - project_id
    - instance_id
    - tag
    - level
clickhouse:
    url: http://localhost:8123/
    user: logger
    password: p@ssw0rd
    database: default
    table: logs
    timestamp_column: timestamp
```