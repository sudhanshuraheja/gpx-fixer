# gpx-fixer

gpx-fixer is a small Go library for reading GPX activity files, building a time series ("stream") of track points, and exporting a cleaned GPX track. It computes distance, elevation-adjusted distance, pace, and aggregates (time, distance, climb), and filters out implausible points (e.g., pace faster than world-record marathon pace).

This repository is primarily a library; `main.go` is currently empty. Sample GPX inputs live in `external/`, and tests validate the calculations.

## How it works

- **Parse GPX**: Uses `github.com/sudhanshuraheja/go-garmin-gpx` to load tracks and track points.
- **Stream building**: Each GPX point becomes a `Drop` in a `Stream`, carrying metrics like timestamp, distance, pace, temperature, heart rate, and cadence.
- **Validation**: Drops with a pace faster than 2.9139 min/km are rejected.
- **Write GPX**: The stream can be exported as a new GPX file (written to `out.gpx` via the `out` base name).

## Packages

- `geo`: Geographic helpers, including 2D and elevation-adjusted distance calculations.
- `stream`: Stream ingestion, derived metrics, validation, and GPX export.

## Example

```go
package main

import (
	"log"

	"github.com/sudhanshuraheja/gpx-fixer/stream"
)

func main() {
	s := stream.Stream{}
	if err := s.Load("external/strava-ADHM-2014.gpx"); err != nil {
		log.Fatal(err)
	}

	if err := s.WriteToGPXFile("fixed"); err != nil {
		log.Fatal(err)
	}
}
```

## Tests

```sh
make test
```

This runs unit tests for the `geo` package and integration-style checks over the sample GPX files in `external/`.
