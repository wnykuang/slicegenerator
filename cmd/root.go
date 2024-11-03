package cmd

import (
	"flag"

	"github.com/wnykuang/slicegenerator/pkg/video"
)

func Execute() {
	// fmt.Print("Enter the video path: ")

	// args := os.Args[1:]
	video_path := flag.String("path", "", "Path to the video file")
	start_time := flag.String("start", "", "Start time of the video")
	end_time := flag.String("end", "", "End time of the video")

	flag.Parse()

	// fmt.Println("Video Path: ", *video_path, "Start Time: ", *start_time, "End Time: ", *end_time)
	// fmt.Println("Arguments: ", args)

	video_clip := video.NewVideoClip(*video_path, *start_time, *end_time)
	video_clip.Print()

}
