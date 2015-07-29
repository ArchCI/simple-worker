package dbutil

import (
	"fmt"
	"os"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	"github.com/ArchCI/archci/models"
)

const (
	ENV_MYSQL_SERVER   = "MYSQL_SERVER"
	ENV_MYSQL_USERNAME = "MYSQL_USERNAME"
	ENV_MYSQL_PASSWORD = "MYSQL_PASSWORD"
	ENV_MYSQL_DATABASE = "MYSQL_DATABASE"

	MYSQL_DRIVER = "mysql"
)

// InitializeModels registries the models of archci.
func InitializeModels() {
	// Registry database models.
	orm.RegisterModel(new(models.Build), new(models.Project), new(models.Worker))

	// Initialize database with environment variables.
	server := ""
	username := "root"
	password := "root"
	database := "mysql"

	if os.Getenv(ENV_MYSQL_SERVER) != "" {
		server = os.Getenv(ENV_MYSQL_SERVER)
	}
	if os.Getenv(ENV_MYSQL_USERNAME) != "" {
		username = os.Getenv(ENV_MYSQL_USERNAME)
	}
	if os.Getenv(ENV_MYSQL_PASSWORD) != "" {
		password = os.Getenv(ENV_MYSQL_PASSWORD)
	}
	if os.Getenv(ENV_MYSQL_DATABASE) != "" {
		database = os.Getenv(ENV_MYSQL_DATABASE)
	}

	// The datasource looks like "root:root@/archci?charset=utf8".
	DATASOURCE := username + ":" + password + "@" + server + "/" + database + "?charset=utf8"
	fmt.Println("Connect to database with " + DATASOURCE)

	orm.RegisterDriver(MYSQL_DRIVER, orm.DR_MySQL)
	orm.RegisterDataBase("default", MYSQL_DRIVER, DATASOURCE, 30)
	orm.RunSyncdb("default", false, true)
}

// GetOneNotStartBuild takes one build whose status is NOT_START.
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

// UpdateBuildStatus takes the status to update the build.
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

// AddWorker inserts worker record in database.
func AddWorker(workerId int64, ip string, lastUpdate time.Time, status int) {
	o := orm.NewOrm()

	var worker models.Worker
	worker.Id = workerId
	worker.Ip = ip
	worker.LastUpdate = lastUpdate
	worker.Status = status

	id, err := o.Insert(&worker)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(id)
	}
}
