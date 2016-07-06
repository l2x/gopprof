package client

import "gopkg.in/robfig/cron.v2"

type job struct {
	c              *Client
	cb             *cron.Cron
	statsEntryID   cron.EntryID
	profileEntryID cron.EntryID
}

func (j *job) reload() {
	j.cb.Remove(j.statsEntryID)
	j.cb.Remove(j.profileEntryID)
	nodeConf := j.c.node.NodeConf
	if nodeConf.EnableStats {
		j.statsEntryID, _ = j.cb.AddFunc(nodeConf.StatsCron, func() {
			eventStats(j.c, nil)
		})
	}
	if nodeConf.EnableProfile && len(nodeConf.Profile) > 0 {
		j.profileEntryID, _ = j.cb.AddFunc(nodeConf.ProfileCron, func() {
			uploadProfile(j.c, nodeConf)
		})
	}
}
