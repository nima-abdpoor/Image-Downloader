package repository

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func DownloadImage(url, filepath string, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to download image from %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to download image from %s: status code %d\n", url, resp.StatusCode)
		return
	}

	file, err := os.Create(filepath)
	if err != nil {
		log.Printf("Failed to create file %s: %v\n", filepath, err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Printf("Failed to save image to file %s: %v\n", filepath, err)
	}

}
