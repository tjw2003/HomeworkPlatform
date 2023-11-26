package service

import (
	"homework_platform/internal/jwt"
	"homework_platform/internal/models"
	"time"

	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserLoginService struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func (service *UserLoginService) Handle(c *gin.Context) (any, error) {
	log.Printf("[UserLoginService]: %v, %v\n", service.Username, service.Password)
	var user models.User
	var err error

	if user, err = models.GetUserByUsername(service.Username); err == gorm.ErrRecordNotFound {
		return nil, err
	}

	if !user.CheckPassword(service.Password) {
		return nil, errors.New("incorrect password")
	}

	var jwtToken string
	jwtToken, err = jwt.CreateToken(user.ID) //根据用id创建jwt
	if err != nil {
		return nil, err
	}

	res := make(map[string]any)
	res["token"] = jwtToken //之后解码token验证和user是否一致
	res["user"] = user
	// res["user_name"] = user.Username
	log.Printf("登陆成功")
	return res, nil
}

// 自己修改密码
type UserselfupdateService struct {
	UserName    string `form:"userName"`
	OldPassword string `form:"oldPassword"` // 旧码
	NewPassword string `form:"newPassword"` //新密码
}

func (service *UserselfupdateService) Handle(c *gin.Context) (any, error) {
	user, err := models.GetUserByUsername(service.UserName)
	if err != nil {
		return nil, errors.New("该用户不存在")
	}
	//验证密码
	passwordCheck := user.CheckPassword(service.OldPassword)
	if !passwordCheck {
		return nil, errors.New("密码错误")
	}
	//修改密码
	result := user.ChangePassword(service.NewPassword)
	if !result {
		return nil, errors.New("修改失败")
	}
	res := make(map[string]any)
	res["msg"] = "修改成功"
	return res, nil
}

type GetUserService struct {
	ID uint `uri:"id" binding:"required"`
}

func (service *GetUserService) Handle(c *gin.Context) (any, error) {
	return models.GetUserByID(service.ID)
}

type UserRegisterService struct {
	Username string `form:"username"` // 用户名
	Password string `form:"password"` // 密码
}

func (service *UserRegisterService) Handle(c *gin.Context) (any, error) {
	_, err := models.CreateUser(service.Username, service.Password)
	return nil, err
}

type GetUserCoursesService struct {
	ID uint `uri:"id" binding:"required"`
}

func (service *GetUserCoursesService) Handle(c *gin.Context) (any, error) {
	user, err := models.GetUserByID(service.ID)
	if err != nil {
		return nil, err
	}
	return user.GetCourses()
}

type UpdateSignature struct {
	Signature string `form:"signature"`
}

func (Service *UpdateSignature) Handle(c *gin.Context) (any, error) {
	id, exist := c.Get("ID")
	if !exist {
		return nil, errors.New("不存在id")
	}
	user, err := models.GetUserByID(id.(uint))
	if err != nil {
		return nil, err
	}
	if err := user.ChangeSignature(Service.Signature); err != nil {
		return nil, err
	}
	return nil, nil
}

type GetUserNotifications struct {
	ID uint `uri:"id" binding:"required"`
}

type Notifications struct {
	Type string `json:"type"`

	TeachingHomeworkListsToFinish  []models.Homework `json:"homeworkInProgress"`
	TeachingHomeworkListsToComment []models.Homework `json:"commentInProgress"`

	ComplaintToBeSolved []models.Complaint `json:"complaintToBeSolved"`
	ComplaintInProgress []models.Complaint `json:"complaintInProgress"`

	LeaningHomeworkListsToFinish  []models.Homework `json:"homeworksToBeCompleted"`
	LeaningHomeworkListsToComment []models.Homework `json:"commentToBeCompleted"`
}

// 返回应该尚未提交的作业,待批阅的作业和每门课最新发布的作业
func (service *GetUserNotifications) Handle(c *gin.Context) (any, error) {
	user, err := models.GetUserByID(service.ID)
	if err != nil {
		return nil, err
	}
	courses, err := user.GetCourses()
	if err != nil {
		return nil, err
	}
	var notifications Notifications
	//得到教的课中进行中和批阅中的作业
	log.Printf("len of homework%d\n", len(courses.LearningCourses))
	//得到学的课中还没完成的作业和还没批阅的作业
	for _, course := range courses.LearningCourses {
		//每门课的作业
		homeworks, err := course.GetHomeworkLists()
		if homeworks == nil {
			continue
		}
		if err != nil {
			return nil, err
		}
		for j := 0; j < len(homeworks); j++ {
			// 在批阅时段中
			if homeworks[j].CommentEndDate.After(time.Now()) {
				// 作业已经开始
				if homeworks[j].BeginDate.Before(time.Now()) {
					// 作业在提交时段内
					if homeworks[j].EndDate.After(time.Now()) {
						homework := models.GetHomeWorkSubmissionByHomeworkIDAndUserID(homeworks[j].ID, user.ID)
						// 没交作业
						if homework == nil {
							notifications.LeaningHomeworkListsToFinish =
								append(notifications.LeaningHomeworkListsToFinish, homeworks[j])
						}
					} else {
						// 评论时段内,获取所有的comment
						comments, err := models.GetCommentListsByUserIDAndHomeworkID(user.ID, homeworks[j].ID)
						if err != nil {
							return nil, err
						}
						// 如果有score==-1就代表尚未完成评论
						for i := 0; i < len(comments); i++ {
							if comments[i].Score == -1 {
								notifications.LeaningHomeworkListsToComment =
									append(notifications.TeachingHomeworkListsToComment, homeworks[j])
								break
							}
						}
					}
				}
			}
		}
	}
	//得到老师的课正在进行的作业
	for _, course := range courses.TeachingCourses {
		// 教的课中的作业
		homeworks, err := course.GetHomeworkLists()
		if homeworks == nil {
			continue
		}
		if err != nil {
			return nil, err
		}
		for j := 0; j < len(homeworks); j++ {
			// comment尚未结束
			if homeworks[j].CommentEndDate.After(time.Now()) {
				//作业已经开始
				if homeworks[j].BeginDate.Before(time.Now()) {
					//在提交时段内
					if homeworks[j].EndDate.After(time.Now()) {
						notifications.TeachingHomeworkListsToFinish =
							append(notifications.TeachingHomeworkListsToFinish, homeworks[j])
					} else {
						notifications.TeachingHomeworkListsToComment =
							append(notifications.TeachingHomeworkListsToComment, homeworks[j])
					}
				}
			}
		}
	}
	//得到老师待审核的complaint
	notifications.ComplaintToBeSolved, err = models.GetComplaintByTeacherID(user.ID)
	if err != nil {
		return nil, err
	}
	//得到学生还未被处理的complaint
	notifications.ComplaintInProgress, err = models.GetComplaintByUserID(user.ID)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}
