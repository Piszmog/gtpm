package run

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/Piszmog/gtpm/tmux"
)

func Clean(ctx context.Context, logger *slog.Logger) error {
	logger.Debug("cleaning plugins")
	confPath, err := tmux.GetConfigFilePath(logger)
	if err != nil {
		return err
	}
	logger.Debug("found tmux conf file", "path", confPath)

	rootDir := filepath.Dir(confPath)
	pluginsPath := filepath.Join(rootDir, "plugins")
	if _, err = os.Stat(pluginsPath); err != nil {
		if os.IsNotExist(err) {
			logger.Debug("plugins directory does not exist, nothing to do", "path", pluginsPath)
			return nil
		} else {
			return fmt.Errorf("failed to check if plugins directory exists: %w", err)
		}
	}

	logger.Debug("finding plugins from conf file", "path", confPath)
	plugins, err := tmux.GetPlugins(confPath)
	if err != nil {
		return err
	}
	logger.Debug("found current plugins", "path", confPath, "plugins", plugins)

	if len(plugins) == 0 {
		logger.Debug("conf has no plugins so removing plugins directory", "dir", pluginsPath)
		if err = os.RemoveAll(pluginsPath); err != nil {
			return fmt.Errorf("failed to remove plugins directory: %w", err)
		}
		return nil
	}

	files, err := os.ReadDir(pluginsPath)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			found := false
			for _, p := range plugins {
				if strings.Contains(p, file.Name()) {
					found = true
					break
				}
			}
			if !found {
				path := filepath.Join(pluginsPath, file.Name())
				if err = os.RemoveAll(path); err != nil {
					return fmt.Errorf("failed to clean plugin "+file.Name()+": %w", err)
				}
			}
		}
	}
	logger.Debug("finished cleaning plugins")

	return nil
}
