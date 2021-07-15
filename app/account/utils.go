package account

import (
	"bytes"
	"hash/fnv"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/modules/util/fileutils"
)

var (
	// fontStore holds already parsed fonts for later use
	fontStore sync.Map
)

var colors = []color.NRGBA{
	{197, 8, 126, 255},
	{227, 207, 18, 255},
	{28, 181, 105, 255},
	{35, 188, 224, 255},
	{116, 49, 196, 255},
	{197, 8, 126, 255},
	{197, 19, 19, 255},
	{250, 134, 6, 255},
	{227, 207, 18, 255},
	{123, 201, 71, 255},
	{28, 181, 105, 255},
	{35, 188, 224, 255},
	{116, 49, 196, 255},
	{197, 8, 126, 255},
	{197, 19, 19, 255},
	{250, 134, 6, 255},
	{227, 207, 18, 255},
	{123, 201, 71, 255},
	{28, 181, 105, 255},
	{35, 188, 224, 255},
	{116, 49, 196, 255},
	{197, 8, 126, 255},
	{197, 19, 19, 255},
	{250, 134, 6, 255},
	{227, 207, 18, 255},
	{123, 201, 71, 255},
}

func CheckEmailDomain(email string, domains string) bool {
	if domains == "" {
		return true
	}

	domainArray := strings.Fields(
		strings.TrimSpace(
			strings.ToLower(
				strings.Replace(
					strings.Replace(domains, "@", " ", -1),
					",", " ", -1,
				),
			),
		),
	)

	for _, d := range domainArray {
		if strings.HasSuffix(strings.ToLower(email), "@"+d) {
			return true
		}
	}

	return false
}

// CheckUserDomain checks that a user's email domain matches a list of space-delimited domains as a string.
func CheckUserDomain(user *account.User, domains string) bool {
	return CheckEmailDomain(user.Email, domains)
}

func getFont(initialFont string) (*truetype.Font, error) {
	// Some people have the old default font still set, so just treat that as if they're using the new default
	if initialFont == "luximbi.ttf" {
		initialFont = "nunito-bold.ttf"
	}

	// try getting font from memory
	if value, ok := fontStore.Load(initialFont); ok {
		return value.(*truetype.Font), nil
	}

	fontDir, _ := fileutils.FindDir("fonts")
	fontBytes, err := ioutil.ReadFile(filepath.Join(fontDir, initialFont))
	if err != nil {
		return nil, err
	}

	parsed, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	// put font into memory
	fontStore.Store(initialFont, parsed)

	return parsed, nil
}

func CreateProfileImage(username string, userID string, initialFont string) ([]byte, *model.AppError) {
	h := fnv.New32a()
	h.Write([]byte(userID))
	seed := h.Sum32()

	initial := string(strings.ToUpper(username)[0])

	font, err := getFont(initialFont)
	if err != nil {
		return nil, model.NewAppError("CreateProfileImage", "api.user.create_profile_image.default_font.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	color := colors[int64(seed)%int64(len(colors))]
	dstImg := image.NewRGBA(image.Rect(0, 0, ImageProfilePixelDimension, ImageProfilePixelDimension))
	srcImg := image.White
	draw.Draw(dstImg, dstImg.Bounds(), &image.Uniform{color}, image.Point{}, draw.Src)
	size := float64(ImageProfilePixelDimension / 2)

	c := freetype.NewContext()
	c.SetFont(font)
	c.SetFontSize(size)
	c.SetClip(dstImg.Bounds())
	c.SetDst(dstImg)
	c.SetSrc(srcImg)

	pt := freetype.Pt(ImageProfilePixelDimension/5, ImageProfilePixelDimension*2/3)
	_, err = c.DrawString(initial, pt)
	if err != nil {
		return nil, model.NewAppError("CreateProfileImage", "api.user.create_profile_image.initial.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	buf := new(bytes.Buffer)

	if imgErr := png.Encode(buf, dstImg); err != nil {
		return nil, model.NewAppError("CreateProfileImage", "api.user.create_profile_image.encode.app_error", nil, imgErr.Error(), http.StatusInternalServerError)
	}

	return buf.Bytes(), nil
}
