package repository

import (
	"context"
	"log-ingestor/internal/core/domain"

	"github.com/labstack/gommon/log"
)

// InsertLog inserts a log into the database
func (repo *DB) InsertLog(log domain.Log) error {
	err := repo.db.WithContext(context.Background()).Create(&log).Error

	if err != nil {
		return err
	}

	return nil
}

// InsertLogWithPreparedStmt inserts a log into the database using a prepared statement
func (repo *DB) InsertLogWithPreparedStmt(log domain.Log) error {

	// Define SQL statement with named parameters
	query := "INSERT INTO log_table (level, message, resource_id, timestamp, trace_id, span_id, commit, parent_resource_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	err := repo.db.WithContext(context.Background()).Exec(query, log.Level, log.Message, log.ResourceID, log.Timestamp, log.TraceID, log.SpanID, log.Commit, log.ParentResourceID).Error

	if err != nil {
		return err
	}

	return nil
}

// InsertBulkLog inserts a log into the database using Create method
func (repo *DB) InsertBulkLog(logs []domain.Log) error {

	err := repo.db.WithContext(context.Background()).Create(&logs).Error
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
