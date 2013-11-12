/**
 * Created with IntelliJ IDEA.
 * User: dustinhiatt
 * Date: 11/5/13
 * Time: 3:13 PM
 * To change this template use File | Settings | File Templates.
 */
package table

import (
	"strconv"
	"regexp"
	"strings"
	"unicode/utf8"
	"math"
	"encoding/json"
//	"log"
)

const (
	letters = "ABCDEFGHIJKLMONQRSTUVWXYZ"
)

type cellrange struct {
	ISerializable
	StartRow		int
	StopRow			int
	StartColumn		int
	StopColumn		int
	TableId			string
}

func (cr *cellrange) ToBytes() []byte {
	res, err := json.Marshal(cr)
	if err != nil {
		return nil
	}

	return res
}

func MakeRange(xrange string) *cellrange {
	cr := new(cellrange)
	rangeParts := strings.Split(xrange, ":")
	if len(rangeParts) == 3 {
		cr.TableId = rangeParts[0]
		rangeParts = rangeParts[1:]
	}
	startParts := getStringPartsFromAlphaNumeric(rangeParts[0])
	startRow, startColumn := parseAlphaNumericParts(startParts)
	stopParts := getStringPartsFromAlphaNumeric(rangeParts[1])
	stopRow, stopColumn := parseAlphaNumericParts(stopParts)
	cr.StartRow = startRow
	cr.StartColumn = startColumn
	cr.StopRow = stopRow + 1
	cr.StopColumn = stopColumn + 1
	return cr
}

func MakeRangeFromBytes(bytes []byte) *cellrange {
	var m cellrange
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		return nil
	}
	return &m
}

func parseFormula(value string) []string {
	funcCall := value[1:len(value)]
	funcParts := strings.Split(funcCall, "(")
	funcParts[1] = funcParts[1][:len(funcParts[1]) - 1] // remove the ')'
	return funcParts
}

func getNumberFromAlpha(alpha string) int {
	sum := 0;
	upperAlpha := strings.ToUpper(alpha)
	la := utf8.RuneCountInString(upperAlpha)
	for i := 0; i < la; i++ {
		index := strings.Index(letters, string([]rune(upperAlpha)[i])) + 1
		sum += index * int(math.Pow(26, float64((la - (i + 1)))))
	}
	return sum - 1
}

func parseAlphaNumericParts(parts []string) (int, int) {
	row, column := 0, -1
	if len(parts) == 2 {
		column = getNumberFromAlpha(parts[0])
		row64, _ := strconv.ParseInt(parts[1], 0, 32)
		row = int(row64)
	} else {
		row64, err := strconv.ParseInt(parts[0], 0, 32)
		if err != nil  {
			column = getNumberFromAlpha(parts[0])
		} else {
			row = int(row64)
		}
	}
	return row - 1, column
}

func getStringPartsFromAlphaNumeric(alpha string) []string {
	re := regexp.MustCompile("[a-zA-Z]+|\\d+")
	return re.FindAllString(alpha, -1)
}

/*
func sum(c *cell, args string) (*cellrange, string) {
	cr := MakeRange(args)
	ch := MakeValueChannel()
	if cr.tableId == "" || cr.tableId == c.table.id {
		c.table.GetRangeByCellRange(cr, ch)
	} else {
		//tableCh := MakeValueChannel()
		//c.table.orchestrator.GetTableById(cr.tableId, tableCh)
		//tableMessage := <- tableCh
		//tableMessage.table.GetRangeByCellRange(cr, ch)
	}
	message := <- ch
	tr := message.tableRange
	sum := 0.0
	for i := cr.startRow; i < cr.stopRow; i++ {
		_, ok := tr.cells[i]
		if ok {
			for j := cr.startColumn; j < cr.stopColumn; j++ {
				cell, ok := tr.cells[i][j]
				if ok {
					if cell.value == "" {
						continue
					}
					if cell.IsFloat() {
						amt, _ := cell.AsFloat()
						sum += amt
					} else if cell.IsInt() {
						amt, _ := cell.AsInt()
						sum += float64(amt)
					}
				}
			}
		}
	}
	value := strconv.FormatFloat(sum, 'f', -1, 64)
	return cr, value
}*/
