package cgroups

import (
	"mydocker/cgroups/subsystems"

	"github.com/sirupsen/logrus"
)

type CgroupManager struct {
	Path     string
	Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

func (c *CgroupManager) Apply(pid int) error {
	for _, subSysInts := range subsystems.SubSystemsIns {
		subSysInts.Apply(c.Path, pid)
	}
	return nil
}
func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subSysInts := range subsystems.SubSystemsIns {
		subSysInts.Set(c.Path, res)
	}
	return nil
}

func (c *CgroupManager) Destory() error {
	for _, subSysInts := range subsystems.SubSystemsIns {
		if err := subSysInts.Remove(c.Path); err != nil {
			logrus.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
