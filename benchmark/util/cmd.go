package bench

type CMDParams struct {
	Concurrency string
	Duration    string
	Target      string
	Method      string
	BodyPath    string
}

var Configs = &CMDParams{}
