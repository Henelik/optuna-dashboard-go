package db

import "time"

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
