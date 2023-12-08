package utils_handle

import (
	"fmt"
	"github.com/streadway/amqp"
	"gopkg.in/gomail.v2"
	"log"
	"strconv"
	"strings"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/proto/models"
)

func SendDiscountMessagesToRabbitMQ(discounts []models.Discount) error {
	conn, err := amqp.Dial("amqps://ithkfqls:w4d4HaTpn_cDiod1r9BOT-CSZON1gYVF@octopus.rmq3.cloudamqp.com/ithkfqls")
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"discount_notification_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	var messages []string

	for _, discount := range discounts {
		valueString := strconv.FormatFloat(discount.Value, 'f', -1, 64)
		message := fmt.Sprintf("Chào bạn, mã giảm giá của bạn là %s có giá trị giảm %s đ. Thưởng thức ưu đãi của bạn!\n", discount.DiscountCode, valueString)
		messages = append(messages, message)
	}

	allMessages := strings.Join(messages, " ")

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(allMessages),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func HandleRabbitMQMessages() {
	conn, err := amqp.Dial("amqps://ithkfqls:w4d4HaTpn_cDiod1r9BOT-CSZON1gYVF@octopus.rmq3.cloudamqp.com/ithkfqls")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"discount_notification_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	for msg := range msgs {
		emailContent := string(msg.Body)
		fmt.Printf("Nhận thông điệp: %s\n", emailContent)

		users, err := dbms.GetUsersByRole("user")
		if err != nil {
			log.Printf("Error getting users: %v\n", err)
			continue
		}

		for _, user := range users {
			err = SendOrderStatusUpdateEmail(user.Email, emailContent)
			if err != nil {
				log.Printf("Error sending email: %v\n", err)
			}
		}

		msg.Ack(false)
	}
}

func SendOrderStatusUpdateEmail(email, emailContent string) error {
	emailAddress := "pau30012002@gmail.com"
	emailPassword := "pljf fqgx yycq ynhq"
	subject := "Order Status Update"
	body := fmt.Sprintf(emailContent)

	m := gomail.NewMessage()
	m.SetHeader("From", emailAddress)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, emailAddress, emailPassword)

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}