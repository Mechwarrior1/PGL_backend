package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Mechwarrior1/PGL_backend/model"
	"github.com/Mechwarrior1/PGL_backend/word2vec"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

// function for the rest api, respond with the slice of all courses
func GetAllListingIndex(c echo.Context, dbHandler *model.DBHandler, embed *word2vec.Embeddings) error {
	// can only return listing results, commentUser and commentItem
	itemName := c.QueryParam("name")
	filterUsername := c.QueryParam("filter")
	filterDate := c.QueryParam("date")
	filterCat := c.QueryParam("cat")
	sortIndex, err := dbHandler.DBController.GetRecordlistingIndex(itemName, filterUsername, filterDate, filterCat, embed)

	if err != nil {
		return newResponseSimple(c, "Bad Request", "false", http.StatusBadRequest)
	}

	return newResponse(c, sortIndex, "nil", "ItemListing", "true", "", http.StatusOK)
}

// function for the rest api, respond with the slice of all courses
func GetAllListing(c echo.Context, dbHandler *model.DBHandler) error {
	// can only return listing results, commentUser and commentItem
	dataPacket1, err1 := readJSONBody(c, dbHandler) // read response JSON
	if err1 != nil {
		return newResponseSimple(c, "Bad Request", "false", http.StatusBadRequest)
	}

	searchString := ""
	indexArr := dataPacket1["DataInfo"].([]interface{})[0].([]interface{})

	// check payload
	if len(indexArr) == 0 {
		return newResponseSimple(c, "Bad Request", "false", http.StatusBadRequest)
	}

	for _, ind := range indexArr {
		searchString = searchString + strconv.Itoa(int(ind.(float64))) + ", "
	}
	searchString = searchString[:len(searchString)-2] //remove last ,

	if err1 != nil {
		return newResponseSimple(c, "Bad Request", "false", http.StatusBadRequest)
	}

	//ReturnDB() returns the database pointer, interface does not contain the db pointer
	results, err := dbHandler.DBController.ReturnDB().Query("Select * FROM ItemListing WHERE ID in (" + searchString + ")")

	if err != nil {
		return newResponseSimple(c, "Bad Request", "false", http.StatusBadRequest)
	}

	sendInfo := []interface{}{}
	for results.Next() {
		newEntry := model.ItemListing{}
		err = results.Scan(&newEntry.ID, &newEntry.Username, &newEntry.Name, &newEntry.ImageLink, &newEntry.DatePosted, &newEntry.CommentItem, &newEntry.ConditionItem, &newEntry.Cat, &newEntry.ContactMeetInfo, &newEntry.Completion)

		//convert string of numbers into int and actual date string
		dateVal, _ := strconv.Atoi(newEntry.DatePosted)
		newEntry.DatePosted = time.Unix(int64(dateVal), 0).Format("02-01-2006")

		if err != nil {
			fmt.Println("logger: error at getRecordlisting:" + err.Error())
		}
		sendInfo = append(sendInfo, newEntry)
	}

	return newResponse(c, sendInfo, "nil", "ItemListing", "true", "", http.StatusOK)
}

func GetAllComment(c echo.Context, dbHandler1 *model.DBHandler, embed *word2vec.Embeddings) error {
	// gets all comments regarding a particular item id
	itemID := c.Param("id")
	sendInfo, _ := dbHandler1.DBController.GetRecord("CommentItem")
	newSendInfo := []interface{}{}

	// loop through the returned results and check against the required itemID
	// append to new array if item id matches
	for i := range sendInfo {
		temp1 := sendInfo[i].(model.CommentItem)

		if temp1.ForItem == itemID {
			// fmt.Println(sendInfo[i])
			newSendInfo = append(newSendInfo, sendInfo[i])
		}
	}

	return newResponse(c, newSendInfo, "nil", "ItemListing", "true", "", http.StatusOK)
}

