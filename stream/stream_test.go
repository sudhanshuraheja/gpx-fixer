package stream_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sudhanshuraheja/gpx-fixer/stream"
)

func Test_LoadStream(t *testing.T) {
	s := stream.Stream{Drops: []stream.Drop{}}
	// s.Load("../external/strava-8kDelhi-May18.gpx")
	// s.Load("../external/strava-BadGPS-Singapore.gpx")
	// s.Load("../external/strava-Mumbai-Jan15.gpx")
	s.Load("../external/strava-ADHM-2014.gpx")

	final := s.Drops[len(s.Drops)-1]
	assert.Equal(t, 21422.215, math.Round(final.Aggregates.Distance*1000)/1000)
	assert.Equal(t, 21440.655, math.Round(final.Aggregates.Distance3D*1000)/1000)
	assert.Equal(t, 7520, int(final.Aggregates.Time))
}
