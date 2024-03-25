package main

import (
	"fmt"
	"log"

	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const outOfRange = 99999
const daysInLastSixMonths = 183
const weeksInLastSixMonths = 26

// 统计过去六个月每天的提交次数
func processRepo(email string) map[int]int {
	filepath := GetSaveFilePath()
	//读取文件
	repos := fileToSlice(filepath)
	dayInmap := daysInLastSixMonths
	commits := make(map[int]int)
	for i := dayInmap; i > 0; i-- {
		commits[i] = 0
	}
	for _, path := range repos {
		commits = getCommits(email, path, commits)
	}
	return commits
}

// fillCommits given a repository found in `path`, gets the commits and
// puts them in the `commits` map, returning it when completed
func getCommits(email, path string, commits map[int]int) map[int]int {
	repo, err := git.PlainOpen(path)
	if err != nil {
		fmt.Println(err)
		log.Panic(err)
	}
	//获取提交历史
	head, err := repo.Head()
	if err != nil {
		fmt.Println("Error:", err)
		log.Panic(err)
	}
	// get the commits history starting from HEAD
	repolog, err := repo.Log(&git.LogOptions{From: head.Hash()})
	if err != nil {
		panic(err)
	}
	//遍历提交历史
	offset := Calloffset()
	repolog.ForEach(func(c *object.Commit) error {
		daysAgo := countDaysSinceDate(c.Author.When) + offset
		if c.Author.Email != email {
			return nil
		}
		if daysAgo != outOfRange {
			commits[daysAgo]++ //距离这周星期日的天数
		}
		return nil
	})
	return commits
}

// 今天离星期天的位置
func Calloffset() int {
	var offset int
	weekday := time.Now().Weekday()

	switch weekday {
	case time.Sunday:
		offset = 7
	case time.Monday:
		offset = 6
	case time.Tuesday:
		offset = 5
	case time.Wednesday:
		offset = 4
	case time.Thursday:
		offset = 3
	case time.Friday:
		offset = 2
	case time.Saturday:
		offset = 1
	}

	return offset
}

// getBeginningOfDay given a time.Time calculates the start time of that day
func getBeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return startOfDay
}

// countDaysSinceDate counts how many days passed since the passed `date`
func countDaysSinceDate(date time.Time) int {
	days := 0
	now := getBeginningOfDay(time.Now())
	for date.Before(now) {
		date = date.Add(time.Hour * 24)
		days++
		if days > daysInLastSixMonths {
			return outOfRange
		}
	}
	return days
}
