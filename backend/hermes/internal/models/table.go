package models

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
