package cmd

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func IntegrationContainerRunner(t *testing.T, dockerfile string, command []string, successMessage []string, debug bool) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../.",
			Dockerfile: dockerfile,
		},
		Entrypoint: []string{"tail", "-f", "/dev/null"},
		WaitingFor: wait.ForExec(successMessage),
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

	_, reader, err := ubuntuC.Exec(ctx, command)
	if err != nil {
		t.Fatalf("genericContainer failed: %v", err)
	}

	if debug {
		buf := new(strings.Builder)
		_, err = io.Copy(buf, reader)
		if err != nil {
			t.Fatalf("genericContainer failed: %v", err)
		}
		// check errors
		fmt.Println(buf.String())
	}
}