func PwCheck(c echo.Context, dbHandler1 *model.DBHandler) error { //works
	dataPacket1, err := readJSONBody(c, dbHandler1) // read response JSON

	if err != nil {
		return newResponseSimple(c, "Bad Request", "false", http.StatusBadRequest)
	}

	receiveInfo := mapInterfaceToString(dataPacket1) // convert received data into map[string]string
	dbData, err1 := dbHandler1.DBController.GetSingleRecord("UserSecret", "WHERE Username", receiveInfo["Username"])
	dbData2, err2 := dbHandler1.DBController.GetSingleRecord("UserInfo", "WHERE Username", receiveInfo["Username"])

	if err1 != nil || err2 != nil || len(dbData) == 0 {
		fmt.Println("error when attempting to retrieve info for logging in, error:", err1.Error(), err2.Error())
		return newResponseSimple(c, "Bad Request", "false", http.StatusBadRequest)
	}

	// update last login
	dbData1 := dbData[0].(model.UserSecret)
	dbData3 := dbData2[0].(model.UserInfo)
	err3 := bcrypt.CompareHashAndPassword([]byte(dbData1.Password), []byte(receiveInfo["Password"]))
	if err3 == nil {
		//update lastLogin if there is no issues

		receiveInfo["CommentItem"] = dbData3.CommentItem
		err := dbHandler1.DBController.EditRecord("UserInfo", receiveInfo)
		if err != nil {
			fmt.Println("error when editing userinfo record for: ", dbData3.Username, "\nerror: ", err)
		}

		payload := make(map[string]string)
		payload["LastLogin"] = dbData3.LastLogin
		if payload["LastLogin"] == "" {
			fmt.Println("empty LastLogin, error: ", err2.Error(), ", data: ", dbData3)
		}
		payload["IsAdmin"] = dbData1.IsAdmin
		// if dbData1.IsAdmin == "" {
		// 	fmt.Println("there seems to be an error with the sql request for user admin and lastlogin info")
		// fmt.Println(dbData1, dbData3)
		// }

		//return response with true if no issues
		return newResponse(c, []interface{}{payload}, "nil", "ItemListing", "true", "", http.StatusOK)
	}

	return newResponseSimple(c, "Bad Request", "false", http.StatusBadRequest)
}

// change map[string]interface to map[string]string
func mapInterfaceToString(dataPacket1 map[string]interface{}) map[string]string {

	receiveInfoRaw := dataPacket1["DataInfo"].([]interface{})[0].(map[string]interface{}) // convert received data into map[string]string
	receiveInfo := make(map[string]string)

	for k, v := range receiveInfoRaw {
		receiveInfo[k] = fmt.Sprintf("%v", v)
	}

	return receiveInfo
}

// checks if the username in the sent info is currently in DB
// returns false if username is not taken
func CheckUsername(c echo.Context, dbHandler1 *model.DBHandler) error {
	username := c.Param("username")

	//check received input
	if username == "" {
		return newResponseSimple(c, "Bad Request", "false", http.StatusBadRequest)
	}

	allData, err := dbHandler1.DBController.GetSingleRecord("UserInfo", " WHERE Username ", username)

	if err != nil || len(allData) == 0 { //err means username not found, ok to proceed
		return newResponseSimple(c, "username not found", "false", http.StatusOK)
	}

	return newResponseSimple(c, "username found", "true", http.StatusOK)
}

// the function that writes the response back
// for when you need to return arrays of entries
func newResponse(c echo.Context, dataInfo []interface{}, errorMsg string, infoType string, resBool string, requestUser string, httpStatus int) error {
	var responseJson model.DataPacket
	responseJson.DataInfo = dataInfo
	responseJson.ErrorMsg = errorMsg       // error msg if any
	responseJson.InfoType = infoType       // to access which db
	responseJson.ResBool = resBool         //
	responseJson.RequestUser = requestUser //request coming from which user

	// fmt.Println(c.Path(), ": \n", responseJson)
	return c.JSON(httpStatus, responseJson) // encode to json and send
}

