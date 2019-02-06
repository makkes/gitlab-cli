package api

import (
	"testing"
	"time"
)

func TestPipelineDetailsDurationRunning(t *testing.T) {
	var zeroTime time.Time

	tests := []struct {
		now              time.Time
		started          time.Time
		updated          time.Time
		finished         time.Time
		recordedDuration *int
		status           string
		out              string
	}{
		{
			time.Unix(1549449133, 0),
			time.Unix(1549448118, 0),
			time.Unix(1549449125, 0),
			time.Unix(1549449129, 0),
			nil,
			"running",
			"8s",
		},
		{
			time.Unix(1549449150, 0),
			time.Unix(1549449125, 0),
			zeroTime,
			zeroTime,
			nil,
			"running",
			"25s",
		},
		{
			time.Unix(1549449150, 0),
			time.Unix(1549449125, 0),
			zeroTime,
			zeroTime,
			nil,
			"failed",
			"-",
		},
		{
			zeroTime,
			zeroTime,
			zeroTime,
			zeroTime,
			func() *int { i := int(2283); return &i }(),
			"success",
			"38m3s",
		},
	}

	for _, tt := range tests {
		pd := PipelineDetails{
			StartedAt:        tt.started,
			UpdatedAt:        tt.updated,
			FinishedAt:       tt.finished,
			RecordedDuration: tt.recordedDuration,
			Status:           tt.status,
		}

		res := pd.Duration(tt.now)
		if res != tt.out {
			t.Errorf("Unexpected duration: %s", res)
		}
	}

}
