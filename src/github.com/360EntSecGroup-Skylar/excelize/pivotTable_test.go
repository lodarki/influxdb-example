package excelize

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddPivotTable(t *testing.T) {
	f := NewFile()
	// Create some data in a sheet
	month := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	year := []int{2017, 2018, 2019}
	types := []string{"Meat", "Dairy", "Beverages", "Produce"}
	region := []string{"East", "West", "North", "South"}
	assert.NoError(t, f.SetSheetRow("Sheet1", "A1", &[]string{"Month", "Year", "Type", "Sales", "Region"}))
	for i := 0; i < 30; i++ {
		assert.NoError(t, f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), month[rand.Intn(12)]))
		assert.NoError(t, f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), year[rand.Intn(3)]))
		assert.NoError(t, f.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), types[rand.Intn(4)]))
		assert.NoError(t, f.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), rand.Intn(5000)))
		assert.NoError(t, f.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), region[rand.Intn(4)]))
	}
	assert.NoError(t, f.AddPivotTable(&PivotTableOption{
		DataRange:       "Sheet1!$A$1:$E$31",
		PivotTableRange: "Sheet1!$G$2:$M$34",
		Rows:            []PivotTableField{{Data: "Month"}, {Data: "Year"}},
		Columns:         []PivotTableField{{Data: "Type"}},
		Data:            []PivotTableField{{Data: "Sales", Subtotal: "Sum", Name: "Summarize by Sum"}},
	}))
	// Use different order of coordinate tests
	assert.NoError(t, f.AddPivotTable(&PivotTableOption{
		DataRange:       "Sheet1!$A$1:$E$31",
		PivotTableRange: "Sheet1!$U$34:$O$2",
		Rows:            []PivotTableField{{Data: "Month"}, {Data: "Year"}},
		Columns:         []PivotTableField{{Data: "Type"}},
		Data:            []PivotTableField{{Data: "Sales", Subtotal: "Average", Name: "Summarize by Average"}},
	}))

	assert.NoError(t, f.AddPivotTable(&PivotTableOption{
		DataRange:       "Sheet1!$A$1:$E$31",
		PivotTableRange: "Sheet1!$W$2:$AC$34",
		Rows:            []PivotTableField{{Data: "Month"}, {Data: "Year"}},
		Columns:         []PivotTableField{{Data: "Region"}},
		Data:            []PivotTableField{{Data: "Sales", Subtotal: "Count", Name: "Summarize by Count"}},
	}))
	assert.NoError(t, f.AddPivotTable(&PivotTableOption{
		DataRange:       "Sheet1!$A$1:$E$31",
		PivotTableRange: "Sheet1!$G$37:$W$50",
		Rows:            []PivotTableField{{Data: "Month"}},
		Columns:         []PivotTableField{{Data: "Region"}, {Data: "Year"}},
		Data:            []PivotTableField{{Data: "Sales", Subtotal: "CountNums", Name: "Summarize by CountNums"}},
	}))
	assert.NoError(t, f.AddPivotTable(&PivotTableOption{
		DataRange:       "Sheet1!$A$1:$E$31",
		PivotTableRange: "Sheet1!$AE$2:$AG$33",
		Rows:            []PivotTableField{{Data: "Month"}, {Data: "Year"}},
		Data:            []PivotTableField{{Data: "Sales", Subtotal: "Max", Name: "Summarize by Max"}},
	}))
	f.NewSheet("Sheet2")
	assert.NoError(t, f.AddPivotTable(&PivotTableOption{
		DataRange:       "Sheet1!$A$1:$E$31",
		PivotTableRange: "Sheet2!$A$1:$AR$15",
		Rows:            []PivotTableField{{Data: "Month"}},
		Columns:         []PivotTableField{{Data: "Region"}, {Data: "Type"}, {Data: "Year"}},
		Data:            []PivotTableField{{Data: "Sales", Subtotal: "Min", Name: "Summarize by Min"}},
	}))
	assert.NoError(t, f.AddPivotTable(&PivotTableOption{
		DataRange:       "Sheet1!$A$1:$E$31",
		PivotTableRange: "Sheet2!$A$18:$AR$54",
		Rows:            []PivotTableField{{Data: "Month"}, {Data: "Type"}},
		Columns:         []PivotTableField{{Data: "Region"}, {Data: "Year"}},
		Data:            []PivotTableField{{Data: "Sales", Subtotal: "Product", Name: "Summarize by Product"}},
	}))

	// Test empty pivot table options
	assert.EqualError(t, f.AddPivotTable(nil), "parameter is required")
	// Test invalid data range
	assert.EqualError(t, f.AddPivotTable(&PivotTableOption{
		DataRange:       "Sheet1!$A$1:$A$1",
		PivotTableRange: "Sheet1!$U$34:$O$2",
		Rows:            []PivotTableField{{Data: "Month"}, {Data: "Year"}},
		Columns:         []PivotTableField{{Data: "Type"}},
		Data:            []PivotTableField{{Data: "Sales"}},
	}), `parameter 'DataRange' parsing error: parameter is invalid`)
	// Test the data range of the worksheet that is not declared
	assert.EqualError(t, f.AddPivotTable(&PivotTableOption{
		DataRange:       "$A$1:$E$31",
		PivotTableRange: "Sheet1!$U$34:$O$2",
		Rows:            []PivotTableField{{Data: "Month"}, {Data: "Year"}},
		Columns:         []PivotTableField{{Data: "Type"}},
		Data:            []PivotTableField{{Data: "Sales"}},
	}), `parameter 'DataRange' parsing error: parameter is invalid`)
	// Test the worksheet declared in the data range does not exist
	assert.EqualError(t, f.AddPivotTable(&PivotTableOption{
		DataRange:       "SheetN!$A$1:$E$31",
		PivotTableRange: "Sheet1!$U$34:$O$2",
		Rows:            []PivotTableField{{Data: "Month"}, {Data: "Year"}},
		Columns:         []PivotTableField{{Data: "Type"}},
		Data:            []PivotTableField{{Data: "Sales"}},
	}), "sheet SheetN is not exist")
	// Test the pivot table range of the worksheet that is not declared
	assert.EqualError(t, f.AddPivotTable(&PivotTableOption{
		DataRange:       "Sheet1!$A$1:$E$31",
		PivotTableRange: "$U$34:$O$2",
		Rows:            []PivotTableField{{Data: "Month"}, {Data: "Year"}},
		Columns:         []PivotTableField{{Data: "Type"}},
		Data:            []PivotTableField{{Data: "Sales"}},
	}), `parameter 'PivotTableRange' parsing error: parameter is invalid`)
	// Test the worksheet declared in the pivot table range does not exist
	assert.EqualError(t, f.AddPivotTable(&PivotTableOption{
		DataRange:       "Sheet1!$A$1:$E$31",
		PivotTableRange: "SheetN!$U$34:$O$2",
		Rows:            []PivotTableField{{Data: "Month"}, {Data: "Year"}},
		Columns:         []PivotTableField{{Data: "Type"}},
		Data:            []PivotTableField{{Data: "Sales"}},
	}), "sheet SheetN is not exist")
	// Test not exists worksheet in data range
	assert.EqualError(t, f.AddPivotTable(&PivotTableOption{
		DataRange:       "SheetN!$A$1:$E$31",
		PivotTableRange: "Sheet1!$U$34:$O$2",
		Rows:            []PivotTableField{{Data: "Month"}, {Data: "Year"}},
		Columns:         []PivotTableField{{Data: "Type"}},
		Data:            []PivotTableField{{Data: "Sales"}},
	}), "sheet SheetN is not exist")
	// Test invalid row number in data range
	assert.EqualError(t, f.AddPivotTable(&PivotTableOption{
		DataRange:       "Sheet1!$A$0:$E$31",
		PivotTableRange: "Sheet1!$U$34:$O$2",
		Rows:            []PivotTableField{{Data: "Month"}, {Data: "Year"}},
		Columns:         []PivotTableField{{Data: "Type"}},
		Data:            []PivotTableField{{Data: "Sales"}},
	}), `parameter 'DataRange' parsing error: cannot convert cell "A0" to coordinates: invalid cell name "A0"`)
	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestAddPivotTable1.xlsx")))
	// Test with field names that exceed the length limit and invalid subtotal
	assert.NoError(t, f.AddPivotTable(&PivotTableOption{
		DataRange:       "Sheet1!$A$1:$E$31",
		PivotTableRange: "Sheet1!$G$2:$M$34",
		Rows:            []PivotTableField{{Data: "Month"}, {Data: "Year"}},
		Columns:         []PivotTableField{{Data: "Type"}},
		Data:            []PivotTableField{{Data: "Sales", Subtotal: "-", Name: strings.Repeat("s", 256)}},
	}))

	// Test adjust range with invalid range
	_, _, err := f.adjustRange("")
	assert.EqualError(t, err, "parameter is required")
	// Test get pivot fields order with empty data range
	_, err = f.getPivotFieldsOrder("")
	assert.EqualError(t, err, `parameter 'DataRange' parsing error: parameter is required`)
	// Test add pivot cache with empty data range
	assert.EqualError(t, f.addPivotCache(0, "", &PivotTableOption{}, nil), "parameter 'DataRange' parsing error: parameter is required")
	// Test add pivot cache with invalid data range
	assert.EqualError(t, f.addPivotCache(0, "", &PivotTableOption{
		DataRange:       "$A$1:$E$31",
		PivotTableRange: "Sheet1!$U$34:$O$2",
		Rows:            []PivotTableField{{Data: "Month"}, {Data: "Year"}},
		Columns:         []PivotTableField{{Data: "Type"}},
		Data:            []PivotTableField{{Data: "Sales"}},
	}, nil), "parameter 'DataRange' parsing error: parameter is invalid")
	// Test add pivot table with empty options
	assert.EqualError(t, f.addPivotTable(0, 0, "", &PivotTableOption{}), "parameter 'PivotTableRange' parsing error: parameter is required")
	// Test add pivot table with invalid data range
	assert.EqualError(t, f.addPivotTable(0, 0, "", &PivotTableOption{}), "parameter 'PivotTableRange' parsing error: parameter is required")
	// Test add pivot fields with empty data range
	assert.EqualError(t, f.addPivotFields(nil, &PivotTableOption{
		DataRange:       "$A$1:$E$31",
		PivotTableRange: "Sheet1!$U$34:$O$2",
		Rows:            []PivotTableField{{Data: "Month"}, {Data: "Year"}},
		Columns:         []PivotTableField{{Data: "Type"}},
		Data:            []PivotTableField{{Data: "Sales"}},
	}), `parameter 'DataRange' parsing error: parameter is invalid`)
	// Test get pivot fields index with empty data range
	_, err = f.getPivotFieldsIndex([]PivotTableField{}, &PivotTableOption{})
	assert.EqualError(t, err, `parameter 'DataRange' parsing error: parameter is required`)
}

func TestInStrSlice(t *testing.T) {
	assert.EqualValues(t, -1, inStrSlice([]string{}, ""))
}

func TestGetPivotTableFieldName(t *testing.T) {
	f := NewFile()
	f.getPivotTableFieldName("-", []PivotTableField{})
}
