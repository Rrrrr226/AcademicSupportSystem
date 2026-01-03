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
func HandleImportStudentSubjectsExcel(r flamego.Render, c flamego.Context, w http.ResponseWriter, req *http.Request) {
	// 解析multipart form
	err := req.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		response.HTTPFail(r, 400001, "文件解析失败")
		return
	}

	file, header, err := req.FormFile("file")
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
func HandleAddManager(r flamego.Render, c flamego.Context, req dto.AddManagerRequest, errs binding.Errors) {
	if errs != nil {
		response.InValidParam(r, errs)
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
func HandleGetManagerList(r flamego.Render, c flamego.Context) {
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
