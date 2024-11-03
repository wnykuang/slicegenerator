//Read the video file and cut the video based on the start and end time

package video

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
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

	err := ffmpeg_go.Input(v.VideoPath, ffmpeg_go.KwArgs{"ss": v.StartTime, "t": v.EndTime}).
		Output(v.ExportPath).OverWriteOutput().WithErrorOutput(&stderr).Run()
	if err != nil {
		return fmt.Errorf("Error in cutting video: %v", stderr.String())
	}

	log.Printf("Video cut successfully saved to " + v.ExportPath)
	return nil
}

func (v *VideoClip) generateScreenShots() ([]image.Image, error) {

	// var screenshots []image.Image
	return nil, nil

}

func (v *VideoClip) CountFrames() (int64, error) {
	output, err := ffmpeg_go.Probe(v.ExportPath, ffmpeg_go.KwArgs{"show_entries": "stream=nb_frames", "select_streams": "v"})
	fmt.Println(output, err)
	return 0, nil
}
