package avatar

import (
	"bytes"
	"fmt"
	"image"
	"image/color/palette"

	// Enable PNG support
	_ "image/png"
	"math/rand"
	"time"

	"github.com/issue9/identicon"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
	"github.com/sitename/sitename/modules/setting"
)

// AvatarSize returns avatar's size
const AvatarSize = 290

// RandomImageSize generates and returns a random avatar image unique to input data
// in custom size (height and width).
func RandomImageSize(size int, data []byte) (image.Image, error) {
	randExtent := len(palette.WebSafe) - 32
	rand.Seed(time.Now().UnixNano())
	colorIndex := rand.Intn(randExtent)
	backColorIndex := colorIndex - 1
	if backColorIndex < 0 {
		backColorIndex = randExtent - 1
	}

	// Define size, background and forecolor
	imgMaker, err := identicon.New(
		size,
		palette.WebSafe[backColorIndex],
		palette.WebSafe[colorIndex:colorIndex+32]...,
	)
	if err != nil {
		return nil, fmt.Errorf("identicon.New: %v", err)
	}
	return imgMaker.Make(data), nil
}

// RandomImage generates and returns a random avatar image unique to input data
// in default size (height and width).
func RandomImage(data []byte) (image.Image, error) {
	return RandomImageSize(AvatarSize, data)
}

// Prepare accepts a byte slice as input, validates it contains an image of an
// acceptable format, and crops and resizes it appropriately.
func Prepare(data []byte) (*image.Image, error) {
	imgCfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("DecodeConfig: %v", err)
	}
	if imgCfg.Width > setting.Avatar.MaxWidth {
		return nil, fmt.Errorf("Image width is too large: %d > %d", imgCfg.Width, setting.Avatar.MaxWidth)
	}
	if imgCfg.Height > setting.Avatar.MaxHeight {
		return nil, fmt.Errorf("Image height is too large: %d > %d", imgCfg.Height, setting.Avatar.MaxHeight)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("Decode: %v", err)
	}

	if imgCfg.Width != imgCfg.Height {
		var newSize, ax, ay int
		if imgCfg.Width > imgCfg.Height {
			newSize = imgCfg.Height
			ax = (imgCfg.Width - imgCfg.Height) / 2
		} else {
			newSize = imgCfg.Height
			ay = (imgCfg.Height - imgCfg.Width) / 2
		}

		img, err = cutter.Crop(img, cutter.Config{
			Width:  newSize,
			Height: newSize,
			Anchor: image.Point{ax, ay},
		})
		if err != nil {
			return nil, err
		}
	}

	img = resize.Resize(AvatarSize, AvatarSize, img, resize.Bilinear)
	return &img, nil
}
