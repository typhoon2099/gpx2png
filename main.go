package main

import (
    "github.com/tkrajina/gpxgo/gpx"
    "flag"
    "fmt"
    "os"
    "io/ioutil"
    "math"
)

type Point struct {
    latitude float64
    longitude float64
    x int
    y int
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    var filename string
    var resolutionX int
    var resolutionY int
    var points = make([]Point, 0)

    // Get the flags
    flag.StringVar(&filename, "i", "", "Filename")
    flag.IntVar(&resolutionX, "width", 1920, "Width")
    flag.IntVar(&resolutionY, "height", 1080, "Height")
    flag.Parse()

    if filename == "" {
        fmt.Println("\nPlease provide an file to parse")
        flag.PrintDefaults()
        os.Exit(1)
    }

    // Attempt to load the file and then read it as GPX data
    gpxBytes, err := ioutil.ReadFile(filename)
    check(err)

    gpxFile, err := gpx.ParseBytes(gpxBytes)
    check(err)

    // Grab the points and stick them in a slice for processing
    // We'll also keep a copy of the original lat/longs, just in case
    // X and Y are decimal shifted and then cast to ints to preserve accuracy
    for _, track := range gpxFile.Tracks {
        for _, segment := range track.Segments {
            for _, point := range segment.Points {
                points = append(points, Point{
                    latitude: point.Latitude,
                    longitude: point.Longitude,
                    x: int(point.Latitude * 100000),
                    y: int((math.Cos(point.Latitude * (math.Pi / 180)) * point.Longitude) * 100000),
                })
            }
        }
    }

    // Loop through the points and get the min values for x and y
    // (I think there's a sexier way to do this, but I don't know it yet
    var minX = points[0].x
    var minY = points[0].y

    for _, point := range points {
        if point.x < minX {
            minX = point.x
        }
        if point.y < minY {
            minY = point.y
        }
    }

    // Now loop through and subtract the mins from each point
    for i := 0; i < len(points); i++ {
        points[i].x = points[i].x - minX
        points[i].y = points[i].y - minY
    }

    // Now find the maxes so we can increase up to our desired image size
    var maxX = points[0].x
    var maxY = points[0].y

    for _, point := range points {
        if point.x > maxX {
            maxX = point.x
        }
        if point.y > maxY {
            maxY = point.y
        }
    }

    var scaleX = float64(resolutionX) / float64(maxX)
    var scaleY = float64(resolutionY) / float64(maxY)

    var finalScale = scaleX

    if(scaleY < finalScale) {
        finalScale = scaleY
    }

    for _, point := range points {
        point.x = int(float64(point.x) * finalScale)
        point.y = int(float64(point.y) * finalScale)
    }
}
