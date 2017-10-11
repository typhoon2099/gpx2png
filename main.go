package main

import (
    "github.com/tkrajina/gpxgo/gpx"
    "flag"
    "fmt"
    "os"
    "io/ioutil"
)

type Point struct {
    x float64
    y float64
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    var filename string
    var points = make([]Point, 0)

    // Get the flags
    flag.StringVar(&filename, "i", "", "Filename")
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
    for _, track := range gpxFile.Tracks {
        for _, segment := range track.Segments {
            for _, point := range segment.Points {
                points = append(points, Point{
                    x: point.Latitude,
                    y: point.Longitude,
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

    fmt.Print(points)

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

    fmt.Print(maxX)
    fmt.Print(maxY)
}
