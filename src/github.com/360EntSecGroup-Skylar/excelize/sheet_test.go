package excelize_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/360EntSecGroup-Skylar/excelize/v2"

	"github.com/mohae/deepcopy"
	"github.com/stretchr/testify/assert"
)

func ExampleFile_SetPageLayout() {
	f := excelize.NewFile()

	if err := f.SetPageLayout(
		"Sheet1",
		excelize.PageLayoutOrientation(excelize.OrientationLandscape),
	); err != nil {
		fmt.Println(err)
	}
	if err := f.SetPageLayout(
		"Sheet1",
		excelize.PageLayoutPaperSize(10),
		excelize.FitToHeight(2),
		excelize.FitToWidth(2),
	); err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleFile_GetPageLayout() {
	f := excelize.NewFile()
	var (
		orientation excelize.PageLayoutOrientation
		paperSize   excelize.PageLayoutPaperSize
		fitToHeight excelize.FitToHeight
		fitToWidth  excelize.FitToWidth
	)
	if err := f.GetPageLayout("Sheet1", &orientation); err != nil {
		fmt.Println(err)
	}
	if err := f.GetPageLayout("Sheet1", &paperSize); err != nil {
		fmt.Println(err)
	}
	if err := f.GetPageLayout("Sheet1", &fitToHeight); err != nil {
		fmt.Println(err)
	}

	if err := f.GetPageLayout("Sheet1", &fitToWidth); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Defaults:")
	fmt.Printf("- orientation: %q\n", orientation)
	fmt.Printf("- paper size: %d\n", paperSize)
	fmt.Printf("- fit to height: %d\n", fitToHeight)
	fmt.Printf("- fit to width: %d\n", fitToWidth)
	// Output:
	// Defaults:
	// - orientation: "portrait"
	// - paper size: 1
	// - fit to height: 1
	// - fit to width: 1
}

func TestNewSheet(t *testing.T) {
	f := excelize.NewFile()
	sheetID := f.NewSheet("Sheet2")
	f.SetActiveSheet(sheetID)
	// delete original sheet
	f.DeleteSheet(f.GetSheetName(f.GetSheetIndex("Sheet1")))
	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestNewSheet.xlsx")))
}

func TestSetPane(t *testing.T) {
	f := excelize.NewFile()
	assert.NoError(t, f.SetPanes("Sheet1", `{"freeze":false,"split":false}`))
	f.NewSheet("Panes 2")
	assert.NoError(t, f.SetPanes("Panes 2", `{"freeze":true,"split":false,"x_split":1,"y_split":0,"top_left_cell":"B1","active_pane":"topRight","panes":[{"sqref":"K16","active_cell":"K16","pane":"topRight"}]}`))
	f.NewSheet("Panes 3")
	assert.NoError(t, f.SetPanes("Panes 3", `{"freeze":false,"split":true,"x_split":3270,"y_split":1800,"top_left_cell":"N57","active_pane":"bottomLeft","panes":[{"sqref":"I36","active_cell":"I36"},{"sqref":"G33","active_cell":"G33","pane":"topRight"},{"sqref":"J60","active_cell":"J60","pane":"bottomLeft"},{"sqref":"O60","active_cell":"O60","pane":"bottomRight"}]}`))
	f.NewSheet("Panes 4")
	assert.NoError(t, f.SetPanes("Panes 4", `{"freeze":true,"split":false,"x_split":0,"y_split":9,"top_left_cell":"A34","active_pane":"bottomLeft","panes":[{"sqref":"A11:XFD11","active_cell":"A11","pane":"bottomLeft"}]}`))
	assert.NoError(t, f.SetPanes("Panes 4", ""))
	assert.EqualError(t, f.SetPanes("SheetN", ""), "sheet SheetN is not exist")
	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestSetPane.xlsx")))
}

func TestPageLayoutOption(t *testing.T) {
	const sheet = "Sheet1"

	testData := []struct {
		container  excelize.PageLayoutOptionPtr
		nonDefault excelize.PageLayoutOption
	}{
		{new(excelize.PageLayoutOrientation), excelize.PageLayoutOrientation(excelize.OrientationLandscape)},
		{new(excelize.PageLayoutPaperSize), excelize.PageLayoutPaperSize(10)},
		{new(excelize.FitToHeight), excelize.FitToHeight(2)},
		{new(excelize.FitToWidth), excelize.FitToWidth(2)},
	}

	for i, test := range testData {
		t.Run(fmt.Sprintf("TestData%d", i), func(t *testing.T) {

			opt := test.nonDefault
			t.Logf("option %T", opt)

			def := deepcopy.Copy(test.container).(excelize.PageLayoutOptionPtr)
			val1 := deepcopy.Copy(def).(excelize.PageLayoutOptionPtr)
			val2 := deepcopy.Copy(def).(excelize.PageLayoutOptionPtr)

			f := excelize.NewFile()
			// Get the default value
			assert.NoError(t, f.GetPageLayout(sheet, def), opt)
			// Get again and check
			assert.NoError(t, f.GetPageLayout(sheet, val1), opt)
			if !assert.Equal(t, val1, def, opt) {
				t.FailNow()
			}
			// Set the same value
			assert.NoError(t, f.SetPageLayout(sheet, val1), opt)
			// Get again and check
			assert.NoError(t, f.GetPageLayout(sheet, val1), opt)
			if !assert.Equal(t, val1, def, "%T: value should not have changed", opt) {
				t.FailNow()
			}
			// Set a different value
			assert.NoError(t, f.SetPageLayout(sheet, test.nonDefault), opt)
			assert.NoError(t, f.GetPageLayout(sheet, val1), opt)
			// Get again and compare
			assert.NoError(t, f.GetPageLayout(sheet, val2), opt)
			if !assert.Equal(t, val1, val2, "%T: value should not have changed", opt) {
				t.FailNow()
			}
			// Value should not be the same as the default
			if !assert.NotEqual(t, def, val1, "%T: value should have changed from default", opt) {
				t.FailNow()
			}
			// Restore the default value
			assert.NoError(t, f.SetPageLayout(sheet, def), opt)
			assert.NoError(t, f.GetPageLayout(sheet, val1), opt)
			if !assert.Equal(t, def, val1) {
				t.FailNow()
			}
		})
	}
}

func TestSearchSheet(t *testing.T) {
	f, err := excelize.OpenFile(filepath.Join("test", "SharedStrings.xlsx"))
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	// Test search in a not exists worksheet.
	_, err = f.SearchSheet("Sheet4", "")
	assert.EqualError(t, err, "sheet Sheet4 is not exist")
	var expected []string
	// Test search a not exists value.
	result, err := f.SearchSheet("Sheet1", "X")
	assert.NoError(t, err)
	assert.EqualValues(t, expected, result)
	result, err = f.SearchSheet("Sheet1", "A")
	assert.NoError(t, err)
	assert.EqualValues(t, []string{"A1"}, result)
	// Test search the coordinates where the numerical value in the range of
	// "0-9" of Sheet1 is described by regular expression:
	result, err = f.SearchSheet("Sheet1", "[0-9]", true)
	assert.NoError(t, err)
	assert.EqualValues(t, expected, result)

	// Test search worksheet data after set cell value
	f = excelize.NewFile()
	assert.NoError(t, f.SetCellValue("Sheet1", "A1", true))
	_, err = f.SearchSheet("Sheet1", "")
	assert.NoError(t, err)
}

func TestSetPageLayout(t *testing.T) {
	f := excelize.NewFile()
	// Test set page layout on not exists worksheet.
	assert.EqualError(t, f.SetPageLayout("SheetN"), "sheet SheetN is not exist")
}

func TestGetPageLayout(t *testing.T) {
	f := excelize.NewFile()
	// Test get page layout on not exists worksheet.
	assert.EqualError(t, f.GetPageLayout("SheetN"), "sheet SheetN is not exist")
}

func TestSetHeaderFooter(t *testing.T) {
	f := excelize.NewFile()
	assert.NoError(t, f.SetCellStr("Sheet1", "A1", "Test SetHeaderFooter"))
	// Test set header and footer on not exists worksheet.
	assert.EqualError(t, f.SetHeaderFooter("SheetN", nil), "sheet SheetN is not exist")
	// Test set header and footer with illegal setting.
	assert.EqualError(t, f.SetHeaderFooter("Sheet1", &excelize.FormatHeaderFooter{
		OddHeader: strings.Repeat("c", 256),
	}), "field OddHeader must be less than 255 characters")

	assert.NoError(t, f.SetHeaderFooter("Sheet1", nil))
	assert.NoError(t, f.SetHeaderFooter("Sheet1", &excelize.FormatHeaderFooter{
		DifferentFirst:   true,
		DifferentOddEven: true,
		OddHeader:        "&R&P",
		OddFooter:        "&C&F",
		EvenHeader:       "&L&P",
		EvenFooter:       "&L&D&R&T",
		FirstHeader:      `&CCenter &"-,Bold"Bold&"-,Regular"HeaderU+000A&D`,
	}))
	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestSetHeaderFooter.xlsx")))
}

func TestDefinedName(t *testing.T) {
	f := excelize.NewFile()
	assert.NoError(t, f.SetDefinedName(&excelize.DefinedName{
		Name:     "Amount",
		RefersTo: "Sheet1!$A$2:$D$5",
		Comment:  "defined name comment",
		Scope:    "Sheet1",
	}))
	assert.NoError(t, f.SetDefinedName(&excelize.DefinedName{
		Name:     "Amount",
		RefersTo: "Sheet1!$A$2:$D$5",
		Comment:  "defined name comment",
	}))
	assert.EqualError(t, f.SetDefinedName(&excelize.DefinedName{
		Name:     "Amount",
		RefersTo: "Sheet1!$A$2:$D$5",
		Comment:  "defined name comment",
	}), "the same name already exists on the scope")
	assert.EqualError(t, f.DeleteDefinedName(&excelize.DefinedName{
		Name: "No Exist Defined Name",
	}), "no defined name on the scope")
	assert.Exactly(t, "Sheet1!$A$2:$D$5", f.GetDefinedName()[1].RefersTo)
	assert.NoError(t, f.DeleteDefinedName(&excelize.DefinedName{
		Name: "Amount",
	}))
	assert.Exactly(t, "Sheet1!$A$2:$D$5", f.GetDefinedName()[0].RefersTo)
	assert.Exactly(t, 1, len(f.GetDefinedName()))
	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestDefinedName.xlsx")))
}

func TestGroupSheets(t *testing.T) {
	f := excelize.NewFile()
	sheets := []string{"Sheet2", "Sheet3"}
	for _, sheet := range sheets {
		f.NewSheet(sheet)
	}
	assert.EqualError(t, f.GroupSheets([]string{"Sheet1", "SheetN"}), "sheet SheetN is not exist")
	assert.EqualError(t, f.GroupSheets([]string{"Sheet2", "Sheet3"}), "group worksheet must contain an active worksheet")
	assert.NoError(t, f.GroupSheets([]string{"Sheet1", "Sheet2"}))
	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestGroupSheets.xlsx")))
}

func TestUngroupSheets(t *testing.T) {
	f := excelize.NewFile()
	sheets := []string{"Sheet2", "Sheet3", "Sheet4", "Sheet5"}
	for _, sheet := range sheets {
		f.NewSheet(sheet)
	}
	assert.NoError(t, f.UngroupSheets())
}

func TestInsertPageBreak(t *testing.T) {
	f := excelize.NewFile()
	assert.NoError(t, f.InsertPageBreak("Sheet1", "A1"))
	assert.NoError(t, f.InsertPageBreak("Sheet1", "B2"))
	assert.NoError(t, f.InsertPageBreak("Sheet1", "C3"))
	assert.NoError(t, f.InsertPageBreak("Sheet1", "C3"))
	assert.EqualError(t, f.InsertPageBreak("Sheet1", "A"), `cannot convert cell "A" to coordinates: invalid cell name "A"`)
	assert.EqualError(t, f.InsertPageBreak("SheetN", "C3"), "sheet SheetN is not exist")
	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestInsertPageBreak.xlsx")))
}

func TestRemovePageBreak(t *testing.T) {
	f := excelize.NewFile()
	assert.NoError(t, f.RemovePageBreak("Sheet1", "A2"))

	assert.NoError(t, f.InsertPageBreak("Sheet1", "A2"))
	assert.NoError(t, f.InsertPageBreak("Sheet1", "B2"))
	assert.NoError(t, f.RemovePageBreak("Sheet1", "A1"))
	assert.NoError(t, f.RemovePageBreak("Sheet1", "B2"))

	assert.NoError(t, f.InsertPageBreak("Sheet1", "C3"))
	assert.NoError(t, f.RemovePageBreak("Sheet1", "C3"))

	assert.NoError(t, f.InsertPageBreak("Sheet1", "A3"))
	assert.NoError(t, f.RemovePageBreak("Sheet1", "B3"))
	assert.NoError(t, f.RemovePageBreak("Sheet1", "A3"))

	f.NewSheet("Sheet2")
	assert.NoError(t, f.InsertPageBreak("Sheet2", "B2"))
	assert.NoError(t, f.InsertPageBreak("Sheet2", "C2"))
	assert.NoError(t, f.RemovePageBreak("Sheet2", "B2"))

	assert.EqualError(t, f.RemovePageBreak("Sheet1", "A"), `cannot convert cell "A" to coordinates: invalid cell name "A"`)
	assert.EqualError(t, f.RemovePageBreak("SheetN", "C3"), "sheet SheetN is not exist")
	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestRemovePageBreak.xlsx")))
}

func TestGetSheetName(t *testing.T) {
	f, _ := excelize.OpenFile(filepath.Join("test", "Book1.xlsx"))
	assert.Equal(t, "Sheet1", f.GetSheetName(1))
	assert.Equal(t, "Sheet2", f.GetSheetName(2))
	assert.Equal(t, "", f.GetSheetName(0))
	assert.Equal(t, "", f.GetSheetName(3))
}

func TestGetSheetMap(t *testing.T) {
	expectedMap := map[int]string{
		1: "Sheet1",
		2: "Sheet2",
	}
	f, _ := excelize.OpenFile(filepath.Join("test", "Book1.xlsx"))
	sheetMap := f.GetSheetMap()
	for idx, name := range sheetMap {
		assert.Equal(t, expectedMap[idx], name)
	}
	assert.Equal(t, len(sheetMap), 2)
}
