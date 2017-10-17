// PT-Server project main.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"mgoSession"
	"models"
	"net/http"
	"os"

	"strconv"
	"strings"
	"time"

	"imgops"
	"ioops"
	"logutil"
	"utils"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type Processer interface {
	Process()
}

type ImageProcesser interface {
	Read([]byte) io.Reader
	CopyToFile(io.Reader, string)
}

var (
	DBConnection, Addr, Root, Quarintine, Invalid, FullImage, MediumImage, Temp, CounterFile string
)

var (
	MaxHeight, MaxWidth, ResieWidth, ResizeHeight uint
	NoOfImages                                    int
)

var (
	err       error
	chPic     chan string
	logger    *logutil.LoggerUtil
	dbSession *mgoSession.Session
)

//init all propertie values from config/app.toml file
func init() {

	viper.SetConfigName("app")    // no need to include file extension
	viper.AddConfigPath("config") // set the path of your config file

	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println("Config file not found...")
	} else {
		Addr = viper.GetString("connections.Address")
		DBConnection = viper.GetString("connections.DBConnection")
		Root = viper.GetString("directories.Root")
		Quarintine = viper.GetString("directories.Quarantune")
		Invalid = viper.GetString("directories.Invalid")
		FullImage = viper.GetString("directories.FullImage")
		MediumImage = viper.GetString("directories.MediumImage")
		CounterFile = viper.GetString("directories.CounterFile")
		Temp = viper.GetString("directories.Temp")
		MaxHeight = uint(viper.GetInt("image.MaxHeight"))
		MaxWidth = uint(viper.GetInt("image.MaxWidth"))
		ResieWidth = uint(viper.GetInt("image.ResizeWidth"))
		ResizeHeight = uint(viper.GetInt("image.ResizeHeight"))
		NoOfImages = viper.GetInt("image.NoOfImages")
	}
}

//init root directory and also initialize a channel that is used for addpic functionality
func init() {
	ioops.CreateDirectory(Root)
	chPic = make(chan string, 1)
}

//init logger
func init() {
	file, err := os.OpenFile("logs/application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file :", err)
	}
	logger = logutil.NewLogger(file, "log: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.SetOutput(file)

}

func main() {
	dbSession, err = mgoSession.New(DBConnection, "local")

	if err != nil {
		log.Fatal("mongodb database is not connected")
	} else {
		log.Println("mongodb session has been created")

		dataServer := http.NewServeMux()

		dataServer.Handle("/AddMovie", MiddlewareHandle(AddMovie()))
		dataServer.Handle("/FetchAllMovies", MiddlewareHandle(GetAll("movies")))
		dataServer.Handle("/GetMovieList", MiddlewareHandle(GetMovieList()))
		dataServer.Handle("/GetMovieByMovieID", MiddlewareHandle(GetMovieByMovieID()))
		dataServer.Handle("/GetMoviesByKeyword", MiddlewareHandle(GetMoviesByKeyword()))

		dataServer.Handle("/UploadImageForMovie", MiddlewareHandle(AddPic()))
		dataServer.Handle("/GetImagesByDelta", MiddlewareHandle(GetImagesByDelta()))
		dataServer.Handle("/GetFullImageByID", MiddlewareHandle(GetFullImageByID()))

		fs := http.FileServer(http.Dir(""))

		dataServer.Handle("/", MiddlewareHandle(fs))

		go ProcessImages()

		http.ListenAndServe(Addr, dataServer)

	}
}

//retunes a http.handler based on middleware conditions..
//if all are passed the argument next is returned. else other handlers are retunred
//for example unathenticated handler etc..
func MiddlewareHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.WriteLog("Call On URI:", r.Host, r.URL, " from the following IP address:", utils.GetIpAddr(r))
		if IsAuthenticated(r) {
			next.ServeHTTP(w, r)
		} else {
			logger.WriteLog("Call On ", r.Host, r.URL, " from the following IP address:", utils.GetIpAddr(r), " is not authenticated")
			Unauthenticated().ServeHTTP(w, r)
		}
	})
}

//when client is unathenticated this handler automatically called
//this is called from the moddlewarehandle.
func Unauthenticated() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		WriteResponseMessage(w, "Unathenticated request", "", 400, true)
	})
}

