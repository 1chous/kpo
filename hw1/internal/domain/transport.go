package domain

import "fmt"

type BaseTransport struct {
	inventoryNumber    string
	name               string
	noviceFriendliness int
}

func newBaseTransport(inventoryNumber, name string, noviceFriendliness int) (BaseTransport, error) {
	if inventoryNumber == "" {
		return BaseTransport{}, ErrEmptyInventoryNumber
	}
	if name == "" {
		return BaseTransport{}, ErrEmptyName
	}
	if noviceFriendliness < 1 || noviceFriendliness > 10 {
		return BaseTransport{}, fmt.Errorf("%w: получено %d", ErrNoviceOutOfRange, noviceFriendliness)
	}
	return BaseTransport{
		inventoryNumber:    inventoryNumber,
		name:               name,
		noviceFriendliness: noviceFriendliness,
	}, nil
}

func (b BaseTransport) InventoryNumber() string { return b.inventoryNumber }
func (b BaseTransport) Name() string            { return b.name }
func (b BaseTransport) NoviceFriendliness() int { return b.noviceFriendliness }

type electricTransport struct {
	BaseTransport
	dailyEnergyKWh float64
}

func newElectricTransport(inventoryNumber, name string, noviceFriendliness int, dailyEnergyKWh float64) (electricTransport, error) {
	if dailyEnergyKWh < 0 {
		return electricTransport{}, fmt.Errorf("%w: получено %.2f", ErrNegativeEnergy, dailyEnergyKWh)
	}
	base, err := newBaseTransport(inventoryNumber, name, noviceFriendliness)
	if err != nil {
		return electricTransport{}, err
	}
	return electricTransport{BaseTransport: base, dailyEnergyKWh: dailyEnergyKWh}, nil
}

func (e electricTransport) DailyEnergyKWh() float64 { return e.dailyEnergyKWh }

type Scooter struct {
	electricTransport
}

func NewScooter(inventoryNumber string, noviceFriendliness int, dailyEnergyKWh float64) (Scooter, error) {
	t, err := newElectricTransport(inventoryNumber, "Самокат", noviceFriendliness, dailyEnergyKWh)
	if err != nil {
		return Scooter{}, err
	}
	return Scooter{electricTransport: t}, nil
}

type EBike struct {
	electricTransport
}

func NewEBike(inventoryNumber string, noviceFriendliness int, dailyEnergyKWh float64) (EBike, error) {
	t, err := newElectricTransport(inventoryNumber, "Электровелосипед", noviceFriendliness, dailyEnergyKWh)
	if err != nil {
		return EBike{}, err
	}
	return EBike{electricTransport: t}, nil
}

type Bicycle struct {
	BaseTransport
}

func NewBicycle(inventoryNumber string, noviceFriendliness int) (Bicycle, error) {
	base, err := newBaseTransport(inventoryNumber, "Велосипед", noviceFriendliness)
	if err != nil {
		return Bicycle{}, err
	}
	return Bicycle{BaseTransport: base}, nil
}

type Monowheel struct {
	electricTransport
}

func NewMonowheel(inventoryNumber string, noviceFriendliness int, dailyEnergyKWh float64) (Monowheel, error) {
	t, err := newElectricTransport(inventoryNumber, "Моноколесо", noviceFriendliness, dailyEnergyKWh)
	if err != nil {
		return Monowheel{}, err
	}
	return Monowheel{electricTransport: t}, nil
}
