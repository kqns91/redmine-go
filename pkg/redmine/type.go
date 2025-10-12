package redmine

type CustomField struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Value any    `json:"value,omitempty"`
}

type Resource struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
