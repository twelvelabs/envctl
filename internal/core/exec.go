package core

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/twelvelabs/termite/run"
)

func NewExecService(config *Config, client *run.Client) *ExecService {
	return &ExecService{
		config: config,
		client: client,
	}
}

type ExecService struct {
	config *Config
	client *run.Client
}

func (s *ExecService) Run(ctx context.Context, args []string, vars EnvVars) (*run.Cmd, error) {
	executable, err := exec.LookPath(args[0])
	if err != nil {
		return nil, err
	}

	cmd := s.client.CommandContext(ctx, executable, args[1:]...)
	env := []string{}
	for k, v := range vars {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the command.
	if err := cmd.Start(); err != nil {
		return cmd, err
	}

	// Forward signals.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals)
	go func() {
		for {
			s := <-signals
			_ = cmd.Process.Signal(s)
		}
	}()

	// Wait for the command to exit.
	if err := cmd.Wait(); err != nil {
		_ = cmd.Process.Signal(syscall.SIGKILL)
		return cmd, err
	}

	return cmd, nil
}
