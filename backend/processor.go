package backend

import (
	"log"
	"os/exec"
	"sync"
)

func RunProcessor(output chan<- string, command []string) {
	data := []string{"--experimental-json-modules", "./node_modules/festuff/cli.mjs"}
	data = append(data, command...)
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
