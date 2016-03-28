package main

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"lagou/util"
	"net/http"
)

func main() {
	cfg, err := ini.Load("./config.ini")
	util.CheckErr(err)
	gdSection := cfg.Section("gaode")

	openDB()
	defer closeDB()

	r := gin.Default()
	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*")

	// 根据时间排序显示
	r.GET("/job", func(c *gin.Context) {
		c.HTML(http.StatusOK, "job_list.tmpl", gin.H{
			"key":  gdSection.Key("key"),
			"jobs": sortedByTime(findAllJobs()),
		})
	})

	r.Run()
}

// 冒泡排序
func sortedByTime(jobs []jobItem) []jobItem {
	for itemCount := len(jobs) - 1; ; itemCount-- {
		swap := false
		for i := 0; i < itemCount; i++ {
			if !util.IsMoreRecent(jobs[i].JobUpTime, jobs[i+1].JobUpTime) {
				swap = true
				jobs[i], jobs[i+1] = jobs[i+1], jobs[i]
			}
		}
		if swap == false {
			break
		}
	}
	return jobs
}
