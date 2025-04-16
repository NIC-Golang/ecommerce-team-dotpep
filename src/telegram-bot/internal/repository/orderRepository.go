package repository

import (
	"context"
	"fmt"
	"log"
	"telegram-bot/internal/db"
	"telegram-bot/internal/helpers"
	"telegram-bot/internal/models"
	"time"
)

func (client *Client) CheckOrders(update models.Update, callbackId int) {
	dbConn, err := db.ConnectToSQL()
	if err != nil {
		log.Printf("Error connecting to db:%v\n", err)
	}
	repo := NewUserRepository(dbConn)

	user, err := repo.findUser(callbackId)
	log.Printf("Error finding user:%v\n", err)

	order, err := helpers.GetRedisOrder(user.NotifierID)
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

func (client *Client) CheckStatus(id string, ctx context.Context, update *models.Update) error {
	order, err := helpers.GetRedisOrder(id)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(30 * time.Second)
	status := order.Status
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			updatedOrder, err := helpers.GetRedisOrder(id)
			if err != nil {
				return err
			}
			if status != updatedOrder.Status {
				status = updatedOrder.Status
				text := helpers.CreateStatusMessage(updatedOrder)
				if text != "" {
					client.SendMessage(update.CallbackQuery.Message.Chat.Id, text)
				} else {
					return fmt.Errorf("error during sending a message")
				}
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func RunNotifications(id int) {
	mode, _ := NotifyMode(id)
	if mode == "on" {
		//TODO
	} else {
		return
	}
}
