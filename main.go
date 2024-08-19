package main

import (
	"fmt"
	"os"
	"time"
)

type Date struct {
	Year  int
	Month int
	Day   int
}

func contributionsOnDay(dates []Date, date Date) int {
	count := 0
	for _, d := range dates {
		if d == date {
			count++
		}
	}
	return count
}

func main() {
	filename := "all_dates"

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	var dates []Date
	for {
		var date Date
		_, err := fmt.Fscanf(file, "%d-%d-%d\n", &date.Year, &date.Month, &date.Day)
		if err != nil {
			break
		}
		dates = append(dates, date)
	}

	for _, date := range dates {
		fmt.Println(date)
	}

	startDate := time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 12, 0, 0, 0, time.UTC)

	for date := startDate; date.Before(endDate); date = date.AddDate(0, 0, 1) {
		d := Date{date.Year(), int(date.Month()), date.Day()}
		count := contributionsOnDay(dates, d)
		fmt.Printf("%d-%d-%d: %d\n", d.Year, d.Month, d.Day, count)
	}
}
