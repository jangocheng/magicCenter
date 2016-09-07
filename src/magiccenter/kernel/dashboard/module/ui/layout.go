package ui

/*
type ApplyModuleSettingResult struct {
	common.Result
	Modules       []model.Module
	DefaultModule string
}

type QueryModuleDetailResult struct {
	common.Result
	Module model.ModuleLayout
}

type DeleteModuleBlockResult struct {
	QueryModuleDetailResult
}

type SaveModuleBlockResult struct {
	QueryModuleDetailResult
}

type SavePageBlockResult struct {
	QueryModuleDetailResult
}

func ApplyModuleSettingHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("ApplyModuleSettingHandler")

	result := ApplyModuleSettingResult{}

	for true {
		err := r.ParseForm()
		if err != nil {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		enableList := strings.Split(r.FormValue("module-enableList"), ",")
		defaultModule := r.FormValue("module-defaultModule")

		ret := false
		result.Modules, ret = bll.EnableModules(enableList)
		if ret {
			ret = configuration.SetOption(configuration.SYS_DEFULTMODULE, defaultModule)
			if ret {
				result.DefaultModule = defaultModule

				result.ErrCode = 0
				result.Reason = "保存设置成功"
			} else {
				result.ErrCode = 1
				result.Reason = "保存设置部分成功"
			}

			break
		}

		result.ErrCode = 1
		result.Reason = "保存设置失败"
		break
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("json.Marshal, failed, err:" + err.Error())
	}

	w.Write(b)
}

func QueryModuleDetailHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("QueryModuleDetailHandler")

	result := QueryModuleDetailResult{}

	params := util.SplitParam(r.URL.RawQuery)
	for true {
		id, found := params["id"]
		if !found {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		result.Module, found = bll.QueryModuleDetail(id)
		if !found {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		result.ErrCode = 0
		result.Reason = "查询成功"
		break
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("marshal failed, err:" + err.Error())
	}

	w.Write(b)
}

func DeleteModuleBlockHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("DeleteModuleBlockHandler")

	result := DeleteModuleBlockResult{}

	for true {
		rawParams := util.SplitParam(r.URL.RawQuery)

		id, found := rawParams["id"]
		if !found {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}
		owner, found := rawParams["owner"]
		if !found {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		idValue, err := strconv.Atoi(id)
		if err != nil {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		result.Module, found = bll.RemoveModuleBlock(idValue, owner)
		if !found {
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
		panic("marshal failed, err:" + err.Error())
	}

	w.Write(b)
}

func SaveModuleBlockHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("SaveModuleBlockHandler")

	result := SaveModuleBlockResult{}
	for true {
		err := r.ParseForm()
		if err != nil {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		owner := r.FormValue("module-id")
		style, err := strconv.Atoi(r.FormValue("block-style"))
		if err != nil {
			result.ErrCode = 1
			result.Reason = "无效请求数据"
			break
		}

		block := r.FormValue("block-name")
		tag := r.FormValue("block-tag")
		ret := false
		result.Module, ret = bll.AddModuleBlock(block, tag, style, owner)
		if !ret {
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
		panic("marshal failed, err:" + err.Error())
	}

	w.Write(b)
}
*/
