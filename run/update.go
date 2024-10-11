package run

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/Piszmog/gtpm/git"
	"github.com/Piszmog/gtpm/tmux"
)

func Update(ctx context.Context, logger *slog.Logger, plugins []string) error {
	if len(plugins) == 0 {
		logger.Debug("updating all plugins")
	} else {
		logger.Debug("updating plugins", "plugins", plugins)
	}

	logger.Debug("checking if git is install")
	if !git.IsInstalled(ctx) {
		return errors.New("git is required to install plugins")
	}

	confPath, err := tmux.GetConfigFilePath(logger)
	if err != nil {
		return err
	}
	logger.Debug("found tmux conf file", "path", confPath)

	existingPlugins, err := tmux.GetPlugins(confPath)
	if err != nil {
		return err
	}
	logger.Debug("configured plugins", "plugins", existingPlugins)

	if len(plugins) > 0 {
		for _, p := range plugins {
			for _, configuredPlugin := range existingPlugins {
				if !strings.Contains(configuredPlugin, p) {
					return errors.New("cannot update plugin " + p + " is it not configured in your tmux conf file")
				}
			}
		}
	}

	rootPath := filepath.Dir(confPath)
	pluginsPath := filepath.Join(rootPath, "plugins")
	files, err := os.ReadDir(pluginsPath)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %w", err)
	}

	pluginsToBeUpdated := plugins
	if len(pluginsToBeUpdated) == 0 {
		pluginsToBeUpdated = existingPlugins
	}

	var pluginPathsToUpdate []string
	for _, p := range plugins {
		for _, configuredPlugin := range existingPlugins {
			if strings.Contains(configuredPlugin, p) {
				plugin, err := tmux.ParsePlugin(configuredPlugin)
				if err != nil {
					return err
				}
				pluginPathsToUpdate = append(pluginPathsToUpdate, filepath.Join(pluginsPath, plugin.Repo))
			}
		}
	}

	logger.Debug("checking if any plugins are not installed")
	err = checkInstalledPlugins(plugins, files)
	if err != nil {
		return err
	}

	for _, p := range pluginPathsToUpdate {
		if err = git.Update(logger, p); err != nil {
			return err
		}
	}

	return nil
}

func checkInstalledPlugins(plugins []string, files []fs.DirEntry) error {
	for _, p := range plugins {
		pluginInstalled := false
		for _, file := range files {
			if file.IsDir() && strings.Contains(p, file.Name()) {
				pluginInstalled = true
				break
			}
		}
		if !pluginInstalled {
			return errors.New("plugin " + p + " is not installed")
		}
	}
	return nil
}
