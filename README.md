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
## output

```
Time,ActivityType,Filename,Latitude,Longitude,Elevation,Cadence,Heartrate,Power
2022-10-01T15:33:08.149Z,Running,8439099964.tcx.gz,1.34324833,103.83843833,6.248,0,104,348
2022-10-01T15:33:09.149Z,Running,8439099964.tcx.gz,1.343325,103.83849167,6.248,0,103,0
2022-10-01T15:33:10.149Z,Running,8439099964.tcx.gz,1.34342333,103.83851833,7.315,0,103,192
2022-10-01T15:33:11.149Z,Running,8439099964.tcx.gz,1.343435,103.838555,7.315,0,102,440
2022-10-01T15:33:12.149Z,Running,8439099964.tcx.gz,1.34348,103.838655,7.315,0,102,760
2022-10-01T15:33:13.149Z,Running,8439099964.tcx.gz,1.34352167,103.83875833,7.315,0,103,908
2022-10-01T15:33:14.149Z,Running,8439099964.tcx.gz,1.34355833,103.838855,7.315,0,103,802
2022-10-01T15:33:15.149Z,Running,8439099964.tcx.gz,1.34363833,103.83907667,7.315,0,103,570
2022-10-01T15:33:16.149Z,Running,8439099964.tcx.gz,0,0,7.315,0,103,282
2022-10-01T15:33:17.149Z,Running,8439099964.tcx.gz,0,0,7.315,0,103,142
```

## perf

I have 1846 GPX / TCX files exported from Strava.

It takes 111.17s to process all of them.

The resulting CSV file is 416MB.
