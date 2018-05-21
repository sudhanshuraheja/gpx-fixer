package stream_test

import (
	"testing"

	"github.com/sudhanshuraheja/gpx-fixer/stream"
)

func Test_LoadStream(t *testing.T) {
	s := stream.Stream{Drops: []stream.Drop{}}
	// s.Load("../external/strava-8kDelhi-May18.gpx")
	// s.Load("../external/strava-BadGPS-Singapore.gpx")
	// s.Load("../external/strava-Mumbai-Jan15.gpx")
	s.Load("../external/strava-ADHM-2014.gpx")
}
