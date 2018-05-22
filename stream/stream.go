package stream

import (
	"fmt"
	"log"
	"math"
	"time"

	gpx "github.com/sudhanshuraheja/go-garmin-gpx"
	"github.com/sudhanshuraheja/gpx-fixer/geo"
)

// Stream of data of a run
type Stream struct {
	Drops []Drop
}

// Drop is an individual point of data
type Drop struct {
	DataPoints  int
	Latitude    float64
	Longitude   float64
	Elevation   float64
	Seconds     float64
	Timestamp   time.Time
	Temperature int
	HeartRate   int
	Cadence     int
	Distance    float64
	Distance3D  float64
	Pace        float64
	Aggregates  struct {
		Climb      float64
		Time       float64
		MovingTime float64
		Distance   float64
		Distance3D float64
	}
	Maximums struct {
		Temperature int
		HeartRate   int
		Cadence     int
	}
}

// Display the drop
func (d *Drop) Display() {
	fmt.Printf("[%v][%v] %v,%v ^%v s%v c%v h%v d%v %v (%v) <%v> [A] ^%v +%v %v (%v) [M] c%v h%v d%v\n", d.Timestamp, d.DataPoints, d.Latitude, d.Longitude, int(d.Elevation), d.Seconds, d.Temperature, d.HeartRate, d.Cadence, d.Distance, d.Distance3D, d.Pace, int(d.Aggregates.Climb), d.Aggregates.Time, int(d.Aggregates.Distance), int(d.Aggregates.Distance3D), d.Maximums.Temperature, d.Maximums.HeartRate, d.Maximums.Cadence)
}

// Load a stream from a file
func (s *Stream) Load(filePath string) error {

	g, err := gpx.ParseFile(filePath)
	if err != nil {
		return err
	}

	for _, track := range g.Tracks {
		for _, trackSegment := range track.TrackSegments {
			for _, trackPoint := range trackSegment.TrackPoint {
				s.AddDrop(trackPoint)
			}
		}
	}

	return nil
}

// AddDrop to the stream
func (s *Stream) AddDrop(point gpx.TrackPoint) {

	var previous Drop
	if len(s.Drops) > 0 {
		previous = s.Drops[len(s.Drops)-1]
	}

	d := Drop{}

	t, err := time.Parse(time.RFC3339, point.Timestamp)
	if err != nil {
		log.Fatalf("Could not parse %v", point.Timestamp)
	}
	d.Timestamp = t

	if point.Latitude != 0.0 && point.Longitude != 0.0 {
		d.DataPoints = previous.DataPoints + 1
		d.Latitude = float64(point.Latitude)
		d.Longitude = float64(point.Longitude)
		d.Elevation = float64(point.Elevation)

		if ((d.Elevation - previous.Elevation) > 0) && (d.DataPoints != 1) {
			d.Aggregates.Climb = previous.Aggregates.Climb + (d.Elevation - previous.Elevation)
		} else {
			d.Aggregates.Climb = previous.Aggregates.Climb
		}

		if len(s.Drops) > 0 {
			d.Distance = geo.DistanceInMetresBetweenPoints(geo.Point{Latitude: d.Latitude, Longitude: d.Longitude, Elevation: d.Elevation}, geo.Point{Latitude: previous.Latitude, Longitude: previous.Longitude, Elevation: previous.Elevation})
			d.Distance3D = geo.DistanceInMetresIncludingElevation(geo.Point{Latitude: d.Latitude, Longitude: d.Longitude, Elevation: d.Elevation}, geo.Point{Latitude: previous.Latitude, Longitude: previous.Longitude, Elevation: previous.Elevation})

			d.Pace = math.Round((1000/(d.Distance3D*60))*10000) / 10000

			d.Aggregates.Distance = previous.Aggregates.Distance + d.Distance
			d.Aggregates.Distance3D = previous.Aggregates.Distance3D + d.Distance3D
		}

		if d.DataPoints != 1 {
			d.Seconds = d.Timestamp.Sub(previous.Timestamp).Seconds()
			d.Aggregates.Time = previous.Aggregates.Time + d.Seconds
		}
	}

	if point.Extensions.TrackPointExtensions.Temperature != 0 {
		d.Temperature = int(point.Extensions.TrackPointExtensions.Temperature)

		if d.Temperature > previous.Maximums.Temperature {
			d.Maximums.Temperature = d.Temperature
		} else {
			d.Maximums.Temperature = previous.Maximums.Temperature
		}
	}

	if point.Extensions.TrackPointExtensions.HeartRate != 0 {
		d.HeartRate = int(point.Extensions.TrackPointExtensions.HeartRate)

		if d.HeartRate > previous.Maximums.HeartRate {
			d.Maximums.HeartRate = d.HeartRate
		} else {
			d.Maximums.HeartRate = previous.Maximums.HeartRate
		}
	}

	if point.Extensions.TrackPointExtensions.Cadence != 0 {
		d.Cadence = int(point.Extensions.TrackPointExtensions.Cadence)

		if d.Cadence > previous.Maximums.Cadence {
			d.Maximums.Cadence = d.Cadence
		} else {
			d.Maximums.Cadence = previous.Maximums.Cadence
		}
	}

	s.Drops = append(s.Drops, d)
}
