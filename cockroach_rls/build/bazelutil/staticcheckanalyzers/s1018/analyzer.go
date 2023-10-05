// Code generated by generate-staticcheck; DO NOT EDIT.

//go:build bazel
// +build bazel

package s1018

import (
	util "github.com/cockroachdb/cockroach/pkg/testutils/lint/passes/staticcheck"
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/simple"
)

var Analyzer *analysis.Analyzer

func init() {
	for _, analyzer := range simple.Analyzers {
		if analyzer.Analyzer.Name == "S1018" {
			Analyzer = analyzer.Analyzer
			break
		}
	}
	util.MungeAnalyzer(Analyzer)
}
