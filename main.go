package main

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/nfnt/resize"
)

var imageFormatsRegex, _ = regexp.Compile(`\.((jpg)|(png))$`)
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

const queueDirectory = "./queue/"
const downscaledFolderPath = "./processed/localDownscaled/"

func RandomString(n int) string {
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func main() {

	queuedFiles, err := ioutil.ReadDir(queueDirectory)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(queuedFiles)

	for iterationNumber, file := range queuedFiles {

		if imageFormatsRegex.MatchString(file.Name()) {
			file, err := os.Open(file.Name())
			if err != nil {
				log.Fatal(err)
			}
			image, err := jpeg.Decode(file)
			if err != nil {
				log.Fatal(err)
			}
			file.Close()

			// Target size: 5mb
			resizedImage := resize.Resize(1000, 0, image, resize.Bicubic)
			fileName := fmt.Sprintf("film-%s-%s-%x", strconv.Itoa(iterationNumber), time.Now().Format("yyyy-MM-dd"), RandomString(6))

			output, err := os.Create(fmt.Sprintf("%s%s", downscaledFolderPath, fileName))
			if err != nil {
				log.Fatal(err)
			}
			defer output.Close()

			png.Encode(output, resizedImage)
		} else {
			log.Fatalf("Failed to parse file: %v, please correct and rerun", file.Name())
		}

		fmt.Println(file.Name(), file.IsDir(), file.Size())
	}

	/*

		// resize to width 1000 using Lanczos resampling
		// and preserve aspect ratio
		m := resize.Resize(1000, 0, img, resize.Lanczos3)

		out, err := os.Create("test_resized.jpg")
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		// write new image to file
		jpeg.Encode(out, m, nil)
	*/
}
