package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"

	"log"
	"os"
	"path/filepath"

	"github.com/Kagami/go-face"
)

const dataDir = "images"
const imageFile = "friends.jpg" // Group Image with more faces

func main() {
	rec, err := face.NewRecognizer(dataDir)
	if err != nil {
		fmt.Println("Cannot initialize go face recognizer")
	}
	defer rec.Close()

	friendsImage := filepath.Join(dataDir, imageFile)

	faces, err := rec.RecognizeFile(friendsImage)
	if err != nil {
		log.Fatalf("Can't recognize: %v", err)
	}
	fmt.Println("Number of Faces in Image: ", len(faces))
	var samples []face.Descriptor
	var friends []int32
	for i, f := range faces {
		// Croping each face and save as image for reference
		err := saveFace(f.Rectangle.Min.X, f.Rectangle.Min.Y, f.Rectangle.Max.X, f.Rectangle.Max.Y, int(i))
		if err != nil {
			fmt.Println(err)
		}
		samples = append(samples, f.Descriptor)
		// Each face is unique on that image so goes to its own category.
		friends = append(friends, int32(i))
	}
	// Name the categories, i.e. people on the image.
	labels := []string{
		"JK",
		"Budha 1",
		"Prasanth",
		"Budha 2",
		"Vyas Reddy",
		"Budha 3",
		"Akhil",
		"Saketh",
		"Prakash",
	}
	// Pass samples to the recognizer.
	rec.SetSamples(samples, friends)

	// Now let's try to classify some not yet known image.
	testVyasReddy := filepath.Join(dataDir, "vyasreddy.png")
	vyasReddy, err := rec.RecognizeSingleFile(testVyasReddy)
	if err != nil {
		log.Fatalf("Can't recognize: %v", err)
	}
	if vyasReddy == nil {
		log.Fatalf("Not a single face on the image")
	}
	faceID := rec.Classify(vyasReddy.Descriptor)
	if faceID < 0 {
		log.Fatalf("Can't classify")
	}

	fmt.Println(faceID)
	fmt.Println(labels[faceID])

}

func saveFace(top int, bottom int, right int, left int, fid int) error {
	friendsImage := filepath.Join(dataDir, imageFile)

	img, err := readImage(friendsImage)
	if err != nil {
		return err
	}
	img, err = cropImage(img, image.Rect(top, bottom, right, left))
	if err != nil {
		return err
	}
	facePath := fmt.Sprintf("images/%d.png", fid)
	return writeImage(img, facePath)

}

// readImage reads a image file from disk.
func readImage(name string) (image.Image, error) {
	friendsImage := filepath.Join(dataDir, imageFile)

	fd, err := os.Open(friendsImage)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	// image.Decode requires that you import the right image package. We've
	// decode jpeg files then we would need to import "image/jpeg".
	img, _, err := image.Decode(fd)
	if err != nil {
		return nil, err
	}

	return img, nil
}
func Newfunc() int {
	return 5
}

// cropImage takes an image and crops it to the specified rectangle.
func cropImage(img image.Image, crop image.Rectangle) (image.Image, error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
		//Newfunc() int
	}

	// method called SubImage. If it does, then we can use SubImage to crop the
	// image.
	simg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}

	return simg.SubImage(crop), nil
}

// writeImage writes an Image back to the disk.
func writeImage(img image.Image, name string) error {
	fd, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fd.Close()

	return png.Encode(fd, img)
}
