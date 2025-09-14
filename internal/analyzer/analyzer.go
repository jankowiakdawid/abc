package analyzer

import (
	"github.com/abc-metrics/abc/internal/metrics"
)

// Analyzer defines the interface for language-specific analyzers
type Analyzer interface {
	// AnalyzeFile analyzes a single file and returns ABC metrics
	AnalyzeFile(filePath string) (metrics.ABCMetrics, error)

	// SupportedExtensions returns a list of file extensions supported by this analyzer
	SupportedExtensions() []string
}

// GetAnalyzerForFile returns the appropriate analyzer for the given file path
// based on the file extension
func GetAnalyzerForFile(filePath string) (Analyzer, error) {
	// Initialize available analyzers
	analyzers := []Analyzer{
		NewGoAnalyzer(),
		// Add more analyzers as they are implemented
		// NewTypeScriptAnalyzer(),
	}

	// Find the first analyzer that supports the file extension
	for _, a := range analyzers {
		for _, ext := range a.SupportedExtensions() {
			if HasExtension(filePath, ext) {
				return a, nil
			}
		}
	}

	return nil, &UnsupportedFileError{FilePath: filePath}
}

// HasExtension checks if a file path has the given extension
func HasExtension(filePath, extension string) bool {
	if len(filePath) < len(extension) {
		return false
	}
	return filePath[len(filePath)-len(extension):] == extension
}

// UnsupportedFileError is returned when no analyzer supports the given file
type UnsupportedFileError struct {
	FilePath string
}

func (e *UnsupportedFileError) Error() string {
	return "unsupported file type: " + e.FilePath
}
