package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"telegram-bot/internal/db"
	"telegram-bot/internal/helpers"
	"telegram-bot/internal/models"
)

func (client *Client) CheckOrders(update models.Update, callbackId int) {
	dbConn, err := db.ConnectToSQL()
	if err != nil {
		log.Printf("Error connecting to db:%v\n", err)
	}
	repo := NewUserRepository(dbConn)

	user, err := repo.findUser(callbackId)
	log.Printf("Error finding user:%v\n", err)

	cart, err := sendToRedis(user.NotifierID)
	log.Printf("Erorr sending to notifier: %v\n", err)
	textItem := "U+1F6CD U+FE0F Order details:\n"
	for _, item := range cart.Items {
		textItem += fmt.Sprintf("• %s - %d \n", item.Description, item.Quantity)
	}

	client.SendMessage(callbackId, textItem+fmt.Sprintf("Order completion date: %s", cart.UpdatedAt))
}

func sendToRedis(id string) (*models.Cart, error) {
	resp, err := http.Post("http://cart-service:8083/cart/checkout", "application/json", strings.NewReader(fmt.Sprintf(`{"id":"%s"}`, id)))
	if err != nil {
		return nil, helpers.ErrorHelper(err, "error sending request to notifier-service")
	}

	defer resp.Body.Close()
	var order *models.Cart
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		return nil, helpers.ErrorHelper(err, "error parsing orders from JSON")
	}
	return order, nil
}
