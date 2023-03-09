package entity

type ItemDetailsContent struct {
	ID       string `json:"id"`       // this is for gopay
	Price    int    `json:"price"`    // this is for gopay
	Quantity int    `json:"quantity"` // this is for gopay
	Name     string `json:"name"`     // this is for gopay
}

type TransactionDetailsContent struct {
	Order_ID     string `json:"order_id"`     // this is for gopay
	Gross_Amount int    `json:"gross_amount"` // this is for gopay
}
type Gopay struct {
	Enable_callback bool   `json:"enable_callback"`
	Callback_url    string `json:"callback_url"`
}

type CustomerDetails struct {
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}

type Payload struct {
	Customer_details    CustomerDetails           `json:"customer_details"`
	Gopay               Gopay                     `json:"gopay"`
	Item_details        []ItemDetailsContent      `json:"item_details"`
	Payment_type        string                    `json:"payment_type"`
	Transaction_details TransactionDetailsContent `json:"transaction_details"`
}

type ValidatePayment struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymenType        string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	ExpiryTime        string `json:"expiry_time"`
	Currency          string `json:"currency"`
}
