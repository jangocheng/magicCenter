package dal

import (
	"database/sql"
	"fmt"
	"log"

	"muidea.com/magicCenter/common/dbhelper"
	"muidea.com/magicCenter/common/resource"
	common_const "muidea.com/magicCommon/common"
	"muidea.com/magicCommon/foundation/util"
	"muidea.com/magicCommon/model"
)

func loadCatalogID(helper dbhelper.DBHelper) int {
	var maxID sql.NullInt64
	sql := fmt.Sprintf(`select max(id) from content_catalog`)
	helper.Query(sql)
	defer helper.Finish()

	if helper.Next() {
		helper.GetValue(&maxID)
	}

	return int(maxID.Int64)
}

// QueryAllCatalog 查询所有分类
func QueryAllCatalog(helper dbhelper.DBHelper) []model.Summary {
	summaryList := []model.Summary{}

	ress := resource.QueryResourceByType(helper, model.CATALOG)
	for _, v := range ress {
		summary := model.Summary{Unit: model.Unit{ID: v.RId(), Name: v.RName()}, Description: v.RDescription(), Type: v.RType(), Catalog: []int{}, CreateDate: v.RCreateDate(), Creater: v.ROwner()}

		for _, r := range v.Relative() {
			summary.Catalog = append(summary.Catalog, r.RId())
		}

		// 如果Catalog没有父分类，则认为其父分类为BuildContentCatalog
		if len(summary.Catalog) == 0 {
			summary.Catalog = append(summary.Catalog, common_const.BuildinContentCatalog.ID)
		}

		summaryList = append(summaryList, summary)
	}

	return summaryList
}

// QueryCatalogs 查询指定分类
func QueryCatalogs(helper dbhelper.DBHelper, ids []int) []model.Catalog {
	catalogList := []model.Catalog{}

	if len(ids) == 0 {
		return catalogList
	}

	sql := fmt.Sprintf(`select id, name from content_catalog where id in(%s)`, util.IntArray2Str(ids))
	helper.Query(sql)
	defer helper.Finish()

	for helper.Next() {
		summary := model.Catalog{}
		helper.GetValue(&summary.ID, &summary.Name)

		catalogList = append(catalogList, summary)
	}

	for _, val := range ids {
		if val == common_const.BuildinContentCatalog.ID {
			catalogList = append(catalogList, model.Catalog{ID: common_const.BuildinContentCatalog.ID, Name: common_const.BuildinContentCatalog.Name})
		}
	}

	return catalogList
}

// QueryCatalogByID 查询指定ID的Catalog
func QueryCatalogByID(helper dbhelper.DBHelper, id int) (model.CatalogDetail, bool) {
	if id == common_const.BuildinContentCatalog.ID {
		return common_const.BuildinContentCatalog, true
	}

	catalog := model.CatalogDetail{Summary: model.Summary{Catalog: []int{}}}
	sql := fmt.Sprintf(`select id, name, description, createdate, creater from content_catalog where id = %d`, id)
	helper.Query(sql)

	result := false
	if helper.Next() {
		helper.GetValue(&catalog.ID, &catalog.Name, &catalog.Description, &catalog.CreateDate, &catalog.Creater)
		result = true
	}
	helper.Finish()

	if result {
		ress := resource.QueryRelativeResource(helper, id, model.CATALOG)

		for _, r := range ress {
			catalog.Catalog = append(catalog.Catalog, r.RId())
		}
		if len(catalog.Catalog) == 0 {
			catalog.Catalog = append(catalog.Catalog, common_const.BuildinContentCatalog.ID)
		}
	}
	return catalog, result
}

// QueryCatalogByName 查询指定Cataloga里名字为Name的Catalog
func QueryCatalogByName(helper dbhelper.DBHelper, name string, parentCatalog int) (model.CatalogDetail, bool) {
	var retCatalog model.CatalogDetail
	sql := fmt.Sprintf(`select id, name, description, createdate, creater from content_catalog where name = '%s'`, name)
	helper.Query(sql)

	catalogList := []*model.CatalogDetail{}
	result := false
	for helper.Next() {
		catalog := &model.CatalogDetail{Summary: model.Summary{Catalog: []int{}}}
		helper.GetValue(&catalog.ID, &catalog.Name, &catalog.Description, &catalog.CreateDate, &catalog.Creater)
		catalogList = append(catalogList, catalog)
	}
	helper.Finish()

	for _, val := range catalogList {
		ress := resource.QueryRelativeResource(helper, val.ID, model.CATALOG)

		found := false
		for _, r := range ress {
			if r.RId() == parentCatalog {
				found = true
			}
			val.Catalog = append(val.Catalog, r.RId())
		}
		if len(val.Catalog) == 0 {
			val.Catalog = append(val.Catalog, common_const.BuildinContentCatalog.ID)
		}
		if parentCatalog == common_const.BuildinContentCatalog.ID {
			found = true
		}

		if found {
			retCatalog = *val
			result = true
			break
		}
	}

	return retCatalog, result
}

// QueryCatalogByCatalog 查询指定分类的子类
func QueryCatalogByCatalog(helper dbhelper.DBHelper, id int) []model.Summary {
	summaryList := []model.Summary{}

	resList := resource.QueryReferenceResource(helper, id, model.CATALOG, model.CATALOG)
	for _, r := range resList {
		summary := model.Summary{Unit: model.Unit{ID: r.RId(), Name: r.RName()}, Description: r.RDescription(), Type: r.RType(), Catalog: []int{}, CreateDate: r.RCreateDate(), Creater: r.ROwner()}
		summaryList = append(summaryList, summary)
	}

	for index, value := range summaryList {
		summary := &summaryList[index]
		ress := resource.QueryRelativeResource(helper, value.ID, value.Type)
		for _, r := range ress {
			summary.Catalog = append(summary.Catalog, r.RId())
		}

		if len(summary.Catalog) == 0 {
			summary.Catalog = append(summary.Catalog, common_const.BuildinContentCatalog.ID)
		}
	}

	return summaryList
}

