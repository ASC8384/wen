package ast

/*
stat            ::=  Assignment | Print | label | break | goto Name | if exp then block "elseif exp then block"* "else block"? end
*/
type Stat interface{}
type EmptyStat struct{} // ‘;’

type AssignStat struct {
	Line     int
	Variable *Variable
	String   string
	Length   int
	Start    int
	Int      bool
}

type Print struct {
	Line     int
	Variable *Variable
	Int      bool
}

type Scanf struct {
	Line     int
	Variable *Variable
	Int      bool
}

type PointerStat struct {
	Line    int
	Pointer int
}

type CellStat struct {
	Line int
	Cell int
}

type LoopStat struct {
	Line  int
	Stats []Stat
}

// 属于
var _ Stat = (*Print)(nil)
var _ Stat = (*AssignStat)(nil)
var _ Stat = (*PointerStat)(nil)
var _ Stat = (*CellStat)(nil)
