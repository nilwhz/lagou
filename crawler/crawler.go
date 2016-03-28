package main

import (
	"encoding/json"
	"fmt"
	gq "github.com/PuerkitoBio/goquery"
	sjs "github.com/bitly/go-simplejson"
	"lagou/util"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	listURL   = "http://www.lagou.com/jobs/positionAjax.json?city=%E5%8C%97%E4%BA%AC"
	itemURL   = "http://www.lagou.com/jobs/"
	pageCount = 3
	kd        = "Go"
)

func main() {
	openDB()

	jobIDs := getJobIDs()
	saveJobDetail(jobIDs)
}

func getJobIDs() []string {
	var jobIDs []string
	for pn := 1; pn <= pageCount; pn++ {
		resp, err := http.PostForm(listURL, url.Values{
			"kd": {kd},
			"pn": {strconv.Itoa(pn)},
		})
		defer resp.Body.Close()
		util.CheckErr(err)

		js, err := sjs.NewFromReader(resp.Body)
		util.CheckErr(err)

		datas, err := js.Get("content").Get("result").Array()
		util.CheckErr(err)

		for _, data := range datas {
			jobIDs = append(jobIDs, data.(map[string]interface{})["positionId"].(json.Number).String())
		}
	}
	return jobIDs
}

var jobList map[string]string = map[string]string{}

func saveJobDetail(ids []string) {
	defer closeDB()

	// 每次爬取新数据的时候，先清空历史数据。
	clearHistoryData()

	counter := 1

	for _, id := range ids {
		url := itemURL + id + ".html"
		doc, err := gq.NewDocument(url)
		util.CheckErr(err)

		// 该工作信息已存在，且数据库中信息更新，放弃执行本次操作。
		job, flag := makeJobItem(doc, url)
		if flag == "useless" {
			continue
		}

		// 执行插入操作。
		if flag == "add" {
			// 插入数据
			insertOneJob(counter, job)

			jobList[job.companyName] = job.jobUpTime

			counter += 1
			continue
		}

		// 更新数据
		updateOneJob(job)
		jobList[job.companyName] = job.jobUpTime
	}
}

func makeJobItem(doc *gq.Document, url string) (jobItem, string) {
	var job jobItem

	companyName := strings.TrimSpace(doc.Find(".job_company dt h2").Contents().Eq(0).Text())
	jobUpTime := strings.TrimSpace(strings.Split(doc.Find(".publish_time").Text(), " ")[0])
	state := checkJobState(companyName, jobUpTime)
	if state == "useless" {
		return job, "useless"
	}

	// company
	job.companyName = companyName
	job.companyArea = strings.TrimSpace(doc.Find(".job_company .c_feature").Eq(0).Find("li").Eq(0).Contents().Eq(2).Text())
	job.companySize = strings.TrimSpace(doc.Find(".job_company .c_feature").Eq(0).Find("li").Eq(1).Contents().Eq(2).Text())
	job.companyUrl = doc.Find(".job_company .c_feature").Eq(0).Find("li").Eq(2).Find("a").AttrOr("href", "")
	job.companyStage = strings.TrimSpace(doc.Find(".job_company .c_feature").Eq(1).Find("li").Contents().Eq(1).Text())
	job.companyAddress = strings.TrimSpace(doc.Find(".job_company dd div").Eq(0).Text())
	// job
	job.jobName = strings.TrimSpace(doc.Find(".job_detail dt h1").First().AttrOr("title", ""))
	job.jobSalary = strings.TrimSpace(doc.Find(".job_request p").Eq(0).Find("span").Eq(0).Text())
	job.jobUpTime = jobUpTime
	job.jobExp = strings.TrimSpace(doc.Find(".job_request p").Eq(0).Find("span").Eq(2).Text())
	job.jobDegree = strings.TrimSpace(doc.Find(".job_request p").Eq(0).Find("span").Eq(3).Text())
	job.jobDesc = doc.Find(".job_bt").Text()
	job.hrName = strings.TrimSpace(doc.Find(".publisher_name a").AttrOr("title", ""))
	job.hrPercent = doc.Find(".publisher_data .data").Eq(0).Text()
	job.hrExecTime = doc.Find(".publisher_data .data").Eq(1).Text()
	job.lagouURL = url

	return job, state
}

var logCounter = 1

func checkJobState(name string, time string) string {
	if len(jobList) == 0 {
		fmt.Printf("插入第%d条数据\n", logCounter)
		logCounter += 1
		return "add"
	}

	for companyName, jobTime := range jobList {
		if companyName == name && util.IsMoreRecent(time, jobTime) {
			fmt.Printf("~~~更新数据\t%s\t历史数据时间:%s\t本次数据时间:%s\n", name, jobTime, time)
			return "update"
		} else if companyName == name && !util.IsMoreRecent(time, jobTime) {
			fmt.Printf("!!!放弃更改\t%s\t历史数据时间:%s\t本次数据时间:%s\n", name, jobTime, time)
			return "useless"
		}
	}

	fmt.Printf("插入第%d条数据\n", logCounter)
	logCounter += 1
	return "add"
}
