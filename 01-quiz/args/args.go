package args

import "flag"

// Args are program arguments.
type Args struct {
	Filename  string
	Timeout   int
	Randomize bool
}

// ParseCommandlineArgs parses program arguments.
func ParseCommandlineArgs() Args {
	a := Args{}
	flag.StringVar(&a.Filename, "filename", "problems.csv", "a name of the file with questions")
	flag.IntVar(&a.Timeout, "timeout", 30, "a timeout for every question")
	flag.BoolVar(&a.Randomize, "randomize", false, "should questions be randomized")
	flag.Parse()
	return a
}
