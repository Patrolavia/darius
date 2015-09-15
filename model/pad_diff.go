package model

func (pad *Pad) tagDiff() (newTags, goneTags []string) {
	newTags = make([]string, 0, len(pad.Tags))
	goneTags = make([]string, 0, len(pad.oldTags))

	oldMap := make(map[string]bool)
	newMap := make(map[string]bool)

	for _, t := range pad.oldTags {
		oldMap[t] = true
	}

	for _, t := range pad.Tags {
		newMap[t] = true
	}

	for _, t := range pad.Tags {
		if !oldMap[t] {
			newTags = append(newTags, t)
		}
	}

	for _, t := range pad.oldTags {
		if !newMap[t] {
			goneTags = append(goneTags, t)
		}
	}

	return
}

func (pad *Pad) coopDiff() (newCoops, goneCoops []int) {
	newCoops = make([]int, 0, len(pad.Cooperators))
	goneCoops = make([]int, 0, len(pad.oldCoops))

	oldMap := make(map[int]bool)
	newMap := make(map[int]bool)

	for _, t := range pad.oldCoops {
		oldMap[t] = true
	}

	for _, t := range pad.Cooperators {
		newMap[t] = true
		if !oldMap[t] {
			newCoops = append(newCoops, t)
		}
	}

	for _, t := range pad.oldCoops {
		if !newMap[t] {
			goneCoops = append(goneCoops, t)
		}
	}

	return
}
