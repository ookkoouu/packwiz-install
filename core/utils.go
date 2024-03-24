package core

import (
	"cmp"
	"slices"
)

func diffSlice[S ~[]E, E cmp.Ordered](old, new S) (added, removed, unchanged S) {
	if !slices.IsSorted(old) {
		slices.Sort(old)
	}
	if !slices.IsSorted(new) {
		slices.Sort(new)
	}

	if len(old) == 0 || len(new) == 0 {
		return new, old, unchanged
	}

	var (
		idxOld = 0
		idxNew = 0
		limOld = len(old) - 1
		limNew = len(new) - 1
	)

	for idxOld <= limOld && idxNew <= limNew {
		oldValue := old[idxOld]
		newValue := new[idxNew]
		cmpr := cmp.Compare(oldValue, newValue)

		// old < new
		if cmpr == -1 {
			// oldValue is removed
			removed = append(removed, oldValue)
			idxOld++
			continue
		}

		// old = new
		if cmpr == 0 {
			unchanged = append(unchanged, newValue)
			idxOld++
			idxNew++
			continue
		}

		// old > new
		if cmpr == 1 {
			// newValue is added
			added = append(added, newValue)
			idxNew++
			continue
		}
	}

	// check remainings
	if idxNew > limNew && idxOld <= limOld {
		for ; idxOld <= limOld; idxOld++ {
			removed = append(removed, old[idxOld])
		}
	}
	if idxOld > limOld && idxNew <= limNew {
		for ; idxNew <= limNew; idxNew++ {
			added = append(added, new[idxNew])
		}
	}

	return
}

func diffSliceFunc[S ~[]E, E any](old S, new S, cmp func(E, E) int) (added, removed, unchanged []E) {
	if !slices.IsSortedFunc(old, cmp) {
		slices.SortFunc(old, cmp)
	}
	if !slices.IsSortedFunc(new, cmp) {
		slices.SortFunc(new, cmp)
	}

	if len(old) == 0 || len(new) == 0 {
		return new, old, unchanged
	}

	var (
		idxOld = 0
		idxNew = 0
		limOld = len(old) - 1
		limNew = len(new) - 1
	)

	for idxOld <= limOld && idxNew <= limNew {
		oldValue := old[idxOld]
		newValue := new[idxNew]
		cmpr := cmp(oldValue, newValue)

		// old < new
		if cmpr == -1 {
			// oldValue is removed
			removed = append(removed, oldValue)
			idxOld++
			continue
		}

		// old = new
		if cmpr == 0 {
			unchanged = append(unchanged, newValue)
			idxOld++
			idxNew++
			continue
		}

		// old > new
		if cmpr == 1 {
			// newValue is added
			added = append(added, newValue)
			idxNew++
			continue
		}
	}

	// check remainings
	if idxNew > limNew && idxOld <= limOld {
		for ; idxOld <= limOld; idxOld++ {
			removed = append(removed, old[idxOld])
		}
	}
	if idxOld > limOld && idxNew <= limNew {
		for ; idxNew <= limNew; idxNew++ {
			added = append(added, new[idxNew])
		}
	}

	return
}