// the function that writes the response back
func newResponseSimple(c echo.Context, msg string, resBool string, httpStatus int) error {
	responseJson := model.DataPacketSimple{
		msg,
		resBool,
	}
	return c.JSON(httpStatus, responseJson) // encode to json and send
}

// function to read the JSON on a request
// maps body into a map[string]interface and checks api key
func readJSONBody(c echo.Context, dbHandler1 *model.DBHandler) (map[string]interface{}, error) {
	// decode JSON body into map
	json_map := make(map[string]interface{})
	err1 := json.NewDecoder(c.Request().Body).Decode(&json_map)

	if err1 == nil {

		if json_map["Key"] != dbHandler1.ApiKey {
			// if api key does not match
			newResponseSimple(c, "Forbidden", "false", http.StatusForbidden)
			return json_map, errors.New("incorrect api key supplied")
		}

		return json_map, nil
	}

	return json_map, errors.New("error while attempting to read body of request")
}

// handler for signing up new users
func Signup(c echo.Context, dbHandler *model.DBHandler) error {
	dataPacket1, err1 := readJSONBody(c, dbHandler) // read response JSON
	if err1 != nil {
		return newResponseSimple(c, "Bad Request", "false", http.StatusBadRequest)
	}

	receiveInfo := mapInterfaceToString(dataPacket1) // convert received data into map[string]string

	err2 := dbHandler.DBController.AddUser(receiveInfo)
	if err2 != nil {
		return newResponseSimple(c, "Bad Request", "false", http.StatusBadRequest)
	}

	return newResponseSimple(c, "nil", "true", http.StatusOK)

}

// general api for posting info into db
func GenInfoPost(c echo.Context, dbHandler1 *model.DBHandler) error {
	dataPacket1, err1 := readJSONBody(c, dbHandler1) // read response JSON
	if err1 != nil {
		fmt.Println(err1)
		return newResponseSimple(c, "Bad Request", "false", http.StatusBadRequest)
	}

	receiveInfoRaw := mapInterfaceToString(dataPacket1) // convert received data into map[string]string
	tarDB := dataPacket1["InfoType"].(string)

	err2 := dbHandler1.DBController.InsertRecord(tarDB, receiveInfoRaw, 0) // deletes if target is found
	if err2 == nil {
		return newResponseSimple(c, "nil", "true", http.StatusOK)
	}

	fmt.Println("logger: insert " + tarDB + " db not found") //reach here only if it is not returned by the switch

	return newResponseSimple(c, "Bad Request", "false", http.StatusBadRequest)

}

// funcs that gets an item based on id
// returns item id info as interface and error if any
func getItem(c echo.Context, tarDB string, tarItemID string, dbHandler1 *model.DBHandler) ([]interface{}, error) {
	dbInfoSlice, err3 := dbHandler1.DBController.GetSingleRecord(tarDB, " WHERE ID", tarItemID)

	if err3 != nil || len(dbInfoSlice) == 0 {

		if tarDB == "UserInfo" {
			dbInfoSlice, err3 = dbHandler1.DBController.GetSingleRecord(tarDB, " WHERE Username", tarItemID)
		}

		if err3 != nil || len(dbInfoSlice) == 0 {
			// fmt.Println("logger: error when looking up id " + tarItemID + " for DB " + tarDB + ", err:" + err3.Error())

			return []interface{}{}, err3
		}

	}
	return dbInfoSlice, err3
}

// handler func, takes item id param and db, and returns item
// returns a bad request if not found
func GenInfoGet(c echo.Context, dbHandler1 *model.DBHandler) error {

	itemID := c.QueryParam("id")
	itemDB := c.QueryParam("db")
	// itemIDInt,_ := strconv.Atoi(itemID)
	dbInfoSlice, err3 := getItem(c, itemDB, itemID, dbHandler1)

	// for userinfo, itemlisting, commentitem and commentuser only
	if itemDB != "UserSecret" && err3 == nil { //prevents any requeset to ask for user secrets
		return newResponse(c, dbInfoSlice, "nil", itemDB, "true", "false", http.StatusCreated)
	}

	return newResponseSimple(c, "Bad Request", "false", http.StatusBadRequest)
}

