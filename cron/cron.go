package cron

import (
	"sync"

	"github.com/owenzhou/ginrbac/contracts"
	"github.com/robfig/cron/v3"
)

func newCron(app contracts.IApplication) *Cron {
	return &Cron{cronList: make(map[string]*cron.Cron)}
}

type Cron struct {
	cronList map[string]*cron.Cron
	sync.Mutex
}

//添加函数任务
func (c *Cron) AddFunc(name, spec string, task func()) (cron.EntryID, error) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.cronList[name]; !ok {
		c.cronList[name] = cron.New()
	}
	defer c.cronList[name].Start()
	return c.cronList[name].AddFunc(spec, task)
}

//添加job任务
func (c *Cron) AddJob(name, spec string, job interface{ Run() }) (cron.EntryID, error) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.cronList[name]; !ok {
		c.cronList[name] = cron.New()
	}
	defer c.cronList[name].Start()
	return c.cronList[name].AddJob(spec, job)
}

//开启任务
func (c *Cron) Start(name string) {
	c.Lock()
	defer c.Unlock()
	if cron, ok := c.cronList[name]; ok {
		cron.Start()
	}
}

//停止任务
func (c *Cron) Stop(name string) {
	c.Lock()
	defer c.Unlock()
	if cron, ok := c.cronList[name]; ok {
		cron.Stop()
	}
}

//移除任务里的子任务
func (c *Cron) Remove(name string, id cron.EntryID) {
	c.Lock()
	defer c.Unlock()
	if cron, ok := c.cronList[name]; ok {
		cron.Remove(id)
	}
}

//清空任务
func (c *Cron) Clear(name string) {
	c.Lock()
	defer c.Unlock()
	if cron, ok := c.cronList[name]; ok {
		cron.Stop()
		delete(c.cronList, name)
	}
}

//查找任务
func (c *Cron) Find(name string) (*cron.Cron, bool) {
	c.Lock()
	defer c.Unlock()
	cron, ok := c.cronList[name]
	return cron, ok
}

//关闭所有任务
func (c *Cron) Close() {
	c.Lock()
	defer c.Unlock()
	for _, v := range c.cronList {
		v.Stop()
	}
}
