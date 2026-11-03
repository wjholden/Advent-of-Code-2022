package main

import "fmt"

type SelfReference struct {
	self *SelfReference
}

func make() SelfReference {
	s := SelfReference{}
	s.self = &s
	return s
}

func main() {
	s := make()
	fmt.Println(s == *s.self, &s == s.self)
}
