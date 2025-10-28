package dto

import "RWB_L0/internal/domain"

// ToDomain - конвертирует CreateOrderInput в domain.Order
func (input *CreateOrderInput) ToDomain() (*domain.Order, error) {
	// Создаём заказ через конструктор Domain
	order, err := domain.NewOrder(input.OrderUID, input.TrackNumber, input.Entry)
	if err != nil {
		return nil, err
	}

	// Delivery
	delivery, err := domain.NewDelivery(
		input.Delivery.Name,
		input.Delivery.Phone,
		input.Delivery.Zip,
		input.Delivery.City,
		input.Delivery.Address,
		input.Delivery.Region,
		input.Delivery.Email,
	)
	if err != nil {
		return nil, err
	}
	order.Delivery = *delivery

	// Payment
	payment, err := domain.NewPayment(
		input.Payment.Transaction,
		input.Payment.RequestID,
		input.Payment.Currency,
		input.Payment.Provider,
		input.Payment.Bank,
		input.Payment.Amount,
		input.Payment.PaymentDt,
		input.Payment.DeliveryCost,
		input.Payment.GoodsTotal,
		input.Payment.CustomFee,
	)
	if err != nil {
		return nil, err
	}
	order.Payment = *payment

	// Items
	for _, itemInput := range input.Items {
		item, err := domain.NewItem(
			itemInput.ChrtID,
			itemInput.TrackNumber,
			itemInput.Name,
			itemInput.Rid,
			itemInput.Size,
			itemInput.Brand,
			itemInput.Price,
			itemInput.Sale,
			itemInput.TotalPrice,
			itemInput.NmID,
			itemInput.Status,
		)
		if err != nil {
			return nil, err
		}
		if err := order.AddItem(*item); err != nil {
			return nil, err
		}
	}

	// Остальные поля
	order.Locale = input.Locale
	order.InternalSignature = input.InternalSignature
	order.CustomerID = input.CustomerID
	order.DeliveryService = input.DeliveryService
	order.Shardkey = input.Shardkey
	order.SmID = input.SmID
	order.DateCreated = input.DateCreated
	order.OofShard = input.OofShard

	return order, nil
}

// FromDomain - конвертирует domain.Order в OrderOutput
func FromDomain(order *domain.Order) *OrderOutput {
	output := &OrderOutput{
		OrderUID:          order.OrderUID,
		TrackNumber:       order.TrackNumber,
		Entry:             order.Entry,
		Locale:            order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerID:        order.CustomerID,
		DeliveryService:   order.DeliveryService,
		Shardkey:          order.Shardkey,
		SmID:              order.SmID,
		DateCreated:       order.DateCreated,
		OofShard:          order.OofShard,
	}

	// Delivery
	output.Delivery = DeliveryOutput{
		Name:    order.Delivery.Name,
		Phone:   order.Delivery.Phone,
		Zip:     order.Delivery.Zip,
		City:    order.Delivery.City,
		Address: order.Delivery.Address,
		Region:  order.Delivery.Region,
		Email:   order.Delivery.Email,
	}

	// Payment
	output.Payment = PaymentOutput{
		Transaction:  order.Payment.Transaction,
		RequestID:    order.Payment.RequestID,
		Currency:     order.Payment.Currency,
		Provider:     order.Payment.Provider,
		Amount:       order.Payment.Amount,
		PaymentDt:    order.Payment.PaymentDt,
		Bank:         order.Payment.Bank,
		DeliveryCost: order.Payment.DeliveryCost,
		GoodsTotal:   order.Payment.GoodsTotal,
		CustomFee:    order.Payment.CustomFee,
	}

	// Items
	output.Items = make([]ItemOutput, len(order.Items))
	for i, item := range order.Items {
		output.Items[i] = ItemOutput{
			ChrtID:      item.ChrtID,
			TrackNumber: item.TrackNumber,
			Price:       item.Price,
			Rid:         item.Rid,
			Name:        item.Name,
			Sale:        item.Sale,
			Size:        item.Size,
			TotalPrice:  item.TotalPrice,
			NmID:        item.NmID,
			Brand:       item.Brand,
			Status:      item.Status,
		}
	}
	return output
}
