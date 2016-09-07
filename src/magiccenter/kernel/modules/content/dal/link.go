package dal

import (
	"fmt"
	"magiccenter/kernel/modules/content/model"
	resdal "magiccenter/resource/dal"
	"magiccenter/util/dbhelper"
)

// QueryAllLink 查询全部Link
func QueryAllLink(helper dbhelper.DBHelper) []model.Link {
	linkList := []model.Link{}
	sql := fmt.Sprintf(`select id, name, url, logo, creater from link`)
	helper.Query(sql)

	for helper.Next() {
		link := model.Link{}
		helper.GetValue(&link.ID, &link.Name, &link.URL, &link.Logo, &link.Creater)

		linkList = append(linkList, link)
	}

	for _, link := range linkList {
		ress := resdal.QueryRelativeResource(helper, link.ID, model.LINK)
		for _, r := range ress {
			link.Catalog = append(link.Catalog, r.RId())
		}
	}

	return linkList
}

// QueryLinkByCatalog 查询指定分类下的Link
func QueryLinkByCatalog(helper dbhelper.DBHelper, id int) []model.Link {
	linkList := []model.Link{}

	resList := resdal.QueryReferenceResource(helper, id, model.CATALOG, model.LINK)
	for _, r := range resList {
		sql := fmt.Sprintf(`select id, name, url, logo, creater from link where id =%d`, r.RId())
		helper.Query(sql)

		if helper.Next() {
			link := model.Link{}
			helper.GetValue(&link.ID, &link.Name, &link.URL, &link.Logo, &link.Creater)
			linkList = append(linkList, link)
		}
	}

	for _, link := range linkList {
		ress := resdal.QueryRelativeResource(helper, link.ID, model.LINK)
		for _, r := range ress {
			link.Catalog = append(link.Catalog, r.RId())
		}
	}

	return linkList
}

// QueryLinkByRang 查询指定范围的Link
func QueryLinkByRang(helper dbhelper.DBHelper, begin int, offset int) []model.Link {
	linkList := []model.Link{}
	sql := fmt.Sprintf(`select id, name, url, logo, creater from link order by id where id >= %d limit %d`, begin, offset)
	helper.Query(sql)

	for helper.Next() {
		link := model.Link{}
		helper.GetValue(&link.ID, &link.Name, &link.URL, &link.Logo, &link.Creater)

		linkList = append(linkList, link)
	}

	for _, link := range linkList {
		ress := resdal.QueryRelativeResource(helper, link.ID, model.LINK)
		for _, r := range ress {
			link.Catalog = append(link.Catalog, r.RId())
		}
	}

	return linkList
}

// QueryLinkByID 查询指定Link
func QueryLinkByID(helper dbhelper.DBHelper, id int) (model.Link, bool) {
	link := model.Link{}
	sql := fmt.Sprintf(`select id, name, url, logo, creater from link where id =%d`, id)
	helper.Query(sql)

	result := false
	if helper.Next() {
		helper.GetValue(&link.ID, &link.Name, &link.URL, &link.Logo, &link.Creater)
		result = true
	}

	if result {
		ress := resdal.QueryRelativeResource(helper, link.ID, model.LINK)
		for _, r := range ress {
			link.Catalog = append(link.Catalog, r.RId())
		}
	}

	return link, result
}

// DeleteLinkByID 删除指定Link
func DeleteLinkByID(helper dbhelper.DBHelper, id int) bool {
	sql := fmt.Sprintf(`delete from link where id =%d`, id)
	num, result := helper.Execute(sql)
	if num > 0 && result {
		lnk := resdal.CreateSimpleRes(id, model.LINK, "")
		result = resdal.DeleteResource(helper, lnk)
	}

	return result
}

// SaveLink 保存Link
func SaveLink(helper dbhelper.DBHelper, link model.Link) bool {
	sql := fmt.Sprintf(`select id from link where id=%d`, link.ID)
	helper.Query(sql)

	result := false
	if helper.Next() {
		var id = 0
		helper.GetValue(&id)
		result = true
	}

	if !result {
		// insert
		sql = fmt.Sprintf(`insert into link (name,url,logo,creater) values ('%s','%s','%s', %d)`, link.Name, link.URL, link.Logo, link.Creater)
		_, result = helper.Execute(sql)
		if result {
			sql = fmt.Sprintf(`select id from link where name='%s' and url ='%s' and creater=%d`, link.Name, link.URL, link.Creater)

			helper.Query(sql)
			if helper.Next() {
				helper.GetValue(&link.ID)
			}
		}
	} else {
		// modify
		sql = fmt.Sprintf(`update link set name ='%s', url ='%s', logo='%s', creater=%d where id=%d`, link.Name, link.URL, link.Logo, link.Creater, link.ID)
		_, result = helper.Execute(sql)
	}

	if result {
		res := resdal.CreateSimpleRes(link.ID, model.LINK, link.Name)
		for _, c := range link.Catalog {
			ca := resdal.CreateSimpleRes(c, model.CATALOG, "")
			res.AppendRelative(ca)
		}
		result = resdal.SaveResource(helper, res)
	}

	return result
}