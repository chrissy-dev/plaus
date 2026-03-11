package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/chrissy-dev/plaus/internal/api"
	"github.com/chrissy-dev/plaus/internal/app"
	"github.com/chrissy-dev/plaus/internal/cli"
	"github.com/chrissy-dev/plaus/internal/config"
)

var version = "dev"

func usage() {
	fmt.Fprintf(os.Stderr, `plaus %s — terminal analytics for Plausible

Usage:
  plaus init            Initialize config
  plaus add <site>      Add a site
  plaus sites           List configured sites
  plaus remove <site>   Remove a site
  plaus <site>          Launch dashboard
  plaus version         Show version
`, version)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "init":
		err = cli.Init()
	case "add":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: plaus add <site>")
			os.Exit(1)
		}
		err = cli.AddSite(os.Args[2])
	case "sites":
		err = cli.ListSites()
	case "remove":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: plaus remove <site>")
			os.Exit(1)
		}
		err = cli.RemoveSite(os.Args[2])
	case "version", "-v", "--version":
		fmt.Println("plaus", version)
	case "help", "-h", "--help":
		usage()
	default:
		err = launch(os.Args[1])
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func launch(domain string) error {
	site, ok := config.GetSite(domain)
	if !ok {
		fmt.Fprintf(os.Stderr, "Site %q not configured.\nRun: plaus add %s\n", domain, domain)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, domain, site.Token)
	graphType := app.GraphTypeFromString(cfg.GraphType)
	period := app.PeriodFromString(cfg.Period)
	m := app.New(domain, client, graphType, period)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}
