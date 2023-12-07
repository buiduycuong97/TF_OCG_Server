package order_handle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	paypalClientID     string
	paypalClientSecret string
	baseURL            string
	port               string
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type CreateOrderResponse struct {
	ID string `json:"id"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	paypalClientID = os.Getenv("PAYPAL_CLIENT_ID")
	paypalClientSecret = os.Getenv("PAYPAL_CLIENT_SECRET")
	baseURL = os.Getenv("PAYPAL_BASE_URL")
	port = os.Getenv("PAYPAL_PORT")
}

func generateAccessToken() (string, error) {
	url := fmt.Sprintf("%s/v1/oauth2/token", baseURL)
	payload := []byte("grant_type=client_credentials")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(paypalClientID, paypalClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp AccessTokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func CreateOrder(grandTotal float64) (CreateOrderResponse, error) {
	accessToken, err := generateAccessToken()
	if err != nil {
		return CreateOrderResponse{}, err
	}

	url := fmt.Sprintf("%s/v2/checkout/orders", baseURL)

	orderAmount := grandTotal
	orderCurrency := "USD"
	orderIntent := "CAPTURE"

	payload := map[string]interface{}{
		"intent": orderIntent,
		"purchase_units": []map[string]interface{}{
			{
				"amount": map[string]interface{}{
					"currency_code": orderCurrency,
					"value":         fmt.Sprintf("%.2f", orderAmount/23500),
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return CreateOrderResponse{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return CreateOrderResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return CreateOrderResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return CreateOrderResponse{}, err
	}

	var orderResp CreateOrderResponse
	err = json.Unmarshal(body, &orderResp)
	if err != nil {
		return CreateOrderResponse{}, err
	}

	return orderResp, nil
}

func captureOrder(orderID string) (map[string]interface{}, error) {
	accessToken, err := generateAccessToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/v2/checkout/orders/%s/capture", baseURL, orderID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var captureResp map[string]interface{}
	err = json.Unmarshal(body, &captureResp)
	if err != nil {
		return nil, err
	}

	return captureResp, nil
}

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var orderData map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&orderData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	grandTotal, ok := orderData["grandTotal"].(float64)
	if !ok {
		http.Error(w, "Invalid grandTotal value", http.StatusBadRequest)
		return
	}

	orderResp, err := CreateOrder(grandTotal)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(orderResp)
}

func CaptureOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["orderID"]

	captureResp, err := captureOrder(orderID)
	if err != nil {
		http.Error(w, "Failed to capture order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(captureResp)
}
