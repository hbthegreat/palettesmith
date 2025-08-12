package config

import (
    "os"
    "os/user"
    "path/filepath"
    "strings"
)

// ExpandPath expands a path that may start with ~ to the current user's home,
// and also expands any ${VAR} or $VAR environment variables.
func ExpandPath(p string) string {
    if p == "" {
        return p
    }
    // env var expansion
    expanded := os.ExpandEnv(p)
    // tilde expansion
    if strings.HasPrefix(expanded, "~") {
        // support ~ and ~/...
        usr, err := user.Current()
        if err == nil && usr.HomeDir != "" {
            if expanded == "~" {
                expanded = usr.HomeDir
            } else if strings.HasPrefix(expanded, "~/") {
                expanded = filepath.Join(usr.HomeDir, strings.TrimPrefix(expanded, "~/"))
            }
        }
    }
    return expanded
}

// ConfigDir returns ~/.config/palettesmith by default (XDG base dir aware).
func ConfigDir() string {
    if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
        return filepath.Join(xdg, "palettesmith")
    }
    home, _ := os.UserHomeDir()
    return filepath.Join(home, ".config", "palettesmith")
}

// ExternalPluginsDir returns the user external plugins dir.
func ExternalPluginsDir() string {
    return filepath.Join(ConfigDir(), "plugins")
}

