// cspell:ignore appimage
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	ci "github.com/megalomania428/go-lib-ci"
)

const roleDir = "ansible-mega-base"

// isCI reports whether the current invocation runs on a CI runner, where
// paths must be unique per job to allow multiple runners on one host.
// Manual/local runs keep the historical shared paths to avoid re-fetching
// galaxy dependencies on every invocation.
func isCI() bool {
	return os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != ""
}

// resolveRepoRoot returns the ansible-mega-base repository root, assuming
// this program runs via `go -C golang run ./cmd/molecule`: the working
// directory is golang/ (the -C target), so the repo root is its parent,
// mirroring manage-dragon's resolveAnsibleDir convention.
func resolveRepoRoot() (string, error) {
	abs, err := filepath.Abs("..")
	if err != nil {
		return "", fmt.Errorf("resolve repo root: %w", err)
	}
	if st, err := os.Stat(filepath.Join(abs, "molecule")); err != nil || !st.IsDir() {
		return "", fmt.Errorf("cannot locate molecule directory under [%s]", abs)
	}
	return abs, nil
}

func run(ctx context.Context) error {
	repoRoot, err := resolveRepoRoot()
	if err != nil {
		return err
	}
	if err := os.Chdir(repoRoot); err != nil {
		return fmt.Errorf("chdir %s: %w", repoRoot, err)
	}
	prepared, err := ci.Prepare(ctx, ci.PrepareOptions{
		RoleDir:      roleDir,
		SourceDir:    repoRoot,
		AppImageName: os.Getenv("APPIMAGE_NAME"),
		Isolated:     isCI(),
		TmpPrefix:    roleDir,
	})
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	if prepared.IsolatedDir != "" {
		defer func() { _ = os.RemoveAll(prepared.IsolatedDir) }()
	}
	const scenario = "default"
	if err := ci.MoleculeCreate(ctx, ci.MoleculeCreateOptions{
		MoleculeBinary: prepared.AppImageBin,
		Scenario:       scenario,
		LogBase:        prepared.LogBase,
	}); err != nil {
		return fmt.Errorf("molecule create: %w", err)
	}
	n := 1
	if err := ci.RunGroup(ctx, ci.RunGroupOptions{
		MoleculeBinary: prepared.AppImageBin,
		Scenario:       scenario,
		LogBase:        prepared.LogBase,
		Counter:        &n,
		RolesPath:      prepared.RolesPath,
	}); err != nil {
		return fmt.Errorf("run group: %w", err)
	}
	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT,
	)
	err := run(ctx)
	stop()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
