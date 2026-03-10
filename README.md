# plaus

A terminal analytics viewer for [Plausible Analytics](https://plausible.io).

## Build

```
go build -o plaus ./cmd/plaus
```

## Usage

### Initialize config

```
plaus init
```

Creates `~/.config/plaus/config.json` and prompts for your Plausible base URL.

### Add a site

```
plaus add example.com
```

Prompts for your API token and stores it in the config.

### List sites

```
plaus sites
```

### Remove a site

```
plaus remove example.com
```

### Launch dashboard

```
plaus example.com
```

Opens the TUI dashboard for the given site.
