package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"proj2/png"
	"strings"
)

func RunSequential(config Config) {

	dirs := strings.Split(config.DataDirs, "+")
	for d := 0; d < len(dirs); d++ {
		effectsPathFile := fmt.Sprintf("../../data/effects.txt")
		effectsFile, _ := os.Open(effectsPathFile)
		var jsonStr png.ImageMetaData
		reader := json.NewDecoder(effectsFile)
		var inputImgDir string = "../../data/in/" + dirs[d] + "/"
		var outputImgDir string = "../../data/out/" + dirs[d] + "_"

		for {
			e := reader.Decode(&jsonStr)
			if e != nil {
				break
			}

			pngImg, err := png.Load(inputImgDir+jsonStr.InPath,
				outputImgDir+jsonStr.OutPath,
				jsonStr.Effects)

			if err != nil {
				panic(err)
			}

			for i := 0; i < len(jsonStr.Effects); i++ {

				switch mode := jsonStr.Effects[i]; mode {
				case "G":
					pngImg.Grayscale(&pngImg.Bounds)
				case "E":
					pngImg.EdgeDetection(&pngImg.Bounds)
				case "B":
					pngImg.Blur(&pngImg.Bounds)
				case "S":
					pngImg.Sharpen(&pngImg.Bounds)
				default:
					fmt.Printf("Undefined mode \n")
				}
				if i < len(jsonStr.Effects)-1 {
					pngImg.Swap()
				}
			}

			err = pngImg.Save(outputImgDir + jsonStr.OutPath)
			if err != nil {
				panic(err)
			}
		}
	}
}
