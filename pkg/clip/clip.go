package clip

import (
	"fmt"
	"image"
	"image/draw"
	"math"
)

type clip struct {
	screenshots []image.Image
}

func NewClip(screenshot []image.Image) *clip {
	return &clip{
		screenshots: screenshot,
	}
}

func (c *clip) GetDifferentFrames() []image.Image {
	//compare the lower 1/5 of the image and if they are similar, then treat them as the same frame
	//since the subtitles are always in the lower 1/5 of the image

	var different_frames []image.Image
	var last_index = 0
	for i := 0; i < len(c.screenshots); i++ {
		if isFrameAllBlack(c.screenshots[i]) {
			last_index = i
			continue
		}
		fmt.Println("Comparing frame: ", i, "with frame: ", last_index)
		if IsFrameDifferent(c.screenshots[last_index], c.screenshots[i]) {
			different_frames = append(different_frames, c.screenshots[i])
			last_index = i
		}
	}

	last_frame := c.screenshots[len(c.screenshots)-1]

	if IsFrameDifferent(c.screenshots[last_index], last_frame) {
		different_frames = append(different_frames, c.screenshots[last_index])
	}

	//append the last frame
	return different_frames
}

func IsFrameDifferent(frame1, frame2 image.Image) bool {
	//get the lower 1/5 of the image, and compare them
	//since they are from same video, the subtitles should be in the same position
	//and the size of the frames are the same
	bounds := frame1.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	//get the lower 1/5 of each image
	lower_height := height / 5
	rect := image.Rect(0, 4*lower_height, width, height)

	lower_frame1 := CopySubImage(frame1, rect)
	lower_frame2 := CopySubImage(frame2, rect)
	// lower_frames := []image.Image{lower_frame1, lower_frame2}
	// util.SaveImages(lower_frames, "./output/lower/")

	//binarizate the images

	binarizeThreshold := 200

	binazed_frame1 := BinarizateImage(lower_frame1, float32(binarizeThreshold))
	binazed_frame2 := BinarizateImage(lower_frame2, float32(binarizeThreshold))
	similarity := CalculateSimilarity(binazed_frame1, binazed_frame2)
	// binazed_frames := []image.Image{binazed_frame1, binazed_frame2}
	// err := util.SaveImages(binazed_frames, "./output/binarizated_frame/")
	// if err != nil {
	// return false
	// }

	if similarity > 90 {
		return false
	}
	return true
}

func isFrameAllBlack(frame image.Image) bool {

	bounds := frame.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	lower_height := height / 5
	rect := image.Rect(0, 4*lower_height, width, height)

	lower_frame1 := CopySubImage(frame, rect)
	binazed_frame := BinarizateImage(lower_frame1, 200)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, g, b, _ := binazed_frame.At(x, y).RGBA()
			if r != 0 || g != 0 || b != 0 {
				return false
			}
		}
	}
	return true
}

func CopySubImage(img image.Image, rect image.Rectangle) image.Image {

	subImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(rect)

	// Create a new image with the dimensions of the sub-image
	newImg := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))

	// Draw the sub-image onto the new image
	draw.Draw(newImg, newImg.Bounds(), subImg, rect.Min, draw.Src)

	return newImg
}

func BinarizateImage(img image.Image, threshold float32) image.Image {
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y
	//create a new image
	new_img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			//get the pixel
			r, g, b, _ := img.At(x, y).RGBA()
			//convert to gray scale
			lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			gray := int(lum / 256)

			//binarizate the image
			if gray > int(threshold) {
				new_img.Set(x, y, image.White)
			} else {
				new_img.Set(x, y, image.Black)
			}
		}
	}
	return new_img
}

func CalculateSimilarity(frame1, frame2 image.Image) int {
	//let me start from naive approach.
	width := frame1.Bounds().Max.X
	height := frame1.Bounds().Max.Y

	start_i_1 := 0
	start_i_2 := 0

	for i := 0; i < width; i++ {
		white_count_1 := 0
		white_count_2 := 0

		for j := 0; j < height; j++ {
			r1, g1, b1, _ := frame1.At(i, j).RGBA()
			r2, g2, b2, _ := frame2.At(i, j).RGBA()

			if r1 != 0 && g1 != 0 && b1 != 0 {
				white_count_1++
			}

			if r2 != 0 && g2 != 0 && b2 != 0 {
				white_count_2++
			}
		}
		if start_i_1 == 0 && float64(white_count_1)/float64(height) > 0.1 {
			start_i_1 = i
		}

		if start_i_2 == 0 && float64(white_count_2)/float64(height) > 0.1 {
			start_i_2 = i
		}
	}

	fmt.Println("start_i_1: ", start_i_1, "start_i_2: ", start_i_2)
	if start_i_1*start_i_2 == 0 {
		return 0
	}

	if math.Abs(float64(start_i_1-start_i_2)) < float64(width)/20 {
		return 100
	}

	return 0
}
