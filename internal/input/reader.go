package input

type InputReader interface {
	// error != nil means stop iteration
	// error == io.EOF means end of input
	// error == nil means read OK, try again for potentially a new result
	ReadOne() (string, error)
}
