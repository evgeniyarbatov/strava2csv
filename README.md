# strava2csv

Convert Strava archive download to CSV file.

Use [bulk export](https://support.strava.com/hc/en-us/articles/216918437-Exporting-your-Data-and-Bulk-Export#h_01GG58HC4F1BGQ9PQZZVANN6WF) on Strava to download the archive.

Run this on `activities` directory.

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
