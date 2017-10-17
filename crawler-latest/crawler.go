//crawler software developed for client
//Author:Jiten Palaparthi

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

var username, password, startdate, enddate, URL, CompleteData, ErrorData, AccountID, filename string

func main() {

	username, password, startdate, enddate, filename, err := GetArgs(os.Args)
	CheckErr(err)

	URL = GetFormatedURL("1", startdate, enddate)

	log.Println("Initiated. Connecting to the server...")
	log.Println("Fetching details for the user ", username, " and date between "+startdate+" and "+enddate)

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	CheckErr(err)
	client := http.Client{Jar: jar}
	resp, err := client.Get("https://customers.salik.ae/connect/login?lang=en")
	CheckErr(err)
	log.Println("Intial URL has been hit.. trying to fetch login callback id")

	urls, err := url.Parse(resp.Request.Referer())
	CheckErr(err)
	log.Println("Redirected url is being fetch.. trying to get all data and cookies and redirected cookies")

	postData := url.Values{}
	postData.Add("UserName", username)
	postData.Add("Password", password)

	req, err := http.NewRequest("POST", resp.Request.URL.String(), strings.NewReader(postData.Encode()))
	CheckErr(err)

	req.AddCookie(jar.Cookies(urls)[0])
	req.AddCookie(jar.Cookies(urls)[1])

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept-Language", "en-GB,en;q=0.8,en-US;q=0.6,te;q=0.4")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	CheckErr(err)
	log.Println("Post username and password along with open id authorization and athentication details")
	resp1, err := client.Do(req)
	CheckErr(err)
	resp2, err := client.Get(resp1.Request.URL.String())
	CheckErr(err)
	root, err := html.Parse(resp2.Body)
	CheckErr(err)
	names := []string{"id_token", "access_token", "token_type", "expires_in", "scope", "state", "session_state"}
	log.Println(names)

	postData1 := getFormValues(names, root)

	log.Println(postData1)

	resp3, err := client.Post("https://customers.salik.ae/Connect/SalikIdCallback", "application/x-www-form-urlencoded", strings.NewReader(postData1.Encode()))
	CheckErr(err)
	log.Println("OpenId authorization token has been posted that fetches from callback url")

	//if startpage == "" && endpage == "" {
	var data []byte
	data, err = GetResponseData(URL, jar.Cookies(resp3.Request.URL)[0], jar.Cookies(resp3.Request.URL)[1])
	//log.Println(string(data))
	i, l := GetScopeOfData(string(data))

	//log.Fatal(i, l)
	if l == -1 {
		log.Fatalln("Server is not responding hence ... exit.")
	}
	//log.Println("Length of data ", len(string(data)))

	CompleteData = CompleteData + string(data)[i:l]

	log.Println("Fetching page-1 records and writng to json file")

	count := GetTotalCount(string(data))

	log.Println("Caluclating number of records")
	pagenumber := int(count / 10)
	if int(count%10) != 0 {
		pagenumber = pagenumber + 1
	}
	log.Println("Caluclating total record count and number of pages to fetch")
	log.Println("Total records found:", count)
	log.Println("Total pages to fetch:", pagenumber)

	for k := 2; k <= pagenumber; k++ {
		var data []byte
		URL = GetFormatedURL(strconv.Itoa(k), startdate, enddate)
		data, err := GetResponseData(URL, jar.Cookies(resp3.Request.URL)[0], jar.Cookies(resp3.Request.URL)[1])
		CheckErr(err)
		if err == nil {
			i, l := GetScopeOfData(string(data))
			if l == -1 {
				ErrorData = ErrorData + "\nPage Number:" + strconv.Itoa(k) + "\n" + string(data)
			} else {
				CompleteData = CompleteData + "," + string(data)[i:l]
				log.Println("Fetching page-" + strconv.Itoa(k) + " records and writng as json file")
			}
		} else {
			k--
		}
	}
	CheckErr(err)
	CompleteData = `{"result":[` + CompleteData + "]}"
	log.Println("Writing to output file")
	err = ioutil.WriteFile(filename+".json", []byte(CompleteData), 0644)
	CheckErr(err)
	if ErrorData != "" {
		log.Println("Writing to error log..")
		err = ioutil.WriteFile("log.txt", []byte(ErrorData), 0644)
		CheckErr(err)
	}
	log.Println("Json file Finished successfully...")
	//-------------------------------------------------------------
	log.Println("Creating Excel file")

	sheetData, err := ioutil.ReadFile(filename + ".json")

	if err != nil {
		log.Println(err)
	}
	s := Records{}
	JSONStringToStructure(string(sheetData), &s)
	if err != nil {
		log.Println(err)
	}

	var file *xlsx.File
	var sheet *xlsx.Sheet

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		log.Printf(err.Error())
	}

	AddRowWithData(sheet, s, "TransactionId", "iAccountID", "TransactionDate", "PostDate", "PlateSource", "PlateCategory", "PlateCode", "PlateNumber", "TagNumber", "Location", "Direction", "Amount")

	err = file.Save(filename + ".xlsx")
	if err != nil {
		log.Printf(err.Error())
	}
	log.Println("Records are stored in excel file successfully")
}

