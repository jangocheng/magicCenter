package ui

import (
	"encoding/json"
	"html/template"
	"log"
	"magiccenter/common"
	"magiccenter/kernel/modules/account/bll"
	"magiccenter/kernel/modules/account/model"
	"net/http"
	"strconv"

	"muidea.com/util"
)

// ManageUserView 用户管理视图数据
type ManageUserView struct {
	Users  []model.UserDetail
	Groups []model.Group
}

// AllUserResult 所有用户结果
type AllUserResult struct {
	common.Result
	Users []model.UserDetail
}

// SingleUserResult 单用户结果
type SingleUserResult struct {
	common.Result
	User model.UserDetail
}

// ManageUserViewHandler 用户管理视图处理器
func ManageUserViewHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("ManageUserViewHandler")

	w.Header().Set("content-type", "text/html")
	w.Header().Set("charset", "utf-8")

	t, err := template.ParseFiles("template/html/admin/account/user.html")
	if err != nil {
		panic("parse files failed")
	}

	view := ManageUserView{}
	view.Users = bll.QueryAllUser()
	view.Groups = bll.QueryAllGroup()

	t.Execute(w, view)
}

// QueryAllUserActionHandler 查询所有用户信息处理器
func QueryAllUserActionHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("queryAllUserActionHandler")

	result := AllUserResult{}
	result.Users = bll.QueryAllUser()
	result.ErrCode = 0
	result.Reason = "查询成功"

	b, err := json.Marshal(result)
	if err != nil {
		panic("json.Marshal, failed, err:" + err.Error())
	}

	w.Write(b)
}

// QueryUserActionHandler 查询单个用户信息处理器
func QueryUserActionHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("QueryUserActionHandler")

	result := SingleUserResult{}

	params := util.SplitParam(r.URL.RawQuery)
	for true {
		id, found := params["id"]
		if !found {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		uid, err := strconv.Atoi(id)
		if err != nil {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		result.User, found = bll.QueryUserByID(uid)
		if !found {
			result.ErrCode = 1
			result.Reason = "指定User不存在"
			break
		}

		result.ErrCode = 0
		result.Reason = "查询成功"
		break
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("json.Marshal, failed, err:" + err.Error())
	}

	w.Write(b)
}

// DeleteUserActionHandler 删除用户处理器
func DeleteUserActionHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("DeleteUserActionHandler")

	result := common.Result{}
	params := util.SplitParam(r.URL.RawQuery)
	for true {
		id, found := params["id"]
		if !found {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		uid, err := strconv.Atoi(id)
		if err != nil {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		ok := bll.DeleteUser(uid)
		if !ok {
			result.ErrCode = 1
			result.Reason = "删除用户失败"
			break
		}

		result.ErrCode = 0
		result.Reason = "删除用户成功"
		break
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("json.Marshal, failed, err:" + err.Error())
	}

	w.Write(b)
}

// CheckAccountActionHandler 检查账号是否可用处理器
func CheckAccountActionHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("CheckAccountActionHandler")

	result := common.Result{}

	params := util.SplitParam(r.URL.RawQuery)
	for true {
		account, found := params["account"]
		if !found {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}
		_, found = bll.QueryUserByAccount(account)
		if !found {
			result.ErrCode = 0
			result.Reason = "该账号可用"
			break
		}

		result.ErrCode = 1
		result.Reason = "该账号不可用"
		break
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("json.Marshal, failed, err:" + err.Error())
	}

	w.Write(b)
}

/*
func sendVerifyMail(user, email, id string) {
	systemInfo := configuration.GetSystemInfo()

	subject := "MagicCenter用户验证"

	content := fmt.Sprintf("<html><head><title>用户信息验证</title></head><body><p>Hi %s</p><p><a href='http://%s/user/verify/?id=%s'>请点击链接继续验证用户信息</a></p><p>该邮件由MagicCenter自动发送，请勿回复该邮件</p></body></html>", user, systemInfo.Domain, id)

	mail.PostMail(email, subject, content)
}*/

// SaveUserActionHandler 保存用户信息处理器
func SaveUserActionHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("SaveUserActionHandler")

	result := SingleUserResult{}
	for {
		err := r.ParseMultipartForm(0)
		if err != nil {
			log.Print("paseform failed")

			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		uid := -1
		id := r.FormValue("user-id")
		if len(id) > 0 {
			uid, err = strconv.Atoi(id)
			if err != nil {
				log.Print("paseform failed")

				result.ErrCode = 1
				result.Reason = "无效请求数据"
				break
			}
		}
		account := r.FormValue("user-account")
		nickName := r.FormValue("user-nickname")
		passWord := r.FormValue("user-password")
		email := r.FormValue("user-email")
		groups := r.MultipartForm.Value["user-group"]
		groupList := []int{}
		for _, g := range groups {
			gid, err := strconv.Atoi(g)
			if err != nil {
				log.Printf("parse group id failed, group:%s", g)

				result.ErrCode = 1
				result.Reason = "无效请求数据"
				break
			}

			groupList = append(groupList, gid)
		}

		usr, found := bll.QueryUserByID(uid)
		if found {
			// 说明是更新用户信息
			usr.Account = account
			usr.Name = nickName
			usr.Email = email
			usr.Groups = groupList
			ok := bll.SaveUser(usr)
			if !ok {
				result.ErrCode = 1
				result.Reason = "保存用户信息失败"
				break
			} else {
				result.ErrCode = 0
				result.Reason = "保存用户信息成功"
				result.User = usr
			}
		} else {
			ok := bll.CreateUser(account, passWord, nickName, email, 0, groupList)
			if !ok {
				result.ErrCode = 1
				result.Reason = "创建用户失败"
			} else {
				usr, ok = bll.QueryUserByAccount(account)
				if ok {
					result.ErrCode = 0
					result.Reason = "创建用户成功"

					result.User = usr
				} else {
					result.ErrCode = 1
					result.Reason = "创建用户失败"
				}
			}
		}

		break
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("json.Marshal, failed, err:" + err.Error())
	}

	w.Write(b)
}