// handler func, takes an id param and change item's completion status to true
// returns bad request if item id is not found
func Completed(c echo.Context, dbHandler *model.DBHandler) error {
	itemID := c.Param("id")
	itemIDInt, err2 := strconv.Atoi(itemID)
	dataPacket1, err1 := readJSONBody(c, dbHandler)

	if err2 != nil || err1 != nil {
		return newResponseSimple(c, "Bad Request", "false", http.StatusBadRequest)
	}
	username, err3 := dbHandler.DBController.GetUsername("ItemListing", itemIDInt)
	if err3 != nil {
		return newResponseSimple(c, "Bad Request", "false", http.StatusBadRequest)
	}

	// check if the username is the same
	if username != dataPacket1["RequestUser"] {
		return newResponseSimple(c, "Not owner", "false", http.StatusBadRequest)
	}

	err := dbHandler.DBController.CompleteItem(itemID)

	if err != nil {
		return newResponseSimple(c, "Bad Request", "false", http.StatusBadRequest)
	}

	return newResponseSimple(c, "nil", "true", http.StatusOK)
}

/// delete is not implemented currently, might change to soft delete instead
// func GenInfoDelete(c echo.Context, dbHandler1 *model.DBHandler) error {

// 	itemID := c.QueryParam("id")
// 	itemDB := c.QueryParam("db")
// 	apiKey := c.QueryParam("key")
// 	if dbHandler1.Apikey != apiKey{
// 				return newErrorResponse(c, "Bad Request", http.StatusBadRequest)
// 	}

// dbInfoSlice, err3 := getItem(c, tarDB, tarItemID, dbHandler1)

// if !checkUser(tarDB, dataPacket1["RequestUser"].(string), dbInfoSlice) || err3 != nil {

// 	return newErrorResponse(c, "Bad Request", http.StatusBadRequest)
// }

// for deleting an entry

// 	if tarDB == "ItemListing" || tarDB == "CommentUser" || tarDB == "CommentItem" { //only delete records for 3 items
// 		err2 := dbHandler1.DBController.DeleteRecord(tarDB, tarItemID) // attempt to delete record
// 		if err2 != nil {

// 			return newResponse(c, []interface{}{}, "nil", "userInfo", "true", "", http.StatusOK)
// 		}

// 	}
// 	//if delete did not occur
// 	fmt.Println("logger:  " + ": " + tarDB + " db not found or not in use for Delete func, err:")
// 	return newErrorResponse(c, "Bad Request", http.StatusBadRequest)

// }

// handler func, updates existing item, based on specified item id

func GenInfoPut(c echo.Context, dbHandler1 *model.DBHandler) error {

	dataPacket1, err1 := readJSONBody(c, dbHandler1) // read response JSON
	if err1 != nil {
		return newResponseSimple(c, "Bad Request", "false", http.StatusBadRequest)
	}

	receiveInfoRaw := mapInterfaceToString(dataPacket1) // convert received data into map[string]string

	err2 := dbHandler1.DBController.EditRecord(dataPacket1["InfoType"].(string), receiveInfoRaw) // deletes if target is found

	if err2 == nil {
		return newResponseSimple(c, "nil", "true", http.StatusCreated)
	}

	return newResponseSimple(c, "Bad Request", "false", http.StatusBadRequest)

}

// for client to check if api is active
func HealthCheckLiveness(c echo.Context) error {
	return newResponseSimple(c, "nil", "true", http.StatusOK)
}

// for client to check if api is ready for traffic
func HealthCheckReadiness(c echo.Context, dbHandler *model.DBHandler) error {

	if dbHandler.ReadyForTraffic {
		return newResponseSimple(c, "nil", "true", http.StatusOK)
	}

	return newResponseSimple(c, "not ready for traffic", "false", 503)

}
