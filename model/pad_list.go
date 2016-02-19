// This file is part of Darius. See License.txt for license information.

package model

// ListPad list all pads in db
func ListPad() (ret []*Pad, err error) {
	data := make(map[int]*Pad)
	rows, err := listPadQuery.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var (
			id    int
			uid   int
			title string
		)
		if err = rows.Scan(&id, &uid, &title); err != nil {
			return
		}
		data[id] = &Pad{id, uid, title, make([]string, 0), nil, make([]int, 0), nil}
	}
	if err = rows.Err(); err != nil {
		return
	}

	// fill tags
	rows, err = listTagQuery.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var name string
		var pid int
		if err = rows.Scan(&name, &pid); err != nil {
			return
		}
		data[pid].Tags = append(data[pid].Tags, name)
	}
	if err = rows.Err(); err != nil {
		return
	}

	// fill coops
	rows, err = listCoopQuery.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var uid, pid int
		if err = rows.Scan(&uid, &pid); err != nil {
			return
		}
		data[pid].Cooperators = append(data[pid].Cooperators, uid)
	}
	if err = rows.Err(); err != nil {
		return
	}

	ret = make([]*Pad, len(data))
	cnt := 0
	for _, p := range data {
		ret[cnt] = p
		cnt++
	}
	return
}
