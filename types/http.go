package types

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type ProjectsResponse struct {
	Projects []string `json:"projects"`
}

type ProjectResponse struct {
	Project  string   `json:"project"`
	CreatedAt int64         `json:"createdAt"`
	Versions []string `json:"versions"`
}

type VersionResponse struct {
	Project string   `json:"project"`
	Version string   `json:"version"`
	CreatedAt int64         `json:"createdAt"`
	Latest  string   `json:"latest"`
	Builds  []string `json:"builds"`
}

type VersionResponseDetailed struct {
	Project string          `json:"project"`
	Version string          `json:"version"`
	CreatedAt int64         `json:"createdAt"`
	Latest  string          `json:"latest"`
	Builds  []BuildResponse `json:"builds"`
}

type BuildResponse struct {
	Project   string         `json:"project"`
	Version   string         `json:"version"`
	Build     string         `json:"build"`
	CreatedAt int64          `json:"createdAt"`
	Result    string   `json:"result"`
	Commits   []Commit `json:"commits"`
	Files     []File   `json:"files"`
}
