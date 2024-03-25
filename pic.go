package main

import (
	"fmt"
	"sort"
	"time"
)

type column []int

// 近六个月提交记录图表
func Pic(email string) {
	commit := processRepo(email)
	//fmt.Println(commit)
	printCommitStats(commit)
}
func printCommitStats(commits map[int]int) {
	cols := generateColumn(commits)
	printCells(cols)
}

// 列表，26~1周 星期几 commit
func generateColumn(m map[int]int) map[int]column {
	keys := sortMapIntoSlice(m)
	cols := make(map[int]column)
	col := column{}

	for _, k := range keys {
		week := int(k / 7) //26,25...1
		dayinweek := k % 7 // 0,1,2,3,4,5,6

		if dayinweek == 0 { //reset星期天是每周第一天
			col = column{}
		}

		col = append(col, m[k])

		if dayinweek == 6 {
			cols[week] = col
		}
	}

	return cols
}

func printMonth() {
	week := getBeginningOfDay(time.Now()).Add(-(daysInLastSixMonths * time.Hour * 24))
	month := week.Month()
	fmt.Printf("         ")
	for {
		if week.Month() != month {
			fmt.Printf("%s ", week.Month().String()[:3])
			month = week.Month()
		} else {
			fmt.Printf("    ")
		}

		week = week.Add(7 * time.Hour * 24)
		if week.After(time.Now()) {
			break
		}
	}
	fmt.Printf("\n")
}
func printDayCol(day int) {
	out := "     "
	switch day {
	case 1:
		out = " Mon "
	case 3:
		out = " Wed "
	case 5:
		out = " Fri "
	}

	fmt.Print(out)
}

func printCell(val int, today bool) {
	escape := "\033[0;37;30m" //黑色
	switch {
	case val > 0 && val < 5:
		escape = "\033[1;30;47m" //白色
	case val >= 5 && val < 10:
		escape = "\033[1;30;43m" //黄色
	case val >= 10:
		escape = "\033[1;30;42m" //绿色
	}

	if today {
		escape = "\033[1;37;45m" //紫色
	}

	if val == 0 {
		fmt.Printf(escape + "  - " + "\033[0m")
		return
	}

	str := "  %d "
	switch {
	case val >= 10:
		str = " %d "
	case val >= 100:
		str = "%d "
	}

	fmt.Printf(escape+str+"\033[0m", val)
}

// printCells prints the cells of the graph
func printCells(cols map[int]column) {
	printMonth()
	for j := 6; j >= 0; j-- { //星期几 星期六~星期天
		for i := weeksInLastSixMonths + 1; i >= 0; i-- { //那一周
			if i == weeksInLastSixMonths+1 {
				printDayCol(j)
			}
			if col, ok := cols[i]; ok { //这周有提交记录
				//special case today
				if i == 0 && j == int(time.Now().Weekday()) {
					printCell(col[j], true)
					continue
				} else {
					if len(col) > j {
						printCell(col[j], false)
						continue
					}
				}
			}
			printCell(0, false)
		}
		fmt.Printf("\n")
	}
}
func sortMapIntoSlice(m map[int]int) []int {
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	return keys
}
