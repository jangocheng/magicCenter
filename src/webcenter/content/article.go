package content

import (
	"fmt"
	"log"
	"html"
	"muidea.com/dao"
	"webcenter/auth"
)

type ArticleInfo struct {
	Id int
	Title string
	CreateDate string
	Catalog Catalog
	Author auth.User
}


type Article struct {
	Id int
	Title string
	Content string
	CreateDate string
	Catalog Catalog
	Author auth.User
}

func newArticleInfo() ArticleInfo {
	articleInfo := ArticleInfo{}
	articleInfo.Id = -1
	articleInfo.Catalog = newCatalog()
	articleInfo.Author = auth.NewUser()
	
	return articleInfo
}

func newArticle() Article {
	article := Article{}
	article.Id = -1
	article.Catalog = newCatalog()
	article.Author = auth.NewUser()
	
	return article
}

func GetAllArticleInfo(dao * dao.Dao) []ArticleInfo {
	articleInfoList := []ArticleInfo{}
	sql := fmt.Sprintf("select id, title, author, createdate, catalog from article")
	if !dao.Query(sql) {
		log.Printf("query article failed, sql:%s", sql)
		return articleInfoList
	}

	for dao.Next() {
		articleInfo := newArticleInfo()
		dao.GetField(&articleInfo.Id, &articleInfo.Title, &articleInfo.Author.Id, &articleInfo.CreateDate, &articleInfo.Catalog.Id)
		
		articleInfoList = append(articleInfoList, articleInfo)
	}
	
	for i:=0; i < len(articleInfoList); i++ {
		articleInfo := &articleInfoList[i]
		if !articleInfo.Author.Query(dao) {
			articleInfo.Author, _ = auth.QueryDefaultUser(dao)
		}
		articleInfo.Catalog.Query(dao)
	}
	
	return articleInfoList
}

func GetArticleByCatalog(id int, dao* dao.Dao) []ArticleInfo {
	articleInfoList := []ArticleInfo{}
	sql := fmt.Sprintf("select id, title, author, createdate, catalog from article where catalog=%d", id)
	if !dao.Query(sql) {
		log.Printf("query article failed, sql:%s", sql)
		return articleInfoList
	}

	for dao.Next() {
		articleInfo := newArticleInfo()
		dao.GetField(&articleInfo.Id, &articleInfo.Title, &articleInfo.Author.Id, &articleInfo.CreateDate, &articleInfo.Catalog.Id)
		
		articleInfoList = append(articleInfoList, articleInfo)
	}
	
	for i:=0; i < len(articleInfoList); i++ {
		articleInfo := &articleInfoList[i]
		if !articleInfo.Author.Query(dao) {
			articleInfo.Author, _ = auth.QueryDefaultUser(dao)
		}
		articleInfo.Catalog.Query(dao)
	}
	
	return articleInfoList	
}

func QueryArticleByRang(begin int,offset int, dao* dao.Dao) []Article {
	articleList := []Article{}
	sql := fmt.Sprintf("select * from (select id, title, content, author, createdate, catalog from article order by id) c where id > %d limit %d", begin, offset)
	if !dao.Query(sql) {
		log.Printf("query article failed, sql:%s", sql)
		return articleList
	}

	for dao.Next() {
		article := newArticle()
		dao.GetField(&article.Id, &article.Title, &article.Content, &article.Author.Id, &article.CreateDate, &article.Catalog.Id)
		
		articleList = append(articleList, article)
	}
	
	for i:=0; i < len(articleList); i++ {
		article := &articleList[i]
		if !article.Author.Query(dao) {
			article.Author, _ = auth.QueryDefaultUser(dao)
		}
		article.Catalog.Query(dao)
		
		article.Content = html.UnescapeString(article.Content)
	}
	
	return articleList	
}


func QueryArticleById(id int, dao* dao.Dao) (Article, bool) {
	article := Article{}
	sql := fmt.Sprintf("select id, title, content, author, createdate, catalog from article where id = %d", id)
	if !dao.Query(sql) {
		log.Printf("query article failed, sql:%s", sql)
		return article, false
	}

	result := false
	for dao.Next() {
		dao.GetField(&article.Id, &article.Title, &article.Content, &article.Author.Id, &article.CreateDate, &article.Catalog.Id)
		result = true
	}
	
	if !article.Author.Query(dao) {
		article.Author, _ = auth.QueryDefaultUser(dao)
	}
	article.Catalog.Query(dao)
	
	article.Content = html.UnescapeString(article.Content)
	
	return article, result	
}


func (this *Article)Query(dao * dao.Dao) bool {
	sql := fmt.Sprintf("select id, title, content, author, createdate, catalog from article where id=%d", this.Id)
	if !dao.Query(sql) {
		log.Printf("query article failed, sql:%s", sql)
		return false
	}

	result := false;
	for dao.Next() {
		result = dao.GetField(&this.Id, &this.Title, &this.Content, &this.Author.Id, &this.CreateDate, &this.Catalog.Id)
	}

	if result {
		result = this.Author.Query(dao)
		if !result {
			this.Author, result = auth.QueryDefaultUser(dao)
		}
		this.Content = html.UnescapeString(this.Content)		
	}
	
	if result {
		result = this.Catalog.Query(dao)
	}
	
	return result		
}



func (this *Article)delete(dao * dao.Dao) bool {
	sql := fmt.Sprintf("delete from article where id=%d", this.Id)
	
	result := dao.Execute(sql)
	
	return result		
}

func (this *Article)save(dao * dao.Dao) bool {
	sql := fmt.Sprintf("select id from article where id=%d", this.Id)
	if !dao.Query(sql) {
		log.Printf("query article failed, sql:%s", sql)
		return false
	}

	result := false;
	for dao.Next() {
		var id = 0
		result = dao.GetField(&id)
		result = true
	}

	if !result {
		// insert
		sql = fmt.Sprintf("insert into article (title,content,author,createdate,catalog) values ('%s','%s',%d,'%s',%d)", this.Title, html.EscapeString(this.Content), this.Author.Id, this.CreateDate, this.Catalog.Id)
	} else {
		// modify
		sql = fmt.Sprintf("update article set title ='%s', content ='%s', author =%d, createdate ='%s', catalog =%d where id=%d", this.Title, html.EscapeString(this.Content), this.Author.Id, this.CreateDate, this.Catalog.Id, this.Id)
	}
	
	result = dao.Execute(sql)
	
	return result
}


