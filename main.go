package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	TaskID    string             `json:"taskid"`
	Body      string             `json:"body"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("GO API")
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading env!")
	}
	PORT := os.Getenv("PORT")
	MongoDB_URI := os.Getenv("MONGODB_URI")

	clientOptions := options.Client().ApplyURI(MongoDB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal("Mongo Err : ", err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal("mongo err : ", err)
	}

	fmt.Println("Connected to mongoDB")

	collection = client.Database("todo_db").Collection("todos")

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",                           // Allow all origins_Test purpose
		AllowMethods: "GET,POST,PUT,DELETE,PATCH",   // Allow specific HTTP methods
		AllowHeaders: "Content-Type, Authorization", // Allow specific headers
	}))

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.Status(200).JSON(fiber.Map{"msg": "Server setup done!"})
	// })

	app.Static("/", "./client/dist")

	app.Get("/api/gettodos", getTodos)
	app.Post("/api/addtodo", addTodo)
	app.Patch("/api/updatetodo/:id", updateTodo)
	app.Delete("/api/deletetod/:id", deleteTodo)

	app.Listen(":" + PORT)
}

func getTodos(c *fiber.Ctx) error {

	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}

	return c.JSON(todos)

}

func addTodo(c *fiber.Ctx) error {
	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
		return err
	}
	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Body can't be empty!"})
	}
	currentTime := time.Now()

	// Get current Unix timestamp in milliseconds
	unixTimestampMillis := currentTime.UnixNano() / int64(time.Millisecond)
	//convert int64 to string
	inString := fmt.Sprintf("%d", unixTimestampMillis)
	fmt.Println(inString)
	todo.TaskID = inString
	insertResult, err := collection.InsertOne(context.Background(), todo)

	if err != nil {
		return err
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(200).JSON(todo)

}

func updateTodo(c *fiber.Ctx) error {

	taskid := c.Params("id")

	fmt.Println(taskid)
	filter := bson.M{"taskid": taskid}
	findTask := collection.FindOne(context.Background(), filter)
	var foundedResult bson.M
	err := findTask.Decode(&foundedResult)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No document found with the given taskid")
		} else {
			log.Fatal(err)
		}
	}

	// Print the found document
	fmt.Printf("Found document: %v\n", foundedResult["completed"])

	var update bson.M

	if completed, ok := foundedResult["completed"].(bool); ok {
		if completed {
			update = bson.M{"$set": bson.M{"completed": false}}
		} else {
			update = bson.M{"$set": bson.M{"completed": true}}
		}

		// Perform the update operation
		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Document updated successfully")
	} else {
		fmt.Println("The 'completed' field is not of type boolean")
	}

	findResult := collection.FindOneAndUpdate(context.Background(), filter, update)
	var updatedRest bson.M

	Updateerr := findResult.Decode(&updatedRest)
	if Updateerr != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No document found with the given taskid")
		} else {
			log.Fatal(err)
		}
	}

	return c.Status(200).JSON(updatedRest)

}

func deleteTodo(c *fiber.Ctx) error {

	taskid := c.Params("id")

	filter := bson.M{"taskid": taskid}

	deleteResult := collection.FindOneAndDelete(context.Background(), filter)

	var result bson.M
	err := deleteResult.Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No document found with the given taskid")
		} else {
			log.Fatal(err)
		}
	}

	return c.JSON(result)

}
