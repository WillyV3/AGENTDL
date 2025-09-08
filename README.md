# AGENTDL

A TUI for finding and downloading Claude agent configurations from GitHub.

<img src="https://vhs.charm.sh/vhs-4654xhf4TTrBrP51jLoo3N.gif" alt="Made with VHS">
<a href="https://vhs.charm.sh">
  <img src="https://stuff.charm.sh/vhs/badge.svg">
</a>

## What it does

Search GitHub for `.claude/agents/*.md` files and download them. Built to learn Bubble Tea and practice Go.

Useful for:
- Finding new agent .md files for your Claude Code setup
- Exploring how others structure their .claude/agents directories
- Batch downloading agent configurations

## Usage

```bash
go build -o agentdl *.go
./agentdl
```

Search, browse repos, select files with space, download to ~/.claude/agents or wherever.

## Built with

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) by Charm
- [Lipgloss](https://github.com/charmbracelet/lipgloss) for styling
- Go

## Why

Learning project to get familiar with Bubble Tea TUI framework and Go development.

## Credits

Built using the excellent TUI libraries from [Charm](https://github.com/charmbracelet).