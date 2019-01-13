package usercontroller

import (
	models "GetItDone-goserver/models"
	variables "GetItDone-goserver/variables"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var store = sessions.NewCookieStore([]byte(variables.MongoSecret))

func Test(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}

func startMongo() *mgo.Database {
	// start database
	session, err := mgo.Dial(variables.MongoURL)
	if err != nil {
		fmt.Println(err)
	}
	db := session.DB("getitdone")
	return db
}

func hashPassword(s string) []byte {
	// creating hash from password
	password := []byte(s)
	hash, err := bcrypt.GenerateFromPassword(password, 11)
	if err != nil {
		fmt.Println(err)
	}
	return hash
}

func CreateUser(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	db := startMongo()

	// Get user data from client request
	userData := &models.User{}
	fmt.Println(userData)
	c.Bind(userData)

	// creates a hash byte
	Hashword := hashPassword(userData.Password)

	// convert hash to string to be stored in user struct which demands a string
	userData.Password = fmt.Sprintf("%s", Hashword)

	// store user information into the mongo lab database
	db.C("userdata").Insert(&userData)

	// print user data to check if it is working
	fmt.Println(userData.Username)
	fmt.Println(userData.Password)
	// send message back to client
	c.JSON(200, gin.H{
		"message": "creation completed",
	})

}

func LoginUser(c *gin.Context) {
	// set headers on request
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	// Get user data from client request
	userData := &models.User{}
	c.Bind(userData)

	// print user data to check if it is working
	fmt.Println(userData.Username)
	fmt.Println(userData.Password)

	// start database
	session, err := mgo.Dial(variables.MongoURL)
	if err != nil {
		fmt.Println(err)
	}
	db := session.DB("getitdone")

	// set up User model
	var users []models.User
	// search for all usera
	db.C("userdata").Find(bson.M{}).All(&users)
	fmt.Println(users)
	for i := 0; i < len(users); i++ {
		// compare entered password and the current hashed password in database
		passwordCheck := bcrypt.CompareHashAndPassword([]byte(users[i].Password), []byte(userData.Password))
		fmt.Println(passwordCheck)
		if users[i].Username == userData.Username && passwordCheck == nil {
			// set max age of store above -1
			store.MaxAge(86400 * 30)
			// set up a session with error handlignto start storing
			session, err := store.Get(c.Request, "session-name")
			if err != nil {
				c.JSON(200, gin.H{
					"message": "Error: Session not found",
				})
				return
			}
			// set session key pair and save it
			session.Values["username"] = userData.Username
			session.Save(c.Request, c.Writer)
			fmt.Println(session.Values["username"])
			// send success statement to client
			c.JSON(200, gin.H{
				"message": "Login Successful",
			})
			return
		}
		// if all users list have been exhausted
		if i == len(users)-1 {
			// send message back to client
			c.JSON(200, gin.H{
				"message": "Login username or password is incorrect",
			})
			return
		}
	}
}

func CheckSession(c *gin.Context) {
	// use gin context in session
	session, err := store.Get(c.Request, "session-name")
	if err != nil {
		c.JSON(200, gin.H{
			"message": "Error: Session not found",
		})
		return
	}
	fmt.Println(session.Values["username"])
	// send back whatever current value of username is currently checking for login state
	c.JSON(200, gin.H{
		"username": session.Values["username"],
	})
	return;
}

func LogOut(c *gin.Context) {
	// set max age to end now ending session
	store.MaxAge(-1)
	c.JSON(200, gin.H{
		"message": "logout successful",
	})
}
