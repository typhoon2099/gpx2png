package main

import (
    "github.com/tkrajina/gpxgo/gpx"
    "flag"
    "fmt"
    "os"
    "io/ioutil"
    "math"
    "github.com/llgcode/draw2d/draw2dimg"
    "image"
    "image/color"
)

type Point struct {
    latitude  float64
    longitude float64
    x         float64
    y         float64
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
    var thickness int
    var outline bool
    var points = make([]Point, 0)

    // Get the flags
    flag.StringVar(&filename, "i", "", "Filename")
    flag.IntVar(&resolutionX, "width", 1920, "Width")
    flag.IntVar(&resolutionY, "height", 1080, "Height")
    flag.IntVar(&thickness, "thickness", 5, "Thickness")
    flag.BoolVar(&outline, "outline", true, "Outline")
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
    // x has some magic applied to get the true width based on latitude
    // y is flipped to match the image output
    for _, track := range gpxFile.Tracks {
        for _, segment := range track.Segments {
            for _, point := range segment.Points {
                points = append(points, Point{
                    latitude: point.Latitude,
                    longitude: point.Longitude,
                    x: math.Cos(point.Latitude * (math.Pi / 180)) * point.Longitude,
                    y: point.Latitude * -1,
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

    // Figure out what scale to apply
    var scaleX = float64(resolutionX) / maxX
    var scaleY = float64(resolutionY) / maxY

    var finalScale = scaleX

    if (scaleY < finalScale) {
        finalScale = scaleY
    }

    // Now we know what scale to use, figure out how much translation need to be done in each direction
    var translateX = (float64(resolutionX) - (maxX * finalScale)) / 2
    var translateY = (float64(resolutionY) - (maxY * finalScale)) / 2

    // Loop through and apply the scaling and translation
    for i := 0; i < len(points); i++ {
        points[i].x = points[i].x * finalScale + translateX
        points[i].y = points[i].y * finalScale + translateY
    }

    dest := image.NewRGBA(image.Rect(0, 0, resolutionX, resolutionY))
    gc := draw2dimg.NewGraphicContext(dest)

    if (outline == true) {
        drawPoints(gc, points, float64(thickness * 2), color.RGBA{0x00, 0x00, 0x00, 0xff})
    }

    drawPoints(gc, points, float64(thickness), color.RGBA{0xff, 0x44, 0xff, 0xff})

    // Save to file
    draw2dimg.SaveToPngFile(filename + ".png", dest)
}

func drawPoints(image *draw2dimg.GraphicContext, points []Point, thickness float64, color color.RGBA) {
    // Set some properties
    image.SetStrokeColor(color)
    image.SetLineWidth(float64(thickness))

    // Move to the first point
    image.MoveTo(float64(points[0].x), float64(points[0].y))

    for _, point := range points {
        image.LineTo(float64(point.x), float64(point.y))
    }

    // Finish drawing the line
    image.Stroke()
}
