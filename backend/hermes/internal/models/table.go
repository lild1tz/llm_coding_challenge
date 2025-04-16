package models

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Line struct {
	Date string `json:"date"`

	Division       string `json:"division"`
	DivisionYellow bool

	Operation       string `json:"operation"`
	OperationYellow bool

	Culture       string `json:"culture"`
	CultureYellow bool

	PerDay       int `json:"per_day"`
	PerOperation int `json:"per_operation"`
	ValDay       int `json:"val_day"`
	ValBeginning int `json:"val_beginning"`
}

type Table []Line

func GetTableName(startedAt time.Time, chatContextName string) string {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Println("failed to load location: %w", err)
	}

	timestamp := startedAt.In(loc)

	fileName := timestamp.Format("04м15ч02/01/2006") + "_AgroScientists"
	if chatContextName != "" {
		fileName += "_" + chatContextName
	}

	return fileName
}

func GetBasicName(name string, number int, timestamp time.Time, chatContextName string) string {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Println("failed to load location: %w", err)
	}

	timestamp = timestamp.In(loc)

	fileName := fmt.Sprintf("%s_%d_%s", strings.ReplaceAll(name, " ", "-"), number, timestamp.Format("04м15ч02/01/2006"))

	if chatContextName != "" {
		fileName += "_" + chatContextName
	}

	return fileName
}

func GetDocxName(name string, number int, timestamp time.Time, chatContextName string) string {
	return GetBasicName(name, number, timestamp, chatContextName) + ".docx"
}

func GetImageName(name string, number int, timestamp time.Time, chatContextName string, postfix string) string {
	return GetBasicName(name, number, timestamp, chatContextName) + postfix
}
