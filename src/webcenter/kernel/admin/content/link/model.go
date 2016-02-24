package link


import (
	"fmt"
	"webcenter/util/modelhelper"
	"webcenter/kernel/admin/common"
	"webcenter/kernel/admin/content/base"
)

type LinkInfo struct {
	Id int
	Name string
	Url string
	Logo string
	Style int
	Creater string
	Catalog []string
}

type Link interface {
	common.Resource
	Url() string
	Logo() string
	Style() int
	Creater() int
	SetId(id int)
	SetName(name string)
	SetUrl(url string)
	SetLogo(logoUrl string)
	SetStyle(style int)
	SetCreater(user int)
	SetCatalog(catalog []int)
}

type link struct {
	id int
	name string
	url string
	logo string
	style int
	creater int
	catalog []int
}


func (this *link) Id() int {
	return this.id
}

func (this *link) Name() string {
	return this.name
}

func (this *link) Type() int {
	return base.LINK
}

func (this *link) Relative() []common.Resource {
	ress := []common.Resource{}
	
	for _, catalog := range this.catalog {
		res := common.NewSimpleRes(catalog,"", base.CATALOG)
		ress = append(ress, res)
	}
	
	return ress
}

func (this *link) Url() string {
	return this.url
}

func (this *link) Logo() string {
	return this.logo
}

func (this *link) Style() int {
	return this.style
}

func (this *link) Creater() int {
	return this.creater
}

func (this *link) SetId(id int) {
	this.id = id
}

func (this *link) SetName(name string) {
	this.name = name
}

func (this *link) SetUrl(url string) {
	this.url = url
}

func (this *link) SetLogo(logo string) {
	this.logo = logo
}

func (this *link) SetStyle(style int) {
	this.style = style
}

func (this *link) SetCreater(user int) {
	this.creater = user
}

func (this *link) SetCatalog(catalog []int) {
	this.catalog = catalog
}

func NewLink() Link {
	link := &link{}

	return link
}

func QueryAllLink(model modelhelper.Model) []LinkInfo {
	linkInfoList := []LinkInfo{}
	sql := fmt.Sprintf(`select l.id,l.name, l.url,l.logo, l.style, u.nickname from link l, user u where l.creater = u.id`)
	model.Query(sql)

	for model.Next() {
		link := LinkInfo{}
		model.GetValue(&link.Id, &link.Name, &link.Url, &link.Logo, &link.Style, &link.Creater)
		
		linkInfoList = append(linkInfoList, link)
	}
	
	for index, info := range linkInfoList {
		sql = fmt.Sprintf(`select r.name from resource r, resource_relative rr where r.id = rr.dst and r.type = rr.dstType and rr.src = %d and rr.srcType=%d`, info.Id, base.LINK)
		name := "-"
		model.Query(sql)
		for model.Next() {
			model.GetValue(&name)
			linkInfoList[index].Catalog = append(linkInfoList[index].Catalog, name)
		}
	}
		
	return linkInfoList
}


func QueryLinkByCatalog(model modelhelper.Model, id int) []LinkInfo {
	linkInfoList := []LinkInfo{}
	sql := fmt.Sprintf(`select l.id,l.name, l.url,l.logo, l.style, u.nickname from link l, user u where l.creater = u.id and l.id in (
		select src from resource_relative where dst = %d and dstType = %d and srcType = %d )`, id, base.CATALOG, base.LINK)
	model.Query(sql)

	for model.Next() {
		link := LinkInfo{}
		model.GetValue(&link.Id, &link.Name, &link.Url, &link.Logo, &link.Style, &link.Creater)
		
		linkInfoList = append(linkInfoList, link)
	}
	
	for index, info := range linkInfoList {
		sql = fmt.Sprintf(`select r.name from resource r, resource_relative rr where r.id = rr.dst and r.type = rr.dstType and rr.src = %d and rr.srcType=%d`, info.Id, base.LINK)
		name := "-"
		model.Query(sql)
		for model.Next() {
			model.GetValue(&name)
			linkInfoList[index].Catalog = append(linkInfoList[index].Catalog, name)
		}
	}
		
	return linkInfoList
}

