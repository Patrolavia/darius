package pad

import (
	"reflect"
	"testing"
)

func (o *Pad) equals(n *Pad) (ret bool) {
	if o.ID != n.ID || o.UID != n.UID || o.Title != n.Title {
		return
	}

	if !reflect.DeepEqual(o.Tags, n.Tags) {
		return
	}

	return reflect.DeepEqual(o.Cooperators, n.Cooperators)
}

func TestList(t *testing.T) {
	pad, err := New(db, u.ID, "list", "content", []string{"tag1", "tag2"}, []int{coop.ID})
	if err != nil {
		t.Fatalf("Error creating pad: %s", err)
	}

	p, err := List()
	if err != nil {
		t.Fatalf("Error listing pads: %s", err)
	}

	for _, pp := range p {
		if pp.equals(pad.Pad) {
			return
		}
	}
	t.Error("WTF did I list?")
}
