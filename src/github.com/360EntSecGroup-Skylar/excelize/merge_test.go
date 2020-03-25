package excelize

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeCell(t *testing.T) {
	f, err := OpenFile(filepath.Join("test", "Book1.xlsx"))
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.EqualError(t, f.MergeCell("Sheet1", "A", "B"), `cannot convert cell "A" to coordinates: invalid cell name "A"`)
	assert.NoError(t, f.MergeCell("Sheet1", "D9", "D9"))
	assert.NoError(t, f.MergeCell("Sheet1", "D9", "E9"))
	assert.NoError(t, f.MergeCell("Sheet1", "H14", "G13"))
	assert.NoError(t, f.MergeCell("Sheet1", "C9", "D8"))
	assert.NoError(t, f.MergeCell("Sheet1", "F11", "G13"))
	assert.NoError(t, f.MergeCell("Sheet1", "H7", "B15"))
	assert.NoError(t, f.MergeCell("Sheet1", "D11", "F13"))
	assert.NoError(t, f.MergeCell("Sheet1", "G10", "K12"))
	assert.NoError(t, f.SetCellValue("Sheet1", "G11", "set value in merged cell"))
	assert.NoError(t, f.SetCellInt("Sheet1", "H11", 100))
	assert.NoError(t, f.SetCellValue("Sheet1", "I11", float64(0.5)))
	assert.NoError(t, f.SetCellHyperLink("Sheet1", "J11", "https://github.com/360EntSecGroup-Skylar/excelize", "External"))
	assert.NoError(t, f.SetCellFormula("Sheet1", "G12", "SUM(Sheet1!B19,Sheet1!C19)"))
	value, err := f.GetCellValue("Sheet1", "H11")
	assert.Equal(t, "0.5", value)
	assert.NoError(t, err)
	value, err = f.GetCellValue("Sheet2", "A6") // Merged cell ref is single coordinate.
	assert.Equal(t, "", value)
	assert.NoError(t, err)
	value, err = f.GetCellFormula("Sheet1", "G12")
	assert.Equal(t, "SUM(Sheet1!B19,Sheet1!C19)", value)
	assert.NoError(t, err)

	f.NewSheet("Sheet3")
	assert.NoError(t, f.MergeCell("Sheet3", "D11", "F13"))
	assert.NoError(t, f.MergeCell("Sheet3", "G10", "K12"))

	assert.NoError(t, f.MergeCell("Sheet3", "B1", "D5")) // B1:D5
	assert.NoError(t, f.MergeCell("Sheet3", "E1", "F5")) // E1:F5

	assert.NoError(t, f.MergeCell("Sheet3", "H2", "I5"))
	assert.NoError(t, f.MergeCell("Sheet3", "I4", "J6")) // H2:J6

	assert.NoError(t, f.MergeCell("Sheet3", "M2", "N5"))
	assert.NoError(t, f.MergeCell("Sheet3", "L4", "M6")) // L2:N6

	assert.NoError(t, f.MergeCell("Sheet3", "P4", "Q7"))
	assert.NoError(t, f.MergeCell("Sheet3", "O2", "P5")) // O2:Q7

	assert.NoError(t, f.MergeCell("Sheet3", "A9", "B12"))
	assert.NoError(t, f.MergeCell("Sheet3", "B7", "C9")) // A7:C12

	assert.NoError(t, f.MergeCell("Sheet3", "E9", "F10"))
	assert.NoError(t, f.MergeCell("Sheet3", "D8", "G12"))

	assert.NoError(t, f.MergeCell("Sheet3", "I8", "I12"))
	assert.NoError(t, f.MergeCell("Sheet3", "I10", "K10"))

	assert.NoError(t, f.MergeCell("Sheet3", "M8", "Q13"))
	assert.NoError(t, f.MergeCell("Sheet3", "N10", "O11"))

	// Test get merged cells on not exists worksheet.
	assert.EqualError(t, f.MergeCell("SheetN", "N10", "O11"), "sheet SheetN is not exist")

	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestMergeCell.xlsx")))

	f = NewFile()
	assert.NoError(t, f.MergeCell("Sheet1", "A2", "B3"))
	f.Sheet["xl/worksheets/sheet1.xml"].MergeCells = &xlsxMergeCells{Cells: []*xlsxMergeCell{nil, nil}}
	assert.NoError(t, f.MergeCell("Sheet1", "A2", "B3"))

	f.Sheet["xl/worksheets/sheet1.xml"].MergeCells = &xlsxMergeCells{Cells: []*xlsxMergeCell{{Ref: "A1"}}}
	assert.EqualError(t, f.MergeCell("Sheet1", "A2", "B3"), `invalid area "A1"`)

	f.Sheet["xl/worksheets/sheet1.xml"].MergeCells = &xlsxMergeCells{Cells: []*xlsxMergeCell{{Ref: "A:A"}}}
	assert.EqualError(t, f.MergeCell("Sheet1", "A2", "B3"), `cannot convert cell "A" to coordinates: invalid cell name "A"`)
}

