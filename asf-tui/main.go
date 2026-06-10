package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	for _, arg := range os.Args[1:] {
		switch arg {
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
			fmt.Println("Place your license key in ~/.asf/license.key")
			os.Exit(1)
		}
	}

	if strings.HasPrefix(os.Args[0], "go") || len(os.Args) > 1 {
		fmt.Printf("ASF v%s — Architecture Security Framework\n", ASFVersion)
		fmt.Println("Run without arguments to open the TUI.")
		if len(os.Args) > 1 {
			fmt.Println()
			fmt.Println("Flags:")
			fmt.Println("  --version, -v    Show version")
			fmt.Println("  --license        Show license status")
		}
		os.Exit(0)
	}

	cfg, err := LoadConfig(ConfigPath())
	if err != nil {
		def := DefaultConfig()
		cfg = &def
	}

	m := newMainModel(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}

	if cfg != nil {
		cfg.Save(ConfigPath())
	}
}
