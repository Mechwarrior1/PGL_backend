package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Mechwarrior1/PGL_backend/encrypt"
	"github.com/Mechwarrior1/PGL_backend/model"
	"github.com/Mechwarrior1/PGL_backend/word2vec"

	_ "github.com/go-sql-driver/mysql" // go mod init api_server.go

	"github.com/joho/godotenv"
)

type DBHandlerMysql struct {
	DB *sql.DB
}

var (
	splitText  = regexp.MustCompile(`\s*,\s*|\s,*\s*`)
	stopWords2 = regexp.MustCompile("^(i|me|my|myself|we|our|ours|ourselves|you|your|yours|yourself|yourselves|he|him|his|himself|she|her|hers|herself|it|its|itself|they|them|their|theirs|themselves|what|which|who|whom|this|that|these|those|am|is|are|was|were|be|been|being|have|has|had|having|do|does|did|doing|a|an|the|and|but|if|or|because|as|until|while|of|at|by|for|with|about|against|between|into|through|during|before|after|above|below|to|from|up|down|in|out|on|off|over|under|again|values|further|then|once|here|there|when|where|why|how|all|any|both|each|few|more|most|other|some|such|no|nor|not|only|own|same|so|than|too|very|s|t|can|will|just|don|should|now|0o|0s|3a|3b|3d|6b|6o|a|a1|a2|a3|a4|ab|able|about|above|abst|ac|accordance|according|accordingly|across|act|actually|ad|added|adj|ae|af|affected|affecting|affects|after|afterwards|ag|again|against|ah|ain|ain't|aj|al|all|allow|allows|almost|alone|along|already|also|although|always|am|among|amongst|amoungst|amount|an|and|announce|another|any|anybody|anyhow|anymore|anyone|anything|anyway|anyways|anywhere|ao|ap|apart|apparently|appear|appreciate|appropriate|approximately|ar|are|aren|arent|aren't|arise|around|as|a's|aside|ask|asking|associated|at|au|auth|av|available|aw|away|awfully|ax|ay|az|b|b1|b2|b3|ba|back|bc|bd|be|became|because|become|becomes|becoming|been|before|beforehand|begin|beginning|beginnings|begins|behind|being|believe|below|beside|besides|best|better|between|beyond|bi|bill|biol|bj|bk|bl|bn|both|bottom|bp|br|brief|briefly|bs|bt|bu|but|bx|by|c|c1|c2|c3|ca|call|came|can|cannot|cant|can't|cause|causes|cc|cd|ce|certain|certainly|cf|cg|ch|changes|ci|cit|cj|cl|clearly|cm|c'mon|cn|co|com|come|comes|con|concerning|consequently|consider|considering|contain|containing|contains|corresponding|could|couldn|couldnt|couldn't|course|cp|cq|cr|cry|cs|c's|ct|cu|currently|cv|cx|cy|cz|d|d2|da|date|dc|dd|de|definitely|describe|described|despite|detail|df|di|did|didn|didn't|different|dj|dk|dl|do|does|doesn|doesn't|doing|don|done|don't|down|downwards|dp|dr|ds|dt|du|due|during|dx|dy|e|e2|e3|ea|each|ec|ed|edu|ee|ef|effect|eg|ei|eight|eighty|either|ej|el|eleven|else|elsewhere|em|empty|en|end|ending|enough|entirely|eo|ep|eq|er|es|especially|est|et|et-al|etc|eu|ev|even|ever|every|everybody|everyone|everything|everywhere|ex|exactly|example|except|ey|f|f2|fa|far|fc|few|ff|fi|fifteen|fifth|fify|fill|find|fire|first|five|fix|fj|fl|fn|fo|followed|following|follows|for|former|formerly|forth|forty|found|four|fr|from|front|fs|ft|fu|full|further|furthermore|fy|g|ga|gave|ge|get|gets|getting|gi|give|given|gives|giving|gj|gl|go|goes|going|gone|got|gotten|gr|greetings|gs|gy|h|h2|h3|had|hadn|hadn't|happens|hardly|has|hasn|hasnt|hasn't|have|haven|haven't|having|he|hed|he'd|he'll|hello|help|hence|her|here|hereafter|hereby|herein|heres|here's|hereupon|hers|herself|hes|he's|hh|hi|hid|him|himself|his|hither|hj|ho|home|hopefully|how|howbeit|however|how's|hr|hs|http|hu|hundred|hy|i|i2|i3|i4|i6|i7|i8|ia|ib|ibid|ic|id|i'd|ie|if|ig|ignored|ih|ii|ij|il|i'll|im|i'm|immediate|immediately|importance|important|in|inasmuch|inc|indeed|index|indicate|indicated|indicates|information|inner|insofar|instead|interest|into|invention|inward|io|ip|iq|ir|is|isn|isn't|it|itd|it'd|it'll|its|it's|itself|iv|i've|ix|iy|iz|j|jj|jr|js|jt|ju|just|k|ke|keep|keeps|kept|kg|kj|km|know|known|knows|ko|l|l2|la|largely|last|lately|later|latter|latterly|lb|lc|le|least|les|less|lest|let|lets|let's|lf|like|liked|likely|line|little|lj|ll|ll|ln|lo|look|looking|looks|los|lr|ls|lt|ltd|m|m2|ma|made|mainly|make|makes|many|may|maybe|me|mean|means|meantime|meanwhile|merely|mg|might|mightn|mightn't|mill|million|mine|miss|ml|mn|mo|more|moreover|most|mostly|move|mr|mrs|ms|mt|mu|much|mug|must|mustn|mustn't|my|myself|n|n2|na|name|namely|nay|nc|nd|ne|near|nearly|necessarily|necessary|need|needn|needn't|needs|neither|never|nevertheless|new|next|ng|ni|nine|ninety|nj|nl|nn|no|nobody|non|none|nonetheless|noone|nor|normally|nos|not|noted|nothing|novel|nowhere|nr|ns|nt|ny|o|oa|ob|obtain|obtained|obviously|oc|od|of|off|often|og|oh|oi|oj|ok|okay|ol|old|om|omitted|on|once|one|ones|only|onto|oo|op|oq|or|ord|os|ot|other|others|otherwise|ou|ought|our|ours|ourselves|out|outside|over|overall|ow|owing|own|ox|oz|p|p1|p2|p3|page|pagecount|pages|par|part|particular|particularly|pas|past|pc|pd|pe|per|perhaps|pf|ph|pi|pj|pk|pl|placed|please|plus|pm|pn|po|poorly|possible|possibly|potentially|pp|pq|pr|predominantly|present|presumably|previously|primarily|probably|promptly|proud|provides|ps|pt|pu|put|py|q|qj|qu|que|quickly|quite|qv|r|r2|ra|ran|rather|rc|rd|re|readily|really|reasonably|recent|recently|ref|refs|regarding|regardless|regards|related|relatively|research|research-articl|respectively|resulted|resulting|results|rf|rh|ri|right|rj|rl|rm|rn|ro|rq|rr|rs|rt|ru|run|rv|ry|s|s2|sa|said|same|saw|say|saying|says|sc|sd|se|sec|second|secondly|section|see|seeing|seem|seemed|seeming|seems|seen|self|selves|sensible|sent|serious|seriously|seven|several|sf|shall|shan|shan't|she|shed|she'd|she'll|shes|she's|should|shouldn|shouldn't|should've|show|showed|shown|showns|shows|si|side|significant|significantly|similar|similarly|since|sincere|six|sixty|sj|sl|slightly|sm|sn|so|some|somebody|somehow|someone|somethan|something|sometime|sometimes|somewhat|somewhere|soon|sorry|sp|specifically|specified|specify|specifying|sq|sr|ss|st|still|stop|strongly|sub|substantially|successfully|such|sufficiently|suggest|sup|sure|sy|system|sz|t|t1|t2|t3|take|taken|taking|tb|tc|td|te|tell|ten|tends|tf|th|than|thank|thanks|thanx|that|that'll|thats|that's|that've|the|their|theirs|them|themselves|then|thence|there|thereafter|thereby|thered|therefore|therein|there'll|thereof|therere|theres|there's|thereto|thereupon|there've|these|they|theyd|they'd|they'll|theyre|they're|they've|thickv|thin|think|third|this|thorough|thoroughly|those|thou|though|thoughh|thousand|three|throug|through|throughout|thru|thus|ti|til|tip|tj|tl|tm|tn|to|together|too|took|top|toward|towards|tp|tq|tr|tried|tries|truly|try|trying|ts|t's|tt|tv|twelve|twenty|twice|two|tx|u|u201d|ue|ui|uj|uk|um|un|under|unfortunately|unless|unlike|unlikely|until|unto|uo|up|upon|ups|ur|us|use|used|useful|usefully|usefulness|uses|using|usually|ut|v|va|value|various|vd|ve|ve|very|via|viz|vj|vo|vol|vols|volumtype|vq|vs|vt|vu|w|wa|want|wants|was|wasn|wasnt|wasn't|way|we|wed|we'd|welcome|well|we'll|well-b|went|were|we're|weren|werent|weren't|we've|what|whatever|what'll|whats|what's|when|whence|whenever|when's|where|whereafter|whereas|whereby|wherein|wheres|where's|whereupon|wherever|whether|which|while|whim|whither|who|whod|whoever|whole|who'll|whom|whomever|whos|who's|whose|why|why's|wi|widely|will|willing|wish|with|within|without|wo|won|wonder|wont|won't|words|world|would|wouldn|wouldnt|wouldn't|www|x|x1|x2|x3|xf|xi|xj|xk|xl|xn|xo|xs|xt|xv|xx|y|y2|yes|yet|yj|yl|you|youd|you'd|you'll|your|youre|you're|yours|yourself|yourselves|you've|yr|ys|yt|z|zero|zi|zz)$")
)

