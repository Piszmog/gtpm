package run

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Piszmog/gtpm/tmux"
)

func Source(ctx context.Context, logger *slog.Logger) error {
	if err := tmux.SetDefaultTPMPath(logger); err != nil {
		return err
	}

	confPath, err := tmux.GetConfigFilePath(logger)
	if err != nil {
		return err
	}
	logger.Debug("found tmux conf file", "path", confPath)

	logger.Debug("binding keys")
	if err = bindKeys(logger); err != nil {
		return err
	}

	logger.Debug("finding plugins from conf file", "path", confPath)
	plugins, err := tmux.GetPlugins(confPath)
	if err != nil {
		return err
	}

	if len(plugins) == 0 {
		logger.Debug("there are no plugins to source", "path", confPath)
		return nil
	} else {
		logger.Debug("there are plugins to source", "plugins", plugins, "path", confPath)
	}

	rootPath := filepath.Dir(confPath)
	pluginsRootPath := filepath.Join(rootPath, "plugins")

	var pluginPaths []string
	for _, p := range plugins {
		plugin, err := tmux.ParsePlugin(p)
		if err != nil {
			return err
		}
		pluginPaths = append(pluginPaths, filepath.Join(pluginsRootPath, plugin.Repo))
	}

	for _, p := range pluginPaths {
		if info, err := os.Stat(p); err != nil {
			if os.IsNotExist(err) {
				return errors.New("plugin " + info.Name() + " has not been installed")
			}
		}
	}

	for _, p := range pluginPaths {
		if strings.HasSuffix(p, "tpm") {
			logger.Debug("skipping tpm plugin")
			continue
		}
		files, err := os.ReadDir(p)
		if err != nil {
			return fmt.Errorf("failed to read directory: %w", err)
		}
		var executableName string
		for _, file := range files {
			if !file.IsDir() {
				if filepath.Ext(file.Name()) == ".tmux" {
					if executableName != "" {
						return errors.New("there are multiple *.tmux files in " + p + " do not know which to source")
					}
					executableName = file.Name()
				}
			}
		}

		if executableName == "" {
			return errors.New("failed to find *.tmux file in " + p)
		}

		cmd := exec.Command("./" + executableName)
		cmd.Dir = p

		out, err := cmd.CombinedOutput()
		logger.Debug("attempted to source plugin", "plugin", p, "output", out)
		if err != nil {
			return fmt.Errorf("failed to source plugin "+p+": %w", err)
		}
	}
	logger.Debug("completed sourcing plugins")

	return nil
}

func bindKeys(logger *slog.Logger) error {
	installKey, err := tmux.GetOption("@tpm-intall")
	if err != nil {
		return err
	}
	if installKey == "" {
		installKey = "I"
	}
	logger.Debug("binding install key", "key", installKey)
	if err = tmux.BindKey(installKey, "gtpm i"); err != nil {
		return err
	}

	updateKey, err := tmux.GetOption("@tpm-update")
	if err != nil {
		return err
	}
	if updateKey == "" {
		updateKey = "U"
	}
	logger.Debug("binding update key", "key", updateKey)
	if err = tmux.BindKey(updateKey, "gtpm u"); err != nil {
		return err
	}

	cleanKey, err := tmux.GetOption("@tpm-clean")
	if err != nil {
		return err
	}
	if cleanKey == "" {
		cleanKey = "M-u"
	}
	logger.Debug("binding clean key", "key", cleanKey)
	if err = tmux.BindKey(cleanKey, "gtpm c"); err != nil {
		return err
	}

	return nil
}
