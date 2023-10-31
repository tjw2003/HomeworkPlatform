package models

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type Homework struct {
	gorm.Model
	CourseID       uint      `json:"courseId" gorm:"type:int(20)"`
	Name           string    `json:"name" gorm:"type:varchar(255)"`
	Description    string    `json:"description"`
	BeginDate      time.Time `json:"beginDate"`
	EndDate        time.Time `json:"endDate"`
	CommentEndDate time.Time `json:"commentEndDate"`
	Assigned       int       `json:"-" gorm:"default:-1"`
	// A homework has many submissions
	// Also check homeworkSubmission.go
	// Check: https://gorm.io/docs/has_many.html
	HomeworkSubmissions []HomeworkSubmission `json:"-"`
	FilePaths           []string             `json:"file_paths" gorm:"-"`
}

func (homework *Homework) UpdateInformation(name string, desciption string, beginDate time.Time, endDate time.Time, commentendate time.Time) bool {
	result := DB.Model(&homework).Updates(Homework{
		Name: name, Description: desciption, BeginDate: beginDate, EndDate: endDate, CommentEndDate: commentendate,
	})
	return result.Error == nil
}

func (homeworkd Homework) Deleteself() error {
	res := DB.Delete(&homeworkd)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func CreateHomework(id uint, name string, description string,
	begindate time.Time, endtime time.Time, commentendate time.Time) (any, error) {
	if begindate.After(endtime) {
		return nil, errors.New("结束时间不可早于开始时间")
	}
	if endtime.After(commentendate) {
		return nil, errors.New("评论开始时间不可早于结束时间")
	}
	newhomework := Homework{
		CourseID:       id,
		Name:           name,
		Description:    description,
		BeginDate:      begindate,
		EndDate:        endtime,
		CommentEndDate: commentendate,
	}
	res := DB.Create(&newhomework)
	if res.Error != nil {
		return nil, errors.New("创建失败")
	}

	return newhomework, nil
}

func GetHomeworkByID(id uint) (Homework, error) {
	log.Printf("正在查找<Homework>(ID = %d)...", id)
	var work Homework

	res := DB.First(&work, id)
	if res.Error != nil {
		log.Printf("查找失败: %s", res.Error)
		return work, res.Error
	}
	log.Printf("查找完成: <Homeworkd>(homeworkName = %s)", work.Name)
	return work, nil
}

func GetHomeworkByIDWithSubmissionLists(id uint) (Homework, error) {
	log.Printf("正在查找<Homework>(ID = %d)...", id)
	var work Homework

	res := DB.Preload("HomeworkSubmissions").First(&work, id)
	if res.Error != nil {
		log.Printf("查找失败: %s", res.Error)
		return work, res.Error
	}
	log.Printf("查找完成: <Homeworkd>(homeworkName = %s)", work.Name)
	return work, nil
}
