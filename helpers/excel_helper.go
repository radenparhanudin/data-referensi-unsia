package helpers

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

/* Auto Size Column Export Excel */
func ExcelAutoSizeColumn(file *excelize.File, sheetName, col string, numRows int) {
	var maxLength float64
	for row := 1; row <= numRows+1; row++ {
		cellValue, _ := file.GetCellValue(sheetName, fmt.Sprintf("%s%d", col, row))
		cellLength := float64(len(cellValue))
		if cellLength > maxLength {
			maxLength = cellLength
		}
	}
	file.SetColWidth(sheetName, col, col, maxLength+2)
}
