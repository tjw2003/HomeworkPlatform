package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Name        string    `json:"name"`
	BeginDate   time.Time `json:"begin_date"`
	EndDate     time.Time `json:"end_date"`
	Description string    `json:"description"`

	// A teacher has many course
	// Also check user.go
	// Check: https://gorm.io/docs/has_many.html
	TeacherID uint

	// A student has many Course, a course has many students
	// Also check user.go
	// Check: https://gorm.io/docs/many_to_many.html
	Students []*User `json:"-" gorm:"many2many:user_courses;"`

	// A course has many homework
	// Also check homework.go
	// Check: https://gorm.io/docs/has_many.html
	Homeworks []Homework `json:"-"`
}

func GetAllStudents(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Model(&User{}).Preload("Course").Find(&users).Error
	return users, err
}

func CreateCourse(name string, begindate time.Time,
	enddate time.Time, description string, teachderID int) error {
	c := Course{
		Name:        name,
		BeginDate:   begindate,
		EndDate:     enddate,
		Description: description,
		TeacherID:   uint(teachderID),
	}
	res := DB.Create(&c)
	if res.Error != nil {
		return errors.New("创建失败")
	}

	return nil
}