type Records struct {
	Result []Record `json:"result"`
}

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

// local functions

func CheckErr(err error) {
	if err != nil {
		err = ioutil.WriteFile("error.txt", []byte(err.Error()), 0644)
		log.Println(err)
	}
}

func getElementByName(name string, n *html.Node) (element *html.Node, ok bool) {
	for _, a := range n.Attr {
		if a.Key == "name" && a.Val == name {
			return n, true
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if element, ok = getElementByName(name, c); ok {
			return
		}
	}
	return
}

func getElmentValue(name string, n *html.Node) string {
	element, ok := getElementByName(name, n)
	if !ok {
		log.Fatal("element not found")
	}
	for _, a := range element.Attr {
		if a.Key == "value" {
			//fmt.Println(a.Val)
			return a.Val
		}
	}
	return ""
}

func getFormValues(names []string, n *html.Node) (postValues url.Values) {
	postValues = url.Values{}
	for _, v := range names {
		value := getElmentValue(v, n)
		if value != "" {
			postValues.Add(v, value)
		}

	}
	return postValues
}

func GetTotalCount(data string) (count int) {

	i := strings.LastIndex(string(data), "TotalCount")
	str := string(data)[i+12:]
	li := 0
	for i, v := range str {
		if string(v) == "}" {
			li = i
			break
		}
	}
	count, err := strconv.Atoi(str[0:li])

	if err != nil {
		return 0
	} else {
		return count
	}
}

func GetScopeOfData(data string) (fi, li int) {
	fi = strings.Index(data, "[")
	li = strings.Index(data, "]")
	return fi + 1, li
}

func GetResponseData(url string, c1, c2 *http.Cookie) (data []byte, err error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(c1)
	req.AddCookie(c2)
	req.Header.Add("Content-Type", "application/json")
	CheckErr(err)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	data, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetArgs(args []string) (username, password, startdate, enddate, filename string, err error) {

	for i, v := range os.Args {
		if v == "-u" {
			username = os.Args[i+1]
		}
		if v == "-p" {
			password = os.Args[i+1]
		}
		if v == "-sd" {
			startdate = os.Args[i+1]
		}
		if v == "-ed" {
			enddate = os.Args[i+1]
		}
		if v == "-fn" {
			filename = os.Args[i+1]
		}
		if v == "-a" {
			AccountID = os.Args[i+1]
		}
	}

	if username == "" || password == "" || startdate == "" || enddate == "" {
		return "", "", "", "", "", errors.New("wrong arguments are passed.")
	}
	_, err = ValidateDateStr(startdate)
	if err != nil {
		return "", "", "", "", "", err
	}

	_, err = ValidateDateStr(enddate)
	if err != nil {
		return "", "", "", "", "", err
	}

	if filename == "" {
		filename = "output.json"
	}

	return username, password, startdate, enddate, filename, nil

}

func ValidateDateStr(date string) (r bool, err error) {

	ss := strings.Split(date, "-")
	if len(ss) != 3 {
		return false, errors.New("Wrong date format.date should be in dd-mm-yyyy")
	}

	dd, err := strconv.Atoi(ss[0])
	if err != nil {
		return false, err
	}
	if dd > 31 {
		return false, errors.New("date cannot be more than 31")
	}
	mm, err := strconv.Atoi(ss[1])
	if err != nil {
		return false, err
	}
	if mm > 12 {
		return false, errors.New("month cannot be more than 12")
	}
	yyyy, err := strconv.Atoi(ss[2])
	if err != nil {
		return false, err
	}
	if yyyy < 2016 || yyyy > 2018 {
		return false, errors.New("year cannot be more than 2018 and less than 2016")
	}

	return true, nil
}

func GetFormatedURL(pid, sd, ed string) string {
	return "https://customers.salik.ae/surface/portal/trips?pageId=" + pid + "&pageSize=10&timePeriod=4&tripType=1&tagPlateDetails=1&tagNumber=&plateNumber=&plateCountryCode=&plateSource=&plateCategory=&plateCode=&startDate=" + sd + "&endDate=" + ed + "&lang=en"
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
		cell2.Value = AccountID

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
