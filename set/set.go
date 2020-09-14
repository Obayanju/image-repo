package set

type StringSet struct {
	items map[string]bool
}

func (set *StringSet) Add(s string) {
	if set.items == nil {
		set.items = make(map[string]bool)
	}

	_, found := set.items[s]
	if !found {
		set.items[s] = true
	}
}
