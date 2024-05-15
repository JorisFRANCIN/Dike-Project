package controllers

import (
	helper "API/helpers"
	"API/models"
	"API/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = utils.OpenCollection(utils.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}

	return check, msg
}

// Register
//	@Summary		Register a new user
//	@Description	Register a new user with the provided information.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body	models.UserRegister	true	"User information to register"
//	@Success		200		"Successfully registered user"
//	@Failure		400		"Bad Request"
//	@Failure		500		"Internal server error"
//	@Router			/users/register [post]
func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
	
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exist"})
			return
		}
		
		password := HashPassword(*user.Password)
		user.Password = &password
		
		count1, err1 := userCollection.CountDocuments(ctx, bson.M{"username": user.Username})
		defer cancel()
		if err1 != nil {
			log.Panic(err1)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the username"})
		}

		if count1 > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this username already exist"})
			return
		}
		
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refresh_token, _ := helper.GenerateAllTokens(*user.Email, *&user.User_id)
		user.Token = &token
		user.Refresh_Token = &refresh_token

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

// Login
//	@Summary		Log in a user
//	@Description	Log in a user with the provided email and password.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body	models.UserLogin	true	"User login information"
//	@Success		200		"Successfully logged in"
//	@Failure		400		"Bad request, invalid credentials"
//	@Failure		500		"Internal server error"
//	@Router			/users/login [post]
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		}

		token, refresh_token, _ := helper.GenerateAllTokens(*foundUser.Email, *&foundUser.User_id)
		helper.UpdateAllTokens(token, refresh_token, foundUser.User_id)
		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}
}

// GetUsers
//	@Summary		Get a list of users
//	@Description	Get a list of users based on specified parameters.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			recordPerPage	query	int	false	"Number of records per page"
//	@Param			page			query	int	false	"Page number"
//	@Param			startIndex		query	int	false	"Start index"
//	@Security		API_Token
//	@Success		200	"List of users"
//	@Failure		400	"Bad request or validation error"
//	@Failure		500	"Internal server error"
//	@Router			/users [get]
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("Page"))
		if err1 != nil || page < 1 {
			page = 10
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}},
			{"total_count", bson.D{{"$sum", 1}}},
			{"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
			}}}
		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		}
		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allUsers[0])
	}
}

// GetUser
//	@Summary		Get user details
//	@Description	Get details of a specific user based on the user ID.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path	string	true	"User ID"
//	@Security		API_Token
//	@Success		200	{object}	models.User	"User details"
//	@Failure		400	"Bad request or validation error"
//	@Failure		500	"Internal server error"
//	@Router			/users/{user_id} [get]
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helper.MatchToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// Cron
//	@Summary		Cron inputs
//	@Description	Create the Cron action
//	@Tags			Services
//	@Accept			json
//	@Produce		json
//	@Param			cron	body	models.Cron	true	"Cron informations"
//	@Security		API_Token
//	@Success		200		"Successfully registered user"
//	@Failure		400		"Bad Request"
//	@Failure		500		"Internal server error"
//	@Router			/users/Cron [post]
func Trigger() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var cron models.Cron

		if err := c.BindJSON(&cron); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(cron)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{"minutes": cron.Minutes, "hours": cron.Hours, "DayOfMonth": cron.Day, "DayOfWeek": cron.Week})
	}
}

// Notify
//	@Summary		Notify inputs
//	@Description	Create the Notify Reaction
//	@Tags			Services
//	@Accept			json
//	@Produce		json
//	@Param			notify	body	models.Notify	true	"Notify informations"
//	@Security		API_Token
//	@Success		200		"Successfully registered user"
//	@Failure		400		"Bad Request"
//	@Failure		500		"Internal server error"
//	@Router			/users/Notify [post]
func ReactionNotify() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var notify models.Notify

		if err := c.BindJSON(&notify); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(notify)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{"title": notify.Title})
	}
}
