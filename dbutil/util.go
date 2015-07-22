package dbutil

import (
	"fmt"
	//"database/sql"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	"github.com/ArchCI/archci/models"
)

func InitializeModels() {
	orm.RegisterDataBase("default", "mysql", "root:wawa316@/archci?charset=utf8")
	orm.RegisterModel(new(models.Build), new(models.Project), new(models.Worker))
}

// Get one build whose status is NOT_START
func GetOneNotStartBuild() (models.Build, error) {
	o := orm.NewOrm()

	build := models.Build{Status: models.BUILD_STATUS_NOT_START}
	err := o.Read(&build, "Status")
	if err == orm.ErrNoRows {
		fmt.Println("No build whose status is 0")
		return build, err
	} else if err == orm.ErrMissPK {
		fmt.Println("No this primary key")
	} else {
		fmt.Println(build.Id)
	}

	return build, nil
}

func UpdateBuildStatus(buildId int64, status int) {
	fmt.Println("Start to update status")
	o := orm.NewOrm()
	build := models.Build{Id: buildId}
	err2 := o.Read(&build)
	if err2 != nil {
		fmt.Println(err2)
	} else {
		fmt.Println("Get this build")
		build.Status = status
		if num, err := o.Update(&build); err == nil {
			fmt.Println(num)
		}
	}
}
