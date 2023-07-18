package converter

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

type Excel struct{}

func (e *Excel) Convert(source []byte) ([]byte, error) {
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	tempFileName := fmt.Sprintf("Sheet1", time.Now().Unix())
	index, err := file.NewSheet(tempFileName)
	if err != nil {
		return nil, err
	}

	file.SetActiveSheet(index)

	file.SetSheetRow("Sheet1", "A1")

	// gen.SetCellValue("Sheet2", "A2", "Hello world.")
	// gen.SetCellValue("Sheet1", "B2", 100)

	// if err := f.SaveAs(dstPath); err != nil {
	// 	fmt.Println(err)
	// }

	buf, err := gen.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func NewExcel() *Excel {
	return &Excel{}
}
