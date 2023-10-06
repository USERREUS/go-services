package model

type Model struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name"`
	Count int    `json:"count"`
	Cost  int    `json:"cost"`
}
