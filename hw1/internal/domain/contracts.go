package domain

const MinNoviceFriendliness = 6

type Inventoriable interface {
	InventoryNumber() string
	Name() string
}

type EnergyConsumer interface {
	DailyEnergyKWh() float64
}

type Rentable interface {
	Inventoriable
	NoviceFriendliness() int
}

func IsNoviceFriendly(r Rentable) bool {
	return r.NoviceFriendliness() >= MinNoviceFriendliness
}
