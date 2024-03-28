package main

import (
	"compress/gzip"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func main() {
	dir := os.Args[1]
	outputFile := os.Args[2]

	var mu sync.Mutex
	var wg sync.WaitGroup

	var extractedPoints []Point

	files, _ := os.ReadDir(dir)
	for _, file := range files {
		wg.Add(1)
		go ParseFiles(
			file,
			dir,
			&extractedPoints,
			&mu,
			&wg,
		)
	}

	wg.Wait()

	WriteCSV(
		extractedPoints,
		outputFile,
	)
}

func ParseFiles(
	file fs.DirEntry,
	dir string,
	extractedPoints *[]Point,
	mu *sync.Mutex,
	wg *sync.WaitGroup,
) {

	defer wg.Done()

	fileName := file.Name()
	filePath := dir + "/" + fileName

	var xmlData []byte

	extension := filepath.Ext(fileName)
	is_compressed := extension == ".gz"

	if is_compressed {
		xmlData = ExtractFile(filePath)
		extension = filepath.Ext(
			strings.TrimSuffix(filePath, extension),
		)
	} else {
		file, _ := os.Open(filePath)
		xmlData, _ = io.ReadAll(file)
	}

	switch extension {
	case ".gpx":
		var gpxData GPX
		xml.Unmarshal(xmlData, &gpxData)

		mu.Lock()
		defer mu.Unlock()
		for _, point := range gpxData.Trk.Trkseg.Trkpt {
			*extractedPoints = append(
				*extractedPoints,
				Point{
					Time:      point.Time,
					Sport:     gpxData.Trk.Type,
					Filename:  fileName,
					Latitude:  point.Lat,
					Longitude: point.Lon,
					Elevation: point.Ele,
					Cadence:   point.Cad,
					Heartrate: point.Hr,
					Power:     point.Pwr,
				},
			)
		}
	case ".tcx":
		var tcxData TCX
		xml.Unmarshal(xmlData, &tcxData)

		for _, activity := range tcxData.Activities {
			mu.Lock()
			defer mu.Unlock()
			for _, lap := range activity.Lap {
				for _, point := range lap.Track {
					*extractedPoints = append(
						*extractedPoints,
						Point{
							Time:      point.Time,
							Sport:     activity.Sport,
							Filename:  fileName,
							Latitude:  point.Position.Latitude,
							Longitude: point.Position.Longitude,
							Elevation: point.Ele,
							Cadence:   point.Cad,
							Heartrate: point.Hr,
							Power:     point.Pwr,
						},
					)
				}
			}
		}
	}
}

func ExtractFile(filePath string) []byte {
	gzFile, _ := os.Open(filePath)
	defer gzFile.Close()

	gzReader, _ := gzip.NewReader(gzFile)
	defer gzReader.Close()

	xmlData, _ := io.ReadAll(gzReader)
	return xmlData
}

func FloatToString(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func WriteCSV(
	points []Point,
	filePath string,
) {
	outputDir := filepath.Dir(filePath)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Println("Error creating output directory:", err)
			return
		}
	}

	outputFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	for _, point := range points {
		entry := []string{
			point.Time,
			point.Sport,
			point.Filename,
			FloatToString(point.Latitude),
			FloatToString(point.Longitude),
			FloatToString(point.Elevation),
			strconv.Itoa(point.Cadence),
			strconv.Itoa(point.Heartrate),
			strconv.Itoa(point.Power),
		}

		writer.Write(entry)
	}
}

type GPX struct {
	Metadata Metadata `xml:"metadata"`
	Trk      Trk      `xml:"trk"`
}

type Metadata struct {
	Time string `xml:"time,omitempty"`
}

type Trk struct {
	Name   string `xml:"name,omitempty"`
	Trkseg Trkseg `xml:"trkseg"`
	Type   string `xml:"type,omitempty"`
}

type Trkseg struct {
	Trkpt []Trkpt `xml:"trkpt"`
}

type Trkpt struct {
	Ele  float64 `xml:"ele,omitempty"`
	Lat  float64 `xml:"lat,attr"`
	Lon  float64 `xml:"lon,attr"`
	Time string  `xml:"time,omitempty"`
	Hr   int     `xml:"extensions>gpxtpx:TrackPointExtension>gpxtpx:hr,omitempty"`
	Cad  int     `xml:"extensions>gpxtpx:TrackPointExtension>gpxtpx:cad,omitempty"`
	Pwr  int     `xml:"extensions>power,omitempty"`
}

type TCX struct {
	XMLName    xml.Name   `xml:"TrainingCenterDatabase"`
	Activities []Activity `xml:"Activities>Activity"`
}

type Activity struct {
	XMLName xml.Name `xml:"Activity"`
	ID      string   `xml:"Id"`
	Sport   string   `xml:"Sport,attr"`
	Lap     []Lap    `xml:"Lap"`
}

type Lap struct {
	StartTime string    `xml:"StartTime,attr"`
	Track     []TrackPt `xml:"Track>Trackpoint"`
}

type TrackPt struct {
	Ele      float64 `xml:"AltitudeMeters"`
	Time     string  `xml:"Time"`
	Position Pos     `xml:"Position"`
	Hr       int     `xml:"HeartRateBpm>Value"`
	Cad      int     `xml:"Cadence"`
	Pwr      int     `xml:"Extensions>TPX>Watts,omitempty"`
}

type Pos struct {
	Latitude  float64 `xml:"LatitudeDegrees"`
	Longitude float64 `xml:"LongitudeDegrees"`
}

type Point struct {
	Time      string
	Sport     string
	Filename  string
	Latitude  float64
	Longitude float64
	Elevation float64
	Cadence   int
	Heartrate int
	Power     int
}
