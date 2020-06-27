package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileInfo struct {
	FileName     string
	FilePath     string
	FileSize     int64
	Permission   os.FileMode
	LastModified time.Time
	IsDir        bool
}

var mu *sync.Mutex

func sendPostRequests(index int, dataChan chan FileInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	for data := range dataChan {
		// Send the http request to the server

		// fmt.Println(" ", data.FilePath)
		bytesData, err := json.Marshal(data)
		if err != nil {
			log.Fatalln(err)
		}

		resp, err := http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(bytesData))
		if err != nil {
			log.Fatalln(err)
		}

		var result map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&result)

		// mu.Lock()
		fmt.Println(index)
		log.Println(result)
		log.Println(result["data"])
		// mu.Unlock()
	}

}

func traverseFiles(filePath string, dataChan chan FileInfo) {

	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		currfileInfo := FileInfo{
			FileName:     info.Name(),
			FilePath:     path,
			FileSize:     info.Size(),
			Permission:   info.Mode(),
			LastModified: info.ModTime(),
			IsDir:        info.IsDir(),
		}

		dataChan <- currfileInfo

		return nil
	})

	close(dataChan)

	if err != nil {
		log.Println(err)
	}
}

func main() {
	filePath := flag.String("path", "C:/Users/bmaheshwari/go/src/github.com/bhuvnesh13396/Blog-Platform/sample-master/sample/comment", "Provide a file path")
	wg := new(sync.WaitGroup)
	dataChan := make(chan FileInfo)
	noOfWorkers := 5
	// mu := new(sync.Mutex)
	go traverseFiles(*filePath, dataChan)

	wg.Add(noOfWorkers)
	for i := 0; i < noOfWorkers; i++ {
		go sendPostRequests(i, dataChan, wg)
	}

	// for data := range dataChan {
	// 	fmt.Println(data.FileName)
	// }

	wg.Wait()

	// timeout := time.Duration(5 * time.Second)
	// client := http.Client{
	// 	Timeout: timeout,
	// }

	// request, err := http.NewRequest("POST", "http://localhost:8081", bytes.NewBuffer(data))
	// request.Header.Set("Content-Type", "application/json")

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// resp, err := client.Do(request)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// log.Println(string(body))

	// http.HandleFunc("/", HelloServer)
	// http.ListenAndServe(":8080", nil)

}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
