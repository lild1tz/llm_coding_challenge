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

	PerDay       int `json:"per_day"`
	PerOperation int `json:"per_operation"`
	ValDay       int `json:"val_day"`
	ValBeginning int `json:"val_beginning"`
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

func GetDocxName(name string, number int, t time.Time, chatContextName string) string {
	return GetBasicName(name, number, t, chatContextName) + ".docx"
}

func GetImageName(name string, number int, timestamp time.Time, chatContextName string, postfix string) string {
	return GetBasicName(name, number, timestamp, chatContextName) + postfix
}
