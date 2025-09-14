package metrics

import (
	"fmt"
	"math"
)

// MetricDetail represents a single item that contributes to a metric
type MetricDetail struct {
	Line    int    // Line number
	Col     int    // Column number
	Text    string // Short description or snippet
	Context string // Additional context
}

// ABCMetrics represents the Assignment, Branch, and Condition metrics
type ABCMetrics struct {
	Assignments    int            // Number of assignments
	Branches       int            // Number of branches (function calls, method calls)
	Conditions     int            // Number of conditions (if, else, switch, case, for, while, etc.)
	AssignmentList []MetricDetail // Details of assignments
	BranchList     []MetricDetail // Details of branches
	ConditionList  []MetricDetail // Details of conditions
}

// Score calculates the ABC score as sqrt(A² + B² + C²)
func (m ABCMetrics) Score() float64 {
	return math.Sqrt(float64(m.Assignments*m.Assignments + m.Branches*m.Branches + m.Conditions*m.Conditions))
}

// String returns a string representation of the ABC metrics
func (m ABCMetrics) String() string {
	return fmt.Sprintf("ABC: %.2f (A=%d, B=%d, C=%d)",
		m.Score(), m.Assignments, m.Branches, m.Conditions)
}

// CombineMetrics combines multiple ABCMetrics into a single metric
func CombineMetrics(metrics ...ABCMetrics) ABCMetrics {
	combined := ABCMetrics{}
	for _, m := range metrics {
		combined.Assignments += m.Assignments
		combined.Branches += m.Branches
		combined.Conditions += m.Conditions

		// Combine detail lists
		combined.AssignmentList = append(combined.AssignmentList, m.AssignmentList...)
		combined.BranchList = append(combined.BranchList, m.BranchList...)
		combined.ConditionList = append(combined.ConditionList, m.ConditionList...)
	}
	return combined
}

// SeverityLevel returns a human-readable severity level based on ABC score
func SeverityLevel(score float64) string {
	switch {
	case score < 10:
		return "Low"
	case score < 20:
		return "Medium"
	case score < 40:
		return "High"
	default:
		return "Very High"
	}
}
