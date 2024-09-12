package main

type Segment struct {
	ID               string `json:"_id"`
	NumberOfSegments int    `json:"number_of_segments"`
	Points           int    `json:"points"`
	ControlGroup     int    `json:"control_group"`
}

type GameBoard struct {
	Board []Segment `json:"board"`
	State Game      `json:"state"`
}

type GameSegment struct {
	ID               string
	NumberOfSegments int
	Segments         []string
	Points           int
	ControlGroup     int
	Player0          []string
	Player1          []string
}
