package yagolib

import (
	"bytes"
	"fmt"
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
		if strings.Contains(line, "RTC in local TZ: yes") {
			status.IsRTCInLocalTZ = true
		}
		if strings.Contains(line, "Network time on: yes") {
			status.IsNetworkTimeOn = true
		}
		if strings.Contains(line, "NTP synchronized: yes") {
			status.IsNTPSynchronized = true
		}
	}
	if status.IsRTCInLocalTZ {
		tz, _ := time.Now().Zone()
		s := status.RTCTime.Format("Mon 2006-01-02 15:04:05 " + tz)
		status.RTCTime, _ = time.Parse("Mon 2006-01-02 15:04:05 MST", s)
	}

	out.Reset()
	cmd = exec.Command("systemctl", "status", "systemd-timesyncd")
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	outStr = out.String()
	outLines = strings.Split(outStr, "\n")
	for _, line := range outLines {
		if strings.Contains(line, "Loaded: loaded") {
			status.IsLoaded = true
		}
		if strings.Contains(line, "Active: active") {
			status.IsActive = true
			if status.StartedAt.IsZero() {
				i := strings.Index(line, "since")
				if i > 0 {
					parseFmt := "Mon 2006-01-02 15:04:05 MST"
					str := strings.TrimSpace(line[i+5 : i+6+len(parseFmt)])
					t, err := time.Parse(parseFmt, str)
					if err == nil {
						status.StartedAt = t
					}
				}
			}
		}
		if strings.Contains(line, "Status") {
			i := strings.Index(line, ":")
			if i > 0 {
				status.Status = strings.Trim(line[i+1:], " \"")
			}
		}
		if status.StartedAt.IsZero() &&
			strings.Contains(line, "Started Network Time Synchronization") {
			now := time.Now()
			loc := now.Location()
			s := fmt.Sprint(line[:16], now.Year())
			t, err := time.ParseInLocation("Jan 02 15:04:05 2006", s, loc)
			if err == nil {
				if t.After(now) {
					t = time.Date(now.Year()-1, t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, loc)
				}
				status.StartedAt = t
			}
		}
		if strings.Contains(line, "Synchronized to time server") {
			now := time.Now()
			loc := now.Location()
			s := fmt.Sprint(line[:16], now.Year())
			t, err := time.ParseInLocation("Jan 02 15:04:05 2006", s, loc)
			if err == nil {
				if t.After(now) {
					t = time.Date(now.Year()-1, t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, loc)
				}
				status.SynchronizedAt = t
			}
		}
	}

	return &status, nil
}
