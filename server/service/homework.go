package service

import (
	"errors"
	"fmt"
	"homework_platform/internal/models"
	"log"
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
)

type AssignHomeworkService struct {
	CourseID    int                     `form:"courseid"`
	Name        string                  `form:"name"`
	Description string                  `form:"description"`
	BeginDate   time.Time               `form:"begindate"`
	EndDate     time.Time               `form:"enddate"`
	Files       []*multipart.FileHeader `form:"files"`
}

func (service *AssignHomeworkService) Handle(c *gin.Context) (any, error) {
	course, err := models.GetCourseByID(service.CourseID)
	if err != nil {
		return nil, err
	}
	id, _ := c.Get("ID")
	if course.TeacherID != id {
		return nil, errors.New("不能发布不是您的课程的作业")
	}
	//CourseID
	homework, err2 := models.CreateHomework(
		service.CourseID,
		service.Name,
		service.Description,
		service.BeginDate,
		service.EndDate,
	)
	if err2 != nil {
		return nil, errors.New("创建失败")
	}
	for _, f := range service.Files {
		log.Println(f.Filename)
		dst := fmt.Sprintf("./homeworkassign/%d//%s", homework.(models.Homework).ID, f.Filename)
		// 上传文件到指定的目录
		c.SaveUploadedFile(f, dst)
	}
	println(service.CourseID)
	return nil, nil
}

type HomeworkLists struct {
	CourseID int `form:"courseid"`
}

func (service *HomeworkLists) Handle(c *gin.Context) (any, error) {
	course, err := models.GetCourseByID(service.CourseID)
	if err != nil {
		return nil, err
	}
	id, _ := c.Get("ID")
	if course.TeacherID != id {
		return nil, errors.New("不能查看不是您的课程的作业")
	}
	homeworks, err2 := course.GetHomeworkLists()
	if err2 != nil {
		return nil, err2
	}
	return homeworks, nil
}
