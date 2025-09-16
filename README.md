# AGENTDL

A TUI for finding and downloading Claude agent and command configurations from GitHub.

<img src="https://vhs.charm.sh/vhs-4654xhf4TTrBrP51jLoo3N.gif" alt="Made with VHS">
<a href="https://vhs.charm.sh">
  <img src="https://stuff.charm.sh/vhs/badge.svg">
</a>

## What it does

Search GitHub for `.claude/agents/*.md` and `.claude/commands/*.md` files and download them. Built to learn Bubble Tea and practice Go.

**New Features:**
- **Toggle between Agents and Commands mode** - Switch between searching `.claude/agents/` and `.claude/commands/` directories
- **Filename-based search** - Find files with keywords in their actual filenames (not just content)
- **Improved search accuracy** - Enhanced filtering and rate limiting for better results

Useful for:
- Finding Claude agent configurations and command files
- Exploring how others structure their .claude directories
- Batch downloading configurations with smart filename filtering

## Installation

```bash
brew install WillyV3/tap/agentdl
```

## Usage

```bash
agentdl
```

### Key Features

- **Search**: Enter keywords to find files with those terms in their filenames
- **Mode Toggle**: Press `tab` to switch between Agents and Commands mode
- **Browse**: Press `v` to browse individual repositories
- **Select**: Use `space` to select/deselect files
- **Download**: Press `enter` to download selected files

### Search Modes

**Agents Mode** (default): Searches `.claude/agents/` directories
- Find agent configurations, assistants, and AI personas
- Example: search "python" to find Python-related agents

**Commands Mode**: Searches `.claude/commands/` directories
- Find command templates, workflows, and automation scripts
- Example: search "hook" to find Git hook commands

### Controls

- `↑/↓` - Navigate results
- `space` - Select/deselect files
- `enter` - Download selected files or view details
- `v` - Browse repository
- `p` - Preview file content
- `tab` - Toggle between Agents/Commands mode
- `q` - Quit

Downloads go to `~/.claude/agents` or `~/.claude/commands` by default.

### Build from source

```bash
go build -o agentdl *.go
./agentdl
```

## Built with

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) by Charm
- [Lipgloss](https://github.com/charmbracelet/lipgloss) for styling
- Go

## Why

Learning project to get familiar with Bubble Tea TUI framework and Go development.

## Credits

Built using the excellent TUI libraries from [Charm](https://github.com/charmbracelet).

Part of [WillyV3](https://github.com/WillyV3) organization.