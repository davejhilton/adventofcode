package set

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(v T) {
	s[v] = struct{}{}
}
func (s Set[T]) AddAll(s2 Set[T]) {
	for v := range s2 {
		s[v] = struct{}{}
	}
}
func (s Set[T]) Remove(v T) {
	delete(s, v)
}
func (s Set[T]) Has(v T) bool {
	_, has := s[v]
	return has
}
func (s Set[T]) Size() int {
	return len(s)
}
func (s Set[T]) ToSlice() []T {
	l := make([]T, 0, len(s))
	for k := range s {
		l = append(l, k)
	}
	return l
}
func (s Set[T]) Filter(f func(T) bool) Set[T] {
	s2 := make(Set[T])
	for v := range s {
		if f(v) {
			s2[v] = struct{}{}
		}
	}
	return s2
}

func (s Set[T]) ToMap(defaultVal any) map[T]any {
	m := make(map[T]any)
	for k := range s {
		m[k] = defaultVal
	}
	return m
}

func New[T comparable](vals ...T) Set[T] {
	var s = make(Set[T])
	for _, v := range vals {
		s[v] = struct{}{}
	}
	return s
}

func FromSlice[T comparable](l []T) Set[T] {
	var s = make(Set[T])
	for _, v := range l {
		s[v] = struct{}{}
	}
	return s
}
func FromMapKeys[T comparable](m map[T]any) Set[T] {
	var s = make(Set[T])
	for k := range m {
		s[k] = struct{}{}
	}
	return s
}

func Union[T comparable](s1, s2 Set[T]) Set[T] {
	u := Set[T]{}
	for v := range s1 {
		u[v] = struct{}{}
	}
	for v := range s2 {
		u[v] = struct{}{}
	}
	return u
}

func Diff[T comparable](s1, s2 Set[T]) Set[T] {
	d := Set[T]{}
	var ok bool
	for v := range s1 {
		if _, ok = s1[v]; !ok {
			d[v] = struct{}{}
		}
	}
	return d
}
