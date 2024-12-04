package bot

import (
	"github.com/vitali-fedulov/images3"
	"os"
)

func MustOpen(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return file
}

func ComparePhotos(path1, path2 string) bool {
	img1, _ := images3.Open(path1)
	img2, _ := images3.Open(path2)
	icon1 := images3.Icon(img1, path1)
	icon2 := images3.Icon(img2, path2)

	if images3.Similar(icon1, icon2) {
		return true
	} else {
		return false
	}
}
