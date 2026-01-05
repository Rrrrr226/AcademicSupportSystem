package handler

import (
	"HelpStudent/core/auth"
	"HelpStudent/core/logx"
	"HelpStudent/core/middleware/response"
	"HelpStudent/internal/app/managers/dao"
	"HelpStudent/internal/app/managers/dto"
	"HelpStudent/internal/app/managers/model"
	subjectDAO "HelpStudent/internal/app/subject/dao"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/flamego/binding"
	"github.com/flamego/flamego"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// HandleImportStudentSubjectsExcel 处理Excel导入学生科目
func HandleImportStudentSubjectsExcel(c flamego.Context, r flamego.Render) {
	// 从 FormFile 获取文件
	file, header, err := c.Request().FormFile("file")
	if err != nil {
		response.HTTPFail(r, 400002, "获取上传文件失败")
		return
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	// 检查文件扩展名
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".xlsx") &&
		!strings.HasSuffix(strings.ToLower(header.Filename), ".xls") {
		response.HTTPFail(r, 400003, "仅支持Excel文件(.xlsx, .xls)")
		return
	}

	// 读取文件内容
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.HTTPFail(r, 400004, "读取文件内容失败")
		return
	}

	// 打开Excel文件
	f, err := excelize.OpenReader(bytes.NewReader(fileBytes))
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.HTTPFail(r, 400005, "打开Excel文件失败")
		return
	}
	defer func() {
		_ = f.Close()
	}()

	// 获取第一个工作表
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		response.HTTPFail(r, 400006, "Excel文件中没有工作表")
		return
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.HTTPFail(r, 400007, "读取Excel内容失败")
		return
	}

	if len(rows) < 2 {
		response.HTTPFail(r, 400008, "Excel文件没有数据行")
		return
	}

	// 解析表头，查找学号和科目名称列
	headerRow := rows[0]
	staffIdColIdx := -1
	subjectColIdx := -1

	for idx, cell := range headerRow {
		cellTrimmed := strings.TrimSpace(cell)
		if cellTrimmed == "学号" || cellTrimmed == "staff_id" || cellTrimmed == "StaffId" {
			staffIdColIdx = idx
		}
		if cellTrimmed == "科目名称" || cellTrimmed == "科目" || cellTrimmed == "subject_name" || cellTrimmed == "SubjectName" {
			subjectColIdx = idx
		}
	}

	if staffIdColIdx == -1 {
		response.HTTPFail(r, 400009, "Excel缺少学号列（列名应为：学号/staff_id/StaffId）")
		return
	}

	if subjectColIdx == -1 {
		response.HTTPFail(r, 400010, "Excel缺少科目名称列（列名应为：科目名称/科目/subject_name/SubjectName）")
		return
	}

	// 解析数据行
	type studentSubject struct {
		StaffId     string
		SubjectName string
	}
	var importData []studentSubject
	var subjectNames []string
	subjectSet := make(map[string]bool)

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) <= staffIdColIdx || len(row) <= subjectColIdx {
			continue
		}

		staffId := strings.TrimSpace(row[staffIdColIdx])
		subjectName := strings.TrimSpace(row[subjectColIdx])

		if staffId == "" || subjectName == "" {
			continue
		}

		importData = append(importData, studentSubject{
			StaffId:     staffId,
			SubjectName: subjectName,
		})

		if !subjectSet[subjectName] {
			subjectSet[subjectName] = true
			subjectNames = append(subjectNames, subjectName)
		}
	}

	if len(importData) == 0 {
		response.HTTPFail(r, 400011, "没有有效的导入数据")
		return
	}

	// 验证所有科目是否存在
	missingSubjects, err := subjectDAO.Subject.SubjectsExist(subjectNames)
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	if len(missingSubjects) > 0 {
		response.HTTPFail(r, 400012, fmt.Sprintf("以下科目不存在：%s", strings.Join(missingSubjects, ", ")))
		return
	}

	// 批量导入 - 构造导入数据
	var importItems []struct {
		StaffId     string
		SubjectName string
	}
	for _, item := range importData {
		importItems = append(importItems, struct {
			StaffId     string
			SubjectName string
		}{
			StaffId:     item.StaffId,
			SubjectName: item.SubjectName,
		})
	}

	successCount, failCount, errorMsgs := subjectDAO.Subject.ImportStudentSubjects(importItems)

	response.HTTPSuccess(r, dto.ImportStudentSubjectsResponse{
		Total:        len(importData),
		SuccessCount: successCount,
		FailCount:    failCount,
		Errors:       errorMsgs,
	})
}

