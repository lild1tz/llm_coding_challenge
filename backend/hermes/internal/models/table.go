package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/lild1tz/llm_coding_challenge/backend/libs/go/loctime"
)

type Line struct {
	Date string `json:"date"`

	Division       string `json:"division"`
	DivisionYellow bool

	Operation       string `json:"operation"`
	OperationYellow bool

	Culture       string `json:"culture"`
	CultureYellow bool

	PerDay       string `json:"per_day"`
	PerOperation string `json:"per_operation"`
	ValDay       string `json:"val_day"`
	ValBeginning string `json:"val_beginning"`
}

type Table []Line

func GetTableName(t time.Time, chatContextName string) string {
	t = loctime.Transfer(t)

	fileName := t.Format("04м15ч02/01/2006") + "_AgroScientists"
	if chatContextName != "" {
		fileName += "_" + chatContextName
	}

	return fileName
}

func GetBasicName(name string, number int, t time.Time, chatContextName string) string {
	t = loctime.Transfer(t)

	fileName := fmt.Sprintf("%s_%d_%s", strings.ReplaceAll(name, " ", "-"), number, t.Format("04м15ч02/01/2006"))

	if chatContextName != "" {
		fileName += "_" + chatContextName
	}

	return fileName
}

func GetFileName(name string, number int, timestamp time.Time, chatContextName string, postfix string) string {
	return GetBasicName(name, number, timestamp, chatContextName) + postfix
}

func GetDocxName(name string, number int, timestamp time.Time, chatContextName string) string {
	return GetFileName(name, number, timestamp, chatContextName, ".docx")
}
