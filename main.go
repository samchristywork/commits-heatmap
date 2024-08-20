package main

import (
	"fmt"
	"github.com/go-pdf/fpdf"
	"os"
	"time"
)

type Date struct {
	Year  int
	Month int
	Day   int
}

type Color struct {
	R int
	G int
	B int
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

func getEarliestDate(dates []Date) Date {
	earliest := dates[0]
	for _, date := range dates {
		if date.Year < earliest.Year {
			earliest = date
		} else if date.Year == earliest.Year {
			if date.Month < earliest.Month {
				earliest = date
			} else if date.Month == earliest.Month {
				if date.Day < earliest.Day {
					earliest = date
				}
			}
		}
	}
	return earliest
}

func getLatestDate(dates []Date) Date {
	latest := dates[0]
	for _, date := range dates {
		if date.Year > latest.Year {
			latest = date
		} else if date.Year == latest.Year {
			if date.Month > latest.Month {
				latest = date
			} else if date.Month == latest.Month {
				if date.Day > latest.Day {
					latest = date
				}
			}
		}
	}
	return latest
}

func mix(a Color, b Color, amount float64) Color {
	r := int(float64(a.R)*(1-amount) + float64(b.R)*amount)
	g := int(float64(a.G)*(1-amount) + float64(b.G)*amount)
	bb := int(float64(a.B)*(1-amount) + float64(b.B)*amount)
	return Color{r, g, bb}
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

	yearOffset := 0.0
	month := -1
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

	pdf := fpdf.New("P", "in", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(127, 127, 127)

	cellSize := 0.15
	xOffset := 0.2
	yOffset := 0.2

	for year := 2022; year <= 2024; year++ {
		pdf.TransformBegin()
		pdf.TransformTranslate(xOffset-0.05, yOffset+yearOffset*1.4-1.4)
		pdf.TransformRotate(90, 1, 1)
		pdf.Text(0, 0, fmt.Sprintf("%d", year))
		pdf.TransformEnd()

		start := time.Date(year, 1, 1, 12, 0, 0, 0, time.UTC)
		row := int(start.Weekday())
		column := 0

		for date := start; year == date.Year(); date = date.AddDate(0, 0, 1) {
			if yearOffset == 0 && month != int(date.Month()) {
				month = int(date.Month())
				pdf.Text(xOffset+float64(column)*cellSize, yOffset-0.05, months[month-1])
			}

			d := Date{date.Year(), int(date.Month()), date.Day()}
			count := contributionsOnDay(dates, d)

			c := float64(count) / 15

			if c > 1 {
				c = 1
			}

			greenColor := Color{0, 255, 0}
			defaultColor := Color{230, 230, 230}
			actual := mix(defaultColor, greenColor, c)
			pdf.SetFillColor(actual.R, actual.G, actual.B)
			pdf.Rect(
				xOffset+float64(column)*cellSize,
				yOffset+(float64(row)+yearOffset*9)*cellSize,
				cellSize*0.9, cellSize*0.9, "F")

			row++
			if row == 7 {
				row = 0
				column++
			}
		}
		yearOffset += 0.85
	}

	startDate := time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 12, 0, 0, 0, time.UTC)

	for date := startDate; date.Before(endDate); date = date.AddDate(0, 0, 1) {
		d := Date{date.Year(), int(date.Month()), date.Day()}
		count := contributionsOnDay(dates, d)
		fmt.Printf("%d-%d-%d: %d\n", d.Year, d.Month, d.Day, count)
	}
}
