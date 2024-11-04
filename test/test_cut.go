package main

import (
	"fmt"
	"image"
	"log"

	"github.com/wnykuang/slicegenerator/pkg/clip"
	"github.com/wnykuang/slicegenerator/pkg/util"
	"github.com/wnykuang/slicegenerator/pkg/video"
)

func main() {
	videoClip := video.NewVideoClip("./example.mp4", "00:01:12", "00:01:17", "")
	videoClip.Print()
	err := videoClip.Cut()
	// str, err := videoClip.CountFrames()

	frames, err := videoClip.GenerateScreenShots()
	if err != nil {
		log.Fatal(err)
	}

	util.SaveImages(frames, "./output/oringinal/")
	c := clip.NewClip(frames)
	different_frames := c.GetDifferentFrames()
	fmt.Printf("length of different frames: %d\n", len(different_frames))

	util.SaveImages(different_frames, "./output/different/")

	//debug use, corp and binarizate the images

	binazied_frames := []image.Image{}
	debug_frames := frames
	for i := 0; i < len(debug_frames); i++ {
		//get the lower 1/5 of the image
		bounds := debug_frames[i].Bounds()
		width := bounds.Max.X
		height := bounds.Max.Y
		lower_height := height / 5
		rect := image.Rect(0, 4*lower_height, width, height)
		lower_frame := clip.CopySubImage(debug_frames[i], rect)
		binazed_frame := clip.BinarizateImage(lower_frame, 200)
		// util.SaveImages([]image.Image{binazed_frame}, "./output/binarizated_frame/")
		binazied_frames = append(binazied_frames, binazed_frame)
	}
	util.SaveImages(binazied_frames, "./output/binarizated_frame/")
	// read image1 from hdd

	// image1_path := "./output/frame68.jpg"
	// file, err := os.Open(image1_path)
	// image1, err := jpeg.Decode(file)

	// image2_path := "./output/frame71.jpg"
	// file, err = os.Open(image2_path)
	// image2, err := jpeg.Decode(file)

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// clip.IsFrameDifferent(image1, image2)
	// fmt.Printf("length of different frames: %d\n", len(different_frames))
	// fmt.Printf("length of frames: %d\n", len(frames))
}
