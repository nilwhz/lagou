package main

type jobItem struct {
	// company
	companyName    string
	companyArea    string
	companySize    string
	companyUrl     string
	companyStage   string
	companyAddress string
	// job
	jobName   string
	jobSalary string
	jobUpTime string
	jobExp    string
	jobDegree string
	jobDesc   string
	//hr
	hrName string
	// 7天内处理完成的简历所占比例
	hrPercent string
	// 完成简历处理的平均用时
	hrExecTime string
	// 拉勾网该工作信息的网址
	lagouURL string
}
