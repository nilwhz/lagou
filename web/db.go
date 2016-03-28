package main

import (
	"database/sql"
	"fmt"
	"github.com/go-ini/ini"
	_ "github.com/go-sql-driver/mysql"
	"lagou/util"
)

var db *sql.DB

func openDB() {
	cfg, err := ini.Load("./config.ini")
	util.CheckErr(err)
	dbSection := cfg.Section("mysql")

	// db, _ = sql.Open("mysql", "user:password@tcp(host:port)/databasename")
	db, _ = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbSection.Key("user"), dbSection.Key("pwd"), dbSection.Key("host"), dbSection.Key("port"), dbSection.Key("db")))

	err = db.Ping()
	util.CheckErr(err)
	fmt.Println("connect db success.")
}

func closeDB() {
	db.Close()
}

func findAllJobs() []jobItem {
	var jobList []jobItem

	rows, err := db.Query("select * from job")
	util.CheckErr(err)
	defer rows.Close()

	for rows.Next() {
		job := jobItem{}
		err = rows.Scan(&job.Id, &job.CompanyName, &job.CompanyArea, &job.CompanySize, &job.CompanyUrl, &job.CompanyStage, &job.CompanyAddress, &job.JobName, &job.JobSalary, &job.JobUpTime, &job.JobExp, &job.JobDegree, &job.JobDesc, &job.HrName, &job.HrPercent, &job.HrExecTime, &job.LagouURL)
		util.CheckErr(err)

		jobList = append(jobList, job)
	}
	err = rows.Err()
	util.CheckErr(err)

	return jobList
}
