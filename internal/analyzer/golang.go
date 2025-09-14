package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/abc-metrics/abc/internal/metrics"
)

// GoAnalyzer implements the Analyzer interface for Go code
type GoAnalyzer struct{}

// NewGoAnalyzer creates a new Go analyzer
func NewGoAnalyzer() *GoAnalyzer {
	return &GoAnalyzer{}
}

// SupportedExtensions returns the list of file extensions supported by this analyzer
func (a *GoAnalyzer) SupportedExtensions() []string {
	return []string{".go"}
}

// AnalyzeFile analyzes a Go file and returns ABC metrics
func (a *GoAnalyzer) AnalyzeFile(filePath string) (metrics.ABCMetrics, error) {
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return metrics.ABCMetrics{}, fmt.Errorf("error reading file: %w", err)
	}

	// Parse the file
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, content, 0)
	if err != nil {
		return metrics.ABCMetrics{}, fmt.Errorf("error parsing file: %w", err)
	}

	// Analyze the AST
	v := &goVisitor{
		metrics: metrics.ABCMetrics{
			AssignmentList: []metrics.MetricDetail{},
			BranchList:     []metrics.MetricDetail{},
			ConditionList:  []metrics.MetricDetail{},
		},
		fset: fset,
	}
	ast.Walk(v, f)

	return v.metrics, nil
}

// goVisitor implements the ast.Visitor interface for Go AST traversal
type goVisitor struct {
	metrics metrics.ABCMetrics
	fset    *token.FileSet
}

// Visit implements the ast.Visitor interface
func (v *goVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	// Assignments
	case *ast.AssignStmt:
		count := len(n.Lhs)
		v.metrics.Assignments += count

		pos := v.fset.Position(n.Pos())
		varNames := make([]string, 0, len(n.Lhs))
		for _, expr := range n.Lhs {
			if ident, ok := expr.(*ast.Ident); ok {
				varNames = append(varNames, ident.Name)
			} else {
				varNames = append(varNames, "expr")
			}
		}

		v.metrics.AssignmentList = append(v.metrics.AssignmentList, metrics.MetricDetail{
			Line:    pos.Line,
			Col:     pos.Column,
			Text:    strings.Join(varNames, ", "),
			Context: fmt.Sprintf("Assignment (%d variables)", count),
		})

	// Branches (function calls)
	case *ast.CallExpr:
		v.metrics.Branches++

		pos := v.fset.Position(n.Pos())
		funcName := "unknown"

		switch fn := n.Fun.(type) {
		case *ast.Ident:
			funcName = fn.Name
		case *ast.SelectorExpr:
			if x, ok := fn.X.(*ast.Ident); ok {
				funcName = x.Name + "." + fn.Sel.Name
			} else {
				funcName = fn.Sel.Name
			}
		}

		v.metrics.BranchList = append(v.metrics.BranchList, metrics.MetricDetail{
			Line:    pos.Line,
			Col:     pos.Column,
			Text:    funcName,
			Context: "Function call",
		})

	// Conditions
	case *ast.IfStmt:
		v.metrics.Conditions++
		pos := v.fset.Position(n.Pos())
		v.metrics.ConditionList = append(v.metrics.ConditionList, metrics.MetricDetail{
			Line:    pos.Line,
			Col:     pos.Column,
			Text:    "if statement",
			Context: "Condition",
		})

	case *ast.ForStmt:
		v.metrics.Conditions++
		pos := v.fset.Position(n.Pos())
		v.metrics.ConditionList = append(v.metrics.ConditionList, metrics.MetricDetail{
			Line:    pos.Line,
			Col:     pos.Column,
			Text:    "for loop",
			Context: "Condition",
		})

	case *ast.RangeStmt:
		v.metrics.Conditions++
		pos := v.fset.Position(n.Pos())
		v.metrics.ConditionList = append(v.metrics.ConditionList, metrics.MetricDetail{
			Line:    pos.Line,
			Col:     pos.Column,
			Text:    "for range loop",
			Context: "Condition",
		})

	case *ast.SwitchStmt:
		v.metrics.Conditions++
		pos := v.fset.Position(n.Pos())
		v.metrics.ConditionList = append(v.metrics.ConditionList, metrics.MetricDetail{
			Line:    pos.Line,
			Col:     pos.Column,
			Text:    "switch statement",
			Context: "Condition",
		})

	case *ast.TypeSwitchStmt:
		v.metrics.Conditions++
		pos := v.fset.Position(n.Pos())
		v.metrics.ConditionList = append(v.metrics.ConditionList, metrics.MetricDetail{
			Line:    pos.Line,
			Col:     pos.Column,
			Text:    "type switch",
			Context: "Condition",
		})

	case *ast.SelectStmt:
		v.metrics.Conditions++
		pos := v.fset.Position(n.Pos())
		v.metrics.ConditionList = append(v.metrics.ConditionList, metrics.MetricDetail{
			Line:    pos.Line,
			Col:     pos.Column,
			Text:    "select statement",
			Context: "Condition",
		})

	case *ast.CaseClause:
		if n.List != nil { // Skip default case
			v.metrics.Conditions++
			pos := v.fset.Position(n.Pos())
			v.metrics.ConditionList = append(v.metrics.ConditionList, metrics.MetricDetail{
				Line:    pos.Line,
				Col:     pos.Column,
				Text:    "case clause",
				Context: "Condition",
			})
		}

	case *ast.BinaryExpr:
		// Count logical operators as conditions
		if n.Op == token.LAND || n.Op == token.LOR {
			v.metrics.Conditions++
			pos := v.fset.Position(n.Pos())
			opText := "&&"
			if n.Op == token.LOR {
				opText = "||"
			}
			v.metrics.ConditionList = append(v.metrics.ConditionList, metrics.MetricDetail{
				Line:    pos.Line,
				Col:     pos.Column,
				Text:    opText,
				Context: "Logical operator",
			})
		}
	}

	return v
}
