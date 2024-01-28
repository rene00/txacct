package dataimporter

import (
	"context"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Cell struct {
	Column string
	Row    int
	Value  string
}

type Row struct {
	Cells []Cell
}

func (r Row) GetCellValueWithCol(col string) string {
	for _, cell := range r.Cells {
		if strings.ToLower(cell.Column) == strings.ToLower(col) {
			return cell.Value
		}
	}
	return ""
}

type ExcelData struct {
	Rows []Row
}

func ProcessExcelData(ctx context.Context, c chan Row, filepath, worksheet string) error {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			return
		}
	}()

	rows, err := f.Rows(worksheet)
	if err != nil {
		return err
	}

	type rowDat map[string]string
	type allDat []rowDat

	headers := []string{}
	rowCount := 0

	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return err
		}
		rowCount++

		r := Row{}

		for idx, colCell := range row {
			if rowCount == 1 {
				headers = append(headers, colCell)
				continue
			}

			cell := Cell{
				Column: headers[idx],
				Row:    rowCount,
				Value:  colCell,
			}

			r.Cells = append(r.Cells, cell)
		}

		if len(r.Cells) == 0 {
			continue
		}

		select {
		case c <- r:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if err = rows.Close(); err != nil {
		return err
	}

	return nil
}
