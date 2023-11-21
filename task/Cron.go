package task

import (
	"github.com/robfig/cron"
)

type (
	Cron struct {
		before func(sign any)
		after  func(sign any)
		*cron.Cron
	}

	Job struct {
		cmd  func()
		sign any
	}
	Schedule = cron.Schedule
)

func (j Job) Run() {
	j.cmd()
}

func (c *Cron) AddFuncSign(spec string, cmd func(), sign any) error {
	return c.AddJob(spec, Job{cmd: cmd, sign: sign})
}

func (c *Cron) AddFunc(spec string, cmd func()) error {
	return c.AddJob(spec, Job{cmd: cmd})
}

func (c *Cron) AddJob(spec string, cmd Job) error {
	schedule, err := cron.Parse(spec)
	if err != nil {
		return err
	}
	c.Schedule(schedule, cmd)
	return nil
}

func (c *Cron) Schedule(schedule Schedule, cmd Job) {
	_cmd := cmd.cmd
	if c.before != nil {
		if c.after != nil {
			cmd.cmd = func() {
				c.before(cmd.sign)
				_cmd()
				c.after(cmd.sign)
			}
		} else {
			cmd.cmd = func() {
				c.before(cmd.sign)
				_cmd()
			}
		}
	} else if c.after != nil {
		cmd.cmd = func() {
			_cmd()
			c.after(cmd.sign)
		}
	} else {
		_cmd = nil
	}
	if c.Cron == nil {
		c.Cron = cron.New()
	}
	c.Cron.Schedule(schedule, cmd)
}
