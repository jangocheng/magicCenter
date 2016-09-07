package cms

import (
	"html/template"
	"log"
	"magiccenter/common"
	"net/http"
)

// PageView 页面视图
type PageView struct {
	View common.PageView
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	log.Printf("indexHandler")

	res.Header().Set("content-type", "text/html")
	res.Header().Set("charset", "utf-8")

	view := PageView{}

	/*
		url := req.URL.Path
		view.View, _ = bll.QueryPageView(ID, url)
	*/
	t, err := template.ParseFiles("template/html/modules/cms/index.html")
	if err != nil {
		panic("ParseFiles failed, err:" + err.Error())
	}

	t.Execute(res, view)
}

func viewContentHandler(res http.ResponseWriter, req *http.Request) {
	log.Printf("viewContentHandler")

	res.Header().Set("content-type", "text/html")
	res.Header().Set("charset", "utf-8")

	view := PageView{}
	/*
		params := util.SplitParam(req.URL.RawQuery)
		str, found := params["id"]
		if !found {
			panic("illegl param")
		}

		id, err := strconv.Atoi(str)
		if err != nil {
			panic("illegl id, err:" + err.Error())
		}


		url := req.URL.Path
		view.View, _ = bll.QueryContentView(ID, url, id)
	*/
	t, err := template.ParseFiles("template/html/modules/cms/view.html")
	if err != nil {
		panic("ParseFiles failed, err:" + err.Error())
	}

	t.Execute(res, view)
}

func viewCatalogHandler(res http.ResponseWriter, req *http.Request) {
	log.Printf("viewCatalogHandler")

	res.Header().Set("content-type", "text/html")
	res.Header().Set("charset", "utf-8")

	view := PageView{}
	/*
		params := util.SplitParam(req.URL.RawQuery)
		str, found := params["id"]
		if !found {
			panic("illegl param")
		}

		id, err := strconv.Atoi(str)
		if err != nil {
			panic("illegl id, err:" + err.Error())
		}

		url := req.URL.Path
		view.View, _ = bll.QueryCatalogView(ID, url, id)
	*/

	t, err := template.ParseFiles("template/html/modules/cms/catalog.html")
	if err != nil {
		panic("ParseFiles failed, err:" + err.Error())
	}

	t.Execute(res, view)
}

func viewLinkHandler(res http.ResponseWriter, req *http.Request) {
	log.Printf("viewLinkHandler")

	/*
		params := util.SplitParam(req.URL.RawQuery)
		str, found := params["id"]
		if !found {
			panic("illegl param")
		}

		id, err := strconv.Atoi(str)
		if err != nil {
			panic("illegl id, err:" + err.Error())
		}

		link, found := contentbll.QueryLinkById(id)
		if !found {
			http.Redirect(res, req, "/", http.StatusNotFound)
			return
		}

		http.Redirect(res, req, link.Url, http.StatusFound)
	*/
}