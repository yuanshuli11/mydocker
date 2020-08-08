package container

import (
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
)

func RunContainerInitProcess(command string, args []string) error {

	//	cmdArray := readuserCommand()
	//	if cmdArray == nil || len(cmdArray) == 0 {
	//		return fmt.Errorf("Run container get user command error,cmdArray is nil")
	//	}
	//	setUpMount()

	logrus.Infof("init command %s", command)
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV

	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}
	return nil
}
