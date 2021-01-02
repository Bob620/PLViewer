package backend

import (
	"log"
	"os/exec"
	"sync"
)

func RunCreator(output chan<- string, uri string, options *CreatorOptions) {
	data := []string{"./node_modules/xes-converter/index.js"}
	data = append(data, options.Serialize()...)
	data = append(data, uri)
	creator := exec.Command("./node.exe", data...)
	creator.Dir = "./node"

	var wg sync.WaitGroup

	stdout, err := creator.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := creator.Start(); err != nil {
		log.Fatal(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		line := ""
		for {
			buf := make([]byte, 100)
			_, err := stdout.Read(buf)
			line += string(buf)

			if line != "" {
				for i, rawChar := range line {
					char := string(rawChar)
					if char == "\n" {
						output <- line[:i]
						line = line[i+1:]
						break
					}
				}
			}
			if err != nil {
				break
			}
		}

		for line != "" {
			for i, rawChar := range line {
				char := string(rawChar)
				if char == "\n" {
					output <- line[:i]
					line = line[i+1:]
					break
				}
			}
		}
		output <- "\b"
	}()

	wg.Wait()
	if err := creator.Wait(); err != nil {
		log.Fatal(err)
	}
}
