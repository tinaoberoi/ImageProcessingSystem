package scheduler

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"proj2/png"
	"strings"
)

type bspWorkerContext struct {
	// Define the necessary fields for your implementation
	workerCount int
	masterId    int
	sem         *Semaphore
	barrier     Barrier
	data        *png.ImageTask
	effect      string
	dirs        []string
	imgArr      []*png.ImageTask
}

func NewBSPContext(config Config) *bspWorkerContext {

	dirs := strings.Split(config.DataDirs, "+")
	arr := make([]*png.ImageTask, 0)
	//Initialize the context
	return &bspWorkerContext{
		config.ThreadCount - 1,
		config.ThreadCount - 1,
		NewSemaphore(0),
		*NewBarrier(config.ThreadCount),
		nil,
		"",
		dirs,
		arr,
	}
}

func threadBounds2(n int, bounds image.Rectangle, threadCount int) *image.Rectangle {
	var tBounds image.Rectangle

	if n == threadCount-1 {
		tBounds = image.Rect(
			bounds.Min.X,
			bounds.Min.Y+((threadCount-1)*(bounds.Max.Y/threadCount)),
			bounds.Max.X,
			bounds.Max.Y,
		)
	} else {
		tBounds = image.Rect(
			bounds.Min.X,
			bounds.Min.Y+(n*(bounds.Max.Y/threadCount)),
			bounds.Max.X,
			bounds.Min.Y+((n+1)*(bounds.Max.Y/threadCount)),
		)
	}

	return &tBounds
}

func RunBSPWorker(id int, ctx *bspWorkerContext) {
	for {
		if id == ctx.masterId {
			// To loop over every dir
			for _, dir := range ctx.dirs {
				effectsPathFile := fmt.Sprintf("../../data/effects.txt")
				effectsFile, _ := os.Open(effectsPathFile)
				reader := json.NewDecoder(effectsFile)
				var jsonStr png.ImageMetaData

				// To read effects.txt file
				for {
					e := reader.Decode(&jsonStr)
					if e != nil {
						break
					}

					pngImg, err := png.Load("../../data/in/"+dir+"/"+jsonStr.InPath,
						"../../data/out/"+dir+"_"+jsonStr.OutPath,
						jsonStr.Effects)

					if err != nil {
						panic(err)
					}
					ctx.data = pngImg

					for j := 0; j < len(jsonStr.Effects); j++ {
						ctx.effect = jsonStr.Effects[j]
						ctx.sem.Up(ctx.workerCount)
						ctx.barrier.Arrive()
						if j < len(jsonStr.Effects)-1 {
							pngImg.Swap()
						}
					}
					ctx.imgArr = append(ctx.imgArr, pngImg)
				}

				saveImgBarrier := NewBarrier(len(ctx.imgArr) + 1)
				for i := 0; i < len(ctx.imgArr); i++ {
					go func(pngImg *png.ImageTask) {
						err := pngImg.Save(pngImg.OutPath)

						if err != nil {
							panic(err)
						}
						saveImgBarrier.Arrive()
					}(ctx.imgArr[i])
				}
				saveImgBarrier.Arrive()
			}
			break
		} else {
			ctx.sem.Down()
			var tBounds *image.Rectangle = threadBounds2(id, ctx.data.Bounds, ctx.workerCount)
			switch ctx.effect {
			case "G":
				ctx.data.Grayscale(tBounds)
			case "E":
				ctx.data.EdgeDetection(tBounds)
			case "S":
				ctx.data.Sharpen(tBounds)
			case "B":
				ctx.data.Blur(tBounds)
			default:
				panic("No effect found!")
			}
			ctx.barrier.Arrive()
		}
	}
}
