package gears

import "fmt"

func ZeroLoc() Location {
	return Location{Line: 1, Col: 1, Offset: 0, End: 0}
}

type Location struct {
	Line   int
	Col    int
	Offset int
	End    int
}

func (l Location) String() string {
	return fmt.Sprintf("line %d col %d", l.Line, l.Col)
}
