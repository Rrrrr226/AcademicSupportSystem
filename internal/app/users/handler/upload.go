package handler

import (
	subjectDAO "HelpStudent/internal/app/subject/dao"
	"HelpStudent/internal/app/users/dao"
	"HelpStudent/internal/app/users/model"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/flamego/flamego"
	"github.com/xuri/excelize/v2"
)

// UserSubjectData 临时结构体用于解析上传数据
type UserSubjectData struct {
	StaffId      string
	Name         string
	NeedSubjects []string
}

// HandleUploadUserXLSX 处理上传的用户信息XLSX文件
func HandleUploadUserXLSX(r flamego.Render, req *http.Request) {
	dbUsers := dao.Users
	db := dbUsers.DB

	// 获取上传的文件
	file, header, err := req.FormFile("user_file")
	if err != nil {
		r.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "获取上传文件失败: " + err.Error(),
		})
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// 检查文件类型
	contentType := header.Header.Get("Content-Type")
	if contentType != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		r.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "仅支持XLSX格式文件",
		})
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		r.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "读取文件内容失败: " + err.Error(),
		})
		return
	}

	reader := bytes.NewReader(fileBytes)

	// 解析Excel文件
	f, err := excelize.OpenReader(reader)
	if err != nil {
		r.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "解析Excel文件失败: " + err.Error(),
		})
		return
	}
	defer func(f *excelize.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	// 获取第一个工作表的数据
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		r.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "读取工作表失败: " + err.Error(),
		})
		return
	}

	// Excel列顺序: StaffId, Name, NeedSubjects
	if len(rows) == 0 || len(rows[0]) < 3 {
		r.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Excel文件格式不正确，必须包含StaffId、Name和NeedSubjects列",
		})
		return
	}

	// 处理每一行数据
	var successCount, failCount int
	var errorMessages []string

	// 收集所有用户的科目数据，用于批量处理
	userSubjectsMap := make(map[string]struct {
		UserId   string
		Subjects []string
	})

	for i, row := range rows {
		if i == 0 { // 跳过表头
			continue
		}

		// 解析用户信息
		userData, err := parseUserRowData(row)
		if err != nil {
			failCount++
			errorMessages = append(errorMessages, fmt.Sprintf("第%d行: %v", i+1, err))
			continue
		}

		// 创建用户对象（不包含科目字段）
		user := model.Users{
			StaffId: userData.StaffId,
			Name:    userData.Name,
		}

		// 创建或更新用户信息
		var existingUser model.Users
		result := db.Where("staff_id = ?", user.StaffId).First(&existingUser)
		if result.Error != nil {
			// 用户不存在，创建新用户
			result = db.Create(&user)
			if result.Error != nil {
				failCount++
				errorMessages = append(errorMessages, fmt.Sprintf("第%d行: 创建用户失败: %v", i+1, result.Error))
				continue
			}
			existingUser = user
		} else {
			// 用户存在，更新姓名
			db.Model(&existingUser).Update("name", userData.Name)
		}

		// 收集用户科目数据
		if len(userData.NeedSubjects) > 0 {
			userSubjectsMap[userData.StaffId] = struct {
				UserId   string
				Subjects []string
			}{
				UserId:   existingUser.ID,
				Subjects: userData.NeedSubjects,
			}
		}

		successCount++
	}

	// 批量设置用户科目
	if len(userSubjectsMap) > 0 {
		if err := subjectDAO.Subject.BatchSetUserSubjects(userSubjectsMap); err != nil {
			// 记录错误但不影响整体结果
			errorMessages = append(errorMessages, fmt.Sprintf("批量设置用户科目失败: %v", err))
		}
	}

	// 返回处理结果
	response := map[string]interface{}{
		"success":      true,
		"total":        len(rows) - 1,
		"successCount": successCount,
		"failCount":    failCount,
		"filename":     header.Filename,
	}

	if failCount > 0 || len(errorMessages) > 0 {
		response["errorDetails"] = errorMessages
	}

	r.JSON(http.StatusOK, response)
}

// parseUserRowData 解析单行用户数据
func parseUserRowData(row []string) (UserSubjectData, error) {
	if len(row) < 3 {
		return UserSubjectData{}, fmt.Errorf("数据列不足")
	}

	staffId := strings.TrimSpace(row[0])
	if staffId == "" {
		return UserSubjectData{}, fmt.Errorf("工号不能为空")
	}

	name := strings.TrimSpace(row[1])
	if name == "" {
		return UserSubjectData{}, fmt.Errorf("姓名不能为空")
	}

	// 解析NeedSubjects(用逗号分隔的字符串)
	var needSubjects []string
	subjectsStr := strings.TrimSpace(row[2])
	if subjectsStr != "" {
		parts := strings.Split(subjectsStr, "，")
		for _, part := range parts {
			subject := strings.TrimSpace(part)
			if subject != "" {
				needSubjects = append(needSubjects, subject)
			}
		}
	}

	return UserSubjectData{
		StaffId:      staffId,
		Name:         name,
		NeedSubjects: needSubjects,
	}, nil
}
