package dal

import (
	"fmt"

	"muidea.com/magicCenter/application/common/dbhelper"
	"muidea.com/magicCenter/application/common/model"
	"muidea.com/magicCenter/foundation/util"
)

//QueryAllUser 查询全部用户信息
func QueryAllUser(helper dbhelper.DBHelper) []model.UserDetail {
	userList := []model.UserDetail{}
	sql := fmt.Sprintf("select id, account, nickname, email, status from user")
	helper.Query(sql)
	for helper.Next() {
		user := model.UserDetail{}
		helper.GetValue(&user.ID, &user.Account, &user.Name, &user.Email, &user.Status)
		userList = append(userList, user)
	}

	return userList
}

// QueryUsers 查询指定用户
func QueryUsers(helper dbhelper.DBHelper, ids []int) []model.UserDetail {
	userList := []model.UserDetail{}
	sql := fmt.Sprintf("select id, account, nickname, email, status from user where id in(%s)", util.IntArray2Str(ids))
	helper.Query(sql)
	for helper.Next() {
		user := model.UserDetail{}
		helper.GetValue(&user.ID, &user.Account, &user.Name, &user.Email, &user.Status)
		userList = append(userList, user)
	}

	return userList
}

// QueryUserByAccount 根据账号查询用户信息
func QueryUserByAccount(helper dbhelper.DBHelper, account, password string) (model.UserDetail, bool) {
	user := model.UserDetail{}

	sql := fmt.Sprintf("select id, account, nickname, email, status from user where account='%s' and password='%s'", account, password)
	helper.Query(sql)

	result := false
	if helper.Next() {
		helper.GetValue(&user.ID, &user.Account, &user.Name, &user.Email, &user.Status)
		result = true
	}

	return user, result
}

// QueryUserByID 根据用户ID查询用户信息
func QueryUserByID(helper dbhelper.DBHelper, id int) (model.UserDetail, bool) {
	user := model.UserDetail{}

	sql := fmt.Sprintf("select id, account, nickname, email, status from user where id=%d", id)
	helper.Query(sql)

	result := false
	if helper.Next() {
		helper.GetValue(&user.ID, &user.Account, &user.Name, &user.Email, &user.Status)
		result = true
	}

	return user, result
}

// DeleteUser 删除用户，根据用户ID
func DeleteUser(helper dbhelper.DBHelper, id int) bool {
	sql := fmt.Sprintf("delete from user where id =%d", id)
	_, ret := helper.Execute(sql)
	return ret
}

// DeleteUserByAccount 删除用户，根据用户账号&密码
func DeleteUserByAccount(helper dbhelper.DBHelper, account, password string) bool {
	sql := fmt.Sprintf("delete from user where account ='%s' and password='%s'", account, password)
	_, ret := helper.Execute(sql)
	return ret
}

// CreateUser 创建新用户，根据用户信息和密码
func CreateUser(helper dbhelper.DBHelper, account, email string) (model.UserDetail, bool) {
	user := model.UserDetail{}
	// insert
	sql := fmt.Sprintf("insert into user(account, password, nickname, email, status) values ('%s', '%s', '%s', '%s', %d)", account, "", "", email, 0)
	_, result := helper.Execute(sql)
	if result {
		sql = fmt.Sprintf("select id from user where account='%s' and email='%s'", account, email)
		helper.Query(sql)

		result = false
		if helper.Next() {
			helper.GetValue(&user.ID)

			user.Name = account

			result = true
		}
	}

	return user, result
}

// SaveUser 保存用户信息
func SaveUser(helper dbhelper.DBHelper, user model.UserDetail) (model.UserDetail, bool) {
	// modify
	sql := fmt.Sprintf("update user set nickname='%s', email='%s', status=%d where id =%d", user.Name, user.Email, user.Status, user.ID)
	num, result := helper.Execute(sql)

	return user, result && num == 1
}

// SaveUserWithPassword 保存用户信息
func SaveUserWithPassword(helper dbhelper.DBHelper, user model.UserDetail, password string) (model.UserDetail, bool) {
	// modify
	sql := fmt.Sprintf("update user set password='%s', nickname='%s', email='%s', status=%d where id =%d", password, user.Name, user.Email, user.Status, user.ID)
	num, result := helper.Execute(sql)

	return user, result && num == 1
}