package repository

import (
	"context"
	"log-ingestor/internal/core/domain"
	"time"

	"github.com/labstack/gommon/log"
)

func ConvertTimestampToTime(timestampString string) (timeStamp time.Time, err error) {

	if timestampString == "" {
		return timeStamp, nil
	}

	timeStamp, err = time.Parse(time.RFC3339, timestampString)
	return timeStamp, err
}

// GetLogs returns all logs from the database
func (repo *DB) GetLogs(logFilter domain.LogFilter) (logs domain.LogResponse, err error) {

	logFilter.Level = "%" + logFilter.Level + "%"

	logFilter.Message = "%" + logFilter.Message + "%"

	logFilter.ResourceId = "%" + logFilter.ResourceId + "%"

	timeStampQueryString := ""

	if logFilter.TimestampStart != "" {
		timeStampQueryString += " AND timestamp >= '" + logFilter.TimestampStart + "'"
	}

	if logFilter.TimestampEnd != "" {
		timeStampQueryString += " AND timestamp <= '" + logFilter.TimestampEnd + "'"
	}

	// count total logs
	err = repo.db.WithContext(context.Background()).
		Model(&domain.Log{}).
		Where("level LIKE ? AND message LIKE ? AND resource_id LIKE ?"+timeStampQueryString, logFilter.Level, logFilter.Message, logFilter.ResourceId).
		Count(&logs.Total).Error

	if err != nil {
		log.Error(err)
		return logs, err
	}

	// get logs using limit and offset
	query := repo.db.WithContext(context.Background()).
		Model(&domain.Log{}).
		Offset(logFilter.Page*logFilter.Limit).
		Where("level LIKE ? AND message LIKE ? AND resource_id LIKE ?"+timeStampQueryString, logFilter.Level, logFilter.Message, logFilter.ResourceId)

	if logFilter.Limit != 0 {
		query = query.Limit(logFilter.Limit)
	}

	err = query.Find(&logs.Logs).Error

	if err != nil {
		log.Error(err)
		return logs, err
	}

	return logs, nil
}
