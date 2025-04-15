package models

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Line struct {
	Date         string `json:"date"`
	Division     string `json:"division"`
	Operation    string `json:"operation"`
	Culture      string `json:"culture"`
	PerDay       int    `json:"per_day"`
	PerOperation int    `json:"per_operation"`
	ValDay       int    `json:"val_day"`
	ValBeginning int    `json:"val_beginning"`
}

type Table []Line

func GetTableName(startedAt time.Time) string {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Println("failed to load location: %w", err)
	}

	timestamp := startedAt.In(loc)
	return timestamp.Format("04м15ч02/01/2006") + "_AgroScientists"
}

func GetDocxName(name string, number int, timestamp time.Time) string {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Println("failed to load location: %w", err)
	}

	timestamp = timestamp.In(loc)

	fileName := fmt.Sprintf("%s_%d_%s.docx", strings.ReplaceAll(name, " ", "-"), number, timestamp.Format("04м15ч02/01/2006"))

	return fileName
}
