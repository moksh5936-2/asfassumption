# Python Engine Discovery

## Old Behavior

The Go TUI hardcoded the Python engine path:

```go
pythonPath: "/Users/moksh/Project/cybersec/.venv/bin/python",
projectDir: "/Users/moksh/Project/cybersec",
```

This meant:
- The Python binary was locked to a specific developer's virtualenv
- The working directory was locked to a specific developer's project folder
- ASF would immediately fail with "no such file or directory" on any other machine
- The error message was confusing (`chdir /Users/moksh/Project/cybersec: no such file or directory`)

## New Behavior

`asf-tui/engine.go:discoverPythonPath()` performs runtime discovery:

### Search Order

1. **Config override** — If `config.yaml` has `engine.python_path` set, use it
2. **Co-located engine** — Look for `{binary_dir}/engine/bin/python3` or `{binary_dir}/asf`
3. **Data directory** — Look for `{XDG_DATA_HOME}/asf/venv/bin/python3`
4. **PATH search** — Look for `asf`, `asf.py`, `asf-cli` in PATH
5. **System Python** — Use `python3` or `python` from PATH (and try `python3 -m asf.cli.main`)
6. **Fallback** — Use `"python3"` (will fail with helpful error)

### Search Implementation

```go
func discoverPythonPath(cfg *Config) string {
    // 1. Config override
    if cfg.Engine.PythonPath != "" {
        return cfg.Engine.PythonPath
    }
    // 2. Binary-relative
    exe, _ := os.Executable()
    binDir := filepath.Dir(exe)
    candidates := []string{
        filepath.Join(binDir, "engine", "bin", "python3"),
        filepath.Join(binDir, "asf"),
    }
    // 3. Data dir
    candidates = append(candidates, 
        filepath.Join(asfDataDir(), "venv", "bin", "python3"))
    // 4-5. PATH search
    candidates = append(candidates, lookPath("asf"), "python3")
    // Return first match
    for _, p := range candidates {
        if info, err := os.Stat(p); err == nil && !info.IsDir() {
            return p
        }
    }
    return "python3"
}
```

### Working Directory

The Python CLI runs with `cmd.Dir = asfCacheDir()` (XDG cache directory) instead of a hardcoded project dir.

## Failure Modes

| Scenario | Error Message | Recovery |
|----------|--------------|----------|
| Python not installed | `ASF engine error: exec: "python3": executable file not found` | Install Python 3.8+ |
| ASF package not installed | `ASF engine error: No module named asf.cli.main` | Run `pip install -e /path/to/asf` |
| Both missing | Same as above | `asf doctor` shows Python + engine status |
| Permission denied | `ASF engine error: permission denied` | Check binary permissions |
| Engine configured but wrong | `ASF engine error: exec: ...` | Check `config.yaml engine.python_path` |

## `asf doctor` Integration

The `doctor` command reports:
- `Python binary:` — what Python was found
- `ASF engine:` — whether the ASF Python package is importable

Example:
```
Python Engine
  Python binary:    /usr/bin/python3
  ASF engine:       python3 -m asf.cli.main (/usr/bin/python3)
```
