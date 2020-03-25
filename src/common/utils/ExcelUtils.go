package utils

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
)

type ExcelHeader struct {
	Name  string
	Alias string
}

// 将列的序号转化为excel的列号
func getColumnNameByColumnIndex(index int) (axis string) {

	var letters = make([]string, 0, 26)
	var start = 'A'
	var iStart = int(start)
	for i := 0; i < 26; i++ {
		letters = append(letters, string(rune(iStart+i)))
	}

	var loopCount = (index + 1) / 26
	var offset = (index + 1) % 26

	axis = letters[offset-1]
	if loopCount > 0 {
		axis = letters[loopCount-1] + axis
	}

	return
}

// 将数据写入excel文件，data的字段数据不能是复合结构，只能是一个单纯的值。
func WriteExcel(headers []ExcelHeader, data []map[string]interface{}) (xFile *excelize.File, err error) {

	xFile = excelize.NewFile()
	sheet := "Sheet1"
	sheetIndex := xFile.NewSheet(sheet)

	for i, h := range headers {
		headerName := h.Alias
		if len(headerName) == 0 {
			headerName = h.Name
		}
		e := xFile.SetCellValue(sheet, getColumnNameByColumnIndex(i)+"1", headerName)
		if e != nil {
			err = e
			return
		}
	}

	for j, subMap := range data {
		for i, h := range headers {
			if v, ok := subMap[h.Name]; ok {
				e := xFile.SetCellValue(sheet, getColumnNameByColumnIndex(i)+strconv.Itoa(j+2), v)
				if e != nil {
					err = e
					return
				}
			}
		}
	}

	xFile.SetActiveSheet(sheetIndex)
	return
}
