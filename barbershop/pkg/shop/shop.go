package shop

import (
	"time"

	"go.uber.org/zap"
)

const (
	barberTag   string = "barber"
	customerTag string = "customer"
)

// Shop is an abstraction of a real Barber Shop.
//
// It has a buffered channel for holding customers
// who enter the shop while the barbers are busy.
// The capacity of this channel is determined by
// the shop's capacity attribute.
//
// It also has a done channel which help to
// determine whether all barbers have finished
// any jobs they are currently doing. This helps
// to ensure that the shop can not be close
// while barbers are busy cutting hair.
//
// A Shop uses its logger to display status information.
type Shop struct {
	capacity        int
	cutDuration     time.Duration
	numberOfBarbers int
	clients         chan string
	done            chan bool
	isOpen          bool
	logger          *zap.Logger
}

// ConfigOptions is a type that make it
// convinient to create values of the Shop type.
// By holding all configuration options on this type,
// the number of parameters to be passed to the NewShop
// function is reduced to 2 instead of 6 which wouldn't
// be an ideal or idiomatic style in Go.
type ConfigOptions struct {
	SeatingCapacity int
	DurationPerCut  int
	Clients         chan string
	Done            chan bool
}

// NewShop is a utility function for creating shops.
//
// It returns a pointer to a newly initialzied shop.
func NewShop(cfg *ConfigOptions, logger *zap.Logger) *Shop {
	return &Shop{
		capacity:        cfg.SeatingCapacity,
		cutDuration:     time.Duration(cfg.DurationPerCut) * time.Millisecond,
		numberOfBarbers: 0,
		clients:         cfg.Clients,
		done:            cfg.Done,
		isOpen:          true,
		logger:          logger,
	}
}

// AddBarber registers a barber in the shop who
// has to participate in cutting hair when customers
// are available.
//
// AddBarber starts a goroutine that immediately kicks the
// barber into the cycle of either taking a nap, cutting hair,
// or going home when the shop is to be closed for the day.
func (s *Shop) AddBarber(nameOfBarber string) {
	s.numberOfBarbers += 1

	go func() {
		isSleeping := false
		for {
			if len(s.clients) == 0 {
				isSleeping = true
				s.logger.Info("barber is sleeping since there are no customers", zap.String(barberTag, nameOfBarber))
			}

			client, opened := <-s.clients
			if opened {
				if isSleeping {
					s.logger.Info("new customer awakens barber", zap.String(barberTag, nameOfBarber), zap.String(customerTag, client))
					isSleeping = false
				}

				s.assignCustomerToBarber(client, nameOfBarber)
			} else {
				s.notifyBarberHasClosed(nameOfBarber)
				return
			}
		}
	}()
}

// AddCustomer similates the situation where a new customer
// enters the shop.
// If all barbers are busy and there's a waiting space, the
// customer waits patiently in the waiting area until attended to.
//
// If all barbers are busy and the waiting area is full, then
// the customer walks away without getting a haircut.
//
// Also, if a customer comes and the shop is closed, then
// the customer returns without a haircut because all barbers
// have closed and gone home.
func (s *Shop) AddCustomer(nameOfCustomer string) {
	s.logger.Info("a new customer has arrived in the shop", zap.String(customerTag, nameOfCustomer))
	if s.isOpen {
		select {
		case s.clients <- nameOfCustomer:
			s.logger.Info("new customer takes a seat in the waiting area", zap.String(customerTag, nameOfCustomer))
		default:
			s.logger.Info("customer leaves because waiting area is full", zap.String(customerTag, nameOfCustomer))
		}
	} else {
		s.logger.Info("shop is closed so customer goes away without a cut", zap.String(customerTag, nameOfCustomer))
	}
}

func (s *Shop) assignCustomerToBarber(customer, barber string) {
	s.logger.Info("barber begins cutting customer's hair", zap.String(barberTag, barber), zap.String(customerTag, customer))
	time.Sleep(s.cutDuration)
	s.logger.Info("barber has finished cutting customer's hair", zap.String(barberTag, barber), zap.String(customerTag, customer))
}

func (s *Shop) notifyBarberHasClosed(barber string) {
	s.logger.Info("barber has closed and is going home now", zap.String(barberTag, barber))
	s.done <- true
}

// Close closes the shop but ensures that
// all customers who are being attended to by
// barbers really finish getting their haircuts
// before closing the shop.
//
// While awaiting busy barbers to finish, no new
// customers are allowed entry into the shop.
func (s *Shop) Close() {
	s.logger.Info("shop is being closed for today")
	close(s.clients)
	s.isOpen = false

	for i := 1; i < s.numberOfBarbers; i++ {
		<-s.done
	}

	close(s.done)
	s.logger.Info("shop has closed for the day, see you tomorrow!!!")
}
