package static

import (
	"embed"
	"html/template"
)

//go:embed otp.html
var files embed.FS

type StaticFiles struct {
	otpHTMLTemplate *template.Template
}

func InitStaticFiles() (*StaticFiles, error) {
	otpHTMLTemplate, err := template.ParseFS(files, "otp.html")
	if err != nil {
		return nil, err
	}

	return &StaticFiles{
		otpHTMLTemplate: otpHTMLTemplate,
	}, nil

}
