package model

type (
	Project struct {
		ProjectId     string `json:"projectid,omitempty"`
		ProjectName   string `json:"projectname"`
		ProjectDesc   string `json:"projectdesc"`
		ProjectOwner  string `json:"projectowner"`
		ProjectStatus string `json:"projectstatus,omitempty"`
	}
)
