package main

import (
	"encoding/json"
	"fmt"
)

type EffectJson struct {
	InPath  string   `json:"inPath"`
	OutPath string   `json:"outPath"`
	Effects []string `json:"effects"`
}

type ImageTask struct {
	Effects     []string
	ImageMatrix [][]int
	RowNumber   int
}

type ImageTaskOP struct {
	Effects     []string
	ImageMatrix [1][3]int
	RowNumber   int
}

var input_arr = [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
var kernel = [][]int{{-1, -2, -1}, {0, 0, 0}, {1, 2, 1}}
var lst = [1]string{"B"}

var num_threads = 3

func conv(rows int, cols int, input [][]int, out [1][3]int) [1][3]int {
	var i, j int
	for i = 0; i < rows; i++ {
		sum := 0
		for j = 0; j < cols; j++ {
			sum += input[i][j] * 2
		}
		out[i][j-1] = sum
	}
	return out
}

func main() {
	imageTaskGenerator := func() <-chan interface{} {
		// read from effects.txt
		// var jsonStr EffectJson
		taskChannel := make(chan interface{})
		go func() {
			defer close(taskChannel)
			// fmt.Println(jsonStr)
			for i := 0; i < 3; i++ {
				temp := ImageTask{[]string{"R"}, [][]int{input_arr[i]}, i}
				tempJson, _ := json.Marshal(temp)
				taskChannel <- string(tempJson)
			}
		}()
		return taskChannel
	}

	workerTask := func(
		done <-chan interface{},
		imageStream <-chan interface{},
	) <-chan interface{} {
		workerStream := make(chan interface{})
		go func() {
			defer close(workerStream)
			for task := range imageStream {
				// fmt.Println(task)
				var inputTask ImageTask
				if str, ok := task.(string); ok {
					/* act on str */
					err := json.Unmarshal([]byte(str), &inputTask)

					if err != nil {
						fmt.Println(err)
					}
				} else {
					fmt.Println("NHP")
				}

				var inputMatrix [][]int = inputTask.ImageMatrix

				var outputMatrix [1][3]int
				outputMatrix = conv(1, 3, inputMatrix, outputMatrix)
				outputImageTask := ImageTaskOP{[]string{"R"}, outputMatrix, -1}
				outputJson, _ := json.Marshal(outputImageTask)
				select {
				case <-done:
					return
				case workerStream <- string(outputJson):
				}
			}
		}()
		return workerStream
	}

	taskChannel := imageTaskGenerator()
	// for val := range taskChannel {
	// 	fmt.Println(val)
	// }

	workers := make([]<-chan interface{}, num_threads)
	done := make(chan interface{})
	defer close(done)

	for i := 0; i < num_threads; i++ {
		workers[i] = workerTask(done, taskChannel)
		fmt.Println("something2")
		for val := range workers[i] {
			fmt.Println(val)
		}

	}

}
