package repository

import (
	"GoogleImageDownloader/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type PixelSearchEngine struct {
	config Config
}

type Config struct {
	ApiKey string     `koanf:apiKey`
	Min    int        `koanf:min`
	File   FileConfig `koanf:file`
}

type FileConfig struct {
	Prefix string `koanf:prefix`
	Type   string `koanf:type`
	Path   string `koanf:path`
}

func New(config Config) SearchEngine {
	return PixelSearchEngine{
		config: config,
	}
}

type PexelsResponse struct {
	TotalResults int     `json:"total_results"`
	Page         int     `json:"page"`
	PerPage      int     `json:"per_page"`
	Photos       []Photo `json:"photos"`
	NextPage     string  `json:"next_page"`
}

type Photo struct {
	ID              int    `json:"id"`
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	URL             string `json:"url"`
	Photographer    string `json:"photographer"`
	PhotographerURL string `json:"photographer_url"`
	PhotographerID  int    `json:"photographer_id"`
	AvgColor        string `json:"avg_color"`
	Src             Src    `json:"src"`
	Liked           bool   `json:"liked"`
	Alt             string `json:"alt"`
	Path            string
}

type Src struct {
	Original  string `json:"original"`
	Large2x   string `json:"large2x"`
	Large     string `json:"large"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Portrait  string `json:"portrait"`
	Landscape string `json:"landscape"`
	Tiny      string `json:"tiny"`
}

func (svc PixelSearchEngine) SearchImages(query string, page, perPage int32, wg *sync.WaitGroup) ([]model.ImageResult, error) {
	url := fmt.Sprintf("https://api.pexels.com/v1/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if len(svc.config.ApiKey) == 0 {
		return nil, fmt.Errorf("please provide API_KEY")
	}
	req.Header.Add("Authorization", svc.config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed getImages for %s: status code %d\n", query, resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse PexelsResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	for _, photo := range apiResponse.Photos {
		photo.Path = fmt.Sprintf("%s%s_%d%s", svc.config.File.Path, query, photo.ID, svc.config.File.Type)
		wg.Add(1)
		go func(photoPath, photoURL string) {
			defer wg.Done()
			wg.Add(1)
			go DownloadImage(photoURL, photoPath, wg)
		}(photo.Path, photo.Src.Medium)

	}

	return apiResponse.pixelResponseToImageModel(), nil
}

func (service PexelsResponse) pixelResponseToImageModel() []model.ImageResult {
	results := make([]model.ImageResult, len(service.Photos))
	for i, photo := range service.Photos {
		results[i].Url = photo.Src.Medium
		results[i].Id = photo.ID
		results[i].Data = photo.Path
	}
	return results
}
