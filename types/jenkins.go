package types

type JenkinsData struct {
	Duration  int64  `json:"duration"`
	Result    string `json:"result"`
	Timestamp int64  `json:"timestamp"`
	ChangeSet struct {
		Items []JenkinsChangeSetItem `json:"items"`
	} `json:"changeSet"`
}

type JenkinsChangeSetItem struct {
	Hash        string                     `json:"commitId"`
	Timestamp   int64                      `json:"timestamp"`
	Author      JenkinsChangeSetItemAuthor `json:"author"`
	AuthorEmail string                     `json:"authorEmail"`
	Summary     string                     `json:"msg"`
	Description string                     `json:"comment"`
}

type JenkinsChangeSetItemAuthor struct {
	Name string `json:"fullName"`
}
