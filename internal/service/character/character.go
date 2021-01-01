package character

import "fmt"

// MultiplayerGrade is a string representation of a Multiplayer grade.
type MultiplayerGrade string

const (
	// GradeS is the S grade.
	GradeS MultiplayerGrade = "grade_s"

	// Grade A is the A grade.
	GradeA MultiplayerGrade = "grade_a"

	// GradeB is the B grade.
	GradeB MultiplayerGrade = "grade_b"

	// GradeC is the C grade.
	GradeC MultiplayerGrade = "grade_c"

	// GradeD is the D grade.
	GradeD MultiplayerGrade = "grade_d"

	// GradeUnknown is the unknown grade.
	GradeUnknown MultiplayerGrade = "grade_unknown"
)

// Grades is a list of valid multiplayer grades and their integer key.
var Grades = map[int]MultiplayerGrade{
	0: GradeS,
	1: GradeA,
	2: GradeB,
	3: GradeC,
	4: GradeD,
}

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

// WorldTendency represents the world tendency values for a character.
type WorldTendency struct {
	Area1 int
	WB1   int
	LR1   int
	Area2 int
	WB2   int
	LR2   int
	Area3 int
	WB3   int
	LR3   int
	Area4 int
	WB4   int
	LR4   int
	Area5 int
	WB5   int
	LR5   int
	Area6 int
	WB6   int
	LR6   int
	Area7 int
	WB7   int
	LR7   int
}

// String implements fmt.Stringer.
func (w WorldTendency) String() string {
	return fmt.Sprintf(
		"area_1: %d wb_1: %d lr_1: %d area_2: %d wb_2: %d lr_2: %d area_3: %d wb_3: %d lr_3: %d area_4: %d wb_4: %d lr_4: %d area_5: %d wb_5: %d lr_5: %d area_6: %d wb_6: %d lr_6: %d area_7: %d wb_7: %d lr_7: %d",
		w.Area1, w.WB1, w.LR1,
		w.Area2, w.WB2, w.LR2,
		w.Area3, w.WB3, w.LR3,
		w.Area4, w.WB4, w.LR4,
		w.Area5, w.WB5, w.LR5,
		w.Area6, w.WB6, w.LR6,
		w.Area7, w.WB7, w.LR7,
	)
}
