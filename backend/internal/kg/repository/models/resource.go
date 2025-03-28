package models

import "time"

type Video struct {
	ID               int64     `gorm:"primaryKey;autoIncrement"`
	Title            string    `gorm:"size:100;not null" json:"title"`
	PlayURL          string    `gorm:"size:255;not null" json:"play_url"`
	CoverURL         string    `gorm:"size:255;not null" json:"cover_url"`
	Description      *string   `gorm:"type:text" json:"description,omitempty"`
	KnowledgePointID int64     `gorm:"index;not null" json:"knowledge_point_id"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type Exercise struct {
	ID               int64     `gorm:"primaryKey;autoIncrement"`
	Title            string    `gorm:"size:100;not null" json:"title"`
	ExerciseURL      string    `gorm:"size:255;not null" json:"exercise_url"`
	Difficulty       string    `gorm:"type:enum('easy', 'medium', 'hard');not null" json:"difficulty"`
	Description      *string   `gorm:"type:text" json:"description,omitempty"`
	KnowledgePointID int64     `gorm:"index;not null" json:"knowledge_point_id"` // Assuming knowledge_point_id is of type BIGINT in MySQL
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type Courseware struct {
	ID               int64     `gorm:"primaryKey;autoIncrement"`
	Title            string    `gorm:"size:100;not null" json:"title"`
	CoursewareURL    string    `gorm:"size:255;not null" json:"courseware_url"`
	Description      *string   `gorm:"type:text" json:"description,omitempty"`
	KnowledgePointID int64     `gorm:"index;not null" json:"knowledge_point_id"` // Assuming knowledge_point_id is of type BIGINT in MySQL
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