func QueryLinkByRang(model modelhelper.Model, begin int,offset int) []LinkInfo {
	linkInfoList := []LinkInfo{}
	sql := fmt.Sprintf(`select l.id,l.name, l.url,l.logo, l.style, u.nickname from link l, user u where l.creater = u.id and l.id in (
		select src from resource_relative where dstType = %d and srcType = %d ) and l.id >= %d limit %d`, base.CATALOG, base.LINK, begin, offset)
	model.Query(sql)

	for model.Next() {
		link := LinkInfo{}
		model.GetValue(&link.Id, &link.Name, &link.Url, &link.Logo, &link.Style, &link.Creater)
		
		linkInfoList = append(linkInfoList, link)
	}
	
	for index, info := range linkInfoList {
		sql = fmt.Sprintf(`select r.name from resource r, resource_relative rr where r.id = rr.dst and r.type = rr.dstType and rr.src = %d and rr.srcType=%d`, info.Id, base.LINK)
		name := "-"
		model.Query(sql)
		for model.Next() {
			model.GetValue(&name)
			linkInfoList[index].Catalog = append(linkInfoList[index].Catalog, name)
		}
	}
		
	return linkInfoList
}

func QueryLinkById(model modelhelper.Model, id int) (Link, bool) {
	link := &link{}
	sql := fmt.Sprintf(`select id, name, url, logo,style, creater from link where id =%d`, id)
	model.Query(sql)

	result := false
	if model.Next() {
		model.GetValue(&link.id, &link.name, &link.url, &link.logo, &link.style, &link.creater)
		result = true
	}
	if !result {
		return link, result
	}
	
	sql = fmt.Sprintf(`select dst from resource_relative where src = %d and srcType = %d and dstType =%d`, link.id, base.LINK, base.CATALOG)
	pid := -1
	model.Query(sql)
	for model.Next() {
		model.GetValue(&pid)
		link.catalog = append(link.catalog, pid)
	}
	
	return link, result
}

func DeleteLink(model modelhelper.Model, id int) bool {
	model.BeginTransaction()
	
	sql := fmt.Sprintf(`delete from link where id =%d`, id)	
	result := model.Execute(sql)
	if result {
		link := link{}
		link.id = id
		result  = common.DeleteResource(model, &link)
	}
		
	if !result {
		model.Rollback()
	} else {
		model.Commit()
	}
		
	return result	
}

func SaveLink(model modelhelper.Model, link Link) bool {
	sql := fmt.Sprintf(`select id from link where id=%d`, link.Id())
	model.Query(sql)

	result := false;
	if model.Next() {
		var id = 0
		model.GetValue(&id)
		result = true
	}

	model.BeginTransaction()
	
	if !result {
		// insert
		sql = fmt.Sprintf(`insert into link (name,url,logo,style, creater) values ('%s','%s','%s',%d, %d)`, link.Name(), link.Url(), link.Logo(), link.Style(), link.Creater())
		result = model.Execute(sql)
		if result {
			sql = fmt.Sprintf(`select id from link where name='%s' and url ='%s' and creater=%d`, link.Name(), link.Url(), link.Creater())
			
			id := -1
			model.Query(sql)
			for model.Next() {
				model.GetValue(&id)
				link.SetId(id)
				break
			}
		}
	} else {
		// modify
		sql = fmt.Sprintf(`update link set name ='%s', url ='%s', logo='%s', style=%d, creater=%d where id=%d`, link.Name(), link.Url(), link.Logo(), link.Style(), link.Creater(), link.Id())
		result = model.Execute(sql)
	}
	
	if result {
		result = common.SaveResource(model, link)
	}
	
	if result {
		model.BeginTransaction()
	} else {
		model.Rollback()
	}
	
	return result	
}
