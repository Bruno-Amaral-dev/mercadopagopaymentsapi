package mercadopago

// MercadoPagoCardTokenRequest represents the request to generate a card token
type MercadoPagoCardTokenRequest struct {
	CardNumber     string `json:"card_number"`
	Cardholder     struct {
		Name     string `json:"name"`
		Identification struct {
			Type   string `json:"type"`
			Number string `json:"number"`
		} `json:"identification"`
	} `json:"cardholder"`
	SecurityCode   string `json:"security_code"`
	ExpirationMonth string `json:"expiration_month"`
	ExpirationYear  string `json:"expiration_year"`
}

// MercadoPagoCardTokenResponse represents the response from card token generation
type MercadoPagoCardTokenResponse struct {
	ID           string `json:"id"`
	FirstSixDigits string `json:"first_six_digits"`
	LastFourDigits string `json:"last_four_digits"`
	ExpirationMonth int    `json:"expiration_month"`
	ExpirationYear  int    `json:"expiration_year"`
	CreationDate    string `json:"creation_date"`
}

// MercadoPagoPaymentRequest represents the request to create a payment
type MercadoPagoPaymentRequest struct {
	TransactionAmount float64 `json:"transaction_amount"`
	Installments      int     `json:"installments"`
	PaymentMethodID   string  `json:"payment_method_id"`
	Token             string  `json:"token"`
	Payer             struct {
		Email string `json:"email"`
		Type  string `json:"type"`
		Identification struct {
			Type   string `json:"type"`
			Number string `json:"number"`
		} `json:"identification"`
	} `json:"payer"`
}

// MercadoPagoPaymentResponse represents the response from payment creation
type MercadoPagoPaymentResponse struct {
	ID                 int     `json:"id"`
	DateCreated        string  `json:"date_created"`
	DateApproved       *string `json:"date_approved"`
	DateLastUpdated    string  `json:"date_last_updated"`
	MoneyReleaseDate   *string `json:"money_release_date"`
	OperationType      string  `json:"operation_type"`
	PaymentMethodID    string  `json:"payment_method_id"`
	PaymentTypeID      string  `json:"payment_type_id"`
	Status             string  `json:"status"`
	StatusDetail       string  `json:"status_detail"`
	CurrencyID         string  `json:"currency_id"`
	TransactionAmount  float64 `json:"transaction_amount"`
	Installments       int     `json:"installments"`
	InstallmentsAmount float64 `json:"installments_amount"`
	DeductionSchema    *string `json:"deduction_schema"`
	CollectorID        int     `json:"collector_id"`
	Payer              struct {
		ID           string `json:"id"`
		Type         string `json:"type"`
		Email        string `json:"email"`
		Identification struct {
			Type   string `json:"type"`
			Number string `json:"number"`
		} `json:"identification"`
	} `json:"payer"`
	ExternalReference string `json:"external_reference"`
	DateOfExpiration   *string `json:"date_of_expiration"`
	DifferentialPricing *struct {
		ID int `json:"id"`
	} `json:"differential_pricing"`
	ApplicationFee   *float64 `json:"application_fee"`
	Order            *struct {
		ID int `json:"id"`
		Type string `json:"type"`
	} `json:"order"`
	CorrelationID    *string `json:"correlation_id"`
	TransactionDetails struct {
		FinancialInstitution   *string `json:"financial_institution"`
		PaymentMethodReferenceID *string `json:"payment_method_reference_id"`
		NetReceivedAmount       float64 `json:"net_received_amount"`
		TotalPaidAmount         float64 `json:"total_paid_amount"`
		OverpaidAmount          float64 `json:"overpaid_amount"`
		ExternalResourceURL     *string `json:"external_resource_url"`
		InstallmentAmount       float64 `json:"installment_amount"`
		FinancialFee            float64 `json:"financial_fee"`
	} `json:"transaction_details"`
	FeeDetails []struct {
		Type string `json:"type"`
		Amount float64 `json:"amount"`
	} `json:"fee_details"`
	Card struct {
		ID           string `json:"id"`
		FirstSixDigits string `json:"first_six_digits"`
		LastFourDigits string `json:"last_four_digits"`
		ExpirationMonth int `json:"expiration_month"`
		ExpirationYear  int `json:"expiration_year"`
		Cardholder struct {
			Name string `json:"name"`
			Identification struct {
				Type   string `json:"type"`
				Number string `json:"number"`
			} `json:"identification"`
		} `json:"cardholder"`
	} `json:"card"`
}
