package model

import "fmt"

// Pad is part of pad record saved in db
type Pad struct {
	ID          int      `json:"id"`
	UID         int      `json:"user"`
	Title       string   `json:"title"`
	Tags        []string `json:"tags"`
	oldTags     []string
	Cooperators []int `json:"cooperator"`
	oldCoops    []int
}

// PadContent have all info stored in db
type PadContent struct {
	*Pad
	Content string `json:"content"`
	HTML    string `json:"html"`
	Version int    `json:"version"`
}

type VersionError int

func (e VersionError) Error() string {
	return fmt.Sprintf("While saving pad into db: Pad version in db differs from provided %d", int(e))
}
