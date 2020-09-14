package set

type StringSet struct {
	items map[string]bool
}

func (s *StringSet) Add(str string) {
	if s.items == nil {
		s.items = make(map[string]bool)
	}

	_, found := s.items[str]
	if !found {
		s.items[str] = true
	}
}

func (s *StringSet) Items() []string {
	items := []string{}
	for i := range s.items {
		items = append(items, i)
	}
	return items
}
