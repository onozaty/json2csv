package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/antchfx/jsonquery"
	"github.com/onozaty/go-customcsv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_File(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "testdata/data/rss.json"

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "//items/*",
		"columns": [
			{
				"header": "title",
				"valuePath": "/title"
			},
			{
				"header": "link",
				"valuePath": "/link"
			}
		]
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)
	require.Empty(t, out.String())

	result := readString(t, outputPath)
	expect := joinRows(
		"title,link",
		"RSS Tutorial,https://www.w3schools.com/xml/xml_rss.asp",
		"XML Tutorial,https://www.w3schools.com/xml",
	)

	assert.Equal(t, expect, result)
}

func TestRun_URL(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "https://github.com/onozaty/json2csv/raw/main/testdata/data/rss.json"

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "//items/*",
		"columns": [
			{
				"header": "title",
				"valuePath": "/title"
			},
			{
				"header": "link",
				"valuePath": "/link"
			}
		]
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)
	require.Empty(t, out.String())

	result := readString(t, outputPath)
	expect := joinRows(
		"title,link",
		"RSS Tutorial,https://www.w3schools.com/xml/xml_rss.asp",
		"XML Tutorial,https://www.w3schools.com/xml",
	)

	assert.Equal(t, expect, result)
}

func TestRun_Dir(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "testdata/data/users" // ディレクトリ指定

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "/*",
		"columns": [
			{
				"header": "id",
				"valuePath": "/id"
			},
			{
				"header": "name",
				"valuePath": "/name"
			},
			{
				"header": "male",
				"valuePath": "boolean(/gender[text()='male'])",
				"useEvaluate": true
			}
		]
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)
	require.Empty(t, out.String())

	result := readString(t, outputPath)
	expect := joinRows(
		"id,name,male",
		"1,Jon,true",
		"2,花子,false",
		"3,Taro,true",
		"4,Kyotaro,true",
		"5,X,false",
	)

	assert.Equal(t, expect, result)
}

func TestRun_WithBom(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "testdata/data/rss.json"

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "//items/*",
		"columns": [
			{
				"header": "title",
				"valuePath": "/title"
			},
			{
				"header": "link",
				"valuePath": "/link"
			}
		]
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
			"-b",
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)
	require.Empty(t, out.String())

	result := readString(t, outputPath)
	expect := joinRows(
		"\uFEFFtitle,link",
		"RSS Tutorial,https://www.w3schools.com/xml/xml_rss.asp",
		"XML Tutorial,https://www.w3schools.com/xml",
	)

	assert.Equal(t, expect, result)
}

func TestRun_CommandParseFailed(t *testing.T) {

	// ARRANGE
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-a", // 存在しないフラグ
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := `json2csv vdev (none)

Usage: json2csv [flags]

Flags
  -i, --input string     JSON input file path or directory or url
  -m, --mapping string   JSON to CSV mapping file path or url
  -o, --output string    CSV output file path
  -b, --bom              CSV with BOM
  -h, --help             Help

unknown shorthand flag: 'a' in -a
`
	assert.Equal(t, expect, out.String())
}

func TestRun_Help(t *testing.T) {

	// ARRANGE
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-h",
		},
		out,
	)

	// ASSERT
	require.Equal(t, OK, exitCode)

	expect := `json2csv vdev (none)

Usage: json2csv [flags]

Flags
  -i, --input string     JSON input file path or directory or url
  -m, --mapping string   JSON to CSV mapping file path or url
  -o, --output string    CSV output file path
  -b, --bom              CSV with BOM
  -h, --help             Help

`
	assert.Equal(t, expect, out.String())
}

func TestRun_NoneInput(t *testing.T) {

	// ARRANGE
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-m", "xxx",
			"-o", "yyy",
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := `json2csv vdev (none)

Usage: json2csv [flags]

Flags
  -i, --input string     JSON input file path or directory or url
  -m, --mapping string   JSON to CSV mapping file path or url
  -o, --output string    CSV output file path
  -b, --bom              CSV with BOM
  -h, --help             Help

`
	assert.Equal(t, expect, out.String())
}

