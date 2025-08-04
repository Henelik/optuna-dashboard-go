package db

import (
	"time"

	"gorm.io/gorm"
)

var DB *gorm.DB

type Study struct {
	ID   uint   `gorm:"column:study_id"`
	Name string `gorm:"column:study_name"`
}

type StudyDirection struct {
	ID        uint   `gorm:"column:study_direction_id"`
	StudyID   uint   `gorm:"column:study_id"`
	Direction string `gorm:"column:direction"`
	Objective uint
}

type Trial struct {
	ID       uint `gorm:"column:trial_id"`
	StudyID  uint `gorm:"column:study_id"`
	Number   uint
	State    string
	Start    time.Time `gorm:"column:datetime_start"`
	Complete time.Time `gorm:"column:datetime_complete"`
}

type TrialParam struct {
	ID               uint   `gorm:"column:param_id"`
	TrialID          uint   `gorm:"column:trial_id"`
	Name             string `gorm:"column:param_name"`
	Value            string `gorm:"column:param_value"`
	DistributionJSON string `gorm:"column:distribution_json"`
}

type TrialValue struct {
	ID        uint    `gorm:"column:trial_value_id"`
	TrialID   uint    `gorm:"column:trial_id"`
	Objective uint    `gorm:"column:objective"`
	Value     float64 `gorm:"column:value"`
	Type      string  `gorm:"column:value_type"`
}
