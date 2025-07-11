package handler

import (
	"HelpStudent/core/auth"
	"HelpStudent/core/middleware/response"
	"HelpStudent/internal/app/permission/dto"
	"HelpStudent/internal/app/permission/service/rbac"
	"errors"
	"github.com/flamego/binding"
	"github.com/flamego/flamego"
	"strings"
)

func HandleManagerList(r flamego.Render, auth auth.Info) {
	if !rbac.CheckStaffSystemPermission(auth.StaffId, "permission", "read") {
		response.Forbidden(r)
		return
	}
	var resp []dto.ProjectManager
	managers := rbac.GetSystemManager()
	for _, manager := range managers {
		if manager == "admin" {
			continue
		}
		ent := dto.ProjectManager{
			StaffId: manager,
		}
		for _, perm := range rbac.GetStaffSystemPermissionArray(manager) {
			ent.Permissions = append(ent.Permissions, dto.KVMap{
				Key:   perm[0] + ":" + perm[1],
				Value: rbac.GetSystemActionName(perm[0], perm[1]),
			})
		}
		resp = append(resp, ent)
	}
	response.HTTPSuccess(r, resp)
}

func HandleAddManager(r flamego.Render, auth auth.Info, req dto.AddProjectManagerRequest, errs binding.Errors) {
	if errs != nil {
		response.InValidParam(r, errs)
		return
	}

	if !rbac.CheckStaffSystemPermission(auth.StaffId, "permission", "write") {
		response.Forbidden(r)
		return
	}
	var filtered [][2]string
	for _, permission := range req.Permissions {
		arr := strings.Split(permission, ":")
		if len(arr) != 2 || !rbac.SystemActionExist(arr[0], arr[1]) {
			response.InValidParam(r, errors.New("permission not exist"), permission)
			return
		} else {
			filtered = append(filtered, [2]string{arr[0], arr[1]})
		}
	}

	err := rbac.UnlinkStaffToSystemAllPermission(req.StaffId)
	if err != nil {
		response.ServiceErr(r, err)
		return
	}

	for _, permission := range filtered {
		err = rbac.LinkStaffToSystemPermission(req.StaffId, permission[0], permission[1])
		if err != nil {
			response.ServiceErr(r, err)
			return
		}
	}

	response.HTTPSuccess(r, nil)

}

func HandleRemoveManager(r flamego.Render, auth auth.Info, req dto.RemoveProjectManagerRequest) {
	if !rbac.CheckStaffSystemPermission(auth.StaffId, "permission", "write") {
		response.Forbidden(r)
		return
	}

	for _, staffId := range req.StaffIds {
		err := rbac.UnlinkStaffToSystemAllPermission(staffId)
		if err != nil {
			response.ServiceErr(r, err)
			return
		}
	}
	response.HTTPSuccess(r, nil)
}
