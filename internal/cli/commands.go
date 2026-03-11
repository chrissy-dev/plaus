package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/chrissy-dev/plaus/internal/config"
)

func prompt(label, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", label, defaultVal)
	} else {
		fmt.Printf("%s: ", label)
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	val := strings.TrimSpace(scanner.Text())
	if val == "" {
		return defaultVal
	}
	return val
}

func Init() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	baseURL := prompt("Plausible base URL", cfg.BaseURL)
	cfg.BaseURL = baseURL

	if err := config.Save(cfg); err != nil {
		return err
	}
	fmt.Println("Config saved.")
	return nil
}

func AddSite(domain string) error {
	token := prompt("API token for "+domain, "")
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}
	if err := config.AddSite(domain, token); err != nil {
		return err
	}
	fmt.Printf("Site %q added.\n", domain)
	return nil
}

func ListSites() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if len(cfg.Sites) == 0 {
		fmt.Println("No sites configured. Run: plaus add <site>")
		return nil
	}
	for domain := range cfg.Sites {
		marker := ""
		if domain == cfg.DefaultSite {
			marker = " (default)"
		}
		fmt.Printf("  %s%s\n", domain, marker)
	}
	return nil
}

func RemoveSite(domain string) error {
	if err := config.RemoveSite(domain); err != nil {
		return err
	}
	fmt.Printf("Site %q removed.\n", domain)
	return nil
}
