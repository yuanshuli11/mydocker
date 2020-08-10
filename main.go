package main

import (
	"fmt"
	"mydocker/cgroups"
	"mydocker/cgroups/subsystems"
	"mydocker/container"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const usage = `mydocker is a simple container runtime implementation.
			   The purpose of this project is to learn how docker works and how to write a docker by ourselves
			   Enjoy it,just for fun.`

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage
	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	app.Before = func(context *cli.Context) error {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}
	fmt.Println("main args=====", os.Args)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var runCommand = cli.Command{
	Name:  "run",
	Usage: "Create a container with namespace and cgroups limit mydocker run -ti [command]",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		tty := context.Bool("ti")
		fmt.Println("=====run args", context.Args())
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: context.String("m"),
			//	CpuSet:      context.String("cpuset"),
			//	CpuShare:    context.String("cpushare"),
		}
		Run(tty, cmdArray, resConf)
		return nil
	},
}
var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container.Do not call it outside",

	Action: func(context *cli.Context) error {
		log.Infof("init come on")
		cmd := context.Args().Get(0)
		fmt.Println("=====ini args", context.Args())
		log.Infof("command 2 %s", cmd)
		cmd = "/proc/self/exe"
		err := container.RunContainerInitProcess(cmd, nil)
		fmt.Println("err:", err)
		return err
	},
}

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
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
