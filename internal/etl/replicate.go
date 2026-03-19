package etl

import (
	"context"
	"etl-pract/internal/db"
	"etl-pract/internal/models"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Replicate(mongoDB *mongo.Database) error {
	ctx := context.Background()
	coll := mongoDB.Collection("customers")

	lastSync := readLastSyncTime()
	now := time.Now()

	rows, err := db.PgPool.Query(ctx,
		`SELECT id, name, email, created_at 
         FROM customers 
         WHERE created_at > $1`, lastSync)
	if err != nil {
		return err
	}
	defer rows.Close()

	newCustomers := make(map[int]*models.Customer)
	for rows.Next() {
		var c models.Customer
		var createdAt time.Time
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &createdAt); err != nil {
			return err
		}
		c.Orders = []models.Order{}
		c.SyncedAt = now
		newCustomers[c.ID] = &c
	}

	orderRows, err := db.PgPool.Query(ctx,
		`SELECT o.id, o.customer_id, o.product, o.amount, o.status, o.created_at
         FROM orders o
         WHERE o.updated_at > $1`, lastSync)
	if err != nil {
		return err
	}
	defer orderRows.Close()

	for orderRows.Next() {
		var o models.Order
		var customerID int
		var createdAt time.Time
		if err := orderRows.Scan(&o.OrderID, &customerID, &o.Product, &o.Amount, &o.Status, &createdAt); err != nil {
			return err
		}
		o.PlacedAt = createdAt

		customer, exists := newCustomers[customerID]
		if !exists {
			customer = &models.Customer{
				ID:       customerID,
				Name:     "",
				Email:    "",
				Orders:   []models.Order{},
				SyncedAt: now,
			}
			newCustomers[customerID] = customer
		}
		customer.Orders = append(customer.Orders, o)
	}

	for _, c := range newCustomers {
		if err := UpsertCustomer(ctx, coll, *c); err != nil {
			fmt.Println("Error upserting customer:", c.ID, err)
		}
	}

	saveLastSyncTime(now)

	fmt.Printf("Replicated %d customers and orders\n", len(newCustomers))
	return nil
}

func UpsertCustomer(ctx context.Context, coll *mongo.Collection, c models.Customer) error {
	filter := bson.M{"_id": c.ID}

	update := bson.M{
		"$setOnInsert": bson.M{
			"name":   c.Name,
			"email":  c.Email,
			"orders": []models.Order{},
		},
		"$set": bson.M{
			"name":      c.Name,
			"email":     c.Email,
			"synced_at": time.Now(),
		},
	}

	if len(c.Orders) > 0 {
		update["$addToSet"] = bson.M{
			"orders": bson.M{
				"$each": c.Orders,
			},
		}
	}

	opts := options.Update().SetUpsert(true)
	_, err := coll.UpdateOne(ctx, filter, update, opts)
	return err
}
