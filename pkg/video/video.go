//Read the video file and cut the video based on the start and end time

package video

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type VideoClip struct {
	VideoPath  string
	StartTime  string
	EndTime    string
	ExportPath string
}

func NewVideoClip(videoPath, startTime, endTime, exportPath string) *VideoClip {
	if exportPath == "" {
		exportPath = "./output/"
	}

	videoName := strings.TrimSuffix(filepath.Base(videoPath), filepath.Ext(videoPath))
	exportPath = filepath.Join(exportPath, videoName+".mp4")

	return &VideoClip{
		VideoPath:  videoPath,
		StartTime:  startTime,
		EndTime:    endTime,
		ExportPath: exportPath,
	}
}

// print out the videoclip information to debug
func (v *VideoClip) Print() {
	fmt.Println("Video Path: ", v.VideoPath, "Start Time: ", v.StartTime, "End Time: ", v.EndTime)
}

func (v *VideoClip) Cut() error {
	var stderr bytes.Buffer

	//create the output folder if it does not exist
	outputFolder := filepath.Dir(v.ExportPath)
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		os.MkdirAll(outputFolder, os.ModePerm)
	}

	//if the file is exited, rename the generated file with a suffix with timestamp
	timestamp := time.Now().Format("2006-01-02_15:04:05")
	if _, err := os.Stat(v.ExportPath); err == nil {
		v.ExportPath = strings.TrimSuffix(v.ExportPath, ".mp4") + "_" + timestamp + ".mp4"
	}

	err := ffmpeg_go.Input(v.VideoPath, ffmpeg_go.KwArgs{"ss": v.StartTime, "to": v.EndTime}).
		Output(v.ExportPath).OverWriteOutput().WithErrorOutput(&stderr).Run()
	if err != nil {
		return fmt.Errorf("Error in cutting video: %v", stderr.String())
	}

	log.Printf("Video cut successfully saved to " + v.ExportPath)
	return nil
}

func (v *VideoClip) GenerateScreenShots() ([]image.Image, error) {

	//get the number of frames in the video
	frameCount, err := v.CountFrames()
	if err != nil {
		return nil, err
	}
	var images []image.Image

	//load the frames into RAM
	for i := 0; i < int(frameCount); i++ {
		// buf := new(bytes.Buffer)
		var buf bytes.Buffer

		//get the ith frame
		var stderr bytes.Buffer
		err := ffmpeg_go.Input(v.ExportPath).
			Filter("select", ffmpeg_go.Args{fmt.Sprintf("eq(n, %d)", i)}).
			Output("pipe:", ffmpeg_go.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
			WithOutput(&buf).
			WithErrorOutput(&stderr).
			Run()

		if err != nil {
			// panic(err + "detail is")
			log.Fatal("Error in generating screenshot: ", stderr.String())
		}

		img, err := jpeg.Decode(&buf)
		if err != nil {
			panic(err)
		}

		images = append(images, img)
	}
	return images, nil
}

func (v *VideoClip) CountFrames() (int64, error) {
	output, err := ffmpeg_go.Probe(v.ExportPath, ffmpeg_go.KwArgs{"show_entries": "stream=nb_frames", "select_streams": "v"})

	if err != nil {
		return 0, err
	}

	//parse the output to get the number of frames
	var output_res map[string]interface{}
	json_err := json.Unmarshal([]byte(output), &output_res)

	if json_err != nil {
		return 0, json_err
	}

	frame_count_st := output_res["streams"].([]interface{})[0].(map[string]interface{})["nb_frames"].(string)
	frame_count, err := strconv.ParseInt(frame_count_st, 10, 64)
	return frame_count, nil
}
