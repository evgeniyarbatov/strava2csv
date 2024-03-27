# strava2csv

Convert Strava archive download to CSV file

## run

```
go run main.go \
../gpx-data/data/strava/ \
~/Downloads/strava2gpx/gpx-file.csv
```

## perf

I have 1846 GPX / TCX files exported from Strava.

It takes 111.17s to process all of them.

The resulting CSV file is 416MB.
