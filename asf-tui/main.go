package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func printUsage() {
	fmt.Printf("ASF v%s — Architecture Security Framework\n", ASFVersion)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  asf                        Launch the TUI")
	fmt.Println("  asf --version, -v          Show version")
	fmt.Println("  asf --license              Show license status")
	fmt.Println("  asf analyze <file>         Run native analysis (JSON output)")
	fmt.Println("  asf analyze <file> -e <ev> ...   With evidence files/dirs")
	fmt.Println("  asf analyze <file> --graph Include graph in JSON output")
	fmt.Println("  asf doctor                 Run system diagnostics")
	fmt.Println("  asf doctor --verbose       Detailed diagnostics")
	fmt.Println("  asf doctor --fix           Clean stale binaries")
	fmt.Println("  asf --help, -h             Show this help")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Printf("  Config:  %s\n", asfConfigPath())
	fmt.Printf("  Cache:   %s\n", asfCacheDir())
	fmt.Printf("  License: %s\n", asfLicensePath())
	fmt.Println()
	fmt.Println("Documentation: https://github.com/moksh5936-2/asfassumption")
}

func main() {
	args := os.Args[1:]

	if len(args) > 0 {
		switch args[0] {
		case "--version", "-v":
			fmt.Printf("ASF v%s\n", ASFVersion)
			os.Exit(0)
		case "--license":
			l := LoadLicense()
			if l != nil && l.Valid {
				fmt.Printf("License: %s\n", l.Message)
				os.Exit(0)
			}
			fmt.Println("No valid license found.")
			fmt.Printf("Place your license key in %s\n", asfLicensePath())
			os.Exit(1)
		case "doctor", "--doctor", "diagnose":
			verbose := false
			fix := false
			for _, a := range args[1:] {
				switch a {
				case "--verbose", "-v":
					verbose = true
				case "--fix":
					fix = true
				}
			}
			if fix {
				runDoctorFix()
			} else {
				runDoctor(verbose)
			}
			os.Exit(0)
		case "analyze":
			runAnalyzeCLI(args[1:])
			os.Exit(0)
		case "--help", "-h":
			printUsage()
			os.Exit(0)
		case "doctor--verbose":
			runDoctor(true)
			os.Exit(0)
		}
	}

	helpFlags := map[string]bool{"--help": true, "-h": true}
	if len(args) > 1 && !helpFlags[args[1]] {
		fmt.Printf("ASF v%s — Architecture Security Framework\n", ASFVersion)
		fmt.Println()
		printUsage()
		os.Exit(0)
	}

	cfg, err := LoadConfig(asfConfigPath())
	if err != nil {
		def := DefaultConfig()
		cfg = &def
	}

	_ = ensureRuntimeDirs()

	m := newMainModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}

	if cfg != nil {
		cfg.Save(asfConfigPath())
	}
}
