package excelize

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumnVisibility(t *testing.T) {
	t.Run("TestBook1", func(t *testing.T) {
		f, err := prepareTestBook1()
		assert.NoError(t, err)

		// Hide/display a column with SetColVisible
		assert.NoError(t, f.SetColVisible("Sheet1", "F", false))
		assert.NoError(t, f.SetColVisible("Sheet1", "F", true))
		visible, err := f.GetColVisible("Sheet1", "F")
		assert.Equal(t, true, visible)
		assert.NoError(t, err)

		// Test hiding a few columns SetColVisible(...false)...
		assert.NoError(t, f.SetColVisible("Sheet1", "F:V", false))
		visible, err = f.GetColVisible("Sheet1", "F")
		assert.Equal(t, false, visible)
		assert.NoError(t, err)
		visible, err = f.GetColVisible("Sheet1", "U")
		assert.Equal(t, false, visible)
		assert.NoError(t, err)
		visible, err = f.GetColVisible("Sheet1", "V")
		assert.Equal(t, false, visible)
		assert.NoError(t, err)
		// ...and displaying them back SetColVisible(...true)
		assert.NoError(t, f.SetColVisible("Sheet1", "V:F", true))
		visible, err = f.GetColVisible("Sheet1", "F")
		assert.Equal(t, true, visible)
		assert.NoError(t, err)
		visible, err = f.GetColVisible("Sheet1", "U")
		assert.Equal(t, true, visible)
		assert.NoError(t, err)
		visible, err = f.GetColVisible("Sheet1", "G")
		assert.Equal(t, true, visible)
		assert.NoError(t, err)

		// Test get column visible on an inexistent worksheet.
		_, err = f.GetColVisible("SheetN", "F")
		assert.EqualError(t, err, "sheet SheetN is not exist")

		// Test get column visible with illegal cell coordinates.
		_, err = f.GetColVisible("Sheet1", "*")
		assert.EqualError(t, err, `invalid column name "*"`)
		assert.EqualError(t, f.SetColVisible("Sheet1", "*", false), `invalid column name "*"`)

		f.NewSheet("Sheet3")
		assert.NoError(t, f.SetColVisible("Sheet3", "E", false))
		assert.EqualError(t, f.SetColVisible("Sheet1", "A:-1", true), "invalid column name \"-1\"")
		assert.EqualError(t, f.SetColVisible("SheetN", "E", false), "sheet SheetN is not exist")
		assert.NoError(t, f.SaveAs(filepath.Join("test", "TestColumnVisibility.xlsx")))
	})

	t.Run("TestBook3", func(t *testing.T) {
		f, err := prepareTestBook3()
		assert.NoError(t, err)
		visible, err := f.GetColVisible("Sheet1", "B")
		assert.Equal(t, true, visible)
		assert.NoError(t, err)
	})
}

func TestOutlineLevel(t *testing.T) {
	f := NewFile()
	level, err := f.GetColOutlineLevel("Sheet1", "D")
	assert.Equal(t, uint8(0), level)
	assert.NoError(t, err)

	f.NewSheet("Sheet2")
	assert.NoError(t, f.SetColOutlineLevel("Sheet1", "D", 4))

	level, err = f.GetColOutlineLevel("Sheet1", "D")
	assert.Equal(t, uint8(4), level)
	assert.NoError(t, err)

	level, err = f.GetColOutlineLevel("Shee2", "A")
	assert.Equal(t, uint8(0), level)
	assert.EqualError(t, err, "sheet Shee2 is not exist")

	assert.NoError(t, f.SetColWidth("Sheet2", "A", "D", 13))
	assert.NoError(t, f.SetColOutlineLevel("Sheet2", "B", 2))
	assert.NoError(t, f.SetRowOutlineLevel("Sheet1", 2, 7))
	assert.EqualError(t, f.SetColOutlineLevel("Sheet1", "D", 8), "invalid outline level")
	assert.EqualError(t, f.SetRowOutlineLevel("Sheet1", 2, 8), "invalid outline level")
	// Test set row outline level on not exists worksheet.
	assert.EqualError(t, f.SetRowOutlineLevel("SheetN", 1, 4), "sheet SheetN is not exist")
	// Test get row outline level on not exists worksheet.
	_, err = f.GetRowOutlineLevel("SheetN", 1)
	assert.EqualError(t, err, "sheet SheetN is not exist")

	// Test set and get column outline level with illegal cell coordinates.
	assert.EqualError(t, f.SetColOutlineLevel("Sheet1", "*", 1), `invalid column name "*"`)
	_, err = f.GetColOutlineLevel("Sheet1", "*")
	assert.EqualError(t, err, `invalid column name "*"`)

	// Test set column outline level on not exists worksheet.
	assert.EqualError(t, f.SetColOutlineLevel("SheetN", "E", 2), "sheet SheetN is not exist")

	assert.EqualError(t, f.SetRowOutlineLevel("Sheet1", 0, 1), "invalid row number 0")
	level, err = f.GetRowOutlineLevel("Sheet1", 2)
	assert.NoError(t, err)
	assert.Equal(t, uint8(7), level)

	_, err = f.GetRowOutlineLevel("Sheet1", 0)
	assert.EqualError(t, err, `invalid row number 0`)

	level, err = f.GetRowOutlineLevel("Sheet1", 10)
	assert.NoError(t, err)
	assert.Equal(t, uint8(0), level)

	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestOutlineLevel.xlsx")))

	f, err = OpenFile(filepath.Join("test", "Book1.xlsx"))
	assert.NoError(t, err)
	assert.NoError(t, f.SetColOutlineLevel("Sheet2", "B", 2))
}

