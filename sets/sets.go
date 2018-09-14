package sets

import "fmt"
import "strings"

var exists = struct{}{}

type Set struct {
	m map[string]struct{}
}

func NewSet() *Set {
	s := &Set{}
	s.m = make(map[string]struct{})
	return s
}

func (s *Set) Add(value string) {
	s.m[value] = exists
}

func (s *Set) Remove(value string) {
	delete(s.m, value)
}

func (s *Set) Contains(value string) bool {
	_, c := s.m[value]
	return c
}

func (s *Set) PrintAll() {
	for v, _ := range s.m {
		fmt.Println(v)
	}
}

func (s *Set) ReturnThird() []string {
	var sThird []string
	var listT []string
	for v, _ := range s.m {
		sThird = strings.Split(v, ",")
		lengthThird := len(sThird)
		if lengthThird > 3 {
			listT = append(listT, sThird[lengthThird-3])

		}
	}
	return listT
}

func (s *Set) PrintThird() {
	var sThird []string
	for v, _ := range s.m {
		sThird = strings.Split(v, ",")
		fmt.Println(sThird[1])
		lengthThird := len(sThird)
		if lengthThird > 4 {
			third := sThird[lengthThird-3]
			fmt.Println(third)
		}
	}
}
