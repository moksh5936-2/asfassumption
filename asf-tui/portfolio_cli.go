package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"asf-tui/intelligence"
)

func runPortfolioCLI(args []string) {
	if len(args) == 0 {
		printPortfolioUsage()
		os.Exit(ExitSuccess)
	}

	switch args[0] {
	case "add":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Usage: asf portfolio add <analysis-result.json>\n")
			os.Exit(ExitInvalidCmd)
		}
		for _, path := range args[1:] {
			if err := portfolioAdd(path); err != nil {
				fmt.Fprintf(os.Stderr, "Error adding %s: %v\n", path, err)
				os.Exit(ExitAnalysisErr)
			}
		}
		fmt.Println("Portfolio updated.")
	case "analyze":
		portfolioAnalyze()
	case "report":
		portfolioReport()
	case "list":
		portfolioList()
	case "remove":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Usage: asf portfolio remove <architecture-id>\n")
			os.Exit(ExitInvalidCmd)
		}
		for _, id := range args[1:] {
			if err := portfolioRemove(id); err != nil {
				fmt.Fprintf(os.Stderr, "Error removing %s: %v\n", id, err)
				os.Exit(ExitAnalysisErr)
			}
		}
		fmt.Println("Architecture(s) removed.")
	case "--help", "-h":
		printPortfolioUsage()
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown portfolio command '%s'\n", args[0])
		printPortfolioUsage()
		os.Exit(ExitInvalidCmd)
	}
}

func printPortfolioUsage() {
	fmt.Println("Usage: asf portfolio <command> [options]")
	fmt.Println()
	fmt.Println("Manage and analyze architecture portfolios.")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  add <file>...           Add architecture analysis result(s) to portfolio")
	fmt.Println("  analyze                 Run portfolio analysis")
	fmt.Println("  report                  Generate portfolio report")
	fmt.Println("  list                    List architectures in portfolio")
	fmt.Println("  remove <id>...          Remove architecture(s) from portfolio")
	fmt.Println("  --help, -h              Show this help")
}

func portfolioDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".asf/portfolio"
	}
	return filepath.Join(home, ".asf", "portfolio")
}

func portfolioPath() string {
	return filepath.Join(portfolioDir(), "portfolio.json")
}

func portfolioAdd(resultPath string) error {
	data, err := os.ReadFile(resultPath)
	if err != nil {
		return err
	}
	var result AnalysisResult
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	id := strings.TrimSuffix(filepath.Base(resultPath), filepath.Ext(resultPath))
	rec := intelligence.ArchitectureRecord{
		ArchitectureID: id,
		Name:           result.ArchitectureName,
		Domain:         result.Domain,
		AnalysisDate:   result.AnalysisDate,
		Version:        "1.0",
		RiskScore:      averageRiskScore(result),
		Metadata: map[string]string{
			"source": resultPath,
		},
	}
	rec.Compliance = append(rec.Compliance, result.Compliance...)
	for _, a := range result.Assumptions {
		rec.Assumptions = append(rec.Assumptions, intelligence.Assumption{
			ID: a.ID, Description: a.Description,
			Category: a.Category, Risk: intelligence.RiskLevel(a.Risk),
		})
	}
	for _, t := range result.Threats {
		rec.Threats = append(rec.Threats, intelligence.Threat{
			ID: t.ID, Name: t.Name,
			Severity:    intelligence.RiskLevel(t.Severity),
			Description: t.Description,
		})
	}
	for _, ap := range result.AttackPaths {
		rec.AttackPaths = append(rec.AttackPaths, intelligence.AttackPath{
			ID: ap.ID, Name: ap.Name,
			Description: ap.Description,
		})
	}
	for _, f := range result.SDRIDesignFindings {
		rec.Findings = append(rec.Findings, intelligence.SDRIFinding{
			ID: f.ID, Title: f.Title,
			Category: f.Category, Severity: f.Severity,
			Description: f.Description,
		})
	}
	for _, c := range result.SDRIControls {
		rec.Controls = append(rec.Controls, intelligence.SDRIControl{
			ID: c.ID, Name: c.Name,
			Category: c.Category, ControlType: intelligence.SDRIControlType(c.ControlType),
		})
	}

	dir := portfolioDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	p := intelligence.NewPortfolio()
	if _, err := os.Stat(portfolioPath()); err == nil {
		if err := p.Load(portfolioPath()); err != nil {
			return err
		}
	}
	p.AddArchitecture(rec)
	return p.Save(portfolioPath())
}

func averageRiskScore(result AnalysisResult) float64 {
	total := 0.0
	count := 0
	if result.CriticalCount > 0 {
		total += 9.0 * float64(result.CriticalCount)
		count += result.CriticalCount
	}
	if result.HighCount > 0 {
		total += 6.0 * float64(result.HighCount)
		count += result.HighCount
	}
	if result.MediumCount > 0 {
		total += 3.0 * float64(result.MediumCount)
		count += result.MediumCount
	}
	if result.LowCount > 0 {
		total += 1.0 * float64(result.LowCount)
		count += result.LowCount
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}

func portfolioAnalyze() {
	p, err := loadPortfolio()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(ExitAnalysisErr)
	}
	engine := intelligence.NewSAMPIEngine()
	result := engine.Run(intelligence.SAMPIInput{Portfolio: p})
	convertResult := convertSAMPIResult(result)
	data, _ := json.MarshalIndent(convertResult, "", "  ")
	fmt.Println(string(data))
}

