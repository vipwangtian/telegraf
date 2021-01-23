package proc_username

// proc_username.go

import (
	// "fmt"
	"strings"
	"strconv"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/processors"
	"github.com/shirou/gopsutil/v3/process"
)

type Procusername struct {
	Log telegraf.Logger `toml:"-"`
}

type ProcInfo struct {
	Pid int32
	Username string
}

var sampleConfig = `
`

func (p *Procusername) SampleConfig() string {
	return sampleConfig
}

func (p *Procusername) Description() string {
	return "Print all metrics that pass through this filter."
}

// Init is for setup, and validating config.
func (p *Procusername) Init() error {
	return nil
}

func (p *Procusername) Apply(in ...telegraf.Metric) []telegraf.Metric {
	for _, metric := range in {
		metricname := metric.Name()
		// p.Log.Debugf(metricname)
		if metricname == "win_proc" {
			if procnameInMetric, ret := metric.GetTag("instance"); ret {
				splitRet := strings.Split(procnameInMetric, "_")
				if len(splitRet) == 2 {
					if pid, err := strconv.ParseInt(splitRet[1], 10, 32); err == nil {
						if proc, err := process.NewProcess(int32(pid)); err == nil {
							if nameWithGroup, err := proc.Username(); err == nil {
								nameArrs := strings.Split(nameWithGroup, "\\")
								if len(nameArrs) == 2 {
									metric.AddTag("username", nameArrs[1])
								} else {
									metric.AddTag("username", nameWithGroup)
								}
							}
						}
					}
				} else {
					p.Log.Errorf("instance tag format must be xxx_pid, but now is %s", procnameInMetric)
				}
			} else {
				p.Log.Debugf("no instance tag in win_proc metric")
			}
		} else if strings.Contains(metricname, "win_dotnet_") {
			if procnameInMetric, ret := metric.GetTag("instance"); ret {
				splitRet := strings.Split(procnameInMetric, "_")
				if len(splitRet) == 3 {
					if pid, err := strconv.ParseInt(strings.Trim(splitRet[1], "p"), 10, 32); err == nil {
						if proc, err := process.NewProcess(int32(pid)); err == nil {
							if nameWithGroup, err := proc.Username(); err == nil {
								nameArrs := strings.Split(nameWithGroup, "\\")
								if len(nameArrs) == 2 {
									metric.AddTag("username", nameArrs[1])
								} else {
									metric.AddTag("username", nameWithGroup)
								}
							}
						}
					}
				} else {
					p.Log.Errorf("instance tag format must be xxx_pid_version, but now is %s", procnameInMetric)
				}
			} else {
				p.Log.Debugf("no instance tag in win_proc metric")
			}
		}
	}
	return in
}

func init() {
	processors.Add("proc_username", func() telegraf.Processor {
		return &Procusername{}
	})
}