package bench

type cmdParams struct {
	Concurrency string
	Duration    string
	Target      string
	Method      string
	BodyPath    string
}

//Configs is cmd configs
var Configs = &cmdParams{}