func TestSetColStyle(t *testing.T) {
	f := NewFile()
	style, err := f.NewStyle(`{"fill":{"type":"pattern","color":["#94d3a2"],"pattern":1}}`)
	assert.NoError(t, err)
	// Test set column style on not exists worksheet.
	assert.EqualError(t, f.SetColStyle("SheetN", "E", style), "sheet SheetN is not exist")
	// Test set column style with illegal cell coordinates.
	assert.EqualError(t, f.SetColStyle("Sheet1", "*", style), `invalid column name "*"`)
	assert.EqualError(t, f.SetColStyle("Sheet1", "A:*", style), `invalid column name "*"`)

	assert.NoError(t, f.SetColStyle("Sheet1", "B", style))
	// Test set column style with already exists column with style.
	assert.NoError(t, f.SetColStyle("Sheet1", "B", style))
	assert.NoError(t, f.SetColStyle("Sheet1", "D:C", style))
	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestSetColStyle.xlsx")))
}

func TestColWidth(t *testing.T) {
	f := NewFile()
	assert.NoError(t, f.SetColWidth("Sheet1", "B", "A", 12))
	assert.NoError(t, f.SetColWidth("Sheet1", "A", "B", 12))
	width, err := f.GetColWidth("Sheet1", "A")
	assert.Equal(t, float64(12), width)
	assert.NoError(t, err)
	width, err = f.GetColWidth("Sheet1", "C")
	assert.Equal(t, float64(64), width)
	assert.NoError(t, err)

	// Test set and get column width with illegal cell coordinates.
	width, err = f.GetColWidth("Sheet1", "*")
	assert.Equal(t, float64(64), width)
	assert.EqualError(t, err, `invalid column name "*"`)
	assert.EqualError(t, f.SetColWidth("Sheet1", "*", "B", 1), `invalid column name "*"`)
	assert.EqualError(t, f.SetColWidth("Sheet1", "A", "*", 1), `invalid column name "*"`)

	// Test set column width on not exists worksheet.
	assert.EqualError(t, f.SetColWidth("SheetN", "B", "A", 12), "sheet SheetN is not exist")

	// Test get column width on not exists worksheet.
	_, err = f.GetColWidth("SheetN", "A")
	assert.EqualError(t, err, "sheet SheetN is not exist")

	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestColWidth.xlsx")))
	convertRowHeightToPixels(0)
}

func TestInsertCol(t *testing.T) {
	f := NewFile()
	sheet1 := f.GetSheetName(1)

	fillCells(f, sheet1, 10, 10)

	assert.NoError(t, f.SetCellHyperLink(sheet1, "A5", "https://github.com/360EntSecGroup-Skylar/excelize", "External"))
	assert.NoError(t, f.MergeCell(sheet1, "A1", "C3"))

	assert.NoError(t, f.AutoFilter(sheet1, "A2", "B2", `{"column":"B","expression":"x != blanks"}`))
	assert.NoError(t, f.InsertCol(sheet1, "A"))

	// Test insert column with illegal cell coordinates.
	assert.EqualError(t, f.InsertCol("Sheet1", "*"), `invalid column name "*"`)

	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestInsertCol.xlsx")))
}

func TestRemoveCol(t *testing.T) {
	f := NewFile()
	sheet1 := f.GetSheetName(1)

	fillCells(f, sheet1, 10, 15)

	assert.NoError(t, f.SetCellHyperLink(sheet1, "A5", "https://github.com/360EntSecGroup-Skylar/excelize", "External"))
	assert.NoError(t, f.SetCellHyperLink(sheet1, "C5", "https://github.com", "External"))

	assert.NoError(t, f.MergeCell(sheet1, "A1", "B1"))
	assert.NoError(t, f.MergeCell(sheet1, "A2", "B2"))

	assert.NoError(t, f.RemoveCol(sheet1, "A"))
	assert.NoError(t, f.RemoveCol(sheet1, "A"))

	// Test remove column with illegal cell coordinates.
	assert.EqualError(t, f.RemoveCol("Sheet1", "*"), `invalid column name "*"`)

	// Test remove column on not exists worksheet.
	assert.EqualError(t, f.RemoveCol("SheetN", "B"), "sheet SheetN is not exist")

	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestRemoveCol.xlsx")))
}
