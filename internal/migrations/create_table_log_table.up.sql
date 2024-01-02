CREATE TABLE log_table (
    id SERIAL PRIMARY KEY,
    level VARCHAR(10),
    message TEXT,
    resource_id VARCHAR(36),
    timestamp TIMESTAMP WITH TIME ZONE,
    trace_id VARCHAR(36),
    span_id VARCHAR(36),
    commit VARCHAR(40),
    parent_resource_id VARCHAR(36)
);

COMMENT ON COLUMN log_table.id IS 'id of the log entry';
COMMENT ON COLUMN log_table.level IS 'log level';
COMMENT ON COLUMN log_table.message IS 'log message';
COMMENT ON COLUMN log_table.resource_id IS 'resource id';
COMMENT ON COLUMN log_table.timestamp IS 'log timestamp';
COMMENT ON COLUMN log_table.trace_id IS 'trace id';
COMMENT ON COLUMN log_table.span_id IS 'span id';
COMMENT ON COLUMN log_table.commit IS 'commit hash';
COMMENT ON COLUMN log_table.parent_resource_id IS 'parent resource id';