package captcha

import (
	"image/color"

	"github.com/mojocn/base64Captcha"
)

type FontsStorage interface{}

type DriverString struct {
	Height          int
	Width           int
	NoiseCount      int
	ShowLineOptions int
	Length          int
	Source          string
	BgColor         *color.RGBA
	FontsStorage    base64Captcha.FontsStorage
	Fonts           []string
}

type DriverMath struct {
	Height          int
	Width           int
	NoiseCount      int
	ShowLineOptions int
	BgColor         *color.RGBA
	FontsStorage    base64Captcha.FontsStorage
	Fonts           []string
}

type DriverAudio struct {
	Length   int
	Language string
}

// Constructors
func NewDriverString(
	height int,
	width int,
	noiseCount int,
	showLineOptions int,
	length int,
	source string,
	bgColor *color.RGBA,
) *DriverString {
	return &DriverString{
		Height:          height,
		Width:           width,
		NoiseCount:      noiseCount,
		ShowLineOptions: showLineOptions,
		Length:          length,
		Source:          source,
		BgColor:         bgColor,
	}
}

func NewDriverMath(
	height int,
	width int,
	noiseCount int,
	showLineOptions int,
	bgColor *color.RGBA,

) *DriverMath {
	return &DriverMath{
		Height:          height,
		Width:           width,
		NoiseCount:      noiseCount,
		ShowLineOptions: showLineOptions,
		BgColor:         bgColor,
	}
}

func NewDriverAudio(length int, language string) *DriverAudio {
	return &DriverAudio{
		Length:   length,
		Language: language,
	}
}

// Suggested Default Instances
var (
	suggestDriverString = DriverString{
		Height:          60,
		Width:           245,
		NoiseCount:      70,
		ShowLineOptions: base64Captcha.OptionShowHollowLine,
		Length:          7,
		Source:          "0123456789bcdfghjklmnpqrstvwxyzBCDFGHJKLMNPQRSTVWXYZ",
		BgColor:         &color.RGBA{R: 93, G: 178, B: 178, A: 200},
		FontsStorage:    nil,
		Fonts:           []string{},
	}

	suggestDriverMath = DriverMath{
		Height:          60,
		Width:           175,
		NoiseCount:      55,
		ShowLineOptions: 2,
		BgColor:         &color.RGBA{R: 93, G: 178, B: 178, A: 200},
		FontsStorage:    nil,
		Fonts:           []string{},
	}

	suggestDriverAudio = DriverAudio{
		Length:   5,
		Language: "en",
	}
)
