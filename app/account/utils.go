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

	"github.com/Masterminds/squirrel"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/sitename/sitename/app/plugin/interfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/util/fileutils"
	"github.com/sitename/sitename/store"
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
func CheckUserDomain(user *model.User, domains string) bool {
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

// StoreUserAddress Add address to user address book and set as default one.
func (s *ServiceAccount) StoreUserAddress(user *model.User, address model.Address, addressType string, manager interfaces.PluginManagerInterface) *model.AppError {
	address_, appErr := manager.ChangeUserAddress(address, addressType, user)
	if appErr != nil {
		return appErr
	}

	addressFilterOptions := squirrel.And{}
	if address_.FirstName != "" {
		addressFilterOptions = append(addressFilterOptions, squirrel.Eq{store.AddressTableName + ".FirstName": address_.FirstName})
	}
	if address_.LastName != "" {
		addressFilterOptions = append(addressFilterOptions, squirrel.Eq{store.AddressTableName + ".LastName": address_.LastName})
	}
	if address_.CompanyName != "" {
		addressFilterOptions = append(addressFilterOptions, squirrel.Eq{store.AddressTableName + ".CompanyName": address_.CompanyName})
	}
	if address_.Phone != "" {
		addressFilterOptions = append(addressFilterOptions, squirrel.Eq{store.AddressTableName + ".Phone": address_.Phone})
	}
	if address_.PostalCode != "" {
		addressFilterOptions = append(addressFilterOptions, squirrel.Eq{store.AddressTableName + ".PostalCode": address_.PostalCode})
	}
	if address_.Country != "" {
		addressFilterOptions = append(addressFilterOptions, squirrel.Eq{store.AddressTableName + ".Country": address_.Country})
	}

	addresses, appErr := s.AddressesByOption(&model.AddressFilterOption{
		UserID: squirrel.Eq{store.UserAddressTableName + ".UserID": user.Id},
		Other:  addressFilterOptions,
	})
	if appErr != nil {
		if appErr.StatusCode == http.StatusInternalServerError {
			return appErr
		}
		// ignore not found error
	}

	if len(addresses) == 0 {
		// create new address
		address_.Id = ""
		address_, appErr = s.UpsertAddress(nil, address_)
		if appErr != nil {
			return appErr
		}

		_, appErr = s.AddUserAddress(&model.UserAddress{
			UserID:    user.Id,
			AddressID: address_.Id,
		})
		if appErr != nil {
			return appErr
		}

	} else {
		address_ = addresses[0]
	}

	if addressType == model.ADDRESS_TYPE_BILLING {
		if user.DefaultBillingAddressID == nil {
			appErr = s.SetUserDefaultBillingAddress(user, address_.Id)
		}
	} else if addressType == model.ADDRESS_TYPE_SHIPPING {
		if user.DefaultShippingAddressID == nil {
			appErr = s.SetUserDefaultShippingAddress(user, address_.Id)
		}
	}

	return appErr
}

// SetUserDefaultBillingAddress sets default billing address for given user
func (s *ServiceAccount) SetUserDefaultBillingAddress(user *model.User, defaultBillingAddressID string) *model.AppError {
	user.DefaultBillingAddressID = &defaultBillingAddressID
	_, appErr := s.UpdateUser(user, false)
	return appErr
}

// SetUserDefaultShippingAddress sets default shipping address for given user
func (s *ServiceAccount) SetUserDefaultShippingAddress(user *model.User, defaultShippingAddressID string) *model.AppError {
	user.DefaultShippingAddressID = &defaultShippingAddressID
	_, appErr := s.UpdateUser(user, false)
	return appErr
}

// ChangeUserDefaultAddress set default address for given user
func (s *ServiceAccount) ChangeUserDefaultAddress(user model.User, address model.Address, addressType string, manager interfaces.PluginManagerInterface) *model.AppError {
	address_, appErr := manager.ChangeUserAddress(address, addressType, &user)
	if appErr != nil {
		return appErr
	}

	if addressType == model.ADDRESS_TYPE_BILLING {
		if user.DefaultBillingAddressID != nil {
			_, appErr := s.AddUserAddress(&model.UserAddress{
				UserID:    user.Id,
				AddressID: *user.DefaultBillingAddressID,
			})
			if appErr != nil {
				return appErr
			}
		}

		appErr := s.SetUserDefaultBillingAddress(&user, address_.Id)
		if appErr != nil {
			return appErr
		}
	} else if addressType == model.ADDRESS_TYPE_SHIPPING {
		if user.DefaultShippingAddressID != nil {
			_, appErr := s.AddUserAddress(&model.UserAddress{
				UserID:    user.Id,
				AddressID: *user.DefaultShippingAddressID,
			})
			if appErr != nil {
				return appErr
			}
		}

		appErr := s.SetUserDefaultShippingAddress(&user, address_.Id)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}
