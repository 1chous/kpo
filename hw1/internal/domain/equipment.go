package domain

import "fmt"

type BaseEquipment struct {
	inventoryNumber string
	name            string
}

func newBaseEquipment(inventoryNumber, name string) (BaseEquipment, error) {
	if inventoryNumber == "" {
		return BaseEquipment{}, ErrEmptyInventoryNumber
	}
	if name == "" {
		return BaseEquipment{}, ErrEmptyName
	}
	return BaseEquipment{inventoryNumber: inventoryNumber, name: name}, nil
}

func (b BaseEquipment) InventoryNumber() string { return b.inventoryNumber }
func (b BaseEquipment) Name() string            { return b.name }

type Helmet struct {
	BaseEquipment
}

func NewHelmet(inventoryNumber string) (Helmet, error) {
	base, err := newBaseEquipment(inventoryNumber, "Шлем")
	if err != nil {
		return Helmet{}, err
	}
	return Helmet{BaseEquipment: base}, nil
}

type DockStation struct {
	BaseEquipment
}

func NewDockStation(inventoryNumber string) (DockStation, error) {
	base, err := newBaseEquipment(inventoryNumber, "Док-станция")
	if err != nil {
		return DockStation{}, err
	}
	return DockStation{BaseEquipment: base}, nil
}

type ChargingCabinet struct {
	BaseEquipment
	dailyEnergyKWh float64
}

func NewChargingCabinet(inventoryNumber string, dailyEnergyKWh float64) (ChargingCabinet, error) {
	if dailyEnergyKWh < 0 {
		return ChargingCabinet{}, fmt.Errorf("%w: получено %.2f", ErrNegativeEnergy, dailyEnergyKWh)
	}
	base, err := newBaseEquipment(inventoryNumber, "Зарядный шкаф")
	if err != nil {
		return ChargingCabinet{}, err
	}
	return ChargingCabinet{BaseEquipment: base, dailyEnergyKWh: dailyEnergyKWh}, nil
}

func (c ChargingCabinet) DailyEnergyKWh() float64 { return c.dailyEnergyKWh }
