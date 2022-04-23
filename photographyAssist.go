package main

import (
	"fmt"
	"image"
	"image/png"
	"io/fs"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/nfnt/resize"
)

type ImageConstraints struct {
	width  int
	height int
}

var (
	imageFormatsRegex, _ = regexp.Compile(`\.((jpg)|(png))$`)
	letters              = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

const (
	queueFolderPath                = "./queue/"
	downscaledFolderPath           = "./processed/localDownscaled/"
	largePictureFolderPath         = "./processed/large/"
	roughTargetFileSize    float64 = 5000000
)

func RandomString(n int) string {
	rand.Seed(time.Now().Unix())
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func createDirectoriesIfNotExist() {
	err := os.MkdirAll(downscaledFolderPath, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	err = os.MkdirAll(queueFolderPath, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	err = os.MkdirAll(largePictureFolderPath, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
}

func assessMaximumScalingForTargetSize(image image.Image, file fs.FileInfo) ImageConstraints {
	roughUncompressedSize := (image.Bounds().Dx() * image.Bounds().Dy()) * 6

	var compressionRatio float64 = float64(file.Size()) / float64(roughUncompressedSize)
	var aspectRatio float64 = float64(image.Bounds().Dy()) / float64((image.Bounds().Dx()))
	assumedScaledWidth := math.Sqrt(((roughTargetFileSize / compressionRatio) / 6) / aspectRatio)

	return ImageConstraints{int(assumedScaledWidth), 0}
}

func main() {
	createDirectoriesIfNotExist()

	queuedFiles, err := ioutil.ReadDir(queueFolderPath)
	if err != nil || len(queuedFiles) == 0 {
		log.Fatal("Failed to run, queue empty")
	}

	for iterationNumber, file := range queuedFiles {

		if imageFormatsRegex.MatchString(file.Name()) {
			fileData, err := os.Open(fmt.Sprintf("%s%s", queueFolderPath, file.Name()))
			if err != nil {
				log.Fatal(err)
			}
			image, err := png.Decode(fileData)
			if err != nil {
				log.Fatal(err)
			}
			fileData.Close()

			constraints := assessMaximumScalingForTargetSize(image, file)
			fmt.Sprintln(constraints)

			resizedImage := resize.Resize(uint(constraints.width), uint(constraints.height), image, resize.Lanczos3)
			fileName := fmt.Sprintf("film-%s-%s-%x.png", strconv.Itoa(iterationNumber), time.Now().Format("2006-02-01"), RandomString(6))

			output, err := os.Create(fmt.Sprintf("%s%s", downscaledFolderPath, fileName))
			if err != nil {
				log.Fatal(err)
			}
			defer output.Close()

			png.Encode(output, resizedImage)
			fmt.Println(file.Name())
			os.Rename(fmt.Sprintf("%s%s", queueFolderPath, file.Name()), fmt.Sprintf("%slarge-%s", largePictureFolderPath, fileName))
		} else {
			log.Fatalf("Failed to parse file: %v, please correct and rerun", file.Name())
		}
	}

}
