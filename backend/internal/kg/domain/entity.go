package domain

type Node struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}
type Link struct {
	Source int    `json:"source"`
	Target int    `json:"target"`
	Type   string `json:"type"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}
