package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Subject struct {
	Filename string
}

func (s *Subject) String() string {
	// convenience function so we can fmt.Printf
	return fmt.Sprintf("File:%q", s.Filename)
}

func (s *Subject) CompareTo(b *Subject) Comparison {
	c := Comparison{SubjectA: s, SubjectB: b}
	// TODO: part 3!
	return c
}

type Comparison struct {
	Distance uint
	SubjectA *Subject
	SubjectB *Subject
}

func (c Comparison) String() string {
	// convenience function so we can fmt.Printf
	return fmt.Sprintf("%s::%s -> %d", c.SubjectA, c.SubjectB, c.Distance)
}

type SubjectList struct {
	list []*Subject
}

func NewSubjectList() *SubjectList {
	// Assume for now we'll have the same number of subjects as arguments
	return &SubjectList{list: make([]*Subject, 0, flag.NArg())}
}

func (s *SubjectList) AddFile(path string) error {
	s.list = append(s.list, &Subject{Filename: path})
	return nil
}

func (s *SubjectList) CompareAll() []Comparison {
	// we know that a list of size N will have ((N-1)*N)/2 comparisons. this
	// relationship is the "Triangular numbers"! https://oeis.org/A000217
	r := make([]Comparison, 0, ((len(s.list)-1)*len(s.list))/2)
	for len(s.list) > 0 {
		// take the last subject
		subjecta := s.list[len(s.list)-1]
		// modify the slice to exclude it
		s.list = s.list[:len(s.list)-1]
		// now, compare all the remaining subjects
		for _, subjectb := range s.list {
			r = append(r, subjecta.CompareTo(subjectb))
		}
	}
	return r
}

func main() {
	flag.Parse()
	sl := NewSubjectList()
	for _, fn := range flag.Args() {
		filepath.Walk(fn, sl.filewalkfunc())
	}
	r := sl.CompareAll()
	fmt.Printf("%+v\n", r) // don't worry, we'll get to this.
}

func (s *SubjectList) filewalkfunc() func(p string, info os.FileInfo, err error) error {
	return func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		return s.AddFile(p)
	}
}
