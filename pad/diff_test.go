package pad

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
		t.Error("empty-empty: why something to new?")
	}
	if len(d) > 0 {
		t.Error("empty-empty: why something to delete?")
	}

	pad.Tags = []string{"tag"}
	i, d = pad.tagDiff()
	if !reflect.DeepEqual(pad.Tags, i) {
		t.Errorf("A-empty: wrong to insert: %v", i)
	}
	if len(d) > 0 {
		t.Error("A-empty: why something to delete?")
	}

	pad.Tags = []string{}
	pad.oldTags = []string{"tag"}
	i, d = pad.tagDiff()
	if len(i) > 0 {
		t.Error("empty-B: why something to new?")
	}
	if !reflect.DeepEqual(pad.oldTags, d) {
		t.Errorf("empty-B: wrong to delete: %v", i)
	}

	pad.Tags = []string{"tag1", "tag2"}
	pad.oldTags = []string{"tag1", "tag3"}
	i, d = pad.tagDiff()
	if !reflect.DeepEqual([]string{"tag2"}, i) {
		t.Errorf("A-B: wrong to insert: %v", i)
	}
	if !reflect.DeepEqual([]string{"tag3"}, d) {
		t.Errorf("A-B: wrong to delete: %v", i)
	}
}
