package osexit_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/sreway/shorturl/pkg/tools/analysis/osexit"
)

func TestFromFileSystem(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, osexit.Analyzer, "a", "b") // loads testdata/src/
}
