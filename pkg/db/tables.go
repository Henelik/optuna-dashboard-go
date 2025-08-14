package db

import (
	"encoding/json"
	"errors"
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
	DistributionJSON []byte `gorm:"column:distribution_json"`
}

type Distribution struct {
	Name       string
	Attributes DistributionAttributes
}

type DistributionAttributes struct {
	Low     *float64
	High    *float64
	Step    *float64
	Log     *bool
	Choices []any
}

type TrialValue struct {
	ID        uint    `gorm:"column:trial_value_id"`
	TrialID   uint    `gorm:"column:trial_id"`
	Objective uint    `gorm:"column:objective"`
	Value     float64 `gorm:"column:value"`
	Type      string  `gorm:"column:value_type"`
}

type TrialUserAttribute struct {
	ID      uint   `gorm:"column:trial_user_attribute_id"`
	TrialID uint   `gorm:"column:trial_id"`
	Key     string `gorm:"column:key"`
	Value   string `gorm:"column:value_json"`
}

type TrialIntermediateValue struct {
	ID      uint    `gorm:"column:trial_intermediate_value_id"`
	TrialID uint    `gorm:"column:trial_id"`
	Step    int     `gorm:"column:step"`
	Value   float64 `gorm:"column:intermediate_value"`
	Type    string  `gorm:"column:intermediate_value_type"`
}

type TrialHeartbeat struct {
	ID      uint      `gorm:"column:trial_heartbeat_id"`
	TrialID uint      `gorm:"column:trial_id"`
	Time    time.Time `gorm:"column:heartbeat"`
}

type TrialSystemAttribute struct {
	ID      uint   `gorm:"column:trial_system_attribute_id"`
	TrialID uint   `gorm:"column:trial_id"`
	Key     string `gorm:"column:key"`
	Value   string `gorm:"column:value_json"`
}

// BestTrialResult represents the best trial with its parameters and value
type BestTrialResult struct {
	Trial       Trial
	TrialParams []TrialParam
	TrialValue  TrialValue
}

func DeleteStudy(studyID uint, tx *gorm.DB) error {
	// delete the study directions
	if err := tx.Where("study_id = ?", studyID).Delete(&StudyDirection{}).Error; err != nil {
		return err
	}

	// delete the study trials
	trialIDs := []uint{}

	if err := tx.Model(&Trial{}).Where("study_id = ?", studyID).Pluck("trial_id", &trialIDs).Error; err != nil {
		return err
	}

	for _, id := range trialIDs {
		if err := DeleteTrial(id, tx); err != nil {
			return err
		}
	}

	// delete the study
	return tx.Delete(&Study{}, studyID).Error
}

func DeleteTrial(trialID uint, tx *gorm.DB) error {
	for _, model := range []any{
		&TrialSystemAttribute{},
		&TrialUserAttribute{},
		&TrialValue{},
		&TrialParam{},
		&TrialIntermediateValue{},
		&TrialHeartbeat{},
	} {
		if err := tx.Where("trial_id = ?", trialID).Delete(model).Error; err != nil {
			return err
		}
	}

	// delete the trial
	return tx.Delete(&Trial{}, trialID).Error
}

// GetBestTrial finds and returns the best trial for a given study ID,
// taking into account the study direction (minimize or maximize).
// It returns the trial, its parameters, and its value.
func GetBestTrial(studyID uint) (*BestTrialResult, error) {
	if DB == nil {
		return nil, errors.New("database connection not initialized")
	}

	// Get the study direction
	var studyDirection StudyDirection
	if err := DB.Where("study_id = ?", studyID).First(&studyDirection).Error; err != nil {
		return nil, err
	}

	// Prepare the query to find the best trial
	var bestTrial Trial
	var bestTrialValue TrialValue

	// Order depends on the study direction
	orderClause := "value ASC" // For minimize
	if studyDirection.Direction == "maximize" {
		orderClause = "value DESC"
	}

	// Join trials and trial_values tables, filter by study_id and objective,
	// order by value according to direction, and get the best one
	query := DB.Table("trials").
		Select("trials.*, trial_values.*").
		Joins("JOIN trial_values ON trials.trial_id = trial_values.trial_id").
		Where("trials.study_id = ? AND trial_values.objective = ? AND trials.state = ?",
			studyID, studyDirection.Objective, "COMPLETE").
		Order(orderClause).
		Limit(1)

	if err := query.First(&bestTrial).Error; err != nil {
		return nil, err
	}

	// Get the trial value for the best trial
	if err := DB.Where("trial_id = ? AND objective = ?",
		bestTrial.ID, studyDirection.Objective).First(&bestTrialValue).Error; err != nil {
		return nil, err
	}

	// Get all parameters for the best trial
	var trialParams []TrialParam
	if err := DB.Where("trial_id = ?", bestTrial.ID).Find(&trialParams).Error; err != nil {
		return nil, err
	}

	return &BestTrialResult{
		Trial:       bestTrial,
		TrialParams: trialParams,
		TrialValue:  bestTrialValue,
	}, nil
}

func GetTrialUserAttributes(trialID uint) (map[string]any, error) {
	var trialUserAttributes []TrialUserAttribute

	if err := DB.Where("trial_id = ?", trialID).Find(&trialUserAttributes).Error; err != nil {
		return nil, err
	}

	attributes := make(map[string]any, len(trialUserAttributes))
	for _, attr := range trialUserAttributes {
		var value any

		if err := json.Unmarshal([]byte(attr.Value), &value); err != nil {
			return nil, err
		}

		attributes[attr.Key] = value
	}

	return attributes, nil
}

func GetUserAttributesList() (attributes []string, err error) {
	err = DB.Model(&TrialUserAttribute{}).Distinct("key").Pluck("key", &attributes).Error
	return
}
