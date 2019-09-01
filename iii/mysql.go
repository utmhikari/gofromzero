package iii

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"log"
	"os"
)

var ctx = context.Background()

func login(cli *client.Client) {
	log.Println("Logging in docker registry...")
	ok, err := cli.RegistryLogin(ctx, types.AuthConfig{
		Username: "用户名",
		Password: "密码",
	})
	if err != nil {
		log.Fatalf("Error while logging in docker registry! %s", err.Error())
	}
	log.Printf("%s --- Token: %s\n", ok.Status, ok.IdentityToken)
}

func closeClient(cli *client.Client) {
	err := cli.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func pullImage(cli *client.Client) {
	log.Println("Pulling MySQL Image...")
	reader, err := cli.ImagePull(
		ctx,
		"docker.io/library/mysql",
		types.ImagePullOptions{})
	if err != nil {
		log.Fatalf("Error while pulling image! %s", err.Error())
	}
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Println("Successfully pulled MySQL Image!")
}

// create and start image
func runImage(cli *client.Client) string {
	log.Println("Running MySQL Image...")
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: "mysql:latest",
			Env: []string{"MYSQL_ROOT_PASSWORD", "123456"},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"3306/tcp": []nat.PortBinding{
					{
						HostIP: "0.0.0.0",
						HostPort: "3306",
					},
				},
			},
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: "E:\\Tools\\MySQL",
					Target: "/var/lib/mysql",
				},
			},
		},
		nil,
		"MySQLDB")
	if err != nil {
		log.Fatalf("Error while creating image! %s", err.Error())
	}
	log.Printf("Successfully created MySQL image: %s!\n", resp.ID)
	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Fatalf("Error while starting image! %s", err.Error())
	}
	log.Println("Successfully ran MySQL image!")
	return resp.ID
}

func logImage(cli *client.Client, containerID string) {
	log.Println("Fetching log on MySQL container...")
	reader, err := cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout:true,
		ShowStderr:true,
		Timestamps:true,
		Follow:true,
		Details:true,
	})
	if err != nil {
		log.Fatalf("Error while logging image! %s", err.Error())
	}
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func LaunchMySQL() {
	log.Println("Creating docker client...")
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("Error while creating docker client! %s", err.Error())
	}
	defer closeClient(cli)
	login(cli)
	pullImage(cli)
	id := runImage(cli)
	logImage(cli, id)
}
