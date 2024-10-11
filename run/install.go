package run

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/Piszmog/gtpm/git"
	"github.com/Piszmog/gtpm/tmux"
)

func Install(ctx context.Context, logger *slog.Logger) error {
	logger.Debug("installing plugins")

	logger.Debug("checking if git is install")
	if !git.IsInstalled(ctx) {
		return errors.New("git is required to install plugins")
	}

	confPath, err := tmux.GetConfigFilePath(logger)
	if err != nil {
		return err
	}
	logger.Debug("found tmux conf file", "path", confPath)

	rootPath := filepath.Dir(confPath)
	if !tmux.HasPermissions(rootPath) {
		return errors.New("do not have write permissions to " + rootPath)
	}

	pluginsPath := filepath.Join(rootPath, "plugins")
	if err = tmux.CreatePluginsDir(pluginsPath); err != nil {
		return err
	}

	plugins, err := tmux.GetPlugins(confPath)
	if err != nil {
		return err
	}
	logger.Debug("plugins to install", "plugins", plugins)

	if len(plugins) == 0 {
		logger.Debug("nothing to install")
		return nil
	}

	files, err := os.ReadDir(pluginsPath)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %w", err)
	}

	pluginsToInstall := []string{}
	for _, p := range plugins {
		pluginInstalled := false
		for _, file := range files {
			if file.IsDir() && strings.Contains(p, file.Name()) {
				logger.Debug("plugin already installed", "plugin", p, "path", pluginsPath)
				pluginInstalled = true
				break
			}
		}
		if !pluginInstalled {
			pluginsToInstall = append(pluginsToInstall, p)
		}
	}

	if len(pluginsToInstall) == 0 {
		logger.Debug("no plugins to install")
		return nil
	}

	logger.Debug("need to install plugins", "plugins", pluginsToInstall)

	for _, p := range pluginsToInstall {
		plugin, err := tmux.ParsePlugin(p)
		if err != nil {
			return err
		}
		logger.Debug("cloning plugin", "plugin", p)
		err = git.Clone(logger, "https://git::@github.com/"+plugin.Owner+"/"+plugin.Repo, plugin.Branch, filepath.Join(pluginsPath, plugin.Repo))
		if err != nil {
			return err
		}
	}
	logger.Debug("completed installing plugins")

	return nil
}
