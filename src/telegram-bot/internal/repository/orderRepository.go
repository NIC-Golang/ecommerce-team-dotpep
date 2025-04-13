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

	order, err := sendToRedis(user.NotifierID)
	log.Printf("Erorr sending to notifier: %v\n", err)
	textItem := "🛍️ Order details:\n"
	for _, item := range order.Items {
		textItem += fmt.Sprintf("• %s - %d \n", item.Description, item.Quantity)
	}
	orderNum := fmt.Sprintf("📦 Your order number: %s\n", order.OrderNumber)
	orderDate := fmt.Sprintf("📅 Order completion date: %s\n", order.CreatedAt.Format("02 Jan 2006 15:04"))
	orderStatus := fmt.Sprintf("🚚 Order status: %s\n", order.Status)
	if len(order.Items) == 0 {
		client.SendMessage(callbackId, "📦 You don't have any items in your order")
		return
	}
	client.SendMessage(callbackId, textItem+orderNum+orderDate+orderStatus+"💬 We will send you a notification as soon as the order is delivered!")
}

func sendToRedis(id string) (*models.Order, error) {
	resp, err := http.Post("http://cart-service:8083/cart/order", "application/json", strings.NewReader(fmt.Sprintf(`{"id":"%s"}`, id)))
	if err != nil {
		return nil, helpers.ErrorHelper(err, "error sending request to notifier-service")
	}

	defer resp.Body.Close()
	var order *models.Order
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		return nil, helpers.ErrorHelper(err, "error parsing orders from JSON")
	} else if order == nil {
		return nil, fmt.Errorf("order is empty")
	}
	return order, nil
}
