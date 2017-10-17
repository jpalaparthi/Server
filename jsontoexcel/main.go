package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/tealeg/xlsx"
)

func main() {

	//var s Records

	sheetData, err := ioutil.ReadFile("output.json")
	//fmt.Println(string(sheetData))
	if err != nil {
		fmt.Println(err)
	}
	s := Records{}
	JSONStringToStructure(string(sheetData), &s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s.Result[4])
	var file *xlsx.File
	var sheet *xlsx.Sheet
	//var row *xlsx.Row
	//var cell1, cell2, cell3, cell4, cell5, cell6, cell7, cell8, cell9 *xlsx.Cell
	//var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	//row = sheet.AddRow()

	AddRowWithData(sheet, s, "TransactionId", "iAccountID", "TransactionDate", "PostDate", "PlateSource", "PlateCategory", "PlateCode", "PlateNumber", "TagNumber", "Location", "Direction", "Amount")

	//  "TransactionID"=="TransactionId"
	//  "iAccountID"
	//  "TripDateTime" == "TransactionDate"
	//   "TransactionPostDate" == "PostDate"
	//  "PlateSource"=="PlateSource"
	//  "PlateCategory" value("Private")
	//  "PlateColor" == "PlateCode"
	//  "PlateNumber" == "PlateNumber"
	//  "TagNumber" == "TagNumber"
	//  "TollGateLocation" =="Location"
	//  "TollGateDirection" =="Direction"
	//  "Amount" == "Amount"

	err = file.Save("MyXLSXFile.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

type Records struct {
	Result []Record `json:"result"`
}

/*type Record struct {
	TransactionId       string `json:"TransactionId"`
	TripDateTime        string `json:"TripDateTime"`
	TransactionPostDate string `json:"TransactionPostDate"`
	TollGateLocation    string `json:"TollGateLocation"`
	TollGateDirection   string `json:"TollGateDirection"`
	Amount              string `json:"Amount"`
	PlateSource         string `json:"PlateSource"`
	PlateColor          string `json:"PlateColor"`
	PlateNumber         string `json:"PlateNumber"`
	TagNumber           int    `json:"TagNumber"`
}*/

type Record struct {
	TransactionId       string `json:"TransactionId"`
	iAccountID          string `json:"AccountiD"`
	TripDateTime        string `json:"TripDateTime"`
	TransactionPostDate string `json:"TransactionPostDate"`
	PlateSource         string `json:"PlateSource"`
	PlateCategory       string `json:PlateCategory`
	PlateColor          string `json:"PlateColor"`
	PlateNumber         string `json:"PlateNumber"`
	TagNumber           int    `json:"TagNumber"`
	TollGateLocation    string `json:"TollGateLocation"`
	TollGateDirection   string `json:"TollGateDirection"`
	Amount              string `json:"Amount"`
}

func JSONStringToStructure(jsonString string, structure interface{}) error {
	jsonBytes := []byte(jsonString)
	return json.Unmarshal(jsonBytes, structure)
}

func AddRowWithData(sheet *xlsx.Sheet, records Records, cells ...string) {
	r := sheet.AddRow()
	for _, v := range cells {
		c := &xlsx.Cell{}
		c = r.AddCell()
		c.Value = v
	}

	for _, v := range records.Result {
		r := sheet.AddRow()
		cell1 := &xlsx.Cell{}
		cell1 = r.AddCell()
		cell1.Value = v.TransactionId

		cell2 := &xlsx.Cell{}
		cell2 = r.AddCell()
		cell2.Value = v.iAccountID

		cell3 := &xlsx.Cell{}
		cell3 = r.AddCell()
		cell3.Value = v.TripDateTime

		cell4 := &xlsx.Cell{}
		cell4 = r.AddCell()
		cell4.Value = v.TransactionPostDate

		cell5 := &xlsx.Cell{}
		cell5 = r.AddCell()
		cell5.Value = v.PlateSource

		cell6 := &xlsx.Cell{}
		cell6 = r.AddCell()
		cell6.Value = "Private"

		cell7 := &xlsx.Cell{}
		cell7 = r.AddCell()
		cell7.Value = v.PlateColor

		cell8 := &xlsx.Cell{}
		cell8 = r.AddCell()
		cell8.Value = v.PlateNumber

		cell9 := &xlsx.Cell{}
		cell9 = r.AddCell()
		cell9.Value = strconv.Itoa(v.TagNumber)

		cell10 := &xlsx.Cell{}
		cell10 = r.AddCell()
		cell10.Value = v.TollGateLocation

		cell11 := &xlsx.Cell{}
		cell11 = r.AddCell()
		cell11.Value = v.TollGateDirection

		cell12 := &xlsx.Cell{}
		cell12 = r.AddCell()
		cell12.Value = v.Amount

	}
}

//return r
