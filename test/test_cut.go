package main

import (
	"fmt"
	"log"

	"github.com/wnykuang/slicegenerator/pkg/video"
)

func main() {
	videoClip := video.NewVideoClip("./example.mp4", "00:00:00", "00:00:5", "")
	videoClip.Print()
	err := videoClip.Cut()
	str, err := videoClip.CountFrames()
	fmt.Println("Frames: ", str)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Video cut successfully")
}
