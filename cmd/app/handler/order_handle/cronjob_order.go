package order_handle

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
	"os"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/proto/models"
	"time"
)

func ScheduleOrderStatusUpdate() {
	cronJob := cron.New()
	cronJob.AddFunc("@every 1m", func() {
		pendingOrders, err := dbms.GetOrdersByStatusNoPage(models.Pending)
		if err != nil {
			fmt.Printf("Error getting pending orders: %v\n", err)
			return
		}

		for _, order := range pendingOrders {
			if time.Since(order.CreatedAt) > 1*time.Minute {
				err := dbms.UpdateOrderStatus(order.OrderID, string(models.OrderBeingDelivered))
				if err != nil {
					fmt.Printf("Error updating order status for ID %d: %v\n", order.OrderID, err)
				} else {
					fmt.Printf("Order status updated to OrderBeingDelivered for ID %d\n", order.OrderID)
					user, _ := dbms.GetUserByID(order.UserID)
					err := sendOrderStatusUpdateEmail(user.Email, fmt.Sprintf("%d", order.OrderID))
					if err != nil {
						fmt.Printf("Error sending email notification for order ID %d: %v\n", order.OrderID, err)
					} else {
						fmt.Printf("Email notification sent for order ID %d\n", order.OrderID)
					}
				}
			}
		}
	})
	cronJob.Start()
}

func sendOrderStatusUpdateEmail(email, orderID string) error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	emailAddress := os.Getenv("EMAIL_ADDRESS")
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	emailhost := os.Getenv("EMAIL_HOST")
	subject := "Order Status Update"
	body := fmt.Sprintf("Your order with ID %s has been confirmed and is currently being shipped. Thank you for shopping with us!", orderID)

	m := gomail.NewMessage()
	m.SetHeader("From", emailAddress)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	dialer := gomail.NewDialer(emailhost, 587, emailAddress, emailPassword)

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
