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

	s.removeInvalidData()

	return nil
}

// AddDrop to the stream
func (s *Stream) AddDrop(point gpx.TrackPoint) {

	var previous Drop
	if len(s.Drops) > 0 {
		previous = s.Drops[len(s.Drops)-1]
	}

	d := Drop{}
	d.DataPoints = previous.DataPoints + 1

	d.addTime(&point, &previous)
	d.addDistance(&point, &previous)
	d.addPace(&point, &previous)
	d.addTemperature(&point, &previous)
	d.addHeartRate(&point, &previous)
	d.addCadence(&point, &previous)

	valid := d.checkValidity(&point, &previous)

	if valid {
		s.Drops = append(s.Drops, d)
	}
}

func (s *Stream) removeInvalidData() {

}

// WriteToGPXFile writes a new compliant GPX file
func (s *Stream) WriteToGPXFile(name string) error {
	g := gpx.GPX{}
	g.Version = "1.1"
	g.Creator = "github/sudhanshuraheja/gpx-fixer"

	ts := gpx.TrackSegment{}
	ts.TrackPoint = []gpx.TrackPoint{}

	for _, d := range s.Drops {
		tp := gpx.TrackPoint{}
		tp.Latitude = gpx.Latitude(d.Latitude)
		tp.Longitude = gpx.Longitude(d.Longitude)
		tp.Elevation = d.Elevation
		tp.Timestamp = d.Timestamp.Format(time.RFC3339)

		tp.Extensions = &gpx.TrackPointExtensions{}
		tp.Extensions.TrackPointExtensions = &gpx.TrackPointExtension{}
		tp.Extensions.TrackPointExtensions.Temperature = gpx.DegreesCelcius(d.Temperature)
		tp.Extensions.TrackPointExtensions.HeartRate = gpx.BeatsPerMinute(d.HeartRate)
		tp.Extensions.TrackPointExtensions.Cadence = gpx.RevolutionsPerMinute(d.Cadence)

		ts.TrackPoint = append(ts.TrackPoint, tp)
	}

	t := gpx.Track{}
	t.Name = name
	t.TrackSegments = []gpx.TrackSegment{ts}

	g.Tracks = []gpx.Track{t}

	return gpx.Write(&g, "out")
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

func (d *Drop) addTime(point *gpx.TrackPoint, previous *Drop) {
	t, err := time.Parse(time.RFC3339, point.Timestamp)
	if err != nil {
		log.Fatalf("Could not parse %v", point.Timestamp)
	}
	d.Timestamp = t

	if d.DataPoints > 1 {
		d.Seconds = d.Timestamp.Sub(previous.Timestamp).Seconds()
		d.Aggregates.Time = previous.Aggregates.Time + d.Seconds
	}
}

func (d *Drop) addDistance(point *gpx.TrackPoint, previous *Drop) {
	if point.Latitude != 0.0 && point.Longitude != 0.0 {

		d.Latitude = float64(point.Latitude)
		d.Longitude = float64(point.Longitude)
		d.Elevation = float64(point.Elevation)

		if ((d.Elevation - previous.Elevation) > 0) && (d.DataPoints > 1) {
			d.Aggregates.Climb = previous.Aggregates.Climb + (d.Elevation - previous.Elevation)
		} else {
			d.Aggregates.Climb = previous.Aggregates.Climb
		}

		if d.DataPoints > 1 {
			d.Distance = geo.DistanceInMetresBetweenPoints(geo.Point{Latitude: d.Latitude, Longitude: d.Longitude, Elevation: d.Elevation}, geo.Point{Latitude: previous.Latitude, Longitude: previous.Longitude, Elevation: previous.Elevation})
			d.Distance3D = geo.DistanceInMetresIncludingElevation(geo.Point{Latitude: d.Latitude, Longitude: d.Longitude, Elevation: d.Elevation}, geo.Point{Latitude: previous.Latitude, Longitude: previous.Longitude, Elevation: previous.Elevation})

			d.Aggregates.Distance = previous.Aggregates.Distance + d.Distance
			d.Aggregates.Distance3D = previous.Aggregates.Distance3D + d.Distance3D
		}
	}
}

func (d *Drop) addPace(point *gpx.TrackPoint, previous *Drop) {
	if point.Latitude != 0.0 && point.Longitude != 0.0 {
		if d.DataPoints > 1 {
			d.Pace = math.Round(((1000*d.Seconds)/(d.Distance3D*60))*10000) / 10000
		}
	}
}

func (d *Drop) addTemperature(point *gpx.TrackPoint, previous *Drop) {
	if point.Extensions.TrackPointExtensions.Temperature != 0 {
		d.Temperature = int(point.Extensions.TrackPointExtensions.Temperature)

		if d.Temperature > previous.Maximums.Temperature {
			d.Maximums.Temperature = d.Temperature
		} else {
			d.Maximums.Temperature = previous.Maximums.Temperature
		}
	}
}

func (d *Drop) addHeartRate(point *gpx.TrackPoint, previous *Drop) {
	if point.Extensions.TrackPointExtensions.HeartRate != 0 {
		d.HeartRate = int(point.Extensions.TrackPointExtensions.HeartRate)

		if d.HeartRate > previous.Maximums.HeartRate {
			d.Maximums.HeartRate = d.HeartRate
		} else {
			d.Maximums.HeartRate = previous.Maximums.HeartRate
		}
	}
}

func (d *Drop) addCadence(point *gpx.TrackPoint, previous *Drop) {
	if point.Extensions.TrackPointExtensions.Cadence != 0 {
		d.Cadence = int(point.Extensions.TrackPointExtensions.Cadence)

		if d.Cadence > previous.Maximums.Cadence {
			d.Maximums.Cadence = d.Cadence
		} else {
			d.Maximums.Cadence = previous.Maximums.Cadence
		}
	}
}

func (d *Drop) checkValidity(point *gpx.TrackPoint, previous *Drop) bool {
	// If pace is under 2.9139 minutes/km, it's faster than world record timing for the marathon
	if d.DataPoints > 1 && d.Pace < 2.9139 {
		return false
	}
	return true
}
