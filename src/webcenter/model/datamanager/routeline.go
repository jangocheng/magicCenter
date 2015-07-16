package datamanager

import (

)

type Routeline struct {
	Id int
	Name string
	Description string
	Creater int
	CreateDate string
}

func (this *Routeline)Valid() bool {
	return this.Id > 0 && len(this.Name) > 0
}
