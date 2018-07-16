package handler

import (
	"muidea.com/magicCenter/common/dbhelper"
	"muidea.com/magicCenter/module/modules/content/dal"
	"muidea.com/magicCommon/model"
)

type mediaActionHandler struct {
	dbhelper dbhelper.DBHelper
}

func (i *mediaActionHandler) getAllMedia() []model.Summary {
	return dal.QueryAllMedia(i.dbhelper)
}

func (i *mediaActionHandler) getMedias(ids []int) []model.Media {
	return dal.QueryMedias(i.dbhelper, ids)
}

func (i *mediaActionHandler) findMediaByID(id int) (model.MediaDetail, bool) {
	return dal.QueryMediaByID(i.dbhelper, id)
}

func (i *mediaActionHandler) findMediaByCatalog(catalog int) []model.Summary {
	return dal.QueryMediaByCatalog(i.dbhelper, catalog)
}

func (i *mediaActionHandler) createMedia(name, desc, fileToken, createDate string, catalog []int, expiration, author int) (model.Summary, bool) {
	return dal.CreateMedia(i.dbhelper, name, desc, fileToken, createDate, expiration, author, catalog)
}

func (i *mediaActionHandler) batchCreateMedia(medias []model.MediaItem, createDate string, creater int) ([]model.Summary, bool) {
	return dal.BatchCreateMedia(i.dbhelper, medias, createDate, creater)
}

func (i *mediaActionHandler) saveMedia(media model.MediaDetail) (model.Summary, bool) {
	return dal.SaveMedia(i.dbhelper, media)
}

func (i *mediaActionHandler) destroyMedia(id int) bool {
	return dal.DeleteMediaByID(i.dbhelper, id)
}