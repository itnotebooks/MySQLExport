// Copyright 2021 EricWinn
// Author:   Eric Winn
// Email:    eng.eric.winn@gmail.com
// Time:     2021/10/12 14:55
// File:     csv.py
// Software: GoLand

package csv

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"
)

type Converter struct {
	Headers         []string  // 列头字段
	WriteHeaders    bool      // 是否写入列头
	TimeFormat      string    // 时间转换
	FloatFormat     string    // Float类型转换
	Delimiter       rune      // 分隔符
	rows            *sql.Rows // SQL查询结果
	rowPreProcessor CsvPreProcessorFunc
}

type CsvPreProcessorFunc func(row []string, columnNames []string) (outputRow bool, processedRow []string)

// WriteFile 将SQL结果写入到CSV文件
func WriteFile(fileName string, rows *sql.Rows) error {
	return New(rows).WriteFile(fileName)
}

func (c *Converter) SetRowPreProcessor(processor CsvPreProcessorFunc) {
	c.rowPreProcessor = processor
}

func (c Converter) WriteFile(csvFileName string) error {
	f, err := os.Create(csvFileName)
	if err != nil {
		return err
	}

	err = c.Write(f)
	if err != nil {
		f.Close()
		return err
	}

	return f.Close()
}

// Write 数据写入到CSV文件的具体实现
func (c Converter) Write(writer io.Writer) error {
	rows := c.rows

	// 新建一个Writer对象
	csvWriter := csv.NewWriter(writer)

	// 分割符
	if c.Delimiter != '\x00' {
		csvWriter.Comma = c.Delimiter
	}

	// 列字段
	columnNames, err := rows.Columns()
	if err != nil {
		return err
	}

	// 是否需要列头
	if c.WriteHeaders {
		var headers []string
		if len(c.Headers) > 0 {
			headers = c.Headers
		} else {
			headers = columnNames
		}
		// 写入列头
		err = csvWriter.Write(headers)
		if err != nil {
			return fmt.Errorf("failed to write headers: %w", err)
		}
	}

	count := len(columnNames)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		row := make([]string, count)

		for i, _ := range columnNames {
			valuePtrs[i] = &values[i]
		}

		if err = rows.Scan(valuePtrs...); err != nil {
			return err
		}

		// 遍历每个字段的类型
		for i, _ := range columnNames {
			var value interface{}
			rawValue := values[i]

			byteArray, ok := rawValue.([]byte)
			if ok {
				value = string(byteArray)
			} else {
				value = rawValue
			}

			float64Value, ok := value.(float64)
			if ok && c.FloatFormat != "" {
				value = fmt.Sprintf(c.FloatFormat, float64Value)
			} else {
				float32Value, ok := value.(float32)
				if ok && c.FloatFormat != "" {
					value = fmt.Sprintf(c.FloatFormat, float32Value)
				}
			}

			timeValue, ok := value.(time.Time)
			if ok && c.TimeFormat != "" {
				value = timeValue.Format(c.TimeFormat)
			}

			if value == nil {
				row[i] = ""
			} else {
				row[i] = fmt.Sprintf("%v", value)
			}
		}

		writeRow := true
		if c.rowPreProcessor != nil {
			writeRow, row = c.rowPreProcessor(row, columnNames)
		}

		// 写入数据
		if writeRow {
			err = csvWriter.Write(row)
			if err != nil {
				return fmt.Errorf("failed to write data row to csv %w", err)
			}
		}
	}
	err = rows.Err()
	// 刷到磁盘
	csvWriter.Flush()

	return err
}

func New(rows *sql.Rows) *Converter {
	return &Converter{
		rows:         rows,
		WriteHeaders: true,
		Delimiter:    ',',
	}
}
