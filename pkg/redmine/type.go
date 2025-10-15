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

type Watcher struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Changeset struct {
	Revision    string   `json:"revision,omitempty"`
	User        Resource `json:"user,omitempty"`
	Comments    string   `json:"comments,omitempty"`
	CommittedOn string   `json:"committed_on,omitempty"`
}
