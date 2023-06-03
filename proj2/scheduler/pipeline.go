package scheduler

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"proj2/png"
	"strconv"
	"strings"
)

func taskCreater(dir string) chan *png.ImageTask {

	taskChannel := make(chan *png.ImageTask)
	effectsPathFile := fmt.Sprintf("../../data/effects.txt")
	effectsFile, err := os.Open(effectsPathFile)

	if err != nil {
		panic(err)
	}

	go func() {
		defer close(taskChannel)

		reader := json.NewDecoder(effectsFile)

		for {
			var data png.ImageMetaData
			e := reader.Decode(&data)
			if e != nil {
				break
			}

			inputImgPath := "../../data/in/" + dir + "/" + data.InPath
			outputImgPath := "../../data/out/" + dir + "_" + data.OutPath
			imgTask, err := png.Load(inputImgPath, outputImgPath, data.Effects)

			if err != nil {
				panic(err)
			}

			taskChannel <- imgTask
		}
	}()
	return taskChannel
}

func threadBounds(n int, bounds image.Rectangle, threadCount int) *image.Rectangle {
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

func FanOut(taskChannel chan *png.ImageTask, threadCount int) chan *png.ImageTask {
	filteredImg := make(chan *png.ImageTask)

	go func() {
		for task := range taskChannel {
			for i := 0; i < len(task.Effects); i++ {
				wg := make(chan bool)
				effect := task.Effects[i]
				for n := 0; n < threadCount; n++ {
					go func(effect string, bounds image.Rectangle, n int, wg chan bool) {
						var tBounds *image.Rectangle = threadBounds(n, bounds, threadCount)

						switch effect {
						case "G":
							task.Grayscale(tBounds)
						case "E":
							task.EdgeDetection(tBounds)
						case "S":
							task.Sharpen(tBounds)
						case "B":
							task.Blur(tBounds)
						default:
							panic("No effect found!")
						}
						wg <- true
					}(effect, task.Bounds, n, wg)

				}
				for t := 0; t < threadCount; t++ {
					<-wg
				}

				if i != len(task.Effects)-1 {
					task.Swap()
				}
				task.Save("test" + strconv.Itoa(i) + ".png")
				filteredImg <- task
			}
		}
		close(filteredImg)
	}()
	return filteredImg
}

func ResultsAggregator(channels ...chan *png.ImageTask) {
	wg := make(chan bool)

	multiplex := func(c <-chan *png.ImageTask) {
		for i := range c {
			err := i.Save(i.OutPath)
			if err != nil {
				panic(err)
			}
		}
		wg <- true
	}

	for _, c := range channels {
		go multiplex(c)
	}

	for i := 0; i < len(channels); i++ {
		<-wg
	}
}

func RunPipeline(config Config) {

	dirs := strings.Split(config.DataDirs, "+")

	for _, dir := range dirs {
		taskChannel := taskCreater(dir)
		done := make(chan bool)
		defer close(done)
		workers := make([]chan *png.ImageTask, config.ThreadCount)
		for i := 0; i < config.ThreadCount; i++ {
			workers[i] = FanOut(taskChannel, config.ThreadCount)
		}

		ResultsAggregator(workers...)
	}

}
