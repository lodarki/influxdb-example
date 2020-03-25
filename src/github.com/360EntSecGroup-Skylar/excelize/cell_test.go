package excelize

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckCellInArea(t *testing.T) {
	f := NewFile()
	expectedTrueCellInAreaList := [][2]string{
		{"c2", "A1:AAZ32"},
		{"B9", "A1:B9"},
		{"C2", "C2:C2"},
	}

	for _, expectedTrueCellInArea := range expectedTrueCellInAreaList {
		cell := expectedTrueCellInArea[0]
		area := expectedTrueCellInArea[1]
		ok, err := f.checkCellInArea(cell, area)
		assert.NoError(t, err)
		assert.Truef(t, ok,
			"Expected cell %v to be in area %v, got false\n", cell, area)
	}

	expectedFalseCellInAreaList := [][2]string{
		{"c2", "A4:AAZ32"},
		{"C4", "D6:A1"}, // weird case, but you never know
		{"AEF42", "BZ40:AEF41"},
	}

	for _, expectedFalseCellInArea := range expectedFalseCellInAreaList {
		cell := expectedFalseCellInArea[0]
		area := expectedFalseCellInArea[1]
		ok, err := f.checkCellInArea(cell, area)
		assert.NoError(t, err)
		assert.Falsef(t, ok,
			"Expected cell %v not to be inside of area %v, but got true\n", cell, area)
	}

	ok, err := f.checkCellInArea("A1", "A:B")
	assert.EqualError(t, err, `cannot convert cell "A" to coordinates: invalid cell name "A"`)
	assert.False(t, ok)

	ok, err = f.checkCellInArea("AA0", "Z0:AB1")
	assert.EqualError(t, err, `cannot convert cell "AA0" to coordinates: invalid cell name "AA0"`)
	assert.False(t, ok)
}

func TestSetCellFloat(t *testing.T) {
	sheet := "Sheet1"
	t.Run("with no decimal", func(t *testing.T) {
		f := NewFile()
		assert.NoError(t, f.SetCellFloat(sheet, "A1", 123.0, -1, 64))
		assert.NoError(t, f.SetCellFloat(sheet, "A2", 123.0, 1, 64))
		val, err := f.GetCellValue(sheet, "A1")
		assert.NoError(t, err)
		assert.Equal(t, "123", val, "A1 should be 123")
		val, err = f.GetCellValue(sheet, "A2")
		assert.NoError(t, err)
		assert.Equal(t, "123.0", val, "A2 should be 123.0")
	})

	t.Run("with a decimal and precision limit", func(t *testing.T) {
		f := NewFile()
		assert.NoError(t, f.SetCellFloat(sheet, "A1", 123.42, 1, 64))
		val, err := f.GetCellValue(sheet, "A1")
		assert.NoError(t, err)
		assert.Equal(t, "123.4", val, "A1 should be 123.4")
	})

	t.Run("with a decimal and no limit", func(t *testing.T) {
		f := NewFile()
		assert.NoError(t, f.SetCellFloat(sheet, "A1", 123.42, -1, 64))
		val, err := f.GetCellValue(sheet, "A1")
		assert.NoError(t, err)
		assert.Equal(t, "123.42", val, "A1 should be 123.42")
	})
	f := NewFile()
	assert.EqualError(t, f.SetCellFloat(sheet, "A", 123.42, -1, 64), `cannot convert cell "A" to coordinates: invalid cell name "A"`)
}

func TestSetCellValue(t *testing.T) {
	f := NewFile()
	assert.EqualError(t, f.SetCellValue("Sheet1", "A", time.Now().UTC()), `cannot convert cell "A" to coordinates: invalid cell name "A"`)
	assert.EqualError(t, f.SetCellValue("Sheet1", "A", time.Duration(1e13)), `cannot convert cell "A" to coordinates: invalid cell name "A"`)
}

func TestSetCellBool(t *testing.T) {
	f := NewFile()
	assert.EqualError(t, f.SetCellBool("Sheet1", "A", true), `cannot convert cell "A" to coordinates: invalid cell name "A"`)
}

func TestGetCellFormula(t *testing.T) {
	// Test get cell formula on not exist worksheet.
	f := NewFile()
	_, err := f.GetCellFormula("SheetN", "A1")
	assert.EqualError(t, err, "sheet SheetN is not exist")

	// Test get cell formula on no formula cell.
	assert.NoError(t, f.SetCellValue("Sheet1", "A1", true))
	_, err = f.GetCellFormula("Sheet1", "A1")
	assert.NoError(t, err)
}

func ExampleFile_SetCellFloat() {
	f := NewFile()
	var x = 3.14159265
	if err := f.SetCellFloat("Sheet1", "A1", x, 2, 64); err != nil {
		fmt.Println(err)
	}
	val, _ := f.GetCellValue("Sheet1", "A1")
	fmt.Println(val)
	// Output: 3.14
}

func BenchmarkSetCellValue(b *testing.B) {
	values := []string{"First", "Second", "Third", "Fourth", "Fifth", "Sixth"}
	cols := []string{"A", "B", "C", "D", "E", "F"}
	f := NewFile()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(values); j++ {
			if err := f.SetCellValue("Sheet1", fmt.Sprint(cols[j], i), values[j]); err != nil {
				b.Error(err)
			}
		}
	}
}

func TestOverflowNumericCell(t *testing.T) {
	f, err := OpenFile(filepath.Join("test", "OverflowNumericCell.xlsx"))
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	val, err := f.GetCellValue("Sheet1", "A1")
	assert.NoError(t, err)
	// GOARCH=amd64 - all ok; GOARCH=386 - actual: "-2147483648"
	assert.Equal(t, "8595602512225", val, "A1 should be 8595602512225")
}
