package dao

import (
	"time"
)

type PipelineBlockStats struct {
	BlockNumber  int       `json:"block_number" gorm:"primaryKey;column:block_number"`
	MetaData     []byte    `json:"meta_data" gorm:"type:jsonb;default:'{}';column:meta_data"`
	ProcessTime  time.Time `json:"process_time" gorm:"column:process_time"`
	State        string    `json:"state" gorm:"column:state"`
	PipelineName string    `json:"pipeline_name" gorm:"column:pipeline_name"`
}
