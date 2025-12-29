package handler

import (
	"HelpStudent/internal/app/users/dao"
	"HelpStudent/internal/app/users/model"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/flamego/flamego"
	"github.com/xuri/excelize/v2"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

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

	for i, row := range rows {
		if i == 0 { // 跳过表头
			continue
		}

		// 解析用户信息
		user, err := parseUserRow(row)
		if err != nil {
			failCount++
			errorMessages = append(errorMessages, fmt.Sprintf("第%d行: %v", i+1, err))
			continue
		}

		// 将NeedSubjects数组转换为JSON字符串存储
		if user.NeedSubjects != nil {
			subjectsJSON, err := json.Marshal(user.NeedSubjects)
			if err != nil {
				failCount++
				errorMessages = append(errorMessages, fmt.Sprintf("第%d行: 科目序列化失败: %v", i+1, err))
				continue
			}
			user.NeedSubjectsDB = string(subjectsJSON)
		}

		// 创建或更新用户信息
		result := db.Where("staff_id = ?", user.StaffId).FirstOrCreate(&user)
		if result.Error != nil {
			failCount++
			errorMessages = append(errorMessages, fmt.Sprintf("第%d行: 保存用户失败: %v", i+1, result.Error))
			continue
		}

		successCount++
	}

	// 返回处理结果
	response := map[string]interface{}{
		"success":      true,
		"total":        len(rows) - 1,
		"successCount": successCount,
		"failCount":    failCount,
		"filename":     header.Filename,
	}

	if failCount > 0 {
		response["errorDetails"] = errorMessages
	}

	r.JSON(http.StatusOK, response)
}

// parseUserRow 解析单行用户数据
func parseUserRow(row []string) (model.Users, error) {
	if len(row) < 3 {
		return model.Users{}, fmt.Errorf("数据列不足")
	}

	staffId := strings.TrimSpace(row[0])
	if staffId == "" {
		return model.Users{}, fmt.Errorf("工号不能为空")
	}

	name := strings.TrimSpace(row[1])
	if name == "" {
		return model.Users{}, fmt.Errorf("姓名不能为空")
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

	return model.Users{
		StaffId:      staffId,
		Name:         name,
		NeedSubjects: needSubjects,
	}, nil
}
