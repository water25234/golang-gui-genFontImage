package main

import "sync"

// Go types that are bound to the UI must be thread-safe, because each binding
// is executed in its own goroutine. In this simple case we may use atomic
// operations, but for more complex cases one should use proper synchronization.
type generator struct {
	sync.Mutex
	pinCodeList string
	folder      string
	fileExt     string
}

type errLog struct {
	errGenCode []string
}

type jobChannel struct {
	index       int
	fileContent string
}
