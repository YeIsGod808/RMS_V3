package model

type Node struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	SectionId   int    `json:"sectionId,omitempty"`
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

type Section struct {
	Graph   Graph `json:"graph"`
	Chapter Node  `json:"chapter"`
}

type Point struct {
	Graph   Graph `json:"graph"`
	Section Node  `json:"section"`
}

type Video struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type Exercise struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Difficulty string `json:"difficulty"`
	URL        string `json:"url"`
}