//adds a movie to the database.
//generates movieid from the GUID function
//reads request body and unmarshel it to Movie struct
//validate whether there is already a movie with that name
//validates Movie object
//creates movieid,full,medium, invalid, quarantine and temp directories.
//temp will be replaced as inmem process will happen later
//adds corrosponding movie data to the mongodb collection
func AddMovie() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header.Get("Forwarded"))

		if r.Method != "POST" {
			WriteResponseMessage(w, "this action works only for Post method", "action works only for Post method", 400, false)
			return
		}
		var buf []byte
		var M models.Movie

		if r.Header.Get("Content-Type") != "application/json" {
			WriteResponseMessage(w, "content type should be a in json", "content type should be a in json", 400, false)
			return
		}

		buf, err = ioutil.ReadAll(r.Body)
		if err != nil {
			WriteResponseMessage(w, "error in reading data from request body", err.Error(), 400, false)
			return
		}

		err = json.Unmarshal(buf, &M)
		if err != nil {
			WriteResponseMessage(w, "error in converting the data to object", err.Error(), 400, false)
			return
		}

		m := make(map[string]interface{}, 1)
		m["title"] = M.Title
		result, _ := dbSession.FindByQuery("movies", m)

		if len(result) > 0 {
			WriteResponseMessage(w, "movie details are alredy existed in the system", "movie details are alredy existed in the system", 400, false)
			return
		}

		if models.ValidateMovie(M) != "" {
			WriteResponseMessage(w, models.ValidateMovie(M), models.ValidateMovie(M), 400, false)
			return
		}

		M.MovieId = utils.GUID()
		err = dbSession.Insert("movies", M)

		if err != nil {
			WriteResponseMessage(w, "error in insering the data", err.Error(), 400, false)
			return
		}

		err = ioops.CreateDirectory(Root + M.MovieId)
		if err != nil {
			WriteResponseMessage(w, "error in creating directory", err.Error(), 400, false)
			return
		}
		err = ioops.CreateFile(Root + M.MovieId + "/" + CounterFile)
		if err != nil {
			WriteResponseMessage(w, "error in creating directory", err.Error(), 400, false)
			return
		}
		err = ioops.CreateDirectory(Root + M.MovieId + "/" + Quarintine)
		if err != nil {
			WriteResponseMessage(w, "error in creating directory", err.Error(), 400, false)
			return
		}
		err = ioops.CreateDirectory(Root + M.MovieId + "/" + Invalid)
		if err != nil {
			WriteResponseMessage(w, "error in creating directory", err.Error(), 400, false)
			return
		}
		err = ioops.CreateDirectory(Root + M.MovieId + "/" + FullImage)
		if err != nil {
			WriteResponseMessage(w, "error in creating directory", err.Error(), 400, false)
			return
		}
		err = ioops.CreateDirectory(Root + M.MovieId + "/" + MediumImage)
		if err != nil {
			WriteResponseMessage(w, "error in creating directory", err.Error(), 400, false)
			return
		}

		WriteResponseMessage(w, "Movie successfully added", M.MovieId, 200, true)
	})
}

//gives movie list by batchno as query string parameter
//for example batchno=1
func GetMovieList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			WriteResponseMessage(w, "this action works only with get method", "action works only for Post method", 400, false)
			return
		}
		vars := r.URL.Query()
		batchno := vars["batchno"]
		if len(batchno) < 1 {
			WriteResponseMessage(w, "query parameter batchno is not supplied", "query parameter batchno is not supplied", 400, false)
			return
		}
		m := make(map[string]interface{}, 1)

		m["batchno"] = batchno[0]
		result, err := dbSession.FindByQueryAll("movies", m)
		if err != nil {
			WriteResponseMessage(w, "error in db connection", err.Error(), 400, false)
			return
		}

		if len(result) < 1 {
			WriteResponseMessage(w, "no data available", "no data available", 400, false)
			return
		}
		jData, err := json.Marshal(result)
		if err != nil {
			WriteResponseMessage(w, "error in reading the result", err.Error(), 400, false)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)

	})
}

//returns movie by its movieid
func GetMovieByMovieID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			WriteResponseMessage(w, "this action works only with get method", "action works only for Post method", 400, false)
			return
		}
		vars := r.URL.Query()
		movieid := vars["movieid"]
		if len(movieid) < 1 {
			WriteResponseMessage(w, "movieid is not provided as querystring parameter", "movieid is not provided as querystring parameter", 400, false)
			return
		}
		m := make(map[string]interface{}, 1)
		m["movieid"] = movieid[0]
		result, err := dbSession.FindByQuery("movies", m)

		if len(result) < 1 {
			WriteResponseMessage(w, "no data available", "no data available", 400, false)
			return
		}
		if err != nil {
			WriteResponseMessage(w, "error in db connection", err.Error(), 400, false)
			return
		}
		jData, err := json.Marshal(result)
		if err != nil {
			WriteResponseMessage(w, "error in reading the result", err.Error(), 400, false)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
	})
}

