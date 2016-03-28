package util

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func IsMoreRecent(time1 string, time2 string) bool {
	iTime, iJobTime := getFlag(time1), getFlag(time2)
	if iTime < iJobTime {
		return true
	} else if iTime > iJobTime {
		return false
	}

	// 当格式相同，且类似： 08:57   1天前
	if iTime == 0 {
		hTime, hJobTime := strings.Split(time1, ":")[0], strings.Split(time2, ":")[0]
		return !isSmaller(hTime, hJobTime)
	} else if iTime == 1 {
		r := regexp.MustCompile(`\d+`)
		dTime, dJobTime := r.FindString(time1), r.FindString(time2)
		return isSmaller(dTime, dJobTime)
	}
	// 当格式相同，且类似：2016-03-17
	yTime, mTime, dTime := strings.Split(time1, "-")[0], strings.Split(time1, "-")[1], strings.Split(time1, "-")[2]
	yJobTime, mJobTime, dJobTime := strings.Split(time2, "-")[0], strings.Split(time2, "-")[1], strings.Split(time2, "-")[2]
	if yTime != yJobTime {
		return !isSmaller(yTime, yJobTime)
	}
	if mTime != mJobTime {
		return !isSmaller(mTime, mJobTime)
	}
	if dTime != dJobTime {
		return !isSmaller(dTime, dJobTime)
	}
	return false
}

func isSmaller(s1 string, s2 string) bool {
	if strings.HasPrefix(s1, "0") && !strings.HasPrefix(s2, "0") {
		return true
	}
	if strings.HasPrefix(s2, "0") && !strings.HasPrefix(s1, "0") {
		return false
	}
	if strings.HasPrefix(s1, "0") && strings.HasPrefix(s2, "0") {
		return compareWithZero(s1, s2)
	}

	return compareWithNoZero(s1, s2)
}

func compareWithZero(s1 string, s2 string) bool {
	var n1, n2 int
	n1, _ = strconv.Atoi(strings.Trim(s1, "0"))
	n2, _ = strconv.Atoi(strings.Trim(s2, "0"))
	return n1 < n2
}

func compareWithNoZero(s1 string, s2 string) bool {
	var n1, n2 int
	n1, _ = strconv.Atoi(s1)
	n2, _ = strconv.Atoi(s2)
	return n1 < n2
}

func getFlag(s string) int {
	// 08:57
	if strings.Contains(s, ":") {
		return 0
	}
	// 2016-03-17
	if strings.Contains(s, "-") {
		return 2
	}
	// 1天前
	return 1
}
