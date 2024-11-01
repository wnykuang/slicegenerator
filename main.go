// import library about cli interface
package main

import (
	"flag"
	"fmt"
	"os"
)

// display the cli interface
// read the flag from the command line
func main() {

	args := os.Args[1:]
	video_path := flag.String("path", "", "Path to the video file")
	start_time := flag.String("start", "", "Start time of the video")
	end_time := flag.String("end", "", "End time of the video")

	flag.Parse()

	fmt.Println("Video Path: ", *video_path, "Start Time: ", *start_time, "End Time: ", *end_time)
	fmt.Println("Arguments: ", args)
}

// Run the program
// go run main.go
// Hello, World!
