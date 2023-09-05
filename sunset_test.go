package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSunset(t *testing.T) {
	j, err := os.Open("testdata/sunrisesunset.json")
	require.NoError(t, err)

	resp := new(Response)
	err = json.NewDecoder(j).Decode(resp)
	require.NoError(t, err)
	require.Equal(t, "OK", resp.Status)

	t.Run("SunsetTime", func(t *testing.T) {
		x, err := resp.Result.SunsetTime()
		require.NoError(t, err)

		now := time.Now()
		s := time.Date(now.Year(), now.Month(), now.Day(), x.Hour(), x.Minute(), x.Second(), 0, time.UTC)
		require.Equal(t, s, x)
	})

	t.Run("SunriseTime", func(t *testing.T) {
		x, err := resp.Result.SunriseTime()
		require.NoError(t, err)

		now := time.Now()
		s := time.Date(now.Year(), now.Month(), now.Day(), x.Hour(), x.Minute(), x.Second(), 0, time.UTC)
		require.Equal(t, s, x)
	})
}
