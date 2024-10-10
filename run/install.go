package run

import (
	"context"
	"errors"
	"log/slog"
	"path/filepath"

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

	plugins, err := tmux.GetPlugins(confPath)
	if err != nil {
		return err
	}
	logger.Debug("plugins to install", "plugins", plugins)

	if len(plugins) == 0 {
		logger.Debug("nothing to install")
		return nil
	}

	return nil
}
