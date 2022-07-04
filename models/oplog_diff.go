package models

import "time"

type OplogDiff struct {
	Optype    string    `json:"optype"`
	TableName string    `json:"table_name"`
	Timestamp time.Time `json:"timestamp"`
	Old       Interest  `json:"old"`
	New       Interest  `json:"new"`
}
