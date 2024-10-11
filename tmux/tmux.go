package tmux

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
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

func HasPermissions(rootPath string) bool {
	tempFile, err := os.CreateTemp(rootPath, "testFile-*")
	if err != nil {
		return false
	}

	tempFile.Close()
	os.Remove(tempFile.Name())

	return true
}

func ParsePlugin(plugin string) (Plugin, error) {
	parts := strings.SplitN(plugin, "#", 2)
	otherParts := strings.SplitN(parts[0], "/", 2)

	if len(otherParts) != 2 {
		return Plugin{}, errors.New("expected plugin to be in format <owner>/<repo>: " + plugin)
	}

	var branch string
	if len(parts) == 2 {
		branch = parts[1]
	}

	return Plugin{
		Owner:  otherParts[0],
		Repo:   otherParts[1],
		Branch: branch,
	}, nil
}

type Plugin struct {
	Owner  string
	Repo   string
	Branch string
}

func CreatePluginsDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed create plugin directory: %w", err)
	}
	return nil
}

func SetDefaultTPMPath(logger *slog.Logger) error {
	isSet := true
	if err := checkIfDefaultEnvSet(); err != nil {
		if errors.Is(err, errDefaultEnvNotSet) {
			logger.Debug("default env path has not been set")
			isSet = false
		} else {
			return err
		}
	}

	if !isSet {
		conf, err := GetConfigFilePath(logger)
		if err != nil {
			return err
		}
		rootDir := filepath.Dir(conf)

		var tpmPath string
		if strings.Contains(rootDir, ".config") {
			tpmPath = filepath.Join(rootDir, "plugins")
		} else {
			tpmPath = filepath.Join(rootDir, ".tmux", "plugins")
		}

		cmd := exec.Command("tmux", "set-environment", "-g", "$DEFAULT_TPM_ENV_VAR_NAME", tpmPath)

		out, err := cmd.CombinedOutput()
		logger.Debug("attempted to set $DEFAULT_TPM_ENV_VAR_NAME", "out", out)
		if err != nil {
			return fmt.Errorf("failed to set $DEFAULT_TPM_ENV_VAR_NAME: %w", err)
		}
	} else {
		logger.Debug("default tpm path has been set already")
	}
	return nil
}

func checkIfDefaultEnvSet() error {
	cmd := exec.Command("tmux", "show-environment", "-g", "$DEFAULT_TPM_ENV_VAR_NAME")
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			switch exitError.ExitCode() {
			case 1:
				return errDefaultEnvNotSet
			default:
				return fmt.Errorf("failed determine if $DEFAULT_TPM_ENV_VAR_NAME is set: %s", exitError.Error())
			}
		}
	}
	return nil

}

var errDefaultEnvNotSet = errors.New("default env is not set")

func GetOption(option string) (string, error) {
	cmd := exec.Command("tmux", "show-option", "-gqv", option)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			switch exitError.ExitCode() {
			case 1:
				return "", nil
			default:
				return "", fmt.Errorf("failed to find option for "+option+": %w", err)
			}
		}
	}
	return strings.TrimSpace(string(out)), nil
}

func BindKey(key string, command string) error {
	cmd := exec.Command("tmux", "bind-key", key, "run-shell", command)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to bind key "+key+" to "+command+": %w", err)
	}
	return nil
}
