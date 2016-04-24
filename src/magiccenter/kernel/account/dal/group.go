package dal

import (
	"fmt"
	"magiccenter/util/modelhelper"
	"magiccenter/kernel/account/model"
)


func QueryAllGroupInfo(helper modelhelper.Model) []model.GroupInfo {
	groupInfoList := []model.GroupInfo{}
	sql := fmt.Sprintf("select id,name,catalog from `group`")
	helper.Query(sql)

	for helper.Next() {
		info := model.GroupInfo{}
		helper.GetValue(&info.Id, &info.Name, &info.Type)
		
		groupInfoList = append(groupInfoList, info)
	}
	
	for i, _ := range groupInfoList {
		info := &groupInfoList[i]
		
		users := QueryUserByGroup(helper,info.Id)		
		info.UserCount = len(users)
	}
	
	return groupInfoList
}

func QueryGroupById(helper modelhelper.Model, id int) (model.Group, bool) {
	group := model.Group{}
	sql := fmt.Sprintf("select id,name,catalog from `group` where id=%d",id)
	helper.Query(sql)

	result := false;
	if helper.Next() {
		helper.GetValue(&group.Id, &group.Name, &group.Type)
		result = true
	}

	return group, result
}

func DeleteGroup(helper modelhelper.Model, id int) bool {
	sql := fmt.Sprintf("delete from `group` where id =%d", id)
	_, result := helper.Execute(sql)
		
	return result	
}

func SaveGroup(helper modelhelper.Model, group model.Group) bool {
	sql := fmt.Sprintf("select id from `group` where id=%d", group.Id)
	helper.Query(sql)

	result := false;
	if helper.Next() {
		var id = 0
		helper.GetValue(&id)
		result = true
	}

	if !result {
		// insert
		group.Type = 1
		sql = fmt.Sprintf("insert into `group` (name,catalog) values ('%s',%d)", group.Name, group.Type)
	} else {
		// modify
		sql = fmt.Sprintf("update `group` set name ='%s' where id=%d", group.Name, group.Id)
	}
	
	_, ret := helper.Execute(sql)
	return ret
}