func (DBHandlerMysql *DBHandlerMysql) ReturnDB() *sql.DB {
	return DBHandlerMysql.DB
}

// Opens db and returns a struct to access it
func OpenDB() *DBHandlerMysql {

	// check environment for the database url
	err := godotenv.Load("go.env")
	if err != nil {
		panic(err.Error())
	}

	databaseURL := os.Getenv("DATABASE_URL_MYSQL")

	//load database connection
	db, err1 := sql.Open("mysql", databaseURL)

	if err1 != nil {
		panic(err.Error())
	} else {
		fmt.Println("no issue")
	}

	dbHandler1 := DBHandlerMysql{db}
	return &dbHandler1
}

// append item if no error
func appendNoError(allData []interface{}, data1 interface{}, err error) []interface{} {
	if err != nil {
		fmt.Println("log error: s" + err.Error())
		return allData
	}
	allData = append(allData, data1)
	return allData
}

// access the DB and get all records
func (DBHandlerMysql DBHandlerMysql) GetRecord(dbTable string) ([]interface{}, error) {
	// allData := []genData{}
	allData := make([]interface{}, 0)
	results, err := DBHandlerMysql.DB.Query("Select * FROM " + dbTable)
	if err != nil {
		return allData, err
	}
	for results.Next() {
		//swtich depending on the database table
		switch dbTable {
		case "UserInfo":
			data1 := model.UserInfo{}
			err = results.Scan(&data1.ID, &data1.Username, &data1.LastLogin, &data1.DateJoin, &data1.CommentItem)
			allData = appendNoError(allData, data1, err)
		case "ItemListing":
			data1 := model.ItemListing{}
			err = results.Scan(&data1.ID, &data1.Username, &data1.Name, &data1.ImageLink, &data1.DatePosted, &data1.CommentItem, &data1.ConditionItem, &data1.Cat, &data1.ContactMeetInfo, &data1.Completion)
			allData = appendNoError(allData, data1, err)
		case "CommentUser":
			data1 := model.CommentUser{}
			err = results.Scan(&data1.ID, &data1.Username, &data1.ForUsername, &data1.Date, &data1.CommentItem)
			allData = appendNoError(allData, data1, err)
		case "CommentItem":
			data1 := model.CommentItem{}
			err = results.Scan(&data1.ID, &data1.Username, &data1.ForItem, &data1.Date, &data1.CommentItem)
			allData = appendNoError(allData, data1, err)
		default:
			return allData, errors.New(dbTable + " not found")
		}
	}
	results.Close()
	return allData, nil
}

