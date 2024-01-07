package domain

// LogResponse is a struct that represents the log response model
type LogResponse struct {
	Logs  []Log `json:"logs"`
	Total int64 `json:"total"`
}
