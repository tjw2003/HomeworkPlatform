package service

import (
	"errors"
	"homework_platform/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentService struct {
	Score                int    `form:"score" gorm:"default:-1"`
	Comment              string `form:"comment"`
	HomeworkSubmissionID uint   `uri:"id" binding:"required"`
}

func (service *CommentService) Handle(c *gin.Context) (any, error) {
	if service.Score < 0 || service.Score > 100 {
		return nil, errors.New("无效分数")
	}
	homewroksubmission := models.GetHomeWorkSubmissionByID(service.HomeworkSubmissionID)
	homework, res1 := models.GetHomeworkByID(homewroksubmission.HomeworkID)
	if res1 != nil {
		return nil, res1
	}
	if homework.CommentEndDate.Before(time.Now()) {
		return nil, errors.New("超时批阅")
	}
	if homewroksubmission == nil {
		return nil, errors.New("没有找到该作业号")
	}
	id, _ := c.Get("ID")
	// comment是预先分配好的,所以不需要自我创建
	comment, res := models.GetCommentByUserIDAndHomeworkSubmissionID(id.(uint), service.HomeworkSubmissionID)
	if res == nil {
		res := comment.(models.Comment).UpdateSelf(service.Comment, service.Score)
		num := models.GetCommentNum(service.HomeworkSubmissionID)
		if num == 3 {
			homewroksubmission.CalculateGrade()
		}
		return nil, res
	}
	return nil, res
}

type GetCommentListsService struct {
	HomeworkID uint `uri:"id" binding:"required"`
}

func (service *GetCommentListsService) Handle(c *gin.Context) (any, error) {
	println("123\n")
	id, _ := c.Get("ID")
	err := models.AssignComment(service.HomeworkID)
	if err != nil {
		return nil, err
	}

	commentLists, res := models.GetCommentListsByUserIDAndHomeworkID(id.(uint), service.HomeworkID)
	if res != nil {
		return nil, res
	}
	var homework_submission []models.HomeworkSubmission
	for _, comment := range commentLists {
		homework_submission = append(homework_submission, *models.GetHomeWorkSubmissionByID(comment.HomeworkSubmissionID))
	}
	m := make(map[string]any)
	m["homework_submission"] = homework_submission
	m["comment_lists"] = commentLists
	return m, nil
}

// type GetCommentHomeworkSubmissionService struct {
// 	HomeworkSubmissionID uint `uri:"id" binding:"required"`
// }

// func (service *GetCommentHomeworkSubmissionService) Handle(c *gin.Context) (any, error) {
// 	homework_submission := models.GetHomeWorkSubmissionByID(service.HomeworkSubmissionID)
// 	return homework_submission, nil
// }
