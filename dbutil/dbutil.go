package dbutil

import (
       "fmt"
        "database/sql"

        _ "github.com/go-sql-driver/mysql"

        "github.com/ArchCI/archci/models"
)

// Get the first build which needs to test
func GetBuildToTest() models.Build {

    db, err := sql.Open("mysql", "root:wawa316@/archci")

    if err != nil {
        panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

    // Prepare statement for reading data
    stmtOut, err := db.Prepare("SELECT id, project_name, repo_url, branch FROM build WHERE status = ?")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtOut.Close()

    var id int64
    var projectName string // we "scan" the result in here
    var repoUrl string
    var branch string

    // Query the square-number of 13
    err2 := stmtOut.QueryRow(0).Scan(&id, &projectName, &repoUrl, &branch) // WHERE number = 13

    if err2 != nil {
        panic(err2.Error()) // proper error handling instead of panic in your app
    }

    fmt.Printf("The project name is: %d", projectName)

    build := models.Build{Id:id, ProjectName: projectName, RepoUrl: repoUrl, Branch: branch}
    return build
}
