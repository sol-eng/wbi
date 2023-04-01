package cmd

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func IntegrationContainerRunner(t *testing.T, dockerfile string, command []string, successMessage string) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../.",
			Dockerfile: dockerfile,
		},
		Entrypoint: []string{"tail", "-f", "/dev/null"},
		WaitingFor: wait.ForExec([]string{successMessage}),
	}
	ubuntuC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}

	defer func() {
		if err := ubuntuC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err.Error())
		}
	}()

	c, _, err := ubuntuC.Exec(ctx, command)
	if err != nil {
		t.Fatalf("genericContainer failed: %v", err)
	}
	if c > 0 {
		t.Fatalf("genericContainer failed: %v", err)
	}
}
