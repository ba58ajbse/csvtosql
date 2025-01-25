package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	filePath := flag.String("csv", "", "Path to the CSV file")
	outputFile := flag.String("sql", "", "Path to output SQL file")
	tableName := flag.String("table", "", "Insert table name")
	flag.Parse()

	// 引数が指定されていない場合のエラーハンドリング
	if *filePath == "" {
		fmt.Println("Error: Please specify the path to the CSV file using the -file flag.")
		flag.Usage() // フラグの使い方を表示
		return
	}
	// CSVファイルを開く
	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// CSVリーダーを作成
	reader := csv.NewReader(file)

	// CSVの内容をすべて読み込む
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	output, err := os.Create(*outputFile)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer output.Close()

	// 読み込んだデータを出力
	var header string
	var query string
	last := len(records) - 1
	for i, record := range records {
		if i == 0 {
			header = makeValue(record)
			query = fmt.Sprintf("INSERT INTO %s (%s)\nVALUES\n", *tableName, header)
			continue
		}
		query += fmt.Sprintf("(%s)", makeValue(record))
		if i != last {
			query += ",\n"
		} else {
			query += ";"
		}
	}

	_, err = output.WriteString(query)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}

	fmt.Println("Finish.")
}

func makeValue(records []string) string {
	var arr []string
	for _, v := range records {
		arr = append(arr, fmt.Sprintf("'%s'", v))
	}
	return strings.Join(arr, ",")
}
