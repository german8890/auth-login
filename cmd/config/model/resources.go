package model

type ArtifactResources struct {
	GroupName string      `json:"groupName"`
	Resources []Resources `json:"resources"`
}

type Resources struct {
	Path      string `json:"path"`
	Method    string `json:"method"`
	Operation string `json:"operation"`
}
