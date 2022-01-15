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
}

type Print struct {
	Line     int
	Variable *Variable
}

// 属于
var _ Stat = (*Print)(nil)
var _ Stat = (*AssignStat)(nil)
