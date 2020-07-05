package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"regexp"
	"time"

	"github.com/jtogrul/comake/util"
	"github.com/spf13/cobra"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	shellwords "github.com/mattn/go-shellwords"
)

var workingDir string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute build steps as described in the build file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting the build")

		buildfile, err := cmd.Flags().GetString("buildfile")
		if err != nil {
			panic(err)
		}

		workingDir, err := cmd.Flags().GetString("workdir")
		if err != nil {
			panic(err)
		}

		if workingDir == "" {
			workingDir, err = os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
		}

		config, err := util.ReadBuildConfig(buildfile)
		if err != nil {
			panic(err)
		}

		cli, err := client.NewEnvClient()
		if err != nil {
			panic(err)
		}

		ctx := context.Background()

		stopRunningContainers(ctx, cli)

		for _, step := range config.Steps {

			pullImageIfDoesntExist(ctx, cli, step.Image)

			containerName := fmt.Sprintf("%s-%d", cleanString(step.Name), time.Now().UnixNano())

			containerID, err := startContainer(ctx, cli, step.Image, containerName, workingDir)
			if err != nil {
				panic(err)
			}

			// TODO wait for the container to be ready?
			// TODO print logs of container starting?

			// Run step commands
			for _, stepCmd := range step.Script {
				runContainerCommand(ctx, cli, containerID, stepCmd)
			}

			// Stop the container
			if err := stopContainer(ctx, cli, containerID); err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringP("workdir", "d", "", "Working directory. Default: current directory")
}

func stopRunningContainers(ctx context.Context, cli *client.Client) {
	// TODO only stop involved containers
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		if err := stopContainer(ctx, cli, container.ID); err != nil {
			panic(err)
		}
	}
}

func pullImageIfDoesntExist(ctx context.Context, cli *client.Client, image string) {
	_, _, err := cli.ImageInspectWithRaw(ctx, image)
	if err != nil {
		if client.IsErrNotFound(err) {
			reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
			if err != nil {
				panic(err)
			}
			io.Copy(os.Stdout, reader)
		} else {
			panic(err)
		}
	} else {
		fmt.Println("Image already exists. Not pulling")
	}
}

func startContainer(ctx context.Context, cli *client.Client, image, containerName, workingDir string) (containerID string, err error) {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	containerCreate, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        image,
			WorkingDir:   "/source",
			User:         user.Uid,
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			OpenStdin:    true,
			Tty:          true,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: workingDir,
					Target: "/source",
				},
			},
			PortBindings: nat.PortMap{
				"8080/tcp": []nat.PortBinding{
					{
						HostIP:   "",
						HostPort: "8080",
					},
				},
			},
		},
		&network.NetworkingConfig{},
		containerName,
	)
	if err != nil {
		panic(err)
	}

	attachResp, err := cli.ContainerAttach(ctx, containerCreate.ID, types.ContainerAttachOptions{
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   true,
	})
	if err != nil {
		panic(err)
	}
	defer attachResp.Close()
	io.Copy(os.Stdout, attachResp.Reader)

	err = cli.ContainerStart(ctx, containerCreate.ID, types.ContainerStartOptions{})
	return containerCreate.ID, err
}

func runContainerCommand(ctx context.Context, cli *client.Client, containerID, stepCmd string) {
	fmt.Printf("$ %s\n", stepCmd)
	cmd, err := shellwords.Parse(stepCmd)
	if err != nil {
		panic(err)
	}
	execConfig := types.ExecConfig{
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          cmd,
	}
	resp, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		panic(err)
	}

	hijackResp, err := cli.ContainerExecAttach(ctx, resp.ID, execConfig)
	if err != nil {
		panic(err)
	}
	defer hijackResp.Close()

	_, err = io.Copy(os.Stdout, hijackResp.Reader)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}

func stopContainer(ctx context.Context, cli *client.Client, containerID string) error {
	fmt.Print("Stopping container ", containerID[:10], "... ")
	return cli.ContainerStop(ctx, containerID, nil)
}

func cleanString(text string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(text, "")
}
