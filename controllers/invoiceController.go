package controllers

import (
	"context"
	"log"
	"net/http"
	"resturant-management/database"
	"resturant-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoiceViewFormat struct {
	Invoice_id           string
	Payment_method       string
	Order_id             string
	Payment_status       *string
	Payment_due          interface{}
	Table_number         interface{}
	Payment_due_date     time.Time
	Order_details        interface{}
}

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")
func GetInvoices() gin.HandlerFunc{
	return func(c *gin.Context) {
		var	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		res, err := invoiceCollection.Find(context.TODO(), bson.M{})
        defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		var allInvoices []bson.M

		if err = res.All(ctx, &allInvoices); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allInvoices)
	}
}

func GetInvoice() gin.HandlerFunc{
	return func(c *gin.Context) {
		var	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		invoiceId := c.Param("invoice_id")
		var invoice models.Invoice

		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		var invoiceView InvoiceViewFormat
		// allOrderItems, err : = ItemsByOrder(invoice.Order_id)
		invoiceView.Order_id = invoice.Order_id
		invoiceView.Payment_due_date = invoice.Payment_due_date
		invoice.Payment_method = nil
		if invoice.Payment_method != nil {
			invoiceView.Payment_method = *invoice.Payment_method
		}

	}
}

func CreateInvoice() gin.HandlerFunc{
	return func(c *gin.Context) {}
}

func UpdateInvoice() gin.HandlerFunc{
	return func(ctx *gin.Context) {}
}