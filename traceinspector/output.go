package traceinspector

type AnalyzerInfoType int

const (
	AnalyzerJoin AnalyzerInfoType = iota
	AnalyzerWiden
	
)

type AnalyzerStatusOutput interface {
	isAnalyzerStatusOutput()
}

type AnalyzerError struct {
	info     string
	line_num int
}
