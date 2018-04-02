package models

type ApiVerion struct {
	Version     string `json:"verion"`
	Status      string `json:"status"`
	BuildDate   string `json:"build_date"`
	BuildCommit string `json:"build_commit"`
}
