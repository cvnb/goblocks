package modules

import (
	"fmt"
	"os"

	"github.com/davidscholberg/go-i3barjson"
)

// Temperature represents the configuration for the CPU temperature block.
// CpuTempPath is the path to the "hwmon" directory of the CPU temperature info.
// e.g. /sys/devices/platform/coretemp.0/hwmon
type Temperature struct {
	BlockConfigBase `yaml:",inline"`
	TempPath        string  `yaml:"cpu_temp_path"`
	CritTemp        float64 `yaml:"crit_temp"`
}

// UpdateBlock updates the CPU temperature info.
// The value output by the block is the average temperature of all cores.
func (c Temperature) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	fullTextFmt := fmt.Sprintf("%s%%s", c.Label)
	temp := 0

	_, err := os.Stat(c.TempPath)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}

	r, err := os.Open(c.TempPath)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	defer r.Close()

	_, err = fmt.Fscanf(r, "%d", &temp)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}

	calcTemp := float64(temp) / float64(1000)
	if calcTemp >= c.CritTemp {
		b.Urgent = true
	} else {
		b.Urgent = false
	}
	b.FullText = fmt.Sprintf("%s%.2f°C", c.Label, calcTemp)
}
