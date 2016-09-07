package ui

import (
	"encoding/json"
	"html"
	"html/template"
	"log"
	"magiccenter/common"
	"magiccenter/configuration"
	accountModel "magiccenter/kernel/modules/account/model"
	"magiccenter/kernel/modules/content/bll"
	"magiccenter/kernel/modules/content/model"
	"magiccenter/session"
	"net/http"
	"strconv"

	"muidea.com/util"
)

// ManageArticleView 文章管理视图
type ManageArticleView struct {
	Articles []model.ArticleSummary
	Catalogs []model.CatalogDetail
}

type QueryAllArticleResult struct {
	Articles []model.ArticleSummary
}

type QueryArticleResult struct {
	common.Result
	Article model.Article
}

type DeleteArticleResult struct {
	common.Result
}

type AjaxArticleResult struct {
	common.Result
}

type EditArticleResult struct {
	common.Result
	Article model.Article
}

// ManageArticleViewHandler 文章管理主界面
// 显示Article列表信息
// 返回html页面
//
func ManageArticleViewHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("ManageArticleViewHandler")

	w.Header().Set("content-type", "text/html")
	w.Header().Set("charset", "utf-8")

	t, err := template.ParseFiles("template/html/admin/content/article.html")
	if err != nil {
		panic("parse files failed, err:" + err.Error())
	}

	view := ManageArticleView{}
	view.Articles = bll.QueryAllArticleSummary()
	view.Catalogs = bll.QueryAllCatalogDetail()

	t.Execute(w, view)
}

//
// 查询Article
// 返回json
//
func QueryAllArticleHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("QueryAllArticleHandler")

	result := QueryAllArticleResult{}
	result.Articles = bll.QueryAllArticleSummary()

	b, err := json.Marshal(result)
	if err != nil {
		panic("json.Marshal, failed, err:" + err.Error())
	}

	w.Write(b)
}

//
// 查询指定Article内容
// 返回json
//
func QueryArticleHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("QueryArticleHandler")

	result := QueryArticleResult{}

	for true {
		err := r.ParseForm()
		if err != nil {
			log.Print("paseform failed")

			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		params := util.SplitParam(r.URL.RawQuery)
		id, found := params["id"]
		if !found {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		aid, err := strconv.Atoi(id)
		if err != nil {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		article, found := bll.QueryArticleByID(aid)
		if !found {
			result.ErrCode = 1
			result.Reason = "操作失败"
			break
		}

		article.Content = html.UnescapeString(article.Content)
		result.Article = article
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

//
// 删除指定Article
// 返回json
//
func DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("DeleteArticleHandler")

	result := DeleteArticleResult{}

	for true {
		err := r.ParseForm()
		if err != nil {
			log.Print("paseform failed")

			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		params := util.SplitParam(r.URL.RawQuery)
		id, found := params["id"]
		if !found {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		aid, err := strconv.Atoi(id)
		if err != nil {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		if !bll.DeleteArticle(aid) {
			result.ErrCode = 1
			result.Reason = "操作失败"
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

//
// 保存Article
// 返回json
//
func AjaxArticleHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("AjaxArticleHandler")

	authID, found := configuration.GetOption(configuration.AuthorithID)
	if !found {
		panic("unexpected, can't fetch authorith id")
	}

	session := session.GetSession(w, r)
	user, found := session.GetOption(authID)
	if !found {
		panic("unexpected, must login system first.")
	}

	result := AjaxArticleResult{}
	for true {
		err := r.ParseMultipartForm(0)
		if err != nil {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		id := r.FormValue("article-id")
		title := r.FormValue("article-title")
		content := html.EscapeString(r.FormValue("article-content"))
		catalog := r.MultipartForm.Value["article-catalog"]

		aid, err := strconv.Atoi(id)
		if err != nil {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		catalogs := []int{}
		for _, ca := range catalog {
			cid, err := strconv.Atoi(ca)
			if err != nil {
				result.ErrCode = 1
				result.Reason = "无效请求数据"
				break
			}

			catalogs = append(catalogs, cid)
		}

		if !bll.SaveArticle(aid, title, content, user.(accountModel.UserDetail).ID, catalogs) {
			result.ErrCode = 1
			result.Reason = "操作失败"
			break
		}

		result.ErrCode = 0
		result.Reason = "操作成功"
		break
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("json.Marshal, failed, err:" + err.Error())
	}

	w.Write(b)
}

//
// 编辑Article
// 返回article内容和当前可用的catalog
//
func EditArticleHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("EditArticleHandler")

	result := EditArticleResult{}

	for true {
		params := util.SplitParam(r.URL.RawQuery)
		id, found := params["id"]
		if !found {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		aid, err := strconv.Atoi(id)
		if err != nil {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		article, found := bll.QueryArticleByID(aid)
		if !found {
			result.ErrCode = 1
			result.Reason = "操作失败"
			break
		}

		article.Content = html.UnescapeString(article.Content)
		result.Article = article
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