// HandleAddManager 添加管理员
func HandleAddManager(r flamego.Render, c flamego.Context, req dto.AddManagerRequest, errs binding.Errors, authInfo auth.Info) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}
	if authInfo.Uid != req.StaffId {
		response.HTTPFail(r, 403001, "不能添加自己")
		return
	}

	// 检查用户名是否已存在
	existing, err := dao.Managers.GetManagerByStaffID(req.StaffId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	if existing != nil {
		response.HTTPFail(r, 401004, "用户名已存在")
		return
	}

	// 创建管理员
	manager := &model.Managers{
		StaffId: req.StaffId, // 使用 userId 作为 StaffId
	}

	if err := dao.Managers.WithContext(c.Request().Context()).Create(manager).Error; err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	response.HTTPSuccess(r, dto.AddManagerResponse{
		Id: manager.StaffId,
	})
}

// HandleDeleteManager 删除管理员
func HandleDeleteManager(r flamego.Render, c flamego.Context, req dto.DeleteManagerRequest, errs binding.Errors, authInfo auth.Info) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	// 防止删除自己
	if req.StaffId == authInfo.Uid {
		response.HTTPFail(r, 403001, "不能删除自己")
		return
	}

	// 检查目标管理员是否存在
	existing, err := dao.Managers.GetManagerByStaffID(req.StaffId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.HTTPFail(r, 404001, "管理员不存在")
			return
		}
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	if existing == nil {
		response.HTTPFail(r, 404001, "管理员不存在")
		return
	}

	// 删除管理员
	if err := dao.Managers.DeleteManagetByStaffId(req.StaffId); err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	response.HTTPSuccess(r, nil)
}

// HandleGetManagerList 获取管理员列表
func HandleGetManagerList(r flamego.Render, c flamego.Context, authInfo auth.Info) {
	if authInfo.Uid == "" {
		response.HTTPFail(r, 403002, "permission denied")
		return
	}
	managers, total, err := dao.Managers.GetAllManagers()
	if err != nil {
		logx.SystemLogger.CtxError(c.Request().Context(), err)
		response.ServiceErr(r, err)
		return
	}

	var list []dto.ManagerItem
	for _, m := range managers {
		list = append(list, dto.ManagerItem{
			StaffId: m.StaffId,
		})
	}

	response.HTTPSuccess(r, dto.ManagerListResponse{
		Managers: list,
		Total:    total,
	})
}

// HandleDownloadTemplate 下载学生科目导入模板
func HandleDownloadTemplate(c flamego.Context, w http.ResponseWriter) {
	// 创建新的Excel文件
	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	sheetName := "学生科目导入模板"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("创建工作表失败"))
		return
	}

	// 设置表头
	headers := []string{"学号", "科目名称"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// 添加示例数据
	exampleData := [][]string{
		{"22050626", "高等数学"},
		{"22050627", "数据结构"},
		{"22050626", "大学物理"},
	}

	for i, row := range exampleData {
		for j, value := range row {
			cell := fmt.Sprintf("%c%d", 'A'+j, i+2)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// 设置列宽
	f.SetColWidth(sheetName, "A", "A", 15)
	f.SetColWidth(sheetName, "B", "B", 20)

	// 设置表头样式
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err == nil {
		f.SetCellStyle(sheetName, "A1", "B1", headerStyle)
	}

	// 设置活动工作表
	f.SetActiveSheet(index)
	// 删除默认的Sheet1
	f.DeleteSheet("Sheet1")

	// 设置响应头
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=student_subject_import_template.xlsx")
	w.Header().Set("Content-Transfer-Encoding", "binary")

	// 将Excel文件写入响应
	if err := f.Write(w); err != nil {
		logx.SystemLogger.Error("写入Excel文件失败", err)
	}
}