// DeleteCatalog 删除指定类
func DeleteCatalog(helper dbhelper.DBHelper, id int) bool {
	result := false
	helper.BeginTransaction()

	for {
		sql := fmt.Sprintf(`delete from content_catalog where id=%d`, id)

		_, result = helper.Execute(sql)
		if result {
			res, ok := resource.QueryResourceByID(helper, id, model.CATALOG)
			if ok {
				result = resource.DeleteResource(helper, res, true)
			} else {
				result = ok
			}
		}

		break
	}

	if result {
		helper.Commit()
	} else {
		helper.Rollback()
	}

	return result
}

// UpdateCatalog 更新Catalog
func UpdateCatalog(helper dbhelper.DBHelper, catalogs []model.Catalog, parentCatalog int, description, updateDate string, updater int) ([]model.Catalog, bool) {
	ids := []int{}
	result := false
	if len(catalogs) > 0 {
		helper.BeginTransaction()
		for _, val := range catalogs {
			result = true
			detail := model.CatalogDetail{}
			existFlag := false
			if val.ID >= 0 {
				detail, existFlag = QueryCatalogByID(helper, val.ID)
			} else {
				detail, existFlag = QueryCatalogByName(helper, val.Name, parentCatalog)
			}

			if existFlag {
				modifyFlag := false
				if detail.Creater != updater {
					detail.Creater = updater
					modifyFlag = true
				}

				if modifyFlag {
					detail.CreateDate = updateDate
					detail.Description = description

					_, result = SaveCatalog(helper, detail, true)
					if !result {
						log.Printf("UpdateCatalog, saveCatalog failed.")
						break
					}
				}

				ids = append(ids, detail.ID)
			} else {
				detail, ok := CreateCatalog(helper, val.Name, description, updateDate, []int{parentCatalog}, updater, true)
				if ok {
					ids = append(ids, detail.ID)
				} else {
					log.Printf("UpdateCatalog, createCatalog failed.")
					result = false
				}
			}
		}

		if result {
			helper.Commit()
			return QueryCatalogs(helper, ids), true
		}

		helper.Rollback()
		return []model.Catalog{}, false
	}

	return []model.Catalog{}, true
}

// CreateCatalog 新建分类
func CreateCatalog(helper dbhelper.DBHelper, name, description, createDate string, parent []int, creater int, enableTransaction bool) (model.Summary, bool) {
	catalog := model.Summary{Unit: model.Unit{Name: name}, Description: description, Type: model.CATALOG, Catalog: parent, CreateDate: createDate, Creater: creater}

	if !enableTransaction {
		helper.BeginTransaction()
	}

	id := allocCatalogID()
	result := false
	for {
		// insert
		sql := fmt.Sprintf(`insert into content_catalog (id, name, description, createdate, creater) values (%d, '%s','%s','%s',%d)`, id, name, description, createDate, creater)
		_, result = helper.Execute(sql)
		if !result {
			log.Printf("insert catalog to db failed,sql:%s", sql)
			break
		}

		catalog.ID = id
		res := resource.CreateSimpleRes(catalog.ID, model.CATALOG, catalog.Name, catalog.Description, catalog.CreateDate, catalog.Creater)
		for _, c := range parent {
			if c != common_const.BuildinContentCatalog.ID {
				ca, ok := resource.QueryResourceByID(helper, c, model.CATALOG)
				if ok {
					res.AppendRelative(ca)
				} else {
					log.Printf("QueryResourceByID failed,%d, catalog:%s", c, model.CATALOG)
					result = false
					break
				}
			}
		}

		if result {
			result = resource.CreateResource(helper, res, true)
		}

		break
	}

	if !enableTransaction {
		if result {
			helper.Commit()
		} else {
			helper.Rollback()
		}
	}

	return catalog, result
}

// SaveCatalog 保存分类
func SaveCatalog(helper dbhelper.DBHelper, catalog model.CatalogDetail, enableTransaction bool) (model.Summary, bool) {
	if !enableTransaction {
		helper.BeginTransaction()
	}
	summary := model.Summary{Unit: model.Unit{ID: catalog.ID, Name: catalog.Name}, Type: model.CATALOG, Catalog: catalog.Catalog, CreateDate: catalog.CreateDate, Creater: catalog.Creater}

	result := false
	for {
		// modify
		sql := fmt.Sprintf(`update content_catalog set description='%s', createdate='%s', creater =%d where id=%d`, catalog.Description, catalog.CreateDate, catalog.Creater, catalog.ID)
		_, result = helper.Execute(sql)

		if result {
			res, ok := resource.QueryResourceByID(helper, catalog.ID, model.CATALOG)
			if !ok {
				result = false
				break
			}

			res.ResetRelative()
			for _, c := range catalog.Catalog {
				if c != common_const.BuildinContentCatalog.ID {
					ca, ok := resource.QueryResourceByID(helper, c, model.CATALOG)
					if ok {
						res.AppendRelative(ca)
					}
				}
			}
			result = resource.SaveResource(helper, res, true)
		}

		break
	}

	if !enableTransaction {
		if result {
			helper.Commit()
		} else {
			helper.Rollback()
		}
	}

	return summary, result
}