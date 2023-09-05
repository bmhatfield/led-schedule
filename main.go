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

	log.Println("fetching ip geo... (specify lat, lon to skip)")
	loc, err := geo()
	if err != nil {
		log.Fatalf("failed to get geo: %v", err)
	}

	end, err := ParseClock(args.EndTime, loc.Timezone)
	if err != nil {
		log.Fatalf("failed to parse end time: %v", err)
	}

	client := Client{
		Address:        args.Address,
		FullBrightness: args.Power,
		FadeInterval:   args.FadeInterval,
	}

	ss := time.Time{}
	for {
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
		down := Next(end)

		switch {
		case now.Before(up):
			if client.IsOn {
				client.FadeDown()
			}

			// wait for fade-up
			wait := time.Until(up)
			log.Printf("sleeping %s until %s before sunset (%s)", wait, args.BeforeSunset, up)
			<-time.After(wait)
			client.FadeUp()
		case now.Before(down):
			if !client.IsOn {
				client.FadeUp()
			}

			// wait for fade down
			wait := time.Until(down)
			log.Printf("sleeping %s until fade-down (%s)", wait, down)
			<-time.After(wait)
			client.FadeDown()
		}
	}
}
