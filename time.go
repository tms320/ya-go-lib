package yagolib

import (
	"bytes"
	"os/exec"
	"strings"
	"time"
)

type TimeCtlStatus struct {
	Time              time.Time
	RTCTime           time.Time
	IsRTCAvailable    bool
	IsRTCInLocalTZ    bool
	IsNetworkTimeOn   bool
	IsNTPSynchronized bool
	IsLoaded          bool
	IsActive          bool
	Status            string
	StartedAt         time.Time
	SynchronizedAt    time.Time
}

func GetTimeCtlStatus() (*TimeCtlStatus, error) {
	cmd := exec.Command("timedatectl", "status")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	var status TimeCtlStatus
	outStr := out.String()
	outLines := strings.Split(outStr, "\n")
	for _, line := range outLines {
		if strings.Contains(line, "Universal time") {
			i := strings.Index(line, ":")
			if i > 0 {
				str := strings.TrimSpace(line[i+1:])
				t, err := time.Parse("Mon 2006-01-02 15:04:05 MST", str)
				if err == nil {
					status.Time = t
				}
			}
		}
		//line = "      RTC time: Sat 2015-11-14 19:20:38"
		if strings.Contains(line, "RTC time") && !strings.Contains(line, "n/a") {
			i := strings.Index(line, ":")
			if i > 0 {
				str := strings.TrimSpace(line[i+1:])
				t, err := time.Parse("Mon 2006-01-02 15:04:05", str)
				if err == nil {
					status.IsRTCAvailable = true
					status.RTCTime = t
				}
			}
		}
		if strings.Contains(line, "RTC in local TZ") && strings.Contains(line, "yes") {
			status.IsRTCInLocalTZ = true
		}
		if strings.Contains(line, "Network time on") && strings.Contains(line, "yes") {
			status.IsNetworkTimeOn = true
		}
		if strings.Contains(line, "NTP synchronized") && strings.Contains(line, "yes") {
			status.IsNTPSynchronized = true
		}
	}

	if status.IsRTCInLocalTZ {
	}

	return &status, nil
}
