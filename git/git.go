package git

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

func IsInstalled(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "git", "--version")
	_, err := cmd.Output()
	return err == nil
}

func Clone(logger *slog.Logger, url string, branch string, targetDir string) error {
	var cmd *exec.Cmd
	if branch != "" {
		cmd = exec.Command("git", "clone", "-b", branch, "--single-branch", "--recursive", url, targetDir)
	} else {
		cmd = exec.Command("git", "clone", "--single-branch", "--recursive", url, targetDir)
	}

	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")

	logger.Debug("cloning repo", "url", url)
	out, err := cmd.CombinedOutput()
	logger.Debug("output of git clone", "url", url, "branch", branch, "output", out)
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	return nil
}

func Update(logger *slog.Logger, path string) error {
	logger.Debug("pulling repository", "path", path)
	if err := pull(path); err != nil {
		return err
	}
	logger.Debug("updating submodule", "path", path)
	if err := update(path); err != nil {
		return err
	}
	return nil
}

func pull(path string) error {
	cmd := exec.Command("git", "-C", path, "pull")
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			switch exitError.ExitCode() {
			case 1:
				return fmt.Errorf("remote repository not found")
			case 128:
				return fmt.Errorf("there is a conflict between remote and local changes")
			default:
				return fmt.Errorf("failed to pull latest changes: %s", exitError.Error())
			}
		}
	}
	return nil
}

func update(path string) error {
	cmd := exec.Command("git", "-C", path, "submodule", "update", "--init", "--recursive")
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to perform submodule update: %w", err)
	}
	return nil
}