//retuns movie key and value.
//key can be any filed in mongo db and value is corrosponding value of the filed.
//value need not to fully given as it fetches as like operator(bson.RegEx in mongodb)
func GetMoviesByKeyword() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			WriteResponseMessage(w, "this action works only with get method", "action works only for Post method", 400, false)
			return
		}
		vars := r.URL.Query()
		key := vars["key"]
		value := vars["value"]
		if len(key) < 1 || len(value) < 1 {
			WriteResponseMessage(w, "key and(or) value query string parameters are not passed", "key and(or) value query string parameters are not passed", 400, false)
			return
		}
		result, err := dbSession.FindByRegEx("movies", key[0], value[0])
		if err != nil {
			WriteResponseMessage(w, "error in db connection", err.Error(), 400, false)
			return
		}
		if len(result) < 1 {
			WriteResponseMessage(w, "no data available", "no data available", 400, false)
			return
		}
		jData, err := json.Marshal(result)
		if err != nil {
			WriteResponseMessage(w, "error in reading the result", err.Error(), 400, false)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
	})
}

//returns all records from the collection provided by the collection name as the parameter
func GetAll(collection string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			WriteResponseMessage(w, "this action works only with get method", "action works only for Post method", 400, false)
			return
		}
		result, err := dbSession.FindAll(collection)

		if err != nil {
			WriteResponseMessage(w, "error in db connection", err.Error(), 400, false)
			return
		}

		if len(result) < 1 {
			WriteResponseMessage(w, "no data available", "no data available", 400, false)
			return
		}
		jData, err := json.Marshal(result)
		if err != nil {
			WriteResponseMessage(w, "error in reading the result", err.Error(), 400, false)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
	})
}

//fetches document from the collection by _id (bson id).
//collection name as the parameter
func GetByID(collection string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			vars := mux.Vars(r)
			_id := vars["id"]

			if _id == "" {
				WriteResponseMessage(w, "id has not passed as a parameter", "id has not passed as a parameter", 400, false)
				return
			}

			result, err := dbSession.ListByID(collection, _id)
			if err != nil {
				w.WriteHeader(400)
				fmt.Fprintln(w, err.Error())
			} else {

				jData, err := json.Marshal(result)
				if err != nil {
					w.WriteHeader(400)
					fmt.Fprintln(w, err.Error())
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(jData)
			}
		}
	})
}

func GetPicByMovieID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			WriteResponseMessage(w, "this action works only with get method", "action works only for Post method", 400, false)
			return
		}
		vars := r.URL.Query()
		movieid := vars["movieid"]
		m := make(map[string]interface{}, 1)
		m["movieid"] = movieid[0]
		result, err := dbSession.FindByQueryAll("pics", m)
		if err != nil {
			WriteResponseMessage(w, "error in db connection", err.Error(), 400, false)
			return
		}

		if len(result) < 1 {
			WriteResponseMessage(w, "no data available", "no data available", 400, false)
			return
		}
		jData, err := json.Marshal(result)
		if err != nil {
			WriteResponseMessage(w, "error in reading the result", err.Error(), 400, false)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
	})
}

func GetPicDetailsByPicID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			WriteResponseMessage(w, "this action works only with get method", "action works only for Post method", 400, false)
			return
		}
		vars := r.URL.Query()
		picid := vars["picid"]
		m := make(map[string]interface{}, 1)
		m["picid"] = picid[0]
		result, err := dbSession.FindByQuery("pics", m)
		if err != nil {
			WriteResponseMessage(w, "error in db connection", err.Error(), 400, false)
			return
		}

		if len(result) < 1 {
			WriteResponseMessage(w, "no data available", "no data available", 400, false)
			return
		}
		jData, err := json.Marshal(result)
		if err != nil {
			WriteResponseMessage(w, "error in reading the result", err.Error(), 400, false)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)

	})
}

