package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

type Employee struct {
	ID     string  `json:"id,omitempty" bson:"_id,omitempty"` // bson is for mongodb
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    float64 `json:"age"`
}

var mg MongoInstance

const dbName = "trial"

func connectMongo(uri string) error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}
	defer context.WithCancel(context.TODO())
	db := client.Database(dbName)
	mg = MongoInstance{
		Client: client,
		Db:     db,
	}

	return nil
}

func main() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	mongoUserPass, exists := os.LookupEnv("MONGO_USER_PASS")
	if !exists {
		return
	}
	uri := fmt.Sprintf("mongodb+srv://cgoyani:%s@trial.cek3scp.mongodb.net/?retryWrites=true&w=majority", mongoUserPass)
	if err := connectMongo(uri); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	app.Get("/employee", func(c *fiber.Ctx) error {
		query := bson.D{{}}
		cursor, err := mg.Db.Collection("employees").Find(c.Context(), query)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		var employees []Employee = make([]Employee, 0)
		if err := cursor.All(c.Context(), &employees); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(employees)
	})
	app.Post("/employee", func(c *fiber.Ctx) error {
		collection := mg.Db.Collection("employees")
		employee := new(Employee)
		if err := c.BodyParser(employee); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		employee.ID = ""
		insertResult, err := collection.InsertOne(c.Context(), employee)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		filter := bson.D{{Key: "_id", Value: insertResult.InsertedID}}
		cursor := collection.FindOne(c.Context(), filter)
		createdEmployee := &Employee{}
		cursor.Decode(createdEmployee)

		return c.Status(201).JSON(createdEmployee)

	})
	app.Put("/employee/:id", func(c *fiber.Ctx) error {
		collection := mg.Db.Collection("employees")
		idParam := c.Params("id")
		employeeId, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return c.SendStatus(400)
		}
		var employee Employee
		if err := c.BodyParser(&employee); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		query := bson.D{{Key: "_id", Value: employeeId}}
		update := bson.D{
			{
				Key: "$set",
				Value: bson.D{
					{Key: "name", Value: employee.Name},
					{Key: "age", Value: employee.Age},
					{Key: "salary", Value: employee.Salary},
				},
			},
		}

		err = collection.FindOneAndUpdate(c.Context(), query, update).Err()
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.SendStatus(404) // 404 is for not found :)
			}
			return c.SendStatus(500)
		}

		employee.ID = idParam
		return c.Status(200).JSON(employee)

		// cursor := collection.FindOne(c.Context(), query)
		// updatedEmployee := &Employee{}
		// cursor.Decode(updatedEmployee)

		// return c.Status(201).JSON(updatedEmployee)
	})

	app.Delete("/employee/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		employeeId, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return c.SendStatus(400)
		}

		query := bson.D{{Key: "_id", Value: employeeId}}
		deleteResult, err := mg.Db.Collection("employees").DeleteOne(c.Context(), query)
		if err != nil {
			return c.SendStatus(500)
		}
		if deleteResult.DeletedCount < 1 {
			return c.SendStatus(404)
		}

		return c.Status(200).JSON("record deleted")

	})

	log.Fatal(app.Listen(":3000"))
	// defer func() { // postpone the execution of a function until the surrounding function has been executed
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()
	// Send a ping to confirm a successful connection
	// if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}
