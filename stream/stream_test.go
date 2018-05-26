package stream_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sudhanshuraheja/gpx-fixer/stream"
)

func Test_LoadStream(t *testing.T) {
	adhm := stream.Stream{Drops: []stream.Drop{}}
	adhm.Load("../external/strava-ADHM-2014.gpx")

	final := adhm.Drops[len(adhm.Drops)-1]
	assert.Equal(t, 21422.16, math.Round(final.Aggregates.Distance*1000)/1000)
	assert.Equal(t, 21440.597, math.Round(final.Aggregates.Distance3D*1000)/1000)
	assert.Equal(t, 7520, int(final.Aggregates.Time))

	bad := stream.Stream{Drops: []stream.Drop{}}
	bad.Load("../external/strava-BadGPS-Singapore.gpx")

	final = bad.Drops[len(bad.Drops)-1]
	assert.Equal(t, 7682.53, math.Round(final.Aggregates.Distance*1000)/1000)
	assert.Equal(t, 7699.526, math.Round(final.Aggregates.Distance3D*1000)/1000)
	assert.Equal(t, 2264, int(final.Aggregates.Time))

	err := bad.WriteToGPXFile("badgps")
	assert.Nil(t, err)

	// s.Load("../external/strava-8kDelhi-May18.gpx")
	// s.Load("../external/strava-Mumbai-Jan15.gpx")

}
