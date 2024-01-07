package domain

// LogRequest is a struct that represents the log request model
type LogRequest struct {
	Level      string   `json:"level"`
	Message    string   `json:"message"`
	ResourceID string   `json:"resourceId"`
	Timestamp  string   `json:"timestamp"`
	TraceID    string   `json:"traceId"`
	SpanID     string   `json:"spanId"`
	Commit     string   `json:"commit"`
	Metadata   MetaData `json:"metadata"`
}

type MetaData struct {
	ParentResourceID string `json:"parentResourceId"`
}

// LogFilter is a struct that represents the log filter model
type LogFilter struct {
	Page  int
	Limit int
}
