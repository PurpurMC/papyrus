package types

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type ProjectsResponse struct {
	Projects []string `json:"projects"`
}

type ProjectResponse struct {
	Project   string   `json:"project"`
	CreatedAt int64    `json:"createdAt"`
	Versions  []string `json:"versions"`
}

type VersionResponse struct {
	Project   string `json:"project"`
	Version   string `json:"version"`
	CreatedAt int64  `json:"createdAt"`
	Builds    struct {
		Latest string   `json:"latest"`
		All    []string `json:"all"`
	} `json:"builds"`
}

type VersionResponseDetailed struct {
	Project   string `json:"project"`
	Version   string `json:"version"`
	CreatedAt int64  `json:"createdAt"`
	Builds    struct {
		Latest string          `json:"latest"`
		All    []BuildResponse `json:"all"`
	} `json:"builds"`
}

type BuildResponse struct {
	Project   string   `json:"project"`
	Versions  []string `json:"versions"`
	Build     string   `json:"build"`
	CreatedAt int64    `json:"createdAt"`
	Result    string   `json:"result"`
	Flags     []string `json:"flags"`
	Commits   []Commit `json:"commits"`
	Files     []File   `json:"files"`
}