func TestGetMergeCells(t *testing.T) {
	wants := []struct {
		value string
		start string
		end   string
	}{{
		value: "A1",
		start: "A1",
		end:   "B1",
	}, {
		value: "A2",
		start: "A2",
		end:   "A3",
	}, {
		value: "A4",
		start: "A4",
		end:   "B5",
	}, {
		value: "A7",
		start: "A7",
		end:   "C10",
	}}

	f, err := OpenFile(filepath.Join("test", "MergeCell.xlsx"))
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	sheet1 := f.GetSheetName(1)

	mergeCells, err := f.GetMergeCells(sheet1)
	if !assert.Len(t, mergeCells, len(wants)) {
		t.FailNow()
	}
	assert.NoError(t, err)

	for i, m := range mergeCells {
		assert.Equal(t, wants[i].value, m.GetCellValue())
		assert.Equal(t, wants[i].start, m.GetStartAxis())
		assert.Equal(t, wants[i].end, m.GetEndAxis())
	}

	// Test get merged cells on not exists worksheet.
	_, err = f.GetMergeCells("SheetN")
	assert.EqualError(t, err, "sheet SheetN is not exist")
}

func TestUnmergeCell(t *testing.T) {
	f, err := OpenFile(filepath.Join("test", "MergeCell.xlsx"))
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	sheet1 := f.GetSheetName(1)

	xlsx, err := f.workSheetReader(sheet1)
	assert.NoError(t, err)

	mergeCellNum := len(xlsx.MergeCells.Cells)

	assert.EqualError(t, f.UnmergeCell("Sheet1", "A", "A"), `cannot convert cell "A" to coordinates: invalid cell name "A"`)

	// unmerge the mergecell that contains A1
	assert.NoError(t, f.UnmergeCell(sheet1, "A1", "A1"))
	if len(xlsx.MergeCells.Cells) != mergeCellNum-1 {
		t.FailNow()
	}

	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestUnmergeCell.xlsx")))

	f = NewFile()
	assert.NoError(t, f.MergeCell("Sheet1", "A2", "B3"))
	// Test unmerged area on not exists worksheet.
	assert.EqualError(t, f.UnmergeCell("SheetN", "A1", "A1"), "sheet SheetN is not exist")

	f.Sheet["xl/worksheets/sheet1.xml"].MergeCells = nil
	assert.NoError(t, f.UnmergeCell("Sheet1", "H7", "B15"))

	f.Sheet["xl/worksheets/sheet1.xml"].MergeCells = &xlsxMergeCells{Cells: []*xlsxMergeCell{nil, nil}}
	assert.NoError(t, f.UnmergeCell("Sheet1", "H15", "B7"))

	f.Sheet["xl/worksheets/sheet1.xml"].MergeCells = &xlsxMergeCells{Cells: []*xlsxMergeCell{{Ref: "A1"}}}
	assert.EqualError(t, f.UnmergeCell("Sheet1", "A2", "B3"), `invalid area "A1"`)

	f.Sheet["xl/worksheets/sheet1.xml"].MergeCells = &xlsxMergeCells{Cells: []*xlsxMergeCell{{Ref: "A:A"}}}
	assert.EqualError(t, f.UnmergeCell("Sheet1", "A2", "B3"), `cannot convert cell "A" to coordinates: invalid cell name "A"`)

}
