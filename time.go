package yagolib

import (
	"time"
	)

type TimeCtlStatus struct {
	LocalTime Time
	UTCTime Time
	RTCTime Time
	TimeZone string
	
}

func Get