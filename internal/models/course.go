package models

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	CourseID    int       `json:"course_id" gorm:"type:int(20)"`
	Name        string    `json:"name" gorm:"type:varchar(255)"`
	TeacherID   string    `json:"teacher_id" gorm:"type:int(20)"`
	BeginDate   time.Time `json:"begin_date"`
	EndDate     time.Time `json:"end_date"`
	Description string    `json:"description"`
}
