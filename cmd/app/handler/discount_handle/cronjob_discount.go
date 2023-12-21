package discount_handle

import (
	"fmt"
	"time"
)

//func ScheduleDiscountCodeGeneration() {
//	cronJob := cron.New()
//
//	cronJob.AddFunc("0 1 * * *", func() {
//		month := time.Now().Month()
//		if int(month) == time.Now().Day() {
//			err := GenerateAndSaveDiscountCodes()
//			if err != nil {
//				fmt.Printf("Error generating and saving discount codes: %v\n", err)
//			}
//		}
//	})
//
//	cronJob.Start()
//}

func ScheduleDiscountCodeGeneration() {
	// Tạo một timer đếm ngược 30 giây
	timer := time.NewTimer(5 * time.Second)

	// Chờ cho timer kết thúc (đã đếm ngược xong)
	<-timer.C

	// Thực hiện công việc
	err := GenerateAndSaveDiscountCodes()
	if err != nil {
		fmt.Printf("Error generating and saving discount codes: %v\n", err)
	}
}
