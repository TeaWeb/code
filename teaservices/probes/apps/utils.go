package apps

import (
	"github.com/iwind/TeaGo/logs"
	"github.com/shirou/gopsutil/process"
	"os/exec"
	"strings"
	"github.com/iwind/TeaGo/types"
	"regexp"
)

func ps(lookup string, matchPatterns []string, onlyParent bool) (result []*process.Process) {
	result = []*process.Process{}

	file, err := exec.LookPath("pgrep")
	if err != nil || len(file) == 0 {
		return
	}

	cmd := exec.Command(file, "-f", lookup)
	data, err := cmd.Output()
	if err != nil {
		return
	}
	dataString := strings.TrimSpace(string(data))
	if len(dataString) == 0 {
		return
	}

	for _, pidString := range strings.Split(dataString, "\n") {
		pid := types.Int32(pidString)
		if pid == 0 {
			continue
		}
		p, err := process.NewProcess(pid)
		if err != nil {
			logs.Error(err)
		} else {
			if onlyParent {
				ppid, err := p.Ppid()

				if err == nil && ppid > 128 {
					continue
				}
			}

			if len(matchPatterns) > 0 {
				cmdLine, err := p.CmdlineSlice()
				if err != nil {
					logs.Error(err)
					continue
				}

				matched := false
				for _, piece := range cmdLine {
					failed := false
					for _, pattern := range matchPatterns {
						reg, err := regexp.Compile(pattern)
						if err != nil {
							logs.Error(err)
							failed = true
							break
						}
						if !reg.MatchString(piece) {
							failed = true
							break
						}
					}
					if !failed {
						matched = true
					}
				}
				if !matched {
					continue
				}
			}

			result = append(result, p)
		}
	}

	return
}
