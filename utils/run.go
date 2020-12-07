package utils

import (
	"mydocker/cgroups"
	"mydocker/cgroups/subsystems"
	"mydocker/container"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

//TODO 3.2.2
func Run(tty bool, comArry []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)

	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	cgroupManager := cgroups.NewCgroupManager("mydocker-group")
	defer cgroupManager.Destory()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)
	sendInitCommand(comArry, writePipe)
	parent.Wait()
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
