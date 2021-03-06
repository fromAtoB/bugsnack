package bugsnack

import (
	"context"
	"errors"
	"net/http"
	"os"
	"runtime"
	"testing"
)

func TestErrorReporter(t *testing.T) {
	if os.Getenv("BUGSNAG_TEST") != "T" {
		t.Skip("not running bugsnag reporter test")
	}

	er := &BugsnagReporter{
		APIKey:       os.Getenv("BUGSNAG_API_KEY"),
		Doer:         http.DefaultClient,
		ReleaseStage: "development",
		Backup:       nil,
	}

	er.Report(context.Background(), errors.New("bugsnag error test"))

	er.Report(context.Background(), errors.New("bugsnag test"), &BugsnagMetadata{
		GroupingHash: "net.timeout",
		EventMetadata: &map[string]interface{}{
			"data": map[string]interface{}{
				"os": runtime.GOOS,
			},
			"key1": "value1",
			"key2": "value2",
			"arbitraryData": map[string]interface{}{
				"goVersion": runtime.Version(),
				"nested": map[string]interface{}{
					"nestedKey": "value",
				},
			},
		},
	})
}

func TestNestedErrorReporter(t *testing.T) {
	if os.Getenv("BUGSNAG_TEST") != "T" {
		t.Skip("not running bugsnag reporter test")
	}

	er := MultiReporter{
		Reporters: []ErrorReporter{&BugsnagReporter{
			APIKey:       os.Getenv("BUGSNAG_API_KEY"),
			Doer:         http.DefaultClient,
			ReleaseStage: "development",
			Backup:       nil,
		}}}

	er.Report(context.Background(), errors.New("bugsnag multireporter test"))
}
