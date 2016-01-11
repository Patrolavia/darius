package model

import (
	"reflect"
	"testing"
)

func TestTagDiff(t *testing.T) {
	pad := &PadContent{Pad: &Pad{}}
	pad.Tags = []string{}
	pad.oldTags = []string{}
	i, d := pad.tagDiff()
	if len(i) > 0 {
		t.Error("tagdiff: empty-empty: why something to new?")
	}
	if len(d) > 0 {
		t.Error("tagdiff: empty-empty: why something to delete?")
	}

	pad.Tags = []string{"tag"}
	i, d = pad.tagDiff()
	if !reflect.DeepEqual(pad.Tags, i) {
		t.Errorf("tagdiff: A-empty: wrong to insert: %v", i)
	}
	if len(d) > 0 {
		t.Error("tagdiff: A-empty: why something to delete?")
	}

	pad.Tags = []string{}
	pad.oldTags = []string{"tag"}
	i, d = pad.tagDiff()
	if len(i) > 0 {
		t.Error("tagdiff: empty-B: why something to new?")
	}
	if !reflect.DeepEqual(pad.oldTags, d) {
		t.Errorf("tagdiff: empty-B: wrong to delete: %v", i)
	}

	pad.Tags = []string{"tag1", "tag2"}
	pad.oldTags = []string{"tag1", "tag3"}
	i, d = pad.tagDiff()
	if !reflect.DeepEqual([]string{"tag2"}, i) {
		t.Errorf("tagdiff: A-B: wrong to insert: %v", i)
	}
	if !reflect.DeepEqual([]string{"tag3"}, d) {
		t.Errorf("tagdiff: A-B: wrong to delete: %v", i)
	}
}

func TestCoopDiff(t *testing.T) {
	pad := &PadContent{Pad: &Pad{}}
	pad.Cooperators = []int{}
	pad.oldCoops = []int{}
	i, d := pad.coopDiff()
	if len(i) > 0 {
		t.Error("coopdiff: empty-empty: why something to new?")
	}
	if len(d) > 0 {
		t.Error("coopdiff: empty-empty: why something to delete?")
	}

	pad.Cooperators = []int{1}
	i, d = pad.coopDiff()
	if !reflect.DeepEqual(pad.Cooperators, i) {
		t.Errorf("coopdiff: A-empty: wrong to insert: %v", i)
	}
	if len(d) > 0 {
		t.Error("coopdiff: A-empty: why something to delete?")
	}

	pad.Cooperators = []int{}
	pad.oldCoops = []int{1}
	i, d = pad.coopDiff()
	if len(i) > 0 {
		t.Error("coopdiff: empty-B: why something to new?")
	}
	if !reflect.DeepEqual(pad.oldCoops, d) {
		t.Errorf("coopdiff: empty-B: wrong to delete: %v", i)
	}

	pad.Cooperators = []int{1, 2}
	pad.oldCoops = []int{1, 3}
	i, d = pad.coopDiff()
	if !reflect.DeepEqual([]int{2}, i) {
		t.Errorf("coopdiff: A-B: wrong to insert: %v", i)
	}
	if !reflect.DeepEqual([]int{3}, d) {
		t.Errorf("coopdiff: A-B: wrong to delete: %v", i)
	}
}