func TestRun_NoneMapping(t *testing.T) {

	// ARRANGE
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", "xxx",
			"-o", "yyy",
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := `json2csv vdev (none)

Usage: json2csv [flags]

Flags
  -i, --input string     JSON input file path or directory or url
  -m, --mapping string   JSON to CSV mapping file path or url
  -o, --output string    CSV output file path
  -b, --bom              CSV with BOM
  -h, --help             Help

`
	assert.Equal(t, expect, out.String())
}

func TestRun_NoneOutput(t *testing.T) {

	// ARRANGE
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", "xxx",
			"-m", "yyy",
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := `json2csv vdev (none)

Usage: json2csv [flags]

Flags
  -i, --input string     JSON input file path or directory or url
  -m, --mapping string   JSON to CSV mapping file path or url
  -o, --output string    CSV output file path
  -b, --bom              CSV with BOM
  -h, --help             Help

`
	assert.Equal(t, expect, out.String())
}

func TestRun_InputFileNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// 作成はしない
	inputPath := filepath.Join(temp, "input.json")

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "//items/*",
		"columns": [
			{
				"header": "title",
				"valuePath": "/title"
			},
			{
				"header": "link",
				"valuePath": "/link"
			}
		]
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := inputPath + " is not found\n"
	assert.Equal(t, expect, out.String())
}

func TestRun_MappingFileNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "testdata/data/rss.json"

	// 作成はしない
	mappingPath := filepath.Join(temp, "mapping.json")

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := mappingPath + " is not found\n"
	assert.Equal(t, expect, out.String())
}

