package main

import (
	"log"
	"time"

	"github.com/alexflint/go-arg"
)

func main() {
	var args struct {
		BeforeSunset time.Duration `arg:"--before" help:"duration before sunset to turn on" default:"10m"`
		EndTime      string        `arg:"--end" help:"time to turn off" default:"00:15:00 AM"`
		Power        int           `arg:"--power" help:"power level" default:"140"`
		FadeInterval time.Duration `arg:"--fade-interval" help:"fade interval" default:"25ms"`
		Address      string        `arg:"--address" help:"address of the light controller" default:"10.3.1.140:9999"`
	}
	arg.MustParse(&args)

	log.Println("fetching ip geo...")
	loc, err := geo()
	if err != nil {
		log.Fatalf("failed to get geo: %v", err)
	}

	client := Client{
		Address:        args.Address,
		FullBrightness: args.Power,
		FadeInterval:   args.FadeInterval,
	}

	ss := time.Time{}
	for {
		end, err := ParseClock(args.EndTime, loc.Timezone)
		if err != nil {
			log.Fatalf("failed to parse end time: %v", err)
		}

		log.Println("fetching sunset...")
		cur, err := sunset(loc.Lat, loc.Lon)
		switch {
		case err != nil && !ss.IsZero():
			// no sunset time available
			log.Fatalf("failed to get initial sunset: %v", err)
		case err != nil && !ss.IsZero():
			// cached sunset time
			ss = Next(ss)
		default:
			// use updated sunset time
			ss, err = cur.SunsetTime()
			if err != nil {
				log.Fatalf("failed to get next sunset: %v", err)
			}
		}

		now := time.Now()
		up := ss.Add(-args.BeforeSunset)
		nextUp := Next(up)
		nextEnd := Next(end)

		switch {
		case now.Before(up):
			err = client.OnAt(up)

		case now.Before(end):
			err = client.OffAt(end)

		case nextUp.Before(nextEnd):
			err = client.OnAt(nextUp)

		case nextEnd.Before(nextUp):
			err = client.OffAt(nextEnd)
		}

		if err != nil {
			log.Printf("failed to perform action: %v", err)
		}
	}
}
