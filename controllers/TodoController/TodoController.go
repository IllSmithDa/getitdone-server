package todocontroller

import (
	models "GetItDone-goserver/models"
	variables "GetItDone-goserver/variables"
	"fmt"
	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var store = sessions.NewCookieStore([]byte(variables.MongoSecret))

func startMongo() *mgo.Database {
	// start database
	session, err := mgo.Dial(variables.MongoURL)
	if err != nil {
		fmt.Println(err)
	}
	db := session.DB("getitdone")
	return db
}

// create Id to indentify a todoitem
func createUniqueID(todolist []models.Todoitem) int {
	// check if ID is unique or if it already exists in database
	isIDNew := true
	uniqueID := 0
	for {
		// set id to true as default;
		isIDNew = true

		// gnerate ID
		uniqueID = 10000000 + rand.Intn(99999999-10000000)

		// iterate through user's to do list
		for i := 0; i < len(todolist); i++ {
			// set variable to false if id exists in datavase
			if todolist[i].ID == uniqueID {
				isIDNew = false
				break
			}
		}
		// repeat until id is original
		if isIDNew == true {
			break
		}
	}
	return uniqueID
}

func AppendTodoList(c *gin.Context) {
	// access the session variables
	
		session, err := store.Get(c.Request, "session-name")
		if err != nil {
			c.JSON(200, gin.H{
				"message": "Error: Session not found",
			})
			return
		}
		username := session.Values["username"]
	

	//create todomodel and append data to it
	todoitem := models.Todoitem{}
	entry := &models.Todoitem{}

	err = c.Bind(entry)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"message": err,
		})
		return
	}

	// setting up item 
	fmt.Println(entry.Todoitem)
	todoitem.Todoitem = entry.Todoitem
	todoitem.Username = entry.Username
	fmt.Println(todoitem.Todoitem)

	// start mongo
	db := startMongo()

	// set up variable which stores array of Todos
	var allTodos []models.Todoitem
	err = db.C("tododata").Find(bson.M{"username": username}).All(&allTodos);
	if err != nil {
		c.JSON(200, gin.H{
			"message": "data not found ",
		})
		return
	}
	
	// create unique id
	newID := createUniqueID(allTodos);
	todoitem.ID = newID;

	// insert entry into database
	db.C("tododata").Insert(todoitem)
	c.JSON(200, gin.H{
		"success": true,
	})
	return

}


func GetTodoList(c *gin.Context) {

	// access the session variables
	session, err := store.Get(c.Request, "session-name")
	if err != nil {
		c.JSON(200, gin.H{
			"message": "Error: Session not found",
		})
		return
	}
	username := session.Values["username"]
	// start mongo
	db := startMongo()
	// set up User model

	// set up variable which stores array of Todos
	var allTodos []models.Todoitem
	err = db.C("tododata").Find(bson.M{"username": username}).All(&allTodos);
	if err != nil {
		c.JSON(200, gin.H{
			"message": "data not found ",
		})
		return
	}
	var userTodos []string
	var todoIDs []int
	for i := 0; i < len(allTodos); i++ {
		userTodos = append(userTodos, allTodos[i].Todoitem)
		todoIDs = append(todoIDs, allTodos[i].ID)
	}
	fmt.Println(todoIDs)
	c.JSON(200, gin.H{
		"todolist": userTodos,
		"todoIDs":  todoIDs,
	})

	return
}

func EditTodoList(c *gin.Context) {

	//create todomodel and append data to it
	todoitem := models.Todoitem{}
	entry := &models.Todoitem{}
	fmt.Println("point 0")
	err := c.Bind(entry)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"message": err,
		})
		return
	}
	fmt.Println("point 1")
	fmt.Println(entry.Todoitem)
	todoitem.Todoitem = entry.Todoitem
	todoitem.ID = entry.ID
	todoitem.Username = entry.Username
	fmt.Println(todoitem.Todoitem)
	fmt.Println(entry.ID)

	// start mongo
	db := startMongo()

	db.C("tododata").Update(bson.M{"id": todoitem.ID}, todoitem)
	c.JSON(200, gin.H{
		"success": true,
	})

}

func DeleteTodoList(c *gin.Context) {
	
	// access the session variables
	//create todomodel and append data to it
	todoitem := models.Todoitem{}
	entry := &models.Todoitem{}
	err := c.Bind(entry)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"message": err,
		})
		return
	}

	// set query id
	todoitem.ID = entry.ID
	fmt.Println(todoitem.ID)

	// start mongo
	db := startMongo()

	// remove data based on matching id in documents
	err = db.C("tododata").Remove(bson.M{"id": todoitem.ID})
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"message": err,
		})
		return
	}
	// return true if successful
	c.JSON(200, gin.H{
		"success": true,
	})
}


