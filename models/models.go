package models

type TraceRequest struct {
	resource     string `json:"resource"`
	misconfig    string `json:"misconfig"`
	account      string `json:"account"`
	organization string `json:"organization"`
}

type GitHubIWebhook struct {
	Installation struct {
		ID int64 `json:"id"`
	} `json:"installation"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}
