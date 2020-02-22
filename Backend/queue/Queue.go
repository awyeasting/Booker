package queue

import (
	"Booker/Backend/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	log "github.com/sirupsen/logrus"

	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetQueueUsersCollection(c context.Context) *mongo.Collection {
	db, ok := c.Value("db").(*mongo.Client)
	if !ok {
		log.Panic("No database context found")
	}

	// Get the handle on the queue users collection based on information from the config
	BookerDB := db.Database(QUEUE_DB_NAME)
	QueueUsersCol := BookerDB.Collection(QUEUE_USERS_COL_NAME)

	return QueueUsersCol
}

func GetBookQueueCollection(c context.Context) *mongo.Collection {
	db, ok := c.Value("db").(*mongo.Client)
	if !ok {
		log.Panic("No database context found")
	}

	// Get the handle on the queue users collection based on information from the config
	BookerDB := db.Database(QUEUE_DB_NAME)
	BookQueueCol := BookerDB.Collection(BOOK_QUEUE_COL_NAME)

	return BookQueueCol
}

// Route all book club queue related API handles to the proper handlers
func QueueRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.AllowContentType("application/json"))

	r.Post("/join", JoinQueueHandler)

	return r
}

func JoinQueueHandler(w http.ResponseWriter, r *http.Request) {

	// Unmarshal json from request body
	var reqModel models.JoinRequest
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &reqModel)
	var res models.ResponseResult
	if err != nil {
		WriteJSONResponse(w, res, http.StatusBadRequest)
		log.Error(err)
		return
	}

	// Set user information as joined in those specific queues
	err = SetQueueUserInfo(reqModel, GetQueueUsersCollection(r.Context()))
	if err != nil {
		WriteJSONResponse(w, res, http.StatusInternalServerError)
		log.Error(err)
		return
	}

	// Join queue for every book
	for _, book := range reqModel.BookIds {
		go JoinQueue(reqModel.UserId, book, GetBookQueueCollection(r.Context()))
	}
}

func SetQueueUserInfo(reqModel models.JoinRequest, QueueUsersCol *mongo.Collection) (err error) {
	// Set queue user info
	_, err = QueueUsersCol.InsertOne(context.TODO(), reqModel)
	if err != nil {
		return err 
	}

	return nil
}

func JoinQueue(userId string, bookId string, BookQueueCol *mongo.Collection) (err error) {
	// Check if book queue already exists
	var bq models.BookQueue
	err = BookQueueCol.FindOne(context.TODO(), bson.D{{"bookId", bookId}}).Decode(&bq)
	if err != nil {
		// Check if the book queue is not started
		if err.Error() == "mongo: no documents in result" {
			// Create book queue
			bq.BookId = bookId
			bq.UserIds = []string{userId}

			// Start book's queue
			_, err = BookQueueCol.InsertOne(context.TODO(), bq)
			if err != nil {
				return err
			}

			return nil
		}

		return err
	}

	bq.UserIds = append(bq.UserIds, userId) 

	// Add user to existing queue
	filter := bson.M{"bookId": bookId}
	update := bson.D{
		{"$set", bson.D{
			{"userIds", bq.UserIds},
		}},
	}
	_, err = BookQueueCol.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	// TODO: Check book queue status

	return nil
}

func CheckQueues(toCheck chan string) {
	// TODO: Select a bookId from the toCheck queue

	// TODO: Check if that queue is ready to pop

}

func CheckQueue(bookId string) {
	// TODO: Given a bookId check if it has enough members
}

func PopQueue(bookId string) {
	// TODO: Pull users on that book queue from all their active queues

	// TODO: Place users in pending status for the book (and generate clubId)

	// TODO: Notify users

}

func ConfirmPop(userId string, bookId string, clubId string) {
	// TODO: Mark user as ready in the pending club

	// TODO: Put it in the checking channel
}

func CheckConfirmations(toCheck chan string) {
	// TODO: Given a clubId check if enough members have confirmed
}

func CheckConfirmation(clubId string) {
	// TODO: If the confirmation window hasn't passed then check if all have confirmed
	// TODO: Else check if enough have confirmed
		// TODO: Re-enter all confirmed users into their queues
		// TODO: Notify all unconfirmed users that they have been removed from queue

	// TODO: Send all confirmed members the chat group
}

// Writes a given interface to an http ResponseWriter with a given status code
func WriteJSONResponse(w http.ResponseWriter, jsonResponse interface{}, header int) {
	json.NewEncoder(w).Encode(jsonResponse)
	w.WriteHeader(header)
}