package models

type nameMapper struct {
	Pattern string `json:"pattern"`
	Target  string `json:"target"`
}

type Mapping struct {
	Name []nameMapper `json:"name"`
}
