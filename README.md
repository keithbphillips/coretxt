# coretxt

A neon console text editor for writing novels and long-form prose. Runs in the terminal on Linux, macOS, and Windows. Built with [Bubbletea](https://github.com/charmbracelet/bubbletea).

![coretxt screenshot](Screenshot.png)

## Features

- Distraction-free full-screen editing
- Four neon color themes: CYBERPUNK, SYNTHWAVE, MATRIX, NEON AMBER
- Spell check via `aspell`
- File browser for opening existing files
- Mouse support

## Install

### Pre-built binaries

Download the latest release for your platform from the [Releases page](https://github.com/keithbphillips/coretxt/releases).

### Install with Go

```sh
go install coretxt@latest
```

### Build from source

**Prerequisites:** [Go 1.24+](https://go.dev/dl/)

```sh
git clone https://github.com/keithbphillips/coretxt
cd coretxt
```

**Linux / macOS**
```sh
go build -o coretxt .
```

**Windows (Command Prompt)**
```cmd
go build -o coretxt.exe .
```

**Windows (PowerShell)**
```powershell
go build -o coretxt.exe .
```

## Usage

**Linux / macOS**
```sh
./coretxt [file]
```

**Windows**
```cmd
coretxt.exe [file]
```

Open an existing file or start a new one. If no filename is given, you'll be prompted when saving.

## Keybindings

| Key | Action |
|-----|--------|
| `Ctrl+S` | Save file |
| `Ctrl+Q` | Quit (confirms if unsaved) |
| `Ctrl+C` | Force quit |
| `Ctrl+A` / `Ctrl+E` | Start / end of paragraph |
| `Ctrl+← / →` | Jump word |
| `Ctrl+Home/End` | Beginning / end of document |
| `PgUp / PgDn` | Scroll page |
| `Enter` | New line |
| `Backspace` | Delete back |
| `Ctrl+W` | Delete word back |
| `Ctrl+K` | Delete to end of line |
| `F1` | Toggle help |
| `F2` | Cycle color theme |
| `F7` / `Ctrl+Space` | Spell check word at cursor |

## Dependencies

- Go 1.24+
- `aspell` (optional, for spell check)
  - **Linux:** install via your package manager, e.g. `sudo apt install aspell` or `sudo dnf install aspell`
  - **macOS:** `brew install aspell`
  - **Windows:** download the installer from [Aspell for Windows (GNU Win32)](https://gnuwin32.sourceforge.net/packages/aspell.htm) and ensure `aspell.exe` is on your `PATH`
