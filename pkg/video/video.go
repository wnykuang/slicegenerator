//Read the video file and cut the video based on the start and end time

package video

import (
	"fmt"
	"path/filepath"
	"strings"

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

func (v *VideoClip) Cut() err {
	path := v.VideoPath
	exportPath := strings.Split(v.ExportPath, "")
	err := ffmpeg_go.Input(v.VideoPath, ffmpeg_go.KwArgs{"ss": v.StartTime, "t": v.EndTime}).Output(v.ExportPath).OverWriteOutput().Run()
}
