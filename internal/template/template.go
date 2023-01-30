package template

type tmplVars map[string]string

type tmplJson struct {
	Version      string      `json:"version"`
	Placeholders tmplVars    `json:"placeholders"`
	Excludes     []string    `json:"excludes"`
	Validation   []validator `json:"validation"`
}
