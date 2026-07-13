package internal

import (
	"context"
	"errors"
	"os"

	"go.podman.io/podman/v6/pkg/bindings"
)

func InitContainerRuntime() error {
	sock_dir := os.Getenv("XDG_RUNTIME_DIR")
	if sock_dir == "" {
		return errors.New("SOCK_DIR not found in env var XDG_RUNTIME_DIR")
	}

	socket := "unix:" + sock_dir + "/podman.sock"

	_, err := bindings.NewConnection(context.Background(), socket)

	if err != nil {
		return err
	}

	return nil
}
