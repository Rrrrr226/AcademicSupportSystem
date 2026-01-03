package dao

import (
	"HelpStudent/core/logx"
	"HelpStudent/internal/app/subject/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type subject struct {
	*gorm.DB
}

func (u *subject) Init(db *gorm.DB) (err error) {
	u.DB = db
	return db.AutoMigrate(&model.Subject{}, &model.UserSubject{})
}

func (d *subject) GetLinksByNames(names []string) (map[string]string, error) {
	if len(names) == 0 {
		return map[string]string{}, nil
	}

	var subjects []model.Subject
	if err := d.Where("subject_name IN ?", names).Find(&subjects).Error; err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, s := range subjects {
		result[s.SubjectName] = s.SubjectLink
	}

	logx.ServiceLogger.Infof("GetLinksByNames 查询结果: %+v", subjects)
	logx.ServiceLogger.Infof("返回的 linkMap: %+v", result)

	return result, nil
}

// GetUserSubjects 获取用户的所有科目
func (d *subject) GetUserSubjects(staffId string) ([]string, error) {
	var userSubjects []model.UserSubject
	if err := d.Where("staff_id = ?", staffId).Find(&userSubjects).Error; err != nil {
		return nil, err
	}

	var subjectNames []string
	for _, us := range userSubjects {
		subjectNames = append(subjectNames, us.SubjectName)
	}
	return subjectNames, nil
}

// GetUserSubjectsByUserId 根据 UserId 获取用户的所有科目
func (d *subject) GetUserSubjectsByUserId(userId string) ([]string, error) {
	var userSubjects []model.UserSubject
	if err := d.Where("user_id = ?", userId).Find(&userSubjects).Error; err != nil {
		return nil, err
	}

	var subjectNames []string
	for _, us := range userSubjects {
		subjectNames = append(subjectNames, us.SubjectName)
	}
	return subjectNames, nil
}

// SetUserSubjects 设置用户的科目（会覆盖原有数据）
func (d *subject) SetUserSubjects(userId, staffId string, subjectNames []string) error {
	return d.Transaction(func(tx *gorm.DB) error {
		// 删除用户原有的科目关联
		if err := tx.Where("user_id = ?", userId).Delete(&model.UserSubject{}).Error; err != nil {
			return err
		}

		// 如果没有新科目，直接返回
		if len(subjectNames) == 0 {
			return nil
		}

		// 批量创建新的关联
		var userSubjects []model.UserSubject
		for _, name := range subjectNames {
			userSubjects = append(userSubjects, model.UserSubject{
				UserId:      userId,
				StaffId:     staffId,
				SubjectName: name,
			})
		}

		return tx.Create(&userSubjects).Error
	})
}

// AddUserSubject 为用户添加一个科目
func (d *subject) AddUserSubject(userId, staffId, subjectName string) error {
	us := model.UserSubject{
		UserId:      userId,
		StaffId:     staffId,
		SubjectName: subjectName,
	}
	// 使用 OnConflict 来处理已存在的情况
	return d.Clauses(clause.OnConflict{DoNothing: true}).Create(&us).Error
}

// RemoveUserSubject 移除用户的一个科目
func (d *subject) RemoveUserSubject(userId, subjectName string) error {
	return d.Where("user_id = ? AND subject_name = ?", userId, subjectName).
		Delete(&model.UserSubject{}).Error
}

// BatchSetUserSubjects 批量设置多个用户的科目
func (d *subject) BatchSetUserSubjects(userSubjectsMap map[string]struct {
	UserId   string
	Subjects []string
}) error {
	return d.Transaction(func(tx *gorm.DB) error {
		for staffId, data := range userSubjectsMap {
			// 删除用户原有的科目关联
			if err := tx.Where("staff_id = ?", staffId).Delete(&model.UserSubject{}).Error; err != nil {
				return err
			}

			// 如果没有新科目，继续下一个用户
			if len(data.Subjects) == 0 {
				continue
			}

			// 批量创建新的关联
			var userSubjects []model.UserSubject
			for _, name := range data.Subjects {
				userSubjects = append(userSubjects, model.UserSubject{
					UserId:      data.UserId,
					StaffId:     staffId,
					SubjectName: name,
				})
			}

			if err := tx.Create(&userSubjects).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// SubjectExists 检查科目是否存在
func (d *subject) SubjectExists(subjectName string) (bool, error) {
	var count int64
	if err := d.Model(&model.Subject{}).Where("subject_name = ?", subjectName).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// SubjectsExist 批量检查科目是否存在，返回不存在的科目列表
func (d *subject) SubjectsExist(subjectNames []string) ([]string, error) {
	if len(subjectNames) == 0 {
		return nil, nil
	}

	var existingSubjects []model.Subject
	if err := d.Where("subject_name IN ?", subjectNames).Find(&existingSubjects).Error; err != nil {
		return nil, err
	}

	existingMap := make(map[string]bool)
	for _, s := range existingSubjects {
		existingMap[s.SubjectName] = true
	}

	var notExist []string
	for _, name := range subjectNames {
		if !existingMap[name] {
			notExist = append(notExist, name)
		}
	}

	return notExist, nil
}

// ImportStudentSubjects 导入学生科目（仅添加，不删除已有的）
// 返回: 成功数, 失败数, 错误列表
func (d *subject) ImportStudentSubjects(items []struct {
	StaffId     string
	SubjectName string
}) (int, int, []string) {
	var successCount, failCount int
	var errors []string

	for _, item := range items {
		us := model.UserSubject{
			StaffId:     item.StaffId,
			SubjectName: item.SubjectName,
		}

		// 使用 OnConflict 来处理已存在的情况（忽略重复）
		result := d.Clauses(clause.OnConflict{DoNothing: true}).Create(&us)
		if result.Error != nil {
			failCount++
			errors = append(errors, "学号 "+item.StaffId+" 科目 "+item.SubjectName+": "+result.Error.Error())
		} else {
			successCount++
		}
	}

	return successCount, failCount, errors
}

// GetAllSubjectNames 获取所有科目名称
func (d *subject) GetAllSubjectNames() ([]string, error) {
	var subjects []model.Subject
	if err := d.Find(&subjects).Error; err != nil {
		return nil, err
	}

	var names []string
	for _, s := range subjects {
		names = append(names, s.SubjectName)
	}
	return names, nil
}
