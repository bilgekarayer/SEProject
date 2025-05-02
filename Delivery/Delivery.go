package Delivery

import (
	"errors"
	"time"
)

// --- Model Tanımları ---

type DeliveryPerson struct {
	ID          string
	Name        string
	PhoneNumber string
	Status      string // "available", "busy", "offline"
	Location    Location
	ActiveOrder *Order // Current order being delivered, nil if none
}

type Location struct {
	Latitude  float64
	Longitude float64
}

// --- Bağımlı Order modeli ve servisi ---

type Order struct {
	ID               string
	CustomerID       string
	RestaurantID     string
	Status           string // "preparing", "assigned_to_delivery", "out_for_delivery", "delivered"
	DeliveryPersonID string
	AssignedAt       time.Time
	PickedUpAt       time.Time
	DeliveredAt      time.Time
}

type OrderService interface {
	GetOrderByID(orderID string) (*Order, error)
	UpdateOrder(order *Order) error
	AddOrder(order *Order) error
}

// --- Arayüz ---

type DeliveryPersonService interface {
	AssignOrder(deliveryPersonID string, orderID string) error
	MarkOrderAsPickedUp(orderID string) error
	MarkOrderAsDelivered(orderID string) error
	GetAvailableDeliveryPersons() ([]*DeliveryPerson, error)
	UpdateLocation(deliveryPersonID string, location Location) error
	GetDeliveryPersonByID(id string) (*DeliveryPerson, error)
	AddDeliveryPerson(dp *DeliveryPerson) error
}

// --- In-memory servis implementasyonu ---

type InMemoryDeliveryPersonService struct {
	deliveryPersons map[string]*DeliveryPerson
	orderService    OrderService
}

func NewInMemoryDeliveryPersonService(orderService OrderService) *InMemoryDeliveryPersonService {
	return &InMemoryDeliveryPersonService{
		deliveryPersons: make(map[string]*DeliveryPerson),
		orderService:    orderService,
	}
}

func (s *InMemoryDeliveryPersonService) AssignOrder(deliveryPersonID string, orderID string) error {
	dp, err := s.GetDeliveryPersonByID(deliveryPersonID)
	if err != nil {
		return err
	}
	if dp.Status != "available" {
		return errors.New("delivery person is not available")
	}

	order, err := s.orderService.GetOrderByID(orderID)
	if err != nil {
		return err
	}
	if order.Status != "preparing" {
		return errors.New("order is not in a state that can be assigned")
	}

	dp.ActiveOrder = order
	dp.Status = "busy"

	order.Status = "assigned_to_delivery"
	order.DeliveryPersonID = deliveryPersonID
	order.AssignedAt = time.Now()
	return s.orderService.UpdateOrder(order)
}

func (s *InMemoryDeliveryPersonService) MarkOrderAsPickedUp(orderID string) error {
	order, err := s.orderService.GetOrderByID(orderID)
	if err != nil {
		return err
	}
	if order.Status != "assigned_to_delivery" {
		return errors.New("order is not in a state that can be picked up")
	}

	dp, err := s.GetDeliveryPersonByID(order.DeliveryPersonID)
	if err != nil {
		return err
	}

	order.Status = "out_for_delivery"
	order.PickedUpAt = time.Now()
	dp.ActiveOrder = order
	return s.orderService.UpdateOrder(order)
}

func (s *InMemoryDeliveryPersonService) MarkOrderAsDelivered(orderID string) error {
	order, err := s.orderService.GetOrderByID(orderID)
	if err != nil {
		return err
	}
	if order.Status != "out_for_delivery" {
		return errors.New("order is not out for delivery")
	}

	dp, err := s.GetDeliveryPersonByID(order.DeliveryPersonID)
	if err != nil {
		return err
	}

	order.Status = "delivered"
	order.DeliveredAt = time.Now()
	dp.Status = "available"
	dp.ActiveOrder = nil
	return s.orderService.UpdateOrder(order)
}

func (s *InMemoryDeliveryPersonService) GetAvailableDeliveryPersons() ([]*DeliveryPerson, error) {
	var available []*DeliveryPerson
	for _, dp := range s.deliveryPersons {
		if dp.Status == "available" {
			available = append(available, dp)
		}
	}
	return available, nil
}

func (s *InMemoryDeliveryPersonService) UpdateLocation(deliveryPersonID string, location Location) error {
	dp, err := s.GetDeliveryPersonByID(deliveryPersonID)
	if err != nil {
		return err
	}
	dp.Location = location
	return nil
}

func (s *InMemoryDeliveryPersonService) GetDeliveryPersonByID(id string) (*DeliveryPerson, error) {
	dp, ok := s.deliveryPersons[id]
	if !ok {
		return nil, errors.New("delivery person not found")
	}
	return dp, nil
}

func (s *InMemoryDeliveryPersonService) AddDeliveryPerson(dp *DeliveryPerson) error {
	if _, exists := s.deliveryPersons[dp.ID]; exists {
		return errors.New("delivery person already exists")
	}
	s.deliveryPersons[dp.ID] = dp
	return nil
}
