package character

import "fmt"

// Stats represents a characters statistics.
type Stats struct {
	GradeS   int
	GradeA   int
	GradeB   int
	GradeC   int
	GradeD   int
	Sessions int
}

// Stats returns the raw values of character stats.
func (s Stats) Vals() []int {
	return []int{
		s.GradeS, s.GradeA, s.GradeB, s.GradeC, s.GradeD, s.Sessions,
	}
}

// String implements fmt.Stringer.
func (s Stats) String() string {
	return fmt.Sprintf(
		"grade_s: %d grade_a: %d grade_b: %d grade_c: %d grade_d: %d sessions: %d",
		s.GradeS, s.GradeA, s.GradeB, s.GradeC, s.GradeD, s.Sessions,
	)
}
