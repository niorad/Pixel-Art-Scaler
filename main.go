package main

import (
	"fmt"
	"html/template"
	"image/gif"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"pixelartscaler/giffactory"
	"pixelartscaler/processing"
	"strconv"
)

var templates = template.Must(template.ParseFiles("./templates/index.html"))

type page struct {
	ErrorMessage string
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	scalingIterationCount, _ := strconv.Atoi(r.FormValue("scalingIterationCount"))
	frameDelay, _ := strconv.Atoi(r.FormValue("frameDelay"))
	frameCount, _ := strconv.Atoi(r.FormValue("frameCount"))
	scalingType := r.FormValue("scalingType")

	if handler.Header.Get("Content-Type") != "image/png" {
		form := page{ErrorMessage: "This is not a PNG file how am I.. what..."}
		templates.ExecuteTemplate(w, "index.html", form)
		return
	}

	if handler.Size > 500000 {
		form := page{ErrorMessage: "Too big a file pal. Keep them below 500 kb."}
		templates.ExecuteTemplate(w, "index.html", form)
		return
	}

	tempDir, err := ioutil.TempDir(".", "_gen")
	if err != nil {
		fmt.Println("Error creating temp-dir")
		fmt.Println(err)
	}

	tempFile, err := ioutil.TempFile(tempDir, "upload-*.png")
	if err != nil {
		fmt.Println("Error Creating Temp File for Upload")
		fmt.Println(err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading from form-file")
		fmt.Println(err)
	}

	_, err = tempFile.Write(fileContents)
	if err != nil {
		fmt.Println("Error writing fileContents into temnpFile")
		fmt.Println(err)
	}

	tempFile.Seek(0, 0)
	imageFile, err := png.Decode(tempFile)
	if err != nil {
		fmt.Println("Error decoding into png-file")
		fmt.Println(err)
	}

	processedImage := imageFile
	var processedGif gif.GIF

	if scalingType == "basic" {

		for i := 0; i < scalingIterationCount; i++ {
			processedImage = processing.BasicScaling(processedImage, false)
		}

	} else if scalingType == "randombasic" {

		for i := 0; i < scalingIterationCount; i++ {
			processedImage = processing.BasicScaling(processedImage, true)
		}

	} else if scalingType == "randombasicanim" {

		processedGif = giffactory.Generate(processedImage, scalingIterationCount, frameDelay, frameCount)

	} else if scalingType == "nn" {

		for i := 0; i < scalingIterationCount; i++ {
			processedImage = processing.NearestNeighbor(processedImage)
		}
	}

	var filenamePattern = "processed-*.png"
	if scalingType == "randombasicanim" {
		filenamePattern = "processed-*.gif"
	}

	tempResponseFile, _ := ioutil.TempFile(tempDir, filenamePattern)

	if scalingType == "randombasicanim" {
		gif.EncodeAll(tempResponseFile, &processedGif)
	} else {
		png.Encode(tempResponseFile, processedImage)
	}
	http.ServeFile(w, r, tempResponseFile.Name())

	file.Close()
	tempFile.Close()
	os.RemoveAll(tempDir)
}

func setupRoutes() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("http://localhost:" + port)
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/", getIndex)
	http.ListenAndServe(":"+port, nil)
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	form := page{ErrorMessage: ""}
	err := templates.ExecuteTemplate(w, "index.html", form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	fmt.Println("Starting PixelArtScaler-Server…")
	setupRoutes()
}
