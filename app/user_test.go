package app

import (
	"bytes"
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckUserDomain(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	user := th.BasicUser

	cases := []struct {
		domains string
		matched bool
	}{
		{"simulator.amazonses.com", true},
		{"gmail.com", false},
		{"", true},
		{"gmail.com simulator.amazonses.com", true},
	}
	for _, c := range cases {
		matched := CheckUserDomain(user, c.domains)
		if matched != c.matched {
			if c.matched {
				t.Logf("'%v' should have matched '%v'", user.Email, c.domains)
			} else {
				t.Logf("'%v' should not have matched '%v'", user.Email, c.domains)
			}
			t.FailNow()
		}
	}
}

func TestCreateProfileImage(t *testing.T) {
	b, err := CreateProfileImage("Corey Hulen", "eo1zkdr96pdj98pjmq8zy35wba", "nunito-bold.ttf")
	require.Nil(t, err)

	rdr := bytes.NewReader(b)
	img, _, err2 := image.Decode(rdr)
	require.NoError(t, err2)

	colorful := color.RGBA{116, 49, 196, 255}

	require.Equal(t, colorful, img.At(1, 1), "Failed to create correct color")
}
