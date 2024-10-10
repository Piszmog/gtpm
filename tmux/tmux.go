package tmux

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GetConfigFilePath(logger *slog.Logger) (string, error) {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")

	var configPath string
	if xdgConfigHome != "" {
		logger.Debug("XDF_CONFIG_HOME is set")
		configPath = filepath.Join(xdgConfigHome, "tmux", "tmux.conf")
	}

	if configPath == "" {
		logger.Debug("HOME is set")
		home := os.Getenv("HOME")
		configPath = filepath.Join(home, ".config", "tmux", "tmux.conf")
	}

	if _, err := os.Stat(configPath); err == nil {
		return configPath, nil
	} else if os.IsNotExist(err) {
		logger.Debug("conf file not found in initial place, checking .tmux", "path", configPath)
		home := os.Getenv("HOME")
		configPath = filepath.Join(home, ".tmux.conf")

		if _, err = os.Stat(configPath); err == nil {
			return configPath, nil
		} else if os.IsNotExist(err) {
			return "", errors.New("failed to find tmux conf file at $XDG_CONFIG_HOME/tmux/tmux.conf, $HOME/tmux/tmux.conf, and $HOME/.tmux.conf")
		} else {
			return "", fmt.Errorf("failed to check if file exists: %w", err)
		}

	} else {
		return "", fmt.Errorf("failed to check if file exists: %w", err)
	}
}

func GetPlugins(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	var plugins []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix("#", line) {
			continue
		}

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		matches := pluginRegex.FindStringSubmatch(line)
		if len(matches) == 2 {
			plugins = append(plugins, matches[1])
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return plugins, nil
}

var pluginRegex = regexp.MustCompile(`^\s*set\s+-g\s+@plugin\s+['"]([^'"]+)['"]`)
