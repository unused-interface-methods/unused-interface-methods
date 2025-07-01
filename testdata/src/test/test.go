package test

type Writer interface {
	Write([]byte) (int, error)
	Close() error // want "method \"Close\" of interface \"Writer\" is declared but not used"
}

type Reader interface {
	Read([]byte) (int, error)
	Unused() error // want "method \"Unused\" of interface \"Reader\" is declared but not used"
}

func test(w Writer, r Reader) {
	w.Write(nil)
	r.Read(nil)
	// Close and Unused are not called
}
