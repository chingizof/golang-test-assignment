package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Transfer struct {
	Name     string    `bson:"name"`
	Price    int       `bson:"price"`
	Date     time.Time `bson:"date"`
	Income   string    `bson:"income"`
	Comment  string    `bson:"comment"`
	Category string    `bson:"category"`
}

var collection *mongo.Collection
var ctx = context.TODO()

func init() {
	clientOptions := options.Client().ApplyURI("connection") //здесь мы подключаемся к MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("transfers").Collection("transfer")
}

func main() {
	app := &cli.App{
		Name:  "transfers_maker",
		Usage: "A simple CLI program to manage your financial literacy",
		Commands: []*cli.Command{
			{
				Name:    "add", //первая функция, добавить трансфер
				Aliases: []string{"a"},
				Usage:   "add a transfer to the list",
				Action: func(c *cli.Context) error {
					var str string
					var integer int
					var income string
					var comment string
					var category string
					fmt.Println("name of your transfer") //запрашиваем заполнить поля
					fmt.Scanf("%s", &str)
					fmt.Println("Price:")
					fmt.Scanf("%d", &integer)
					fmt.Println("type only 'income' or 'outcome'")
					fmt.Scan(&income)
					fmt.Println("any comments?")
					fmt.Scan(&comment)
					fmt.Println("What category? only 'food', 'transport' for outcome and 'salary' for income")
					fmt.Scan(&category)
					if str == "" {
						return errors.New("cannot add an empty transfer") //нельзя проводить пустой трансфер
					}
					if income != "income" && income != "outcome" {
						return errors.New("is it income or outcome?") //только income или outcome
					}
					if category != "food" && category != "transport" && category != "salary" { //только ограниченные категории
						return errors.New("please type right category")
					}
					transfer := &Transfer{ //заполняем структуру
						Name:     str,
						Price:    integer,
						Date:     time.Now(),
						Income:   income,
						Comment:  comment,
						Category: category,
					}

					return createTransfer(transfer)
				},
			},
			{
				Name:    "all", //вторая функция, увидеть список трансферов.
				Aliases: []string{"l"},
				Usage:   "list all transfers",
				Action: func(c *cli.Context) error {
					transfers, err := getAll()
					if err != nil {
						if err == mongo.ErrNoDocuments {
							fmt.Print("Nothing to see here.\nRun `add 'transfer'` to add a transfer")
							return nil
						}

						return err
					}

					printTransfers(transfers)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func createTransfer(transfer *Transfer) error { //добавляем трансфер в  таблицу базы данных
	_, err := collection.InsertOne(ctx, transfer)
	return err
}

func getAll() ([]*Transfer, error) { //запрашиваем все таблицы в формате bson
	filter := bson.D{{}}
	return filterTransfers(filter)
}

func filterTransfers(filter interface{}) ([]*Transfer, error) { //проводим расшифровку
	var transfers []*Transfer

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return transfers, err
	}

	for cur.Next(ctx) {
		var t Transfer
		err := cur.Decode(&t) //расшифровываем dson в string
		if err != nil {
			return transfers, err
		}

		transfers = append(transfers, &t)
	}

	if err := cur.Err(); err != nil {
		return transfers, err
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	if len(transfers) == 0 {
		return transfers, mongo.ErrNoDocuments
	}

	return transfers, nil
}

func printTransfers(transfers []*Transfer) { //печатаем трансферы в консоль
	for i, v := range transfers {
		fmt.Println("     ITEM", i+1)
		fmt.Println("Name:", v.Name)
		fmt.Println("Price:", v.Price)
		fmt.Println("Date:", v.Date)
		fmt.Println("Income:", v.Income)
		fmt.Println("Comments:", v.Comment)
		fmt.Println("Category:", v.Category)
		fmt.Println("===================")
	}
}
