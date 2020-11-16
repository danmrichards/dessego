package player

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
