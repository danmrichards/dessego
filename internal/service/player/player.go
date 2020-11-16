package player

import "fmt"

// Stats represents a players statistics.
type Stats struct {
	GradeS   int
	GradeA   int
	GradeB   int
	GradeC   int
	GradeD   int
	Sessions int
}

// Stats returns the raw values of player stats.
func (s Stats) Vals() []int {
	return []int{
		s.GradeS, s.GradeA, s.GradeB, s.GradeC, s.GradeD, s.Sessions,
	}
}

// String implements fmt.Stringer.
func (s Stats) String() string {
	return fmt.Sprintf(
		"grade_s: %d grade_a: %d grade_b: %s grade_c: %s grade_d: %d sessions: %d",
		s.GradeS, s.GradeA, s.GradeB, s.GradeC, s.GradeD, s.Sessions,
	)
}
