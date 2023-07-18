package converter

import (
	"bytes"

	pdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type Pdf struct{}

func (p *Pdf) Convert(source []byte) ([]byte, error) {
	gen, err := pdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	gen.Dpi.Set(300)
	gen.Orientation.Set(pdf.OrientationLandscape)
	gen.Grayscale.Set(true)

	// page := pdf.NewPage(sourcePath)
	// page.FooterRight.Set("[page]")
	// page.FooterFontSize.Set(10)
	// page.Zoom.Set(0.95)

	page := pdf.NewPageReader(bytes.NewReader(source))
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(0.95)

	gen.AddPage(page)

	gen.AddPage(page)

	err = gen.Create()
	if err != nil {
		return nil, err
	}

	// err = gen.WriteFile(dstPath)
	// if err != nil {
	// 	return nil, err
	// }

	return gen.Bytes(), nil
}

func NewPdf() *Pdf {
	return &Pdf{}
}