func TestRun_OutputFileDirNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "testdata/data/rss.json"

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "//items/*",
		"columns": [
			{
				"header": "title",
				"valuePath": "/title"
			},
			{
				"header": "link",
				"valuePath": "/link"
			}
		]
	}`)

	// 存在しないディレクトリを親に指定
	outputPath := filepath.Join(temp, "xxx", "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	// OSによってエラーメッセージが異なるのでファイル名部分だけチェック
	expect := "open " + outputPath
	assert.Contains(t, out.String(), expect)
}

func TestRun_InvalidXPath_RowPath(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "testdata/data/rss.json"

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "item[",
		"columns": [
			{
				"header": "title",
				"valuePath": "/title"
			},
			{
				"header": "link",
				"valuePath": "/link"
			}
		]
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := "xpath 'item[' is failed: expression must evaluate to a node-set\n"
	assert.Equal(t, expect, out.String())
}

func TestRun_InvalidXPath_ValuePath(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "testdata/data/rss.json"

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "//items/*",
		"columns": [
			{
				"header": "title",
				"valuePath": "/title["
			},
			{
				"header": "link",
				"valuePath": "/link"
			}
		]
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := "xpath '/title[' is failed: expression must evaluate to a node-set\n"
	assert.Equal(t, expect, out.String())
}

func TestRun_InvalidXPath_ValuePath_UseEvaluate(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "testdata/data/rss.json"

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "//items/*",
		"columns": [
			{
				"header": "title",
				"valuePath": "/title"
			},
			{
				"header": "link",
				"valuePath": "boolean(/link",
				"useEvaluate": true
			}
		]
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := "xpath 'boolean(/link' is failed: boolean(/link has an invalid token\n"
	assert.Equal(t, expect, out.String())
}

func TestRun_InvalidJSON(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := createFile(t, temp, "input.json", `
	{
		"id": 1,
		"name"
	}
	`)

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "//items/*",
		"columns": [
			{
				"header": "name",
				"valuePath": "/name"
			},
			{
				"header": "value",
				"valuePath": "/value"
			}
		]
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := inputPath + " is failed: invalid character '}' after object key\n"
	assert.Equal(t, expect, out.String())
}

func TestRun_InvalidMappingJson(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	inputPath := "testdata/data/rss.json"

	mappingPath := createFile(t, temp, "mapping.json", `
	{
		"rowsPath": "//items/*",
	}`)

	outputPath := filepath.Join(temp, "output.csv")
	out := new(bytes.Buffer)

	// ACT
	exitCode := run(
		[]string{
			"-i", inputPath,
			"-m", mappingPath,
			"-o", outputPath,
		},
		out,
	)

	// ASSERT
	require.Equal(t, NG, exitCode)

	expect := "invalid mapping format: invalid character '}' looking for beginning of object key string\n"
	assert.Equal(t, expect, out.String())
}

func TestConvertOne(t *testing.T) {

	// ARRANGE
	json := `
	{
		"items": [
			{
				"id": 1,
				"name": "name1",
				"value": "value1"
			},
			{
				"id": 2,
				"name": "name2",
				"value": "value2,xx"
			},
			{
				"id": 3,
				"name": "name3"
			}
		]
	}`

	doc, err := jsonquery.Parse(strings.NewReader(json))
	require.NoError(t, err)

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	csv := customcsv.NewWriter(writer)

	mapping := Mapping{
		RowsPath: "//items/*",
		Columns: []Column{
			{Header: "id", ValuePath: "/id"},
			{Header: "name", ValuePath: "/name"},
			{Header: "value", ValuePath: "/value"},
			{Header: "has value", ValuePath: "boolean(/value)", UseEvaluate: true},
		},
	}

	// ACT
	err = convertOne(doc, &mapping, csv)
	csv.Flush()

	// ASSERT
	require.NoError(t, err)

	result := b.String()

	expect := joinRows(
		"1,name1,value1,true",
		"2,name2,\"value2,xx\",true",
		"3,name3,,false",
	)

	assert.Equal(t, expect, result)
}

func TestLoadMapping_File(t *testing.T) {

	// ARRANGE/ACT
	result, err := loadMapping("testdata/mapping/rss.json")

	// ASSERT
	require.NoError(t, err)

	expect := &Mapping{
		RowsPath: "//items/*",
		Columns: []Column{
			{Header: "title", ValuePath: "/title"},
			{Header: "link", ValuePath: "/link"},
			{Header: "description", ValuePath: "/description"},
		},
	}

	assert.Equal(t, expect, result)
}

func TestLoadMapping_URL(t *testing.T) {

	// ARRANGE/ACT
	result, err := loadMapping("https://github.com/onozaty/json2csv/raw/main/testdata/mapping/rss.json")

	// ASSERT
	require.NoError(t, err)

	expect := &Mapping{
		RowsPath: "//items/*",
		Columns: []Column{
			{Header: "title", ValuePath: "/title"},
			{Header: "link", ValuePath: "/link"},
			{Header: "description", ValuePath: "/description"},
		},
	}

	assert.Equal(t, expect, result)
}

func TestFindJSON_Dir(t *testing.T) {

	// ARRANGE/ACT
	result, err := findJSON("testdata/data/users")

	// ASSERT
	require.NoError(t, err)

	expect := []string{
		filepath.Join("testdata", "data", "users", "users1.json"),
		filepath.Join("testdata", "data", "users", "users2.json"),
	}

	assert.Equal(t, expect, result)
}

func TestFindJSON_Dir_Nest(t *testing.T) {

	// ARRANGE/ACT
	result, err := findJSON("testdata/data")

	// ASSERT
	require.NoError(t, err)

	expect := []string{filepath.Join("testdata", "data", "rss.json")}

	assert.Equal(t, expect, result)
}

func TestFindJSON_File(t *testing.T) {

	// ARRANGE/ACT
	result, err := findJSON("testdata/data/rss.json")

	// ASSERT
	require.NoError(t, err)

	expect := []string{"testdata/data/rss.json"}

	assert.Equal(t, expect, result)
}

func TestFindJSON_URL(t *testing.T) {

	// ARRANGE/ACT
	result, err := findJSON("https://github.com/onozaty/json2csv/raw/main/testdata/data/rss.json")

	// ASSERT
	require.NoError(t, err)

	expect := []string{"https://github.com/onozaty/json2csv/raw/main/testdata/data/rss.json"}

	assert.Equal(t, expect, result)
}

func createFile(t *testing.T, dir string, name string, content string) string {

	file, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		t.Fatal("craete file failed\n", err)
	}

	_, err = file.Write([]byte(content))
	if err != nil {
		t.Fatal("write file failed\n", err)
	}

	err = file.Close()
	if err != nil {
		t.Fatal("write file failed\n", err)
	}

	return file.Name()
}

func readBytes(t *testing.T, name string) []byte {

	bo, err := os.ReadFile(name)
	if err != nil {
		t.Fatal("read failed\n", err)
	}

	return bo
}

func readString(t *testing.T, name string) string {

	bo := readBytes(t, name)
	return string(bo)
}

func joinRows(rows ...string) string {
	return strings.Join(rows, "\r\n") + "\r\n"
}
