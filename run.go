package main

import (
	"mydocker/cgroups"
	"mydocker/cgroups/subsystems"
	"mydocker/container"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Run(tty bool, comArray []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("new parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	//use mydocker cgroup
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destory()
	//设置容器资源限制
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)
	//对容器设置完限制后，初始化容器
	sendInitCommand(comArray, writePipe)
	parent.Wait()
	os.Exit(0)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
