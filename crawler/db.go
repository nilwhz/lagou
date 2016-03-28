package main

import (
	"database/sql"
	"fmt"
	"github.com/go-ini/ini"
	_ "github.com/go-sql-driver/mysql"
	"lagou/util"
)

// create table job(
//     id int(10) AUTO_INCREMENT not null,
//     companyName varchar(64),
//     companyArea    varchar(64),
//     companySize    varchar(64),
//     companyUrl     varchar(64),
//     companyStage   varchar(64),
//     companyAddress varchar(64),
//     jobName   varchar(64),
//     jobSalary varchar(64),
//     jobUpTime varchar(64),
//     jobExp    varchar(64),
//     jobDegree varchar(64),
//     jobDesc   text,
//     hrName varchar(64),
//     hrPercent varchar(64),
//     hrExecTime varchar(64),
//     primary key(id),
//     lagouURL varchar(64)
// ) charset=utf8;

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

func clearHistoryData() {
	// 每次执行爬虫脚本，都清空历史数据。
	res, err := db.Exec("delete from job")
	util.CheckErr(err)
	count, err := res.RowsAffected()
	util.CheckErr(err)
	fmt.Printf("清空历史数据%d条\n", count)
	fmt.Println()
}

func insertOneJob(id int, job jobItem) {
	stmt, err := db.Prepare("insert job set id=?, lagouURL=?, companyName=?, companyArea=?, companySize=?, companyUrl=?, companyStage=?, companyAddress=?, jobName=?, jobSalary=?, jobUptime=?, jobExp=?, jobDegree=?, jobDesc=?, hrName=?, hrPercent=?, hrExecTime=?")
	util.CheckErr(err)

	_, err = stmt.Exec(id, job.lagouURL, job.companyName, job.companyArea, job.companySize, job.companyUrl, job.companyStage, job.companyAddress, job.jobName, job.jobSalary, job.jobUpTime, job.jobExp, job.jobDegree, job.jobDesc, job.hrName, job.hrPercent, job.hrExecTime)
	util.CheckErr(err)
}

func updateOneJob(job jobItem) {
	// 执行更新操作
	stmt, err := db.Prepare("update job set lagouURL=?, companyName=?, companyArea=?, companySize=?, companyUrl=?, companyStage=?, companyAddress=?, jobName=?, jobSalary=?, jobUptime=?, jobExp=?, jobDegree=?, jobDesc=?, hrName=?, hrPercent=?, hrExecTime=? where companyName=?")
	util.CheckErr(err)

	_, err = stmt.Exec(job.lagouURL, job.companyName, job.companyArea, job.companySize, job.companyUrl, job.companyStage, job.companyAddress, job.jobName, job.jobSalary, job.jobUpTime, job.jobExp, job.jobDegree, job.jobDesc, job.hrName, job.hrPercent, job.hrExecTime, job.companyName)
	util.CheckErr(err)
}
