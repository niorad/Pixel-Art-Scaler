package main

import (
	"fmt"
	"html/template"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
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

	if handler.Header.Get("Content-Type") != "image/png" {
		form := page{ErrorMessage: "This is not a PNG file how am I.. what..."}
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

	processedImage := processing.BasicScaling(imageFile)

	for i := 1; i < scalingIterationCount; i++ {
		processedImage = processing.BasicScaling(processedImage)
	}

	tempResponseFile, _ := ioutil.TempFile(tempDir, "processed-*.png")
	png.Encode(tempResponseFile, processedImage)

	http.ServeFile(w, r, tempResponseFile.Name())

	file.Close()
	tempFile.Close()
	os.RemoveAll(tempDir)
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/", getIndex)
	http.ListenAndServe(":8080", nil)
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
