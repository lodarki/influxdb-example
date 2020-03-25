package excelize

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStyleFill(t *testing.T) {
	cases := []struct {
		label      string
		format     string
		expectFill bool
	}{{
		label:      "no_fill",
		format:     `{"alignment":{"wrap_text":true}}`,
		expectFill: false,
	}, {
		label:      "fill",
		format:     `{"fill":{"type":"pattern","pattern":1,"color":["#000000"]}}`,
		expectFill: true,
	}}

	for _, testCase := range cases {
		xl := NewFile()
		styleID, err := xl.NewStyle(testCase.format)
		if err != nil {
			t.Fatal(err)
		}

		styles := xl.stylesReader()
		style := styles.CellXfs.Xf[styleID]
		if testCase.expectFill {
			assert.NotEqual(t, style.FillID, 0, testCase.label)
		} else {
			assert.Equal(t, style.FillID, 0, testCase.label)
		}
	}
}

func TestSetConditionalFormat(t *testing.T) {
	cases := []struct {
		label  string
		format string
		rules  []*xlsxCfRule
	}{{
		label: "3_color_scale",
		format: `[{
			"type":"3_color_scale",
			"criteria":"=",
			"min_type":"num",
			"mid_type":"num",
			"max_type":"num",
			"min_value": "-10",
			"mid_value": "0",
			"max_value": "10",
			"min_color":"ff0000",
			"mid_color":"00ff00",
			"max_color":"0000ff"
		}]`,
		rules: []*xlsxCfRule{{
			Priority: 1,
			Type:     "colorScale",
			ColorScale: &xlsxColorScale{
				Cfvo: []*xlsxCfvo{{
					Type: "num",
					Val:  "-10",
				}, {
					Type: "num",
					Val:  "0",
				}, {
					Type: "num",
					Val:  "10",
				}},
				Color: []*xlsxColor{{
					RGB: "FFFF0000",
				}, {
					RGB: "FF00FF00",
				}, {
					RGB: "FF0000FF",
				}},
			},
		}},
	}, {
		label: "3_color_scale default min/mid/max",
		format: `[{
			"type":"3_color_scale",
			"criteria":"=",
			"min_type":"num",
			"mid_type":"num",
			"max_type":"num",
			"min_color":"ff0000",
			"mid_color":"00ff00",
			"max_color":"0000ff"
		}]`,
		rules: []*xlsxCfRule{{
			Priority: 1,
			Type:     "colorScale",
			ColorScale: &xlsxColorScale{
				Cfvo: []*xlsxCfvo{{
					Type: "num",
					Val:  "0",
				}, {
					Type: "num",
					Val:  "50",
				}, {
					Type: "num",
					Val:  "0",
				}},
				Color: []*xlsxColor{{
					RGB: "FFFF0000",
				}, {
					RGB: "FF00FF00",
				}, {
					RGB: "FF0000FF",
				}},
			},
		}},
	}, {
		label: "2_color_scale default min/max",
		format: `[{
			"type":"2_color_scale",
			"criteria":"=",
			"min_type":"num",
			"max_type":"num",
			"min_color":"ff0000",
			"max_color":"0000ff"
		}]`,
		rules: []*xlsxCfRule{{
			Priority: 1,
			Type:     "colorScale",
			ColorScale: &xlsxColorScale{
				Cfvo: []*xlsxCfvo{{
					Type: "num",
					Val:  "0",
				}, {
					Type: "num",
					Val:  "0",
				}},
				Color: []*xlsxColor{{
					RGB: "FFFF0000",
				}, {
					RGB: "FF0000FF",
				}},
			},
		}},
	}}

	for _, testCase := range cases {
		xl := NewFile()
		const sheet = "Sheet1"
		const cellRange = "A1:A1"

		err := xl.SetConditionalFormat(sheet, cellRange, testCase.format)
		if err != nil {
			t.Fatalf("%s", err)
		}

		xlsx, err := xl.workSheetReader(sheet)
		assert.NoError(t, err)
		cf := xlsx.ConditionalFormatting
		assert.Len(t, cf, 1, testCase.label)
		assert.Len(t, cf[0].CfRule, 1, testCase.label)
		assert.Equal(t, cellRange, cf[0].SQRef, testCase.label)
		assert.EqualValues(t, testCase.rules, cf[0].CfRule, testCase.label)
	}
}

func TestUnsetConditionalFormat(t *testing.T) {
	f := NewFile()
	assert.NoError(t, f.SetCellValue("Sheet1", "A1", 7))
	assert.NoError(t, f.UnsetConditionalFormat("Sheet1", "A1:A10"))
	format, err := f.NewConditionalStyle(`{"font":{"color":"#9A0511"},"fill":{"type":"pattern","color":["#FEC7CE"],"pattern":1}}`)
	assert.NoError(t, err)
	assert.NoError(t, f.SetConditionalFormat("Sheet1", "A1:A10", fmt.Sprintf(`[{"type":"cell","criteria":">","format":%d,"value":"6"}]`, format)))
	assert.NoError(t, f.UnsetConditionalFormat("Sheet1", "A1:A10"))
	// Test unset conditional format on not exists worksheet.
	assert.EqualError(t, f.UnsetConditionalFormat("SheetN", "A1:A10"), "sheet SheetN is not exist")
	// Save xlsx file by the given path.
	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestUnsetConditionalFormat.xlsx")))
}

func TestNewStyle(t *testing.T) {
	f := NewFile()
	styleID, err := f.NewStyle(`{"font":{"bold":true,"italic":true,"family":"Times New Roman","size":36,"color":"#777777"}}`)
	assert.NoError(t, err)
	styles := f.stylesReader()
	fontID := styles.CellXfs.Xf[styleID].FontID
	font := styles.Fonts.Font[fontID]
	assert.Contains(t, *font.Name.Val, "Times New Roman", "Stored font should contain font name")
	assert.Equal(t, 2, styles.CellXfs.Count, "Should have 2 styles")
	_, err = f.NewStyle(&Style{})
	assert.NoError(t, err)
}

func TestGetDefaultFont(t *testing.T) {
	f := NewFile()
	s := f.GetDefaultFont()
	assert.Equal(t, s, "Calibri", "Default font should be Calibri")
}

func TestSetDefaultFont(t *testing.T) {
	f := NewFile()
	f.SetDefaultFont("Ariel")
	styles := f.stylesReader()
	s := f.GetDefaultFont()
	assert.Equal(t, s, "Ariel", "Default font should change to Ariel")
	assert.Equal(t, *styles.CellStyles.CellStyle[0].CustomBuiltIn, true)
}

func TestStylesReader(t *testing.T) {
	f := NewFile()
	// Test read styles with unsupport charset.
	f.Styles = nil
	f.XLSX["xl/styles.xml"] = MacintoshCyrillicCharset
	assert.EqualValues(t, new(xlsxStyleSheet), f.stylesReader())
}

func TestThemeReader(t *testing.T) {
	f := NewFile()
	// Test read theme with unsupport charset.
	f.XLSX["xl/theme/theme1.xml"] = MacintoshCyrillicCharset
	assert.EqualValues(t, new(xlsxTheme), f.themeReader())
}

func TestSetCellStyle(t *testing.T) {
	f := NewFile()
	// Test set cell style on not exists worksheet.
	assert.EqualError(t, f.SetCellStyle("SheetN", "A1", "A2", 1), "sheet SheetN is not exist")
}