// removes stop words and split into array of words with regexp
func CleanWord(input1 string, splitText *regexp.Regexp, stopWords *regexp.Regexp) []string {
	newArr := []string{}

	for _, word1 := range splitText.Split(input1, -1) {
		aa := stopWords.Match([]byte(word1))
		// fmt.Println(word1, aa)

		if !aa {
			newArr = append(newArr, word1)
		}
	}
	return newArr
}

// access the DB and get all records for itemListing
// computes similarity between vectors for each word embedding
// puts the similarity into the struct and returns array of struct to the request
func (DBHandlerMysql *DBHandlerMysql) GetRecordlistingIndex(requestWords string, filterUsername string, filterDate string, filterCat string, embed *word2vec.Embeddings) ([]interface{}, error) {
	// fmt.Printf("getrecordlistingindex: requestWords %v,filterUsername %v,filterDate %v,filterCat %v\n", requestWords, filterUsername, filterDate, filterCat)

	allData := []interface{}{}
	requestWords2 := CleanWord(requestWords, splitText, stopWords2)               //clean and split the words for embeding
	requestWordsEmbed := embed.GetWordEmbeddingCombine(requestWords2, []string{}) //get combined word embedding
	results, err := DBHandlerMysql.DB.Query("Select * FROM ItemListing")

	//calculate the cut off date
	timenow := time.Now().Unix()
	switch filterDate {
	case "7days":
		timenow = timenow - (7 * 24 * 60 * 60)
	case "30days":
		timenow = timenow - (30 * 24 * 60 * 60)
	default:
		timenow = timenow - (365 * 24 * 60 * 60)
	}

	//
	i := 0
	sortIndex := []int{}
	sortArr := []float32{}
	sortInd := []interface{}{}

	if err != nil {
		return sortInd, err
	}

	for results.Next() {
		data1 := model.ItemListing{} //needs a seperate call due to different output
		err = results.Scan(&data1.ID, &data1.Username, &data1.Name, &data1.ImageLink, &data1.DatePosted, &data1.CommentItem, &data1.ConditionItem, &data1.Cat, &data1.ContactMeetInfo, &data1.Completion)
		if err != nil {
			fmt.Println("logger: error at getRecordlisting:" + err.Error())
		}

		//convert string of numbers into int and actual date string
		dateVal, _ := strconv.Atoi(data1.DatePosted)
		data1.DatePosted = time.Unix(int64(dateVal), 0).Format("02-01-2006")

		//checks entry with filter
		if (timenow < int64(dateVal) || filterDate == "All" || filterDate == "") &&
			(filterCat == data1.Cat || filterCat == "All" || filterCat == "") &&
			(filterUsername == data1.Username || filterUsername == "All" || filterUsername == "") &&
			(data1.Completion != "true" || filterUsername != "") {

			requestWordsEmbed2 := embed.GetWordEmbeddingCombine(CleanWord(data1.Name, splitText, stopWords2), []string{})
			addVal := float32(0)
			addVal2 := float32(0)

			for _, word := range requestWords2 { // checks for any similar words in the name string
				if strings.Contains(data1.Name, word) {
					addVal += 0.05
				}
				if strings.Contains(data1.CommentItem, word) { // checks for any similar words in the description string
					addVal2 += 0.005
				}
			}

			addVal3 := math.Min(0.15, math.Max(float64(addVal2), 0))
			addVal4 := math.Min(0.2, math.Max(float64(addVal), 0))
			cosSim := word2vec.CosineSimilarity(requestWordsEmbed, requestWordsEmbed2) // computes similarity of words using their vectors
			data1.Similarity = cosSim + float32(addVal3+addVal4)                       // fmt.Sprintf("%f",  //puts the score into struct

			// fmt.Println(requestWords, data1.Name, cosSim+float32(addVal3+addVal4))
			allData = append(allData, data1)

			sortArr = append(sortArr, cosSim+float32(addVal3+addVal4))
			sortIndex = append(sortIndex, i)
			i++
		}
	}
	results.Close()

	// fmt.Println(i, sortIndex, sortArr, sortInd, filterUsername == "")

	_, sortArr2 := encrypt.MergeSort(sortArr, sortIndex)
	maxLen := len(sortArr)
	newSorted := []interface{}{}

	for idx := maxLen - 1; idx >= 0; idx-- { //sorts results in descending order
		newRow := allData[sortArr2[idx]]
		newSorted = append(newSorted, newRow)
		sortInd = append(sortInd, newRow.(model.ItemListing).ID)
	}

	// variable not used, but might be in the future
	allData = newSorted

	return sortInd, nil
}