func GetImagesByDelta() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			WriteResponseMessage(w, "this action works only for GET method", "action works only for Post method", 400, false)
			return
		}
		vars := r.URL.Query()
		min := vars["min"]
		max := vars["max"]

		movieid := vars["movieid"]
		if len(movieid) < 1 {
			WriteResponseMessage(w, "movieid is not provided as querystring parameter", "movieid is not provided as querystring parameter", 400, false)
			return
		}
		if len(min) < 1 {
			WriteResponseMessage(w, "min image number is not provided as querystring parameter", "movieid is not provided as querystring parameter", 400, false)
			return
		}
		if len(max) < 1 {
			WriteResponseMessage(w, "max image number is not provided as querystring parameter", "movieid is not provided as querystring parameter", 400, false)
			return
		}

		m := make(map[string]interface{}, 1)
		m["movieid"] = movieid[0]
		result, err := dbSession.FindByQuery("movies", m)

		if len(result) < 1 {
			WriteResponseMessage(w, "no data available", "no data available", 400, false)
			return
		}
		if err != nil {
			WriteResponseMessage(w, "error in db connection", err.Error(), 400, false)
			return
		}

		currentCount, err := ioops.GetFileCount(Root + movieid[0] + "/" + CounterFile)
		if err != nil {
			WriteResponseMessage(w, "wrong count of images", err.Error(), 400, false)
			return
		}

		mx, _ := strconv.Atoi(max[0])
		mn, _ := strconv.Atoi(min[0])

		magicNos, err := utils.GetMagicNumbers(mx, mn, currentCount, NoOfImages)

		if err != nil {
			WriteResponseMessage(w, err.Error(), err.Error(), 400, false)
			return
		}

		m1 := make(map[string]interface{}, 2)
		m1["movieid"] = movieid[0]
		m2 := make(map[string]interface{}, 1)
		m2["$in"] = magicNos
		m1["picid"] = m2

		data, err := dbSession.FindByQueryAll("pics", m1)
		if err != nil {
			WriteResponseMessage(w, "error in db connection", err.Error(), 400, false)
			return
		}
		if len(data) < 1 {
			WriteResponseMessage(w, "no data available", "no data available", 400, false)
			return
		}
		jData, err := json.Marshal(data)
		if err != nil {
			WriteResponseMessage(w, "error in reading the result", err.Error(), 400, false)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)

		//WriteResponseMessage(w, message, "", 200, true)

	})
}

func GetFullImageByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			WriteResponseMessage(w, "this action works only for GET method", "action works only for Post method", 400, false)
			return
		}
		vars := r.URL.Query()

		movieid := vars["movieid"]
		picid := vars["picid"]
		imagesize := vars["imagesize"]
		//fmt.Println(imagesize, len(imagesize))
		fmt.Println(imagesize)
		if len(movieid) < 1 {
			WriteResponseMessage(w, "movieid is not provided as querystring parameter", "movieid is not provided as querystring parameter", 400, false)
			return
		}
		if len(picid) < 1 {
			WriteResponseMessage(w, "picid is not provided as querystring parameter", "movieid is not provided as querystring parameter", 400, false)
			return
		}
		if len(imagesize) < 1 {
			WriteResponseMessage(w, "imagesize is not provided as querystring parameter", "movieid is not provided as querystring parameter", 400, false)
			return
		}
		if strings.ToLower(imagesize[0]) != "full" && strings.ToLower(imagesize[0]) != "medium" {
			WriteResponseMessage(w, "imagesize either full or medium", "movieid is not provided as querystring parameter", 400, false)
			return
		}
		m := make(map[string]interface{}, 1)
		m["movieid"] = movieid[0]
		m["picid"] = picid[0]

		result, err := dbSession.FindByQuery("pics", m)

		if len(result) < 1 {
			WriteResponseMessage(w, "no data available", "no data available", 400, false)
			return
		}
		if err != nil {
			WriteResponseMessage(w, "error in db connection", err.Error(), 400, false)
			return
		}

		resp, err := http.Get("http://localhost" + Addr + Root + movieid[0] + "/" + imagesize[0] + "/" + picid[0])
		if err != nil {
			WriteResponseMessage(w, "error in db connection", err.Error(), 400, false)
			return
		}
		defer resp.Body.Close()
		content, err := ioutil.ReadAll(resp.Body)
		w.Header().Add("Content-Type", "image/jpeg")
		w.Write(content)
		if err != nil {
			WriteResponseMessage(w, "error in db connection", err.Error(), 400, false)
			return
		}

		//WriteResponseMessage(w, message, "", 200, true)

	})
}

