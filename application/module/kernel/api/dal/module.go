package dal

import (
	"fmt"
	"magiccenter/common"
	"magiccenter/common/model"
)

// DeleteModule 删除指定Module
func DeleteModule(helper common.DBHelper, id string) bool {
	result := false
	sql := fmt.Sprintf("delete from module where id='%s'", id)
	num, ret := helper.Execute(sql)
	if num == 1 && ret {
		sql = fmt.Sprintf("delete from block where owner='%s'", id)
		helper.Execute(sql)

		result = true
	}

	return result
}

// QueryModule 查询指定Module
func QueryModule(helper common.DBHelper, id string) (model.Module, bool) {
	m := model.Module{}
	sql := fmt.Sprintf("select id, name, description, uri, type, status from module where id='%s'", id)
	helper.Query(sql)

	result := false
	if helper.Next() {
		helper.GetValue(&m.ID, &m.Name, &m.Description, &m.URL, &m.Type, &m.Status)
		result = true
	}

	return m, result
}

// QueryAllModule 查询所有Module
func QueryAllModule(helper common.DBHelper) []model.Module {
	moduleList := []model.Module{}

	sql := fmt.Sprintf("select id, name, description, uri, type, status from module")
	helper.Query(sql)

	for helper.Next() {
		m := model.Module{}
		helper.GetValue(&m.ID, &m.Name, &m.Description, &m.URL, &m.Type, &m.Status)

		moduleList = append(moduleList, m)
	}

	return moduleList
}

// CreateModule 新建Module
func CreateModule(helper common.DBHelper, id, name, description, url string, mType, mStatus int) (model.Module, bool) {
	result := false
	m := model.Module{}

	_, found := QueryModule(helper, id)
	if !found {
		sql := fmt.Sprintf("insert into module(id, name, description, uri, type, status) values ('%s','%s','%s','%s',%d,%d)", id, name, description, url, mType, mStatus)
		num, ret := helper.Execute(sql)
		result = (num == 1 && ret)
	}

	return m, result
}

// SaveModule 保存Module
func SaveModule(helper common.DBHelper, m model.Module) (model.Module, bool) {
	result := false
	_, found := QueryModule(helper, m.ID)
	if found {
		sql := fmt.Sprintf("update module set name ='%s', description ='%s', uri='%s', type =%d, status =%d where Id='%s'", m.Name, m.Description, m.URL, m.Type, m.Status, m.ID)
		num, ret := helper.Execute(sql)
		result = (num == 1 && ret)
	}

	return m, result
}