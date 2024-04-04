package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type DemoBlockStats struct {
	BlockNumber  int       `json:"block_number" gorm:"primaryKey;column:block_number"`
	MetaData     []byte    `json:"meta_data" gorm:"type:jsonb;default:'{}';column:meta_data"`
	ProcessTime  time.Time `json:"process_time" gorm:"column:process_time"`
	State        string    `json:"state" gorm:"column:state"`
	PipelineName string    `json:"pipeline_name" gorm:"column:pipeline_name"`
}

func (m *DemoBlockStats) TableName() string {
	return "public.demo_block_stats"
}

type DemoBlockStatsDao struct {
	sourceDB  *gorm.DB
	replicaDB []*gorm.DB
	m         *DemoBlockStats
}

func NewDemoBlockStatsDao(ctx context.Context, dbs ...*gorm.DB) *DemoBlockStatsDao {
	dao := new(DemoBlockStatsDao)
	switch len(dbs) {
	case 0:
		panic("database connection required")
	case 1:
		dao.sourceDB = dbs[0]
		dao.replicaDB = []*gorm.DB{dbs[0]}
	default:
		dao.sourceDB = dbs[0]
		dao.replicaDB = dbs[1:]
	}
	return dao
}

func (d *DemoBlockStatsDao) Upsert(ctx context.Context, obj *DemoBlockStats) error {
	err := d.sourceDB.Save(&obj).Error
	if err != nil {
		return fmt.Errorf("DemoBlockStatsDao: %w", err)
	}
	return nil
}

func (d *DemoBlockStatsDao) Create(ctx context.Context, obj *DemoBlockStats) error {
	err := d.sourceDB.Model(d.m).Create(&obj).Error
	if err != nil {
		return fmt.Errorf("DemoBlockStatsDao: %w", err)
	}
	return nil
}

func (d *DemoBlockStatsDao) Get(ctx context.Context, fields, where string) (*DemoBlockStats, error) {
	items, err := d.List(ctx, fields, where, 0, 1)
	if err != nil {
		return nil, fmt.Errorf("DemoBlockStatsDao: Get where=%s: %w", where, err)
	}
	if len(items) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &items[0], nil
}

func (d *DemoBlockStatsDao) List(ctx context.Context, fields, where string, offset, limit int) ([]DemoBlockStats, error) {
	var results []DemoBlockStats
	err := d.replicaDB[rand.Intn(len(d.replicaDB))].Model(d.m).
		Select(fields).Where(where).Offset(offset).Limit(limit).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("DemoBlockStatsDao: List where=%s: %w", where, err)
	}
	return results, nil
}

func (d *DemoBlockStatsDao) Update(ctx context.Context, where string, update map[string]interface{}, args ...interface{}) error {
	err := d.sourceDB.Model(d.m).Where(where, args...).
		Updates(update).Error
	if err != nil {
		return fmt.Errorf("DemoBlockStatsDao:Update where=%s: %w", where, err)
	}
	return nil
}

func (d *DemoBlockStatsDao) Delete(ctx context.Context, where string, args ...interface{}) error {
	if len(where) == 0 {
		return fmt.Errorf("DemoBlockStatsDao: Delete where=%s, which is empty", where)
	}
	if err := d.sourceDB.Where(where, args...).Delete(d.m).Error; err != nil {
		return fmt.Errorf("DemoBlockStatsDao: Delete where=%s: %w", where, err)
	}
	return nil
}
