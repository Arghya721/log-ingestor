package repository

import (
	"context"
	"log-ingestor/internal/core/domain"

	"github.com/labstack/gommon/log"
)

// GetLogs returns all logs from the database
func (repo *DB) GetLogs(logFilter domain.LogFilter) (logs domain.LogResponse, err error) {

	// count total logs
	err = repo.db.WithContext(context.Background()).Model(&domain.Log{}).Count(&logs.Total).Error
	if err != nil {
		log.Error(err)
		return logs, err
	}

	// get logs using limit and offset
	query := repo.db.WithContext(context.Background()).Model(&domain.Log{}).Limit(logFilter.Limit).Offset(logFilter.Page)

	err = query.Find(&logs.Logs).Error

	if err != nil {
		log.Error(err)
		return logs, err
	}

	return logs, nil
}
