package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mydocker/container"
	"strconv"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var StopCommand = cli.Command{
	Name:  "stop",
	Usage: "stop a container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container name")
		}
		containerName := context.Args().Get(0)
		stopContainer(containerName)
		return nil
	},
}

func stopContainer(containerName string) {
	pid, err := GetContainerPidByName(containerName)
	if err != nil {
		log.Errorf("Get contaienr pid by name %s error %v", containerName, err)
		return
	}
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		log.Errorf("Conver pid from string to int error %v", err)
		return
	}
	if err := syscall.Kill(pidInt, syscall.SIGTERM); err != nil {
		log.Errorf("Stop container %s error %v", containerName, err)
		return
	}
	containerInfo, err := getContainerInfoByName(containerName)
	if err != nil {
		log.Errorf("Get container %s info error %v", containerName, err)
		return
	}
	containerInfo.Status = container.STOP
	containerInfo.Pid = " "
	newContentBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Json marshal %s error %v", containerName, err)
		return
	}
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath := dirURL + container.ConfigName
	if err := ioutil.WriteFile(configFilePath, newContentBytes, 0622); err != nil {
		log.Errorf("Write file %s error", configFilePath, err)
	}
}

func getContainerInfoByName(containerName string) (*container.ContainerInfo, error) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath := dirURL + container.ConfigName
	contentBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Errorf("Read file %s error %v", configFilePath, err)
		return nil, err
	}
	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		log.Errorf("GetContainerInfoByName unmarshal error %v", err)
		return nil, err
	}
	return &containerInfo, nil
}
