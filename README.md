# json2csv

[![GitHub license](https://img.shields.io/github/license/onozaty/json2csv)](https://github.com/onozaty/json2csv/blob/main/LICENSE)
[![Test](https://github.com/onozaty/json2csv/actions/workflows/test.yaml/badge.svg)](https://github.com/onozaty/json2csv/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/onozaty/json2csv/branch/main/graph/badge.svg?token=DPGXS4UDAP)](https://codecov.io/gh/onozaty/json2csv)

json2csv converts JSON to CSV.  
You can easily define mappings for converts using XPath.

## Usage

```
$ json2csv -i input.json -m mapping.json -o output.csv
```

The arguments are as follows.

```
Usage: json2csv [flags]

Flags
  -i, --input string     JSON input file path or directory or url
  -m, --mapping string   JSON to CSV mapping file path or url
  -o, --output string    CSV output file path
  -b, --bom              CSV with BOM
  -h, --help             Help
```

JSON and mapping files can be specified by URL.

```
json2csv -i https://github.com/onozaty/json2csv/raw/main/testdata/data/rss.json -m https://github.com/onozaty/json2csv/raw/main/testdata/mapping/rss.json -o output.csv
```

## Mapping

The conversion mapping definition is written in JSON.    
Specify the position on the JSON with XPath.

```json
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
        },
        {
            "header": "description",
            "valuePath": "/description"
        }
    ]
}
```

* `rowsPath` : XPath to get as a rows.
* `columns` : Definition of each column.
    * `header` : CSV header.
    * `valuePath` : XPath to get as a value.
    * `useEvaluate` : Specify `true` when using an expression with `valuePath`. For example, when using `sum()` or `not()`, `boolean()`.

[antchfx/xpath](https://github.com/antchfx/xpath) is used in json2csv.  
See below for supported XPath.

* https://github.com/antchfx/xpath#supported-features

For XPath in JSON, please refer to the following.

* https://github.com/antchfx/jsonquery#xpath-tests

Please refer to the sample below.

* https://github.com/onozaty/json2csv/tree/main/testdata/mapping

## Install

You can download the binary from the following.

* https://github.com/onozaty/json2csv/releases/latest

## License

MIT

## Author

[onozaty](https://github.com/onozaty)
