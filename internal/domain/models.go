package domain

import "time"

// ========================================
// Delivery - информация о доставке
// ========================================

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

func NewDelivery(name, phone, zip, city, address, region, email string) (*Delivery, error) {
	if name == "" {
		return nil, ErrEmptyDeliveryName
	}
	if phone == "" {
		return nil, ErrEmptyDeliveryPhone
	}

	return &Delivery{
		Name:    name,
		Phone:   phone,
		Zip:     zip,
		City:    city,
		Address: address,
		Region:  region,
		Email:   email,
	}, nil
}

func (d *Delivery) Validate() error {
	if d.Name == "" {
		return ErrEmptyDeliveryName
	}
	if d.Phone == "" {
		return ErrEmptyDeliveryPhone
	}
	return nil
}

// ========================================
// Payment - информация о платеже
// ========================================

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

func NewPayment(transaction, requestID, currency, provider, bank string,
	amount int, paymentDt int64, deliveryCost, goodsTotal, customFee int) (*Payment, error) {

	if transaction == "" {
		return nil, ErrEmptyPaymentTransaction
	}
	if amount <= 0 {
		return nil, ErrInvalidPaymentAmount
	}

	return &Payment{
		Transaction:  transaction,
		RequestID:    requestID,
		Currency:     currency,
		Provider:     provider,
		Amount:       amount,
		PaymentDt:    paymentDt,
		Bank:         bank,
		DeliveryCost: deliveryCost,
		GoodsTotal:   goodsTotal,
		CustomFee:    customFee,
	}, nil
}

func (p *Payment) Validate() error {
	if p.Transaction == "" {
		return ErrEmptyPaymentTransaction
	}
	if p.Amount <= 0 {
		return ErrInvalidPaymentAmount
	}
	return nil
}

// ========================================
// Item - товар в заказе
// ========================================

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func NewItem(chrtID int, trackNumber, name, rid, size, brand string,
	price, sale, totalPrice, nmID, status int) (*Item, error) {

	if name == "" {
		return nil, ErrEmptyItemName
	}
	if price < 0 {
		return nil, ErrInvalidItemPrice
	}

	return &Item{
		ChrtID:      chrtID,
		TrackNumber: trackNumber,
		Price:       price,
		Rid:         rid,
		Name:        name,
		Sale:        sale,
		Size:        size,
		TotalPrice:  totalPrice,
		NmID:        nmID,
		Brand:       brand,
		Status:      status,
	}, nil
}

func (i *Item) Validate() error {
	if i.Name == "" {
		return ErrEmptyItemName
	}
	if i.Price < 0 {
		return ErrInvalidItemPrice
	}
	return nil
}

// ========================================
// Order - главная доменная модель заказа
// ========================================

type Order struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

func NewOrder(orderUID, trackNumber, entry string) (*Order, error) {
	if orderUID == "" {
		return nil, ErrEmptyOrderUID
	}
	if trackNumber == "" {
		return nil, ErrEmptyTrackNumber
	}

	return &Order{
		OrderUID:    orderUID,
		TrackNumber: trackNumber,
		Entry:       entry,
		DateCreated: time.Now(),
		Items:       make([]Item, 0),
	}, nil
}

func (o *Order) Validate() error {
	if o.OrderUID == "" {
		return ErrEmptyOrderUID
	}
	if o.TrackNumber == "" {
		return ErrEmptyTrackNumber
	}

	if err := o.Delivery.Validate(); err != nil {
		return err
	}
	if err := o.Payment.Validate(); err != nil {
		return err
	}

	for _, item := range o.Items {
		if err := item.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (o *Order) AddItem(item Item) error {
	if err := item.Validate(); err != nil {
		return err
	}
	o.Items = append(o.Items, item)
	return nil
}

func (o *Order) GetTotal() int {
	return o.Payment.Amount
}

func (o *Order) GetItemsCount() int {
	return len(o.Items)
}