// access the DB and get a single record, search using courseName
// based on requested database, it will be marshalled into the struct
func (DBHandlerMysql *DBHandlerMysql) GetSingleRecord(dbTable string, queryString string, queryString2 string) ([]interface{}, error) {
	//queryString examples, " WHERE ID = 1" or "WHERE Username = alvin"
	allData := make([]interface{}, 0)
	results, err := DBHandlerMysql.DB.Query("Select * FROM "+dbTable+" "+queryString+" = ?", queryString2)
	if err != nil {
		return allData, err
	}
	results.Next()
	switch dbTable {
	case "UserSecret":
		data1 := model.UserSecret{}
		err = results.Scan(&data1.ID, &data1.Username, &data1.Password, &data1.IsAdmin, &data1.CommentItem)
		allData = appendNoError(allData, data1, err)
	case "UserInfo":
		data1 := model.UserInfo{}
		err = results.Scan(&data1.ID, &data1.Username, &data1.LastLogin, &data1.DateJoin, &data1.CommentItem)
		allData = appendNoError(allData, data1, err)
	case "ItemListing":
		data1 := model.ItemListing{}
		err = results.Scan(&data1.ID, &data1.Username, &data1.Name, &data1.ImageLink, &data1.DatePosted, &data1.CommentItem, &data1.ConditionItem, &data1.Cat, &data1.ContactMeetInfo, &data1.Completion)
		allData = appendNoError(allData, data1, err)
	case "CommentUser":
		data1 := model.CommentUser{}
		err = results.Scan(&data1.ID, &data1.Username, &data1.ForUsername, &data1.Date, &data1.CommentItem)
		allData = appendNoError(allData, data1, err)
	case "CommentItem":
		data1 := model.CommentItem{}
		err = results.Scan(&data1.ID, &data1.Username, &data1.ForItem, &data1.Date, &data1.CommentItem)
		allData = appendNoError(allData, data1, err)
	default:
		return allData, errors.New(dbTable + " not found in switch")
	}
	results.Close()
	return allData, err
}