func portfolioReport() {
	p, err := loadPortfolio()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(ExitAnalysisErr)
	}
	engine := intelligence.NewSAMPIEngine()
	result := engine.Run(intelligence.SAMPIInput{Portfolio: p})
	convertResult := convertSAMPIResult(result)

	fmt.Printf("=== Portfolio Report ===\n")
	fmt.Printf("Generated: %s\n", time.Now().Format(time.RFC3339))
	fmt.Printf("Architectures: %d\n", convertResult.Dashboard.TotalArchitectures)
	fmt.Printf("Total Findings: %d\n", convertResult.Dashboard.TotalFindings)
	fmt.Printf("Total Threats: %d\n", convertResult.Dashboard.TotalThreats)
	fmt.Printf("Total Controls: %d\n", convertResult.Dashboard.TotalControls)
	fmt.Printf("Average Risk Score: %.1f\n", convertResult.Dashboard.AverageRiskScore)
	fmt.Printf("Average Coverage: %.1f\n", convertResult.Dashboard.AverageCoverage)
	fmt.Println()

	if len(convertResult.RepeatedWeaknesses) > 0 {
		fmt.Printf("=== Repeated Weaknesses (%d) ===\n", len(convertResult.RepeatedWeaknesses))
		for _, rw := range convertResult.RepeatedWeaknesses {
			mark := " "
			if rw.Systemic {
				mark = "!"
			}
			fmt.Printf("  %s %s [%s] appears in %d architectures\n", mark, rw.FindingTitle, rw.Category, rw.OccurrenceCount)
		}
		fmt.Println()
	}

	if len(convertResult.EnterpriseThemes) > 0 {
		fmt.Printf("=== Enterprise Risk Themes ===\n")
		for _, th := range convertResult.EnterpriseThemes {
			fmt.Printf("  %s (%d) — %s\n", th.Name, th.RiskCount, th.Severity)
		}
		fmt.Println()
	}

	if len(convertResult.ControlCoverage) > 0 {
		fmt.Printf("=== Control Coverage (lowest first) ===\n")
		for _, cc := range convertResult.ControlCoverage {
			mark := " "
			if cc.CoveragePercent < 50 {
				mark = "!"
			}
			fmt.Printf("  %s %s: %.1f%% (%d/%d architectures)\n", mark, cc.ControlName, cc.CoveragePercent, cc.ArchitecturesWith, cc.ArchitecturesTotal)
		}
		fmt.Println()
	}

	if len(convertResult.Heatmaps) > 0 {
		fmt.Printf("=== Executive Heatmap ===\n")
		for _, h := range convertResult.Heatmaps {
			fmt.Printf("  %s [%s] score=%.1f findings=%d controls=%d\n", h.ArchitectureName, h.RiskBand, h.RiskScore, h.FindingCount, h.ControlCount)
		}
		fmt.Println()
	}

	if convertResult.SecurityDebt.Score > 0 {
		fmt.Printf("=== Security Debt ===\n")
		fmt.Printf("  Score: %.1f\n", convertResult.SecurityDebt.Score)
		fmt.Printf("  Longstanding: %d\n", convertResult.SecurityDebt.LongstandingCount)
		fmt.Printf("  Repeated: %d\n", convertResult.SecurityDebt.RepeatedCount)
		fmt.Println()
	}

	if len(convertResult.ProgramInsights) > 0 {
		fmt.Printf("=== Security Program Insights ===\n")
		for _, pi := range convertResult.ProgramInsights {
			fmt.Printf("  [%s] %s: %s\n", pi.Priority, pi.Area, pi.Insight)
		}
		fmt.Println()
	}
}

func portfolioList() {
	p, err := loadPortfolio()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(ExitAnalysisErr)
	}
	fmt.Printf("Portfolio (%d architectures):\n", len(p.Architectures))
	for _, a := range p.Architectures {
		fmt.Printf("  %s (%s) — %s, risk=%.1f, findings=%d\n",
			a.ArchitectureID, a.Name, a.Domain, a.RiskScore, len(a.Findings))
	}
}

func portfolioRemove(id string) error {
	p, err := loadPortfolio()
	if err != nil {
		return err
	}
	p.RemoveArchitecture(id)
	return p.Save(portfolioPath())
}

func loadPortfolio() (*intelligence.Portfolio, error) {
	path := portfolioPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("no portfolio found at %s — use 'asf portfolio add' first", path)
	}
	p := intelligence.NewPortfolio()
	if err := p.Load(path); err != nil {
		return nil, err
	}
	if len(p.Architectures) == 0 {
		return nil, fmt.Errorf("portfolio is empty")
	}
	return p, nil
}
