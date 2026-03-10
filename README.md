# plaus

[![CI](https://github.com/chrissy-dev/plaus/actions/workflows/ci.yml/badge.svg)](https://github.com/chrissy-dev/plaus/actions/workflows/ci.yml)
[![Release](https://github.com/chrissy-dev/plaus/actions/workflows/release.yml/badge.svg)](https://github.com/chrissy-dev/plaus/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/chriswk/plaus)](https://goreportcard.com/report/github.com/chriswk/plaus)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A terminal analytics dashboard for [Plausible Analytics](https://plausible.io). View your website stats without leaving the terminal.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lip Gloss](https://github.com/charmbracelet/lipgloss), and [ntcharts](https://github.com/NimbleMarkets/ntcharts).

## Features

- Top-line metrics: visitors, visits, pageviews, bounce rate, visit duration
- Visitor trend sparkline chart
- Top pages and top sources panels
- Time period switching: today, yesterday, 7 days, 30 days
- Live visitor count with realtime indicator (today view)
- Self-hosted Plausible support

## Install

### From release

Download the latest binary from [Releases](https://github.com/chrissy-dev/plaus/releases).

### From source

```
go install github.com/chriswk/plaus/cmd/plaus@latest
```

### Build locally

```
git clone https://github.com/chrissy-dev/plaus.git
cd plaus
go build -o plaus ./cmd/plaus
```

## Setup

### 1. Initialize config

```
plaus init
```

Creates `~/.config/plaus/config.json` and prompts for your Plausible base URL (defaults to `https://plausible.io`).

### 2. Add a site

```
plaus add example.com
```

Prompts for your [Stats API key](https://plausible.io/docs/stats-api) and stores it in the config.

### 3. Launch dashboard

```
plaus example.com
```

## Keyboard shortcuts

| Key | Action |
|-----|--------|
| `1` | Today |
| `2` | Yesterday |
| `3` | Last 7 days |
| `4` | Last 30 days |
| `r` | Refresh data |
| `q` | Quit |

## Other commands

```
plaus sites             # List configured sites
plaus remove example.com  # Remove a site
plaus version           # Show version
```

## Config

Config is stored at `~/.config/plaus/config.json`:

```json
{
  "base_url": "https://plausible.io",
  "default_site": "example.com",
  "sites": {
    "example.com": {
      "token": "your-api-key"
    }
  }
}
```

## License

MIT
