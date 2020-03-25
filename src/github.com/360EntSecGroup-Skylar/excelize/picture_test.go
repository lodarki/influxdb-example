package excelize

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/tiff"

	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkAddPictureFromBytes(b *testing.B) {
	f := NewFile()
	imgFile, err := ioutil.ReadFile(filepath.Join("test", "images", "excel.png"))
	if err != nil {
		b.Error("unable to load image for benchmark")
	}
	b.ResetTimer()
	for i := 1; i <= b.N; i++ {
		if err := f.AddPictureFromBytes("Sheet1", fmt.Sprint("A", i), "", "excel", ".png", imgFile); err != nil {
			b.Error(err)
		}
	}
}

func TestAddPicture(t *testing.T) {
	f, err := OpenFile(filepath.Join("test", "Book1.xlsx"))
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	// Test add picture to worksheet with offset and location hyperlink.
	assert.NoError(t, f.AddPicture("Sheet2", "I9", filepath.Join("test", "images", "excel.jpg"),
		`{"x_offset": 140, "y_offset": 120, "hyperlink": "#Sheet2!D8", "hyperlink_type": "Location"}`))
	// Test add picture to worksheet with offset, external hyperlink and positioning.
	assert.NoError(t, f.AddPicture("Sheet1", "F21", filepath.Join("test", "images", "excel.jpg"),
		`{"x_offset": 10, "y_offset": 10, "hyperlink": "https://github.com/360EntSecGroup-Skylar/excelize", "hyperlink_type": "External", "positioning": "oneCell"}`))

	file, err := ioutil.ReadFile(filepath.Join("test", "images", "excel.png"))
	assert.NoError(t, err)

	// Test add picture to worksheet from bytes.
	assert.NoError(t, f.AddPictureFromBytes("Sheet1", "Q1", "", "Excel Logo", ".png", file))
	// Test add picture to worksheet from bytes with illegal cell coordinates.
	assert.EqualError(t, f.AddPictureFromBytes("Sheet1", "A", "", "Excel Logo", ".png", file), `cannot convert cell "A" to coordinates: invalid cell name "A"`)

	assert.NoError(t, f.AddPicture("Sheet1", "Q8", filepath.Join("test", "images", "excel.gif"), ""))
	assert.NoError(t, f.AddPicture("Sheet1", "Q15", filepath.Join("test", "images", "excel.jpg"), ""))
	assert.NoError(t, f.AddPicture("Sheet1", "Q22", filepath.Join("test", "images", "excel.tif"), ""))

	// Test write file to given path.
	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestAddPicture.xlsx")))
}

func TestAddPictureErrors(t *testing.T) {
	xlsx, err := OpenFile(filepath.Join("test", "Book1.xlsx"))
	assert.NoError(t, err)

	// Test add picture to worksheet with invalid file path.
	err = xlsx.AddPicture("Sheet1", "G21", filepath.Join("test", "not_exists_dir", "not_exists.icon"), "")
	if assert.Error(t, err) {
		assert.True(t, os.IsNotExist(err), "Expected os.IsNotExist(err) == true")
	}

	// Test add picture to worksheet with unsupport file type.
	err = xlsx.AddPicture("Sheet1", "G21", filepath.Join("test", "Book1.xlsx"), "")
	assert.EqualError(t, err, "unsupported image extension")

	err = xlsx.AddPictureFromBytes("Sheet1", "G21", "", "Excel Logo", "jpg", make([]byte, 1))
	assert.EqualError(t, err, "unsupported image extension")

	// Test add picture to worksheet with invalid file data.
	err = xlsx.AddPictureFromBytes("Sheet1", "G21", "", "Excel Logo", ".jpg", make([]byte, 1))
	assert.EqualError(t, err, "image: unknown format")
}

