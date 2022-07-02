package main

import (
	"bytes"
	"fmt"
	"github.com/lib/pq"
	"github.com/noelyahan/impexp"
	"github.com/noelyahan/mergi"
	"image"
	"image/png"
)

func MergeItemImages(items pq.Int64Array) (imageUrl *string, err error) {

	var images []image.Image
	for _, item := range items {
		if item > 0 {
			url := fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/%s/img/item/%d.png", Configs.GameVersion, item)
			image, _ := mergi.Import(impexp.NewURLImporter(url))
			resizedImage, _ := mergi.Resize(image, uint(32), uint(32))
			images = append(images, resizedImage)
		}
	}

	var template string
	for i := 0; i < len(images); i++ {
		template += "T"
	}

	itemListImage, err := mergi.Merge(template, images)
	if err != nil {
		return nil, err
	}

	b := make([]byte, 0)
	buf := bytes.NewBuffer(b)
	png.Encode(buf, itemListImage)

	imgUrl, err := UploadImageToImgur(buf.Bytes())
	if err != nil {
		return nil, err
	}

	return imgUrl, nil
}
