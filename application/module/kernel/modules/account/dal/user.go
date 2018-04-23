package dal

import (
	"database/sql"
	"fmt"
	"time"

	"muidea.com/magicCenter/application/common/dbhelper"
	"muidea.com/magicCenter/foundation/util"
	"muidea.com/magicCommon/model"
)

func loadUserID(helper dbhelper.DBHelper) int {
	var maxID sql.NullInt64
	sql := fmt.Sprintf(`select max(id) from account_user`)
	helper.Query(sql)
	defer helper.Finish()

	if helper.Next() {
		helper.GetValue(&maxID)
	}

	return int(maxID.Int64)
}

//QueryAllUser 查询全部用户信息
func QueryAllUser(helper dbhelper.DBHelper) []model.UserDetail {
	userList := []model.UserDetail{}
	sql := fmt.Sprintf("select id, account, nickname, email, groups, status, registertime from account_user")
	helper.Query(sql)
	defer helper.Finish()

	for helper.Next() {
		user := model.UserDetail{}
		groups := ""
		helper.GetValue(&user.ID, &user.Account, &user.Name, &user.Email, &groups, &user.Status, &user.RegisterTime)
		user.Group, _ = util.Str2IntArray(groups)

		userList = append(userList, user)
	}

	return userList
}

// QueryUsers 查询指定用户
func QueryUsers(helper dbhelper.DBHelper, ids []int) []model.User {
	userList := []model.User{}
	sql := fmt.Sprintf("select id, nickname from account_user where id in(%s)", util.IntArray2Str(ids))
	helper.Query(sql)
	defer helper.Finish()

	for helper.Next() {
		user := model.User{}
		helper.GetValue(&user.ID, &user.Name)
		userList = append(userList, user)
	}

	return userList
}

// QueryUserByAccount 根据账号查询用户信息
func QueryUserByAccount(helper dbhelper.DBHelper, account, password string) (model.UserDetail, bool) {
	user := model.UserDetail{}

	sql := fmt.Sprintf("select id, account, nickname, email, groups, status, registertime from account_user where account='%s' and password='%s'", account, password)
	helper.Query(sql)
	defer helper.Finish()

	result := false
	if helper.Next() {
		groups := ""
		helper.GetValue(&user.ID, &user.Account, &user.Name, &user.Email, &groups, &user.Status, &user.RegisterTime)
		user.Group, _ = util.Str2IntArray(groups)
		result = true
	}

	return user, result
}

// QueryUserByGroup 查询指定分组下的用户信息
func QueryUserByGroup(helper dbhelper.DBHelper, groupID int) []model.User {
	userList := []model.User{}

	sql := fmt.Sprintf("select id, nickname from `account_user` where groups like '%d' union select id, nickname from `account_user` where groups like '%d,%%' union select id, nickname from `account_user` where groups like '%%,%d,%%' union select id, nickname from `account_user` where groups like '%%,%d'", groupID, groupID, groupID, groupID)
	helper.Query(sql)
	defer helper.Finish()

	for helper.Next() {
		user := model.User{}
		helper.GetValue(&user.ID, &user.Name)
		userList = append(userList, user)
	}

	return userList
}

// QueryUserByID 根据用户ID查询用户信息
func QueryUserByID(helper dbhelper.DBHelper, id int) (model.UserDetail, bool) {
	user := model.UserDetail{}

	sql := fmt.Sprintf("select id, account, nickname, email, groups, status, registertime from account_user where id=%d", id)
	helper.Query(sql)
	defer helper.Finish()

	result := false
	if helper.Next() {
		groups := ""
		helper.GetValue(&user.ID, &user.Account, &user.Name, &user.Email, &groups, &user.Status, &user.RegisterTime)
		user.Group, _ = util.Str2IntArray(groups)
		result = true
	}

	return user, result
}

// DeleteUser 删除用户，根据用户ID
func DeleteUser(helper dbhelper.DBHelper, id int) bool {
	sql := fmt.Sprintf("delete from account_user where id =%d", id)
	_, ret := helper.Execute(sql)
	return ret
}

// DeleteUserByAccount 删除用户，根据用户账号&密码
func DeleteUserByAccount(helper dbhelper.DBHelper, account, password string) bool {
	sql := fmt.Sprintf("delete from account_user where account ='%s' and password='%s'", account, password)
	_, ret := helper.Execute(sql)
	return ret
}

// CreateUser 创建新用户，根据用户信息和密码
func CreateUser(helper dbhelper.DBHelper, account, password, email string, groups []int) (model.UserDetail, bool) {
	user := model.UserDetail{Account: account, Email: email, Group: groups}
	sql := fmt.Sprintf("select id from account_user where account='%s'", account)
	helper.Query(sql)

	if helper.Next() {
		helper.Finish()
		return user, false
	}
	helper.Finish()

	gVal := util.IntArray2Str(groups)
	createTime := time.Now().Format("2006-01-02 15:04:05")
	user.RegisterTime = createTime

	id := allocUserID()
	// insert
	sql = fmt.Sprintf("insert into account_user(id, account, password, nickname, email, groups, status, registertime) values (%d, '%s', '%s', '%s', '%s', '%s', %d, '%s')", id, account, password, "", email, gVal, 0, createTime)
	_, result := helper.Execute(sql)
	if !result {
		return user, false
	}
	user.ID = id

	return user, result
}

// SaveUser 保存用户信息
func SaveUser(helper dbhelper.DBHelper, user model.UserDetail) (model.UserDetail, bool) {
	gVal := util.IntArray2Str(user.Group)
	// modify
	sql := fmt.Sprintf("update account_user set nickname='%s', email='%s', groups='%s', status=%d where id =%d", user.Name, user.Email, gVal, user.Status, user.ID)
	num, result := helper.Execute(sql)

	return user, result && num == 1
}

// SaveUserWithPassword 保存用户信息
func SaveUserWithPassword(helper dbhelper.DBHelper, user model.UserDetail, password string) (model.UserDetail, bool) {
	gVal := util.IntArray2Str(user.Group)
	// modify
	sql := fmt.Sprintf("update account_user set password='%s', nickname='%s', email='%s', groups='%s', status=%d where id =%d", password, user.Name, user.Email, gVal, user.Status, user.ID)
	num, result := helper.Execute(sql)

	return user, result && num == 1
}

// GetLastRegisterUser 获取最新的用户
func GetLastRegisterUser(helper dbhelper.DBHelper, count int) []model.UserDetail {
	userList := []model.UserDetail{}
	sql := fmt.Sprintf("select id, account, nickname, email, groups, status, registertime from account_user order by registertime desc limit %d", count)
	helper.Query(sql)
	defer helper.Finish()

	for helper.Next() {
		user := model.UserDetail{}
		groups := ""
		helper.GetValue(&user.ID, &user.Account, &user.Name, &user.Email, &groups, &user.Status, &user.RegisterTime)
		user.Group, _ = util.Str2IntArray(groups)

		userList = append(userList, user)
	}

	return userList
}