func TestGetPicture(t *testing.T) {
	f, err := prepareTestBook1()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	file, raw, err := f.GetPicture("Sheet1", "F21")
	assert.NoError(t, err)
	if !assert.NotEmpty(t, filepath.Join("test", file)) || !assert.NotEmpty(t, raw) ||
		!assert.NoError(t, ioutil.WriteFile(filepath.Join("test", file), raw, 0644)) {

		t.FailNow()
	}

	// Try to get picture from a worksheet with illegal cell coordinates.
	_, _, err = f.GetPicture("Sheet1", "A")
	assert.EqualError(t, err, `cannot convert cell "A" to coordinates: invalid cell name "A"`)

	// Try to get picture from a worksheet that doesn't contain any images.
	file, raw, err = f.GetPicture("Sheet3", "I9")
	assert.EqualError(t, err, "sheet Sheet3 is not exist")
	assert.Empty(t, file)
	assert.Empty(t, raw)

	// Try to get picture from a cell that doesn't contain an image.
	file, raw, err = f.GetPicture("Sheet2", "A2")
	assert.NoError(t, err)
	assert.Empty(t, file)
	assert.Empty(t, raw)

	f.getDrawingRelationships("xl/worksheets/_rels/sheet1.xml.rels", "rId8")
	f.getDrawingRelationships("", "")
	f.getSheetRelationshipsTargetByID("", "")
	f.deleteSheetRelationships("", "")

	// Try to get picture from a local storage file.
	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestGetPicture.xlsx")))

	f, err = OpenFile(filepath.Join("test", "TestGetPicture.xlsx"))
	assert.NoError(t, err)

	file, raw, err = f.GetPicture("Sheet1", "F21")
	assert.NoError(t, err)
	if !assert.NotEmpty(t, filepath.Join("test", file)) || !assert.NotEmpty(t, raw) ||
		!assert.NoError(t, ioutil.WriteFile(filepath.Join("test", file), raw, 0644)) {

		t.FailNow()
	}

	// Try to get picture from a local storage file that doesn't contain an image.
	file, raw, err = f.GetPicture("Sheet1", "F22")
	assert.NoError(t, err)
	assert.Empty(t, file)
	assert.Empty(t, raw)

	// Test get picture from none drawing worksheet.
	f = NewFile()
	file, raw, err = f.GetPicture("Sheet1", "F22")
	assert.NoError(t, err)
	assert.Empty(t, file)
	assert.Empty(t, raw)
}

func TestAddDrawingPicture(t *testing.T) {
	// testing addDrawingPicture with illegal cell coordinates.
	f := NewFile()
	assert.EqualError(t, f.addDrawingPicture("sheet1", "", "A", "", 0, 0, 0, 0, nil), `cannot convert cell "A" to coordinates: invalid cell name "A"`)
}

func TestAddPictureFromBytes(t *testing.T) {
	f := NewFile()
	imgFile, err := ioutil.ReadFile("logo.png")
	assert.NoError(t, err, "Unable to load logo for test")
	assert.NoError(t, f.AddPictureFromBytes("Sheet1", fmt.Sprint("A", 1), "", "logo", ".png", imgFile))
	assert.NoError(t, f.AddPictureFromBytes("Sheet1", fmt.Sprint("A", 50), "", "logo", ".png", imgFile))
	imageCount := 0
	for fileName := range f.XLSX {
		if strings.Contains(fileName, "media/image") {
			imageCount++
		}
	}
	assert.Equal(t, 1, imageCount, "Duplicate image should only be stored once.")
	assert.EqualError(t, f.AddPictureFromBytes("SheetN", fmt.Sprint("A", 1), "", "logo", ".png", imgFile), "sheet SheetN is not exist")
}

func TestDeletePicture(t *testing.T) {
	f, err := OpenFile(filepath.Join("test", "Book1.xlsx"))
	assert.NoError(t, err)
	assert.NoError(t, f.DeletePicture("Sheet1", "A1"))
	assert.NoError(t, f.AddPicture("Sheet1", "P1", filepath.Join("test", "images", "excel.jpg"), ""))
	assert.NoError(t, f.DeletePicture("Sheet1", "P1"))
	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestDeletePicture.xlsx")))
	// Test delete picture on not exists worksheet.
	assert.EqualError(t, f.DeletePicture("SheetN", "A1"), "sheet SheetN is not exist")
	// Test delete picture with invalid coordinates.
	assert.EqualError(t, f.DeletePicture("Sheet1", ""), `cannot convert cell "" to coordinates: invalid cell name ""`)
	// Test delete picture on no chart worksheet.
	assert.NoError(t, NewFile().DeletePicture("Sheet1", "A1"))
}
