package bugsnack

import (
	"context"
	"fmt"
	"io"
	"sync"
)

// An ErrorReporter is used to Report errors
type ErrorReporter interface {
	Report(ctx context.Context, err error, metadata ...interface{})
}

// A MultiReporter is capable of sending a single error
// to multiple ErrorReporters
type MultiReporter struct {
	Reporters []ErrorReporter
}

// Report sends the same error to all underlying Reporters
// concurrently.
func (mr *MultiReporter) Report(ctx context.Context, err error, metadata ...interface{}) {
	var wg sync.WaitGroup
	for _, er := range mr.Reporters {
		wg.Add(1)
		go func(wg *sync.WaitGroup, er ErrorReporter) {
			defer wg.Done()
			er.Report(ctx, err)
		}(&wg, er)
	}
	wg.Wait()
}

// A WriterReporter writes errors to an io.Writer
type WriterReporter struct {
	Writer io.Writer
}

// Report printf's the error using %s, then writes it to the
// underlying writer
func (wr *WriterReporter) Report(_ context.Context, err error, metadata ...interface{}) {
	if wr.Writer == nil {
		return
	}
	fmt.Fprintf(wr.Writer, "%s\n", err)
}
