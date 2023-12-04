package order_handle

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/proto/models"
)

var orderUpdateCron *cron.Cron

func updateOrderStatus(orderID int32) {
	order, err := dbms.GetOrderByID(orderID)
	if err != nil {
		fmt.Printf("Error getting order with ID %d: %v\n", orderID, err)
		return
	}

	if order.Status == models.Pending {
		if err := dbms.UpdateOrderStatus(orderID, string(models.OrderBeingDelivered)); err != nil {
			fmt.Printf("Error updating order status for ID %d: %v\n", orderID, err)
		} else {
			fmt.Printf("Order status updated to OrderBeingDelivered for ID %d\n", orderID)
		}
	}
}

func startOrderUpdateCron(orderID int32) {
	if orderUpdateCron == nil {
		orderUpdateCron = cron.New()
		orderUpdateCron.AddFunc("@every 3h", func() {
			updateOrderStatus(orderID)
		})
		orderUpdateCron.Start()
	}
}