func AddPic() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var picid string
		var Buffer []byte
		if r.Method != "POST" {
			WriteResponseMessage(w, "this action works only for Post method", "action works only for Post method", 400, false)
			return
		}
		err = r.ParseMultipartForm(32 << 20)

		if err != nil {
			WriteResponseMessage(w, "error in ParseMultipartForm", err.Error(), 400, false)
			return
		}
		movieid := r.FormValue("movieid")
		title := r.FormValue("title")
		relevance := r.FormValue("relevance")
		option := r.FormValue("option")

		if movieid == "" {
			WriteResponseMessage(w, "movieid cannot be empty", "movieid cannot be empty", 400, false)
			return
		}

		if title == "" {
			WriteResponseMessage(w, "title formdata is empty", "title formdata is empty", 400, false)
			return
		}
		pic := models.Pic{}
		pic.MovieId = movieid
		pic.Title = title
		pic.Option = option
		pic.Relevance = relevance
		pic.Status = "Active"
		pic.Timestamp = time.Now().String()

		if models.ValidatePic(pic) != "" {
			WriteResponseMessage(w, models.ValidatePic(pic), models.ValidatePic(pic), 400, false)
			return
		}

		file, handller, err := r.FormFile("uploadfile")

		if err != nil {
			WriteResponseMessage(w, "error in reading file", err.Error(), 400, false)
			return
		}
		defer file.Close() // when error no need to even deffered close since response is written at the above step
		Buffer, err = ioutil.ReadAll(file)

		imp := imgops.New(Buffer)

		if err != nil {
			WriteResponseMessage(w, "error in reading file", err.Error(), 400, false)
			return
		}

		fileExt := ioops.GetFileExt(handller.Filename)

		//picid = picid + fileExt //tempararly extension has been removed.
		pic.FileType = fileExt

		if err != nil {
			WriteResponseMessage(w, "upload failure", err.Error(), 400, false)
			return
		}

		width, height := imp.GetImageDimensionBy()
		fmt.Println(width, height)

		//isnude, err := imageops.IsNude(ioops.Read(Buffer))

		isnude, err := imp.IsNude()

		if err != nil {

			WriteResponseMessage(w, "error in nude functionality", err.Error(), 400, true)
		}

		if uint(width) < MaxWidth || uint(height) < MaxHeight {
			picid = utils.GUID()
			pic.PicId = picid
			ioops.CopyToFile(ioops.Read(Buffer), Root+movieid+"/"+Invalid+picid)
			pic.PicPath = Root + movieid + "/" + Invalid + picid
			pic.PicStatus = "short"
			err = dbSession.Insert("pics", pic)
			if err != nil {
				WriteResponseMessage(w, "file is not copied to the database", err.Error(), 400, false)
				return
			}
			WriteResponseMessage(w, "short", picid, 400, true)
			return
		} else if isnude {
			picid = utils.GUID()
			pic.PicId = picid
			ioops.CopyToFile(ioops.Read(Buffer), Root+movieid+"/"+Quarintine+picid)
			pic.PicPath = Root + movieid + "/" + Quarintine + picid
			pic.PicStatus = "nude"
			err = dbSession.Insert("pics", pic)
			if err != nil {
				WriteResponseMessage(w, "file is not copied to the database", err.Error(), 400, false)
				return
			}
			WriteResponseMessage(w, "flagged", picid, 400, true)
			return
		} else {
			count, err := ioops.FileCountIncrement(Root + movieid + "/" + CounterFile)
			picid = strconv.Itoa(count)
			pic.PicId = picid
			if err != nil {
				WriteResponseMessage(w, "error in reading file", err.Error(), 400, false)
				return
			}
			ioops.CopyToFile(ioops.Read(Buffer), Root+movieid+"/"+FullImage+picid)
			pic.PicPath = Root + movieid + "/" + FullImage + picid
			pic.PicStatus = "valid"
			err = dbSession.Insert("pics", pic)
			if err != nil {
				WriteResponseMessage(w, "file is not copied to the database", err.Error(), 400, false)
				return
			}
			chPic <- Root + movieid + "/" + FullImage + picid
		}

		WriteResponseMessage(w, "Pic successfully uploaded", picid, 200, true)
		return
	})
}

//write response to the responsewrite in json fomat.
func WriteResponseMessage(w http.ResponseWriter, message string, trace string, status int, success bool) {
	msg := models.Message{MSG: message, Success: success, Trace: trace, Status: status}
	go logger.LogString(msg)
	w.Header().Set("Content-Type", "application/json")
	w.Write(msg.Bytes())
}

func WriteLog(status, user, source, message string) {
	l := models.Log{TimeStamp: time.Now().String(), Status: status, Source: source, Message: message}
	go logger.LogString(l)
}

//Processes list of the defined methods asynchronously
//Channel filepath is channeled from AddPic method
func ProcessImages() {
	for v := range chPic {
		r, m, _, f := ioops.SplitPath(v)
		go imgops.ResizeAndMoveGo(v, r+"/"+m+"/"+MediumImage+f, ResieWidth, ResizeHeight)
	}
}

//returns true when athenticated false when not
//athentication headers are obtained from r
func IsAuthenticated(r *http.Request) bool {
	//todo write logic here
	return true
}
