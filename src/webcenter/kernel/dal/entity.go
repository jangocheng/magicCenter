package dal

import (
	"fmt"
	"webcenter/util/modelhelper"
)

type Entity struct {
	Id string
	Name string
	Description string
	EnableFlag int
	DefaultFlag int
	Module string
}

func NewEntity(helper modelhelper.Model, id, name, description, module string) (Entity, bool) {
	e := Entity{}
	e.Id = id
	e.Name = name
	e.Description = description
	e.Module = module
	e.EnableFlag = 0
	e.DefaultFlag = 0

	return SaveEntity(helper, e)
}

func DeleteEntity(helper modelhelper.Model, id string) bool {
	sql := fmt.Sprintf("delete from entity where id='%s'", id)
	num, ret := helper.Execute(sql)
	return num == 1 && ret
}

func QueryEntity(helper modelhelper.Model, id string) (Entity, bool) {
	e := Entity{}
	sql := fmt.Sprintf("select id, name, description, enableflag, defaultflag, module from entity where id='%s'", id)	
	helper.Query(sql)
	
	result := false
	if helper.Next() {
		helper.GetValue(&e.Id, &e.Name, &e.Description, &e.EnableFlag, &e.DefaultFlag, &e.Module)
		result = true
	}
	
	return e, result	
}

func QueryEntities(helper modelhelper.Model, module string) []Entity {
	entities := []Entity{}
	
	sql := fmt.Sprintf("select id, name, description, enableflag, defaultlfag, module from entity where module='%s'", module)
	helper.Query(sql)
	
	for helper.Next() {
		e := Entity{}
		helper.GetValue(&e.Id, &e.Name, &e.Description, &e.EnableFlag, &e.DefaultFlag, &e.Module)
		
		entities = append(entities, e)
	}
	
	return entities
}

func SaveEntity(helper modelhelper.Model, e Entity) (Entity, bool) {
	result := false
	_, found := QueryEntity(helper, e.Id)
	if found {
		sql := fmt.Sprintf("update entity set Name ='%s', Description ='%s', enableflag =%d, defaultflag =%d where Id='%s'", e.Name, e.Description, e.EnableFlag, e.DefaultFlag, e.Id)
		_, result = helper.Execute(sql)
	} else {
		sql := fmt.Sprintf("insert into entity(id, name, description, enableflag, defaultflag, module) values ('%s','%s','%s',%d,%d,'%s')", e.Id, e.Name, e.Description, e.EnableFlag, e.DefaultFlag, e.Module)
		_, result = helper.Execute(sql)
	}
	
	return e, result
}
