package git

import (
	"context"
	"log/slog"
	"os/exec"
)

func IsInstalled(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "git", "--version")
	_, err := cmd.Output()
	return err == nil
}

func Clone(logger *slog.Logger, owner string, repo string) error {
	return nil
}
