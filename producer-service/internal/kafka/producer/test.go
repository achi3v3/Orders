package producer

import (
	"fmt"
	"math/rand"
	"producer-service/internal/models"
	"time"
)

func createRandomOrder() models.Order {
	orderUID := fmt.Sprintf("order-%s", randomString(8))
	trackNumber := fmt.Sprintf("TRACK%s", randomString(6))
	return models.Order{
		OrderUID:    orderUID,
		TrackNumber: trackNumber,
		Entry:       "ENT" + randomString(3),
		Delivery: models.Delivery{
			Name:    "Name " + randomString(5),
			Phone:   "+1" + randomDigits(10),
			Zip:     randomDigits(6),
			City:    "City" + randomString(4),
			Address: "Address " + randomString(6),
			Region:  "Region" + randomString(3),
			Email:   "email" + randomString(4) + "@test.com",
		},
		Payment: models.Payment{
			Transaction:  orderUID,
			RequestID:    "",
			Currency:     "USD",
			Provider:     "provider" + randomString(3),
			Amount:       rand.Intn(5000),
			PaymentDT:    int64(time.Now().Unix()),
			Bank:         "bank" + randomString(3),
			DeliveryCost: rand.Intn(1000),
			GoodsTotal:   rand.Intn(3000),
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      rand.Intn(1000000),
				TrackNumber: trackNumber,
				Price:       rand.Intn(1000),
				RID:         "rid-" + randomString(8),
				Name:        "Item " + randomString(5),
				Sale:        rand.Intn(100),
				Size:        "M",
				TotalPrice:  rand.Intn(1000),
				NmID:        rand.Intn(100000),
				Brand:       "Brand" + randomString(4),
				Status:      200,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "customer" + randomString(4),
		DeliveryService:   "service" + randomString(3),
		ShardKey:          "1",
		SmID:              rand.Intn(100),
		DateCreated:       time.Now().Format(time.RFC3339),
		OofShard:          "1",
	}
}
func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
func randomDigits(n int) string {
	digits := []rune("0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = digits[rand.Intn(len(digits))]
	}
	return string(s)
}
