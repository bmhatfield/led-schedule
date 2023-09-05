package main

import (
	"log"
	"time"

	"github.com/alexflint/go-arg"
)

func main() {
	var args struct {
		BeforeSunset time.Duration `arg:"--before" help:"duration before sunset to turn on" default:"10m"`
		Power        int           `arg:"--power" help:"power level" default:"140"`
		Lat          float64       `arg:"--lat" help:"latitude"`
		Lon          float64       `arg:"--lon" help:"longitude"`
		FadeInterval time.Duration `arg:"--fade-interval" help:"fade interval" default:"25ms"`
		Address      string        `arg:"--address" help:"address of the light controller" default:"10.3.1.140:9999"`
	}
	arg.MustParse(&args)

	lat, lon := args.Lat, args.Lon
	if lat == 0 && lon == 0 {
		log.Println("fetching ip geo... (specify lat, lon to skip)")
		loc, err := geo()
		if err != nil {
			log.Fatalf("failed to get geo: %v", err)
		}

		lat, lon = loc.Lat, loc.Lon
	}

	client := Client{
		Address:        args.Address,
		FullBrightness: args.Power,
		FadeInterval:   args.FadeInterval,
	}

	ss := time.Time{}
	for {
		log.Println("fetching sunset...")
		cur, err := sunset(lat, lon)
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
		down := Next(Clock(0, 15, 0))

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
