package commands

import (
	"fmt"
	"os"

	"github.com/abc-metrics/abc/internal/analyzer"
	"github.com/abc-metrics/abc/internal/metrics"
	"github.com/spf13/cobra"
)

var (
	// RootCmd is the entry point for the ABC metrics CLI
	RootCmd = &cobra.Command{
		Use:   "abc",
		Short: "ABC is a tool for measuring Assignment, Branch, Condition metrics in code",
		Long: `ABC metrics tool analyzes source code and calculates the Assignment, Branch, Condition
complexity metrics. These metrics provide an indication of code complexity
based on the number of assignments, branches, and conditions in the code.

The ABC score is calculated as sqrt(A² + B² + C²) where:
- A: number of assignments
- B: number of branches (function calls, method calls)
- C: number of conditions (if, else, switch, case, for, while, etc.)`,
		Run: func(cmd *cobra.Command, args []string) {
			// If no subcommand is provided, print help
			cmd.Help()
		},
	}

	// Flags
	verbose     bool
	filePath    string
	showDetails bool
)

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	RootCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "Path to the file for analysis")
	RootCmd.PersistentFlags().BoolVar(&showDetails, "show", false, "Show detailed list of assignments, branches, and conditions")

	// Add the analyze command
	RootCmd.AddCommand(analyzeCmd)
}

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze a file for ABC metrics",
	Long:  `Analyze a single file and calculate its ABC metrics.`,
	Run: func(cmd *cobra.Command, args []string) {
		if filePath == "" {
			if len(args) == 0 {
				fmt.Println("Error: file path is required")
				cmd.Help()
				os.Exit(1)
			}
			filePath = args[0]
		}

		fmt.Printf("Analyzing file: %s\n", filePath)

		// Get analyzer for file
		analyzer, err := analyzer.GetAnalyzerForFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Analyze file
		abcMetrics, err := analyzer.AnalyzeFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error analyzing file: %v\n", err)
			os.Exit(1)
		}

		// Print results
		fmt.Println(abcMetrics.String())
		fmt.Printf("Complexity: %s\n", metrics.SeverityLevel(abcMetrics.Score()))

		// If show details flag is set, print detailed metrics
		if showDetails {
			fmt.Println("\nAssignments:")
			for i, assignment := range abcMetrics.AssignmentList {
				fmt.Printf("  %d. Line %d: %s (%s)\n", i+1, assignment.Line, assignment.Text, assignment.Context)
			}

			fmt.Println("\nBranches:")
			for i, branch := range abcMetrics.BranchList {
				fmt.Printf("  %d. Line %d: %s (%s)\n", i+1, branch.Line, branch.Text, branch.Context)
			}

			fmt.Println("\nConditions:")
			for i, condition := range abcMetrics.ConditionList {
				fmt.Printf("  %d. Line %d: %s (%s)\n", i+1, condition.Line, condition.Text, condition.Context)
			}
		}
	},
}