// post a record into the DB
// based on requested database, it will be marshalled into the respective struct
func (DBHandlerMysql *DBHandlerMysql) InsertRecord(dbTable string, receiveInfo map[string]string, maxIDInt1 int) error {
	maxIDInt := 1
	if maxIDInt1 == 0 {
		maxID, err1 := DBHandlerMysql.GetMaxID(dbTable)

		if err1 == nil {
			maxIDInt = maxID + 1
		}

	} else {
		maxIDInt = maxIDInt1
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch dbTable {
	case "UserSecret":
		stmt, err := DBHandlerMysql.DB.PrepareContext(ctx, "INSERT INTO "+dbTable+" VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, maxIDInt, receiveInfo["Username"], receiveInfo["Password"], receiveInfo["IsAdmin"], receiveInfo["CommentItem"])
		return err

	case "UserInfo":

		stmt, err := DBHandlerMysql.DB.PrepareContext(ctx, "INSERT INTO "+dbTable+" VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, maxIDInt, receiveInfo["Username"], receiveInfo["LastLogin"], receiveInfo["DateJoin"], receiveInfo["CommentItem"])
		return err

	case "ItemListing":

		stmt, err := DBHandlerMysql.DB.PrepareContext(ctx, "INSERT INTO "+dbTable+" VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, maxIDInt,
			receiveInfo["Username"],
			receiveInfo["Name"],
			receiveInfo["ImageLink"],
			receiveInfo["DatePosted"],
			receiveInfo["CommentItem"],
			receiveInfo["ConditionItem"],
			receiveInfo["Cat"],
			receiveInfo["ContactMeetInfo"],
			receiveInfo["Completion"])
		return err

	case "CommentUser":
		stmt, err := DBHandlerMysql.DB.PrepareContext(ctx, "INSERT INTO "+dbTable+" VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, maxIDInt, receiveInfo["Username"], receiveInfo["ForUsername"], receiveInfo["Date"], receiveInfo["CommentItem"])
		return err

	case "CommentItem":

		stmt, err := DBHandlerMysql.DB.PrepareContext(ctx, "INSERT INTO "+dbTable+" VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, maxIDInt, receiveInfo["Username"], receiveInfo["ForItem"], receiveInfo["Date"], receiveInfo["CommentItem"])
		return err

	default:
		return errors.New(dbTable + " not found in switch")
	}
}

// get the current max ID in the server
func (DBHandlerMysql *DBHandlerMysql) GetMaxID(dbTable string) (int, error) {
	results, err := DBHandlerMysql.DB.Query("SELECT MAX(ID) FROM " + dbTable)
	results.Next()
	var maxID int
	results.Scan(&maxID)
	results.Close()
	if err != nil {
		maxID = 0
	}
	// defer recover() //recover if error from no entry
	return maxID, err
}

// get the current max ID in the server
func (DBHandlerMysql *DBHandlerMysql) GetUsername(dbTable string, id int) (string, error) {
	results, err := DBHandlerMysql.DB.Query("SELECT Username FROM "+dbTable+" WHERE ID=?", id)
	results.Next()
	var username string
	results.Scan(&username)
	results.Close()
	return username, err
}

// edit a single record on DB, chosen record based on ID
// based on requested database, it will be marshalled into the respective struct
func (DBHandlerMysql *DBHandlerMysql) EditRecord(dbTable string, receiveInfo map[string]string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch dbTable {
	// case "UserSecret":
	// 	_, err := DBHandlerMysql .DB.Exec("UPDATE "+dbTable+" SET Username=?, Password=?, IsAdmin=?, CommentItem=? WHERE ID=?", values[0], values[1], values[2], values[3], values[4])receiveInfo["Username"], receiveInfo["LastLogin"], receiveInfo["DateJoin"], receiveInfo["CommentItem"]
	case "UserInfo":

		if _, ok := receiveInfo["ID"]; !ok {
			stmt, err := DBHandlerMysql.DB.PrepareContext(ctx, "UPDATE "+dbTable+" SET LastLogin=?, CommentItem=? WHERE Username=?")

			if err != nil {
				return err
			}
			defer stmt.Close()

			_, err = stmt.ExecContext(ctx, receiveInfo["LastLogin"], receiveInfo["CommentItem"], receiveInfo["Username"])
			return err
		}
		stmt, err := DBHandlerMysql.DB.PrepareContext(ctx, "UPDATE "+dbTable+" SET LastLogin=?, CommentItem=? WHERE ID=?")

		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, receiveInfo["LastLogin"], receiveInfo["CommentItem"], receiveInfo["ID"])

		return err

	case "ItemListing":

		stmt, err := DBHandlerMysql.DB.PrepareContext(ctx, "UPDATE "+dbTable+" SET ImageLink=?, CommentItem=?, ConditionItem=?, Cat=?, ContactMeetInfo=?, Completion=? WHERE ID=?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, receiveInfo["ImageLink"], receiveInfo["CommentItem"], receiveInfo["ConditionItem"], receiveInfo["Cat"], receiveInfo["ContactMeetInfo"], "false", receiveInfo["ID"])
		return err

	case "CommentUser":

		stmt, err := DBHandlerMysql.DB.PrepareContext(ctx, "UPDATE "+dbTable+" SET CommentItem=? WHERE ID=?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, receiveInfo["CommentItem"], receiveInfo["ID"])
		return err

	case "CommentItem":

		stmt, err := DBHandlerMysql.DB.PrepareContext(ctx, "UPDATE "+dbTable+" SET CommentItem=? WHERE ID=?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, receiveInfo["CommentItem"], receiveInfo["ID"])
		return err

	default:
		return errors.New(dbTable + " not found in switch")
	}
}

// delete a single record, chosen based on provided ID
func (DBHandlerMysql *DBHandlerMysql) DeleteRecord(dbTable string, id int) error {
	_, err := DBHandlerMysql.DB.Exec("DELETE FROM "+dbTable+" WHERE ID=?", id)
	return err
}

func (DBHandlerMysql *DBHandlerMysql) AddUser(receiveInfo map[string]string) error {
	// get current max number
	maxIDInput := 0
	maxID, err := DBHandlerMysql.GetMaxID("UserSecret")

	if err == nil {
		maxIDInput = maxID + 1
	}

	// begin ticket
	tx, err := DBHandlerMysql.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			fmt.Println(err.Error())
			fmt.Println("rolled back registration of user: ", receiveInfo["Username"])
			tx.Rollback()
		}
	}()

	//first statement
	stmt, err1 := tx.Prepare("INSERT INTO UserSecret VALUES (?, ?, ?, ?, ?)")

	if err1 != nil {
		err = err1

		return err
	}

	_, err2 := stmt.Exec(maxIDInput, receiveInfo["Username"], receiveInfo["Password"], receiveInfo["IsAdmin"], receiveInfo["CommentItem"])
	if err2 != nil {
		err = err2
		return err
	}

	stmt.Close()
	//2nd statement
	stmt, err3 := tx.Prepare("INSERT INTO UserInfo VALUES (?, ?, ?, ?, ?)")

	if err3 != nil {
		err = err3
		return err
	}

	_, err = stmt.Exec(maxIDInput, receiveInfo["Username"], receiveInfo["LastLogin"], receiveInfo["DateJoin"], receiveInfo["CommentItem"])
	stmt.Close()

	//check error
	return err
}

func (DBHandlerMysql *DBHandlerMysql) CompleteItem(itemID string) error {
	tx, err := DBHandlerMysql.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	//first statement
	stmt, err1 := tx.Prepare("UPDATE ItemListing SET Completion =? WHERE ID=?")
	if err1 != nil {
		err = err1
	}

	_, err = stmt.Exec("true", itemID)
	stmt.Close()

	return err
}
