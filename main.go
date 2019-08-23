package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kshvakov/clickhouse"
	"go/format"
	"log"
	"os"
	"strings"
	"text/template"
)

//default
const defaultModelName = "AutoGenerateModel"

//init
func init() {
	flag.StringVar(&Path, "path", "./", "path to output folder. example --path=./model")
	flag.StringVar(&Pack, "pack", "model", "go package. example --pack=model")
	flag.StringVar(&Table, "table", "simple", "target table. example --table=simple_table")
	flag.StringVar(&SettingFile, "c", "./config.yaml", "path to config file. example --c=config.yaml")
	flag.Parse()

	if err := LoadSettings(SettingFile); err != nil {
		log.Fatal(err)
	}

	fmt.Println("==========================")
	fmt.Printf("%vPath: %v\nPacket: %s\nTable: %s%v\n",
		COLOR_GREEN,
		Path,
		Pack,
		Table,
		COLOR_RESET)
	fmt.Println("==========================")
}

//main
func main() {
	connect, err := sqlx.Open("clickhouse", buildConnect())
	if err != nil {
		log.Fatal(err)
	}
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			printErr("clickhouse", err)
			printErr("please check you connection data in file", SettingFile)
		}
		return
	}

	var col []column

	query := fmt.Sprintf("SELECT database, table, name, type, default_kind, default_expression FROM system.columns where table = '%s'", Table)
	if err := connect.Select(&col, query); err != nil {
		printErr(err)
		return
	}

	if len(col) == 0 {
		fmt.Printf("%v✘ Table %s not found in %s database :( %v\n", COLOR_RED, Table, ClickHouse.BaseName, COLOR_RESET)
		return
	}

	generate(col)
}

//column
type column struct {
	Database    string `db:"database"`
	Table       string `db:"table"`
	Name        string `db:"name"`
	Typed       string `db:"type"`
	DefaultKind string `db:"default_kind"`
	DefaultExp  string `db:"default_expression"`
}

//Opt
type Opt struct {
	Name       string
	Type       string
	ColumnName string
}

//data
type data struct {
	TableName   string
	BaseName    string
	ModelName   string
	Package     string
	Opts        []Opt
	Imports     []string
	InsertQuery string
}

//generate
func generate(col []column) {
	var d data
	d.ModelName = defaultModelName
	d.Package = Pack
	d.Opts = make([]Opt, 0, len(col))
	impr := map[string]struct{}{}

	for k, v := range col {
		if k == 0 {
			d.BaseName = v.Database
			d.TableName = v.Table
			d.ModelName = name(v.Table)
		}
		tcast, ok := cast(v.Typed)
		if !ok {
			continue
		}

		if len(strings.TrimSpace(v.DefaultKind)) != 0 && strings.TrimSpace(v.Typed) != "DateTime" {
			continue
		}

		if strings.EqualFold(v.Typed, "Date") || strings.EqualFold(v.Typed, "DateTime") {
			impr["time"] = struct{}{}
		}

		d.Opts = append(d.Opts, Opt{
			Name:       name(v.Name),
			Type:       tcast,
			ColumnName: v.Name,
		})
	}

	for k, _ := range impr {
		d.Imports = append(d.Imports, k)
	}

	colsQuery := make([]string, 0, len(d.Opts))
	colsArgs := make([]string, 0, len(d.Opts))
	for _, v := range d.Opts {
		colsQuery = append(colsQuery, v.ColumnName)
		colsArgs = append(colsArgs, "?")
	}

	d.InsertQuery = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", d.TableName, strings.Join(colsQuery, ","), strings.Join(colsArgs, ","))

	t := template.Must(template.New("model").Parse(temp))

	var buff bytes.Buffer
	err := t.Execute(&buff, d)
	if err != nil {
		printErr("execute: ", err)
		return
	}
	content, err := format.Source(buff.Bytes())
	if err != nil {
		printErr("Format error: ", err)
		return
	}

	path, err := mkdir()
	if err != nil {
		printErr(err)
		return
	}
	path = fmt.Sprintf("%s%s.model.go", path, strings.ToLower(strings.Replace(d.TableName, "_", ".", -1)))
	f, err := os.Create(path)
	if err != nil {
		printErr("create file: ", err)
		return
	}

	if _, err := f.Write(content); err != nil {
		printErr(err)
		return
	}
	defer f.Close()

	fmt.Printf("%vGenerate completed √: %s%v\n", COLOR_GREEN, path, COLOR_RESET)
}

func mkdir() (string, error) {
	if !strings.HasPrefix(Path, "/") && !strings.HasPrefix(Path, "./") {
		Path = "./" + Path
	}
	if err := os.MkdirAll(Path, os.ModePerm); err != nil {
		return Path, fmt.Errorf("mkdir %s error: %s", Path, err)
	}
	if !strings.HasSuffix(Path, "/") {
		Path += "/"
	}
	return Path, nil
}

func name(t string) string {
	return strings.Replace(strings.Title(strings.Replace(t, "_", ".", -1)), ".", "", -1)
}

func printErr(a ...interface{}) {
	fmt.Printf("%v✘ %v%v\n", COLOR_RED, a, COLOR_RESET)
}
