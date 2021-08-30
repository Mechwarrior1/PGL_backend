package model

import (
	"database/sql"

	"github.com/Mechwarrior1/PGL_backend/word2vec"
)

type (
	UserSecret struct {
		ID          int    `json:"ID"`
		Username    string `json:"Username"`
		Password    string `json:"Password"`
		IsAdmin     string `json:"IsAdmin"`
		CommentItem string `json:"CommentItem"`
	}

	UserInfo struct {
		ID          int    `json:"ID"`
		Username    string `json:"Username"`
		LastLogin   string `json:"LastLogin"`
		DateJoin    string `json:"DateJoin"`
		CommentItem string `json:"CommentItem"`
	}

	ItemListing struct {
		ID              int     `json:"ID"`
		Username        string  `json:"Username"`
		Name            string  `json:"Name"`
		ImageLink       string  `json:"ImageLink"`
		DatePosted      string  `json:"DatePosted"`
		CommentItem     string  `json:"CommentItem"`
		ConditionItem   string  `json:"ConditionItem"`
		Cat             string  `json:"Cat"`
		ContactMeetInfo string  `json:"ContactMeetInfo"`
		Similarity      float32 `json:"Similarity"`
		Completion      string  `json:"Completion"`
	}

	CommentUser struct {
		ID          int    `json:"ID"`
		Username    string `json:"Username"`
		ForUsername string `json:"ForUsername"`
		Date        string `json:"Date"`
		CommentItem string `json:"CommentItem"`
	}

	CommentItem struct {
		ID          int    `json:"ID"`
		Username    string `json:"Username"`
		ForItem     string `json:"ForItem"`
		Date        string `json:"Date"`
		CommentItem string `json:"CommentItem"`
	}

	DataPacket struct {
		// key to access rest api
		Key         string        `json:"Key"`
		ErrorMsg    string        `json:"ErrorMsg"`
		InfoType    string        `json:"InfoType"`
		ResBool     string        `json:"ResBool"`
		RequestUser string        `json:"RequestUser"`
		DataInfo    []interface{} `json:"DataInfo"`
	}

	DataPacketSimple struct {
		ErrorMsg string `json:"ErrorMsg"`
		ResBool  string `json:"ResBool"`
	}

	// interface for the mysql/postgres dbhandler
	dbController interface {
		GetRecord(dbTable string) ([]interface{}, error)
		GetRecordlistingIndex(requestWords string, filterUsername string, filterDate string, filterCat string, embed *word2vec.Embeddings) ([]interface{}, error)
		GetSingleRecord(dbTable string, queryString string, queryString2 string) ([]interface{}, error)
		InsertRecord(dbTable string, receiveInfo map[string]string, maxIDInt1 int) error
		GetUsername(dbTable string, id int) (string, error)
		GetMaxID(dbTable string) (int, error)
		EditRecord(dbTable string, receiveInfo map[string]string) error
		DeleteRecord(dbTable string, id int) error
		ReturnDB() *sql.DB
		AddUser(receiveInfo map[string]string) error
		CompleteItem(itemID string) error
	}

	DBHandler struct {
		DBController    dbController
		ApiKey          string
		ReadyForTraffic bool
	}
)
