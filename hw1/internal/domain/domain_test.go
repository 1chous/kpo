package domain_test

import (
	"errors"
	"testing"

	"github.com/1chous/kpo/hw1/internal/domain"
)

func TestConstructorsRejectInvalidData(t *testing.T) {
	tests := []struct {
		name    string
		build   func() error
		wantErr error
	}{
		{
			name:    "самокат без инвентарного номера",
			build:   func() error { _, err := domain.NewScooter("", 5, 0.5); return err },
			wantErr: domain.ErrEmptyInventoryNumber,
		},
		{
			name:    "простота ниже допустимой",
			build:   func() error { _, err := domain.NewScooter("SC-1", 0, 0.5); return err },
			wantErr: domain.ErrNoviceOutOfRange,
		},
		{
			name:    "простота выше допустимой",
			build:   func() error { _, err := domain.NewScooter("SC-1", 11, 0.5); return err },
			wantErr: domain.ErrNoviceOutOfRange,
		},
		{
			name:    "отрицательный расход энергии",
			build:   func() error { _, err := domain.NewEBike("EB-1", 5, -0.1); return err },
			wantErr: domain.ErrNegativeEnergy,
		},
		{
			name:    "велосипед без инвентарного номера",
			build:   func() error { _, err := domain.NewBicycle("", 5); return err },
			wantErr: domain.ErrEmptyInventoryNumber,
		},
		{
			name:    "шлем без инвентарного номера",
			build:   func() error { _, err := domain.NewHelmet(""); return err },
			wantErr: domain.ErrEmptyInventoryNumber,
		},
		{
			name:    "док-станция без инвентарного номера",
			build:   func() error { _, err := domain.NewDockStation(""); return err },
			wantErr: domain.ErrEmptyInventoryNumber,
		},
		{
			name:    "зарядный шкаф с отрицательным расходом",
			build:   func() error { _, err := domain.NewChargingCabinet("CH-1", -5); return err },
			wantErr: domain.ErrNegativeEnergy,
		},
		{
			name:    "зарядный шкаф без инвентарного номера",
			build:   func() error { _, err := domain.NewChargingCabinet("", 1.0); return err },
			wantErr: domain.ErrEmptyInventoryNumber,
		},
		{
			name:    "моноколесо с отрицательным расходом",
			build:   func() error { _, err := domain.NewMonowheel("MW-1", 2, -0.5); return err },
			wantErr: domain.ErrNegativeEnergy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.build()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("получено %v, ожидалось %v", err, tt.wantErr)
			}
		})
	}
}

func TestConstructorsAcceptValidData(t *testing.T) {
	scooter, err := domain.NewScooter("SC-1", 8, 0.45)
	if err != nil {
		t.Fatalf("самокат не создан: %v", err)
	}
	if scooter.InventoryNumber() != "SC-1" {
		t.Errorf("инвентарный номер: получено %q", scooter.InventoryNumber())
	}
	if scooter.Name() != "Самокат" {
		t.Errorf("наименование: получено %q", scooter.Name())
	}
	if scooter.NoviceFriendliness() != 8 {
		t.Errorf("простота: получено %d", scooter.NoviceFriendliness())
	}
	if scooter.DailyEnergyKWh() != 0.45 {
		t.Errorf("расход энергии: получено %v", scooter.DailyEnergyKWh())
	}

	dock, err := domain.NewDockStation("DS-1")
	if err != nil {
		t.Fatalf("док-станция не создана: %v", err)
	}
	if dock.Name() != "Док-станция" || dock.InventoryNumber() != "DS-1" {
		t.Errorf("док-станция: получено %q %q", dock.Name(), dock.InventoryNumber())
	}

	bike, err := domain.NewEBike("EB-1", 6, 0.9)
	if err != nil {
		t.Fatalf("электровелосипед не создан: %v", err)
	}
	if bike.Name() != "Электровелосипед" {
		t.Errorf("электровелосипед: получено %q", bike.Name())
	}
}

func TestBicycleIsNotEnergyConsumer(t *testing.T) {
	bicycle, err := domain.NewBicycle("BC-1", 9)
	if err != nil {
		t.Fatalf("велосипед не создан: %v", err)
	}

	var item domain.Inventoriable = bicycle
	if _, ok := item.(domain.EnergyConsumer); ok {
		t.Fatal("велосипед не должен реализовывать EnergyConsumer")
	}
}

func TestChargingCabinetIsConsumerButNotRentable(t *testing.T) {
	cabinet, err := domain.NewChargingCabinet("CH-1", 1.2)
	if err != nil {
		t.Fatalf("шкаф не создан: %v", err)
	}

	var item domain.Inventoriable = cabinet
	if _, ok := item.(domain.EnergyConsumer); !ok {
		t.Error("зарядный шкаф должен реализовывать EnergyConsumer")
	}
	if _, ok := item.(domain.Rentable); ok {
		t.Error("зарядный шкаф не должен реализовывать Rentable")
	}
	if cabinet.Name() != "Зарядный шкаф" {
		t.Errorf("наименование: получено %q", cabinet.Name())
	}
	if cabinet.DailyEnergyKWh() != 1.2 {
		t.Errorf("расход энергии: получено %v", cabinet.DailyEnergyKWh())
	}
}

func TestHelmetIsOnlyInventoriable(t *testing.T) {
	helmet, err := domain.NewHelmet("HM-1")
	if err != nil {
		t.Fatalf("шлем не создан: %v", err)
	}

	var item domain.Inventoriable = helmet
	if _, ok := item.(domain.EnergyConsumer); ok {
		t.Error("шлем не должен реализовывать EnergyConsumer")
	}
	if _, ok := item.(domain.Rentable); ok {
		t.Error("шлем не должен реализовывать Rentable")
	}
	if helmet.Name() != "Шлем" || helmet.InventoryNumber() != "HM-1" {
		t.Errorf("шлем: получено %q %q", helmet.Name(), helmet.InventoryNumber())
	}
}

func TestIsNoviceFriendly(t *testing.T) {
	tests := []struct {
		name     string
		friendly int
		want     bool
	}{
		{"ниже порога", 5, false},
		{"ровно порог", 6, true},
		{"выше порога", 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bicycle, err := domain.NewBicycle("BC-1", tt.friendly)
			if err != nil {
				t.Fatalf("велосипед не создан: %v", err)
			}
			if got := domain.IsNoviceFriendly(bicycle); got != tt.want {
				t.Fatalf("получено %v, ожидалось %v", got, tt.want)
			}
		})
	}
}

func TestMonowheelImplementsAllContracts(t *testing.T) {
	monowheel, err := domain.NewMonowheel("MW-1", 2, 0.30)
	if err != nil {
		t.Fatalf("моноколесо не создано: %v", err)
	}

	var item domain.Inventoriable = monowheel
	if _, ok := item.(domain.EnergyConsumer); !ok {
		t.Error("моноколесо должно реализовывать EnergyConsumer")
	}
	if _, ok := item.(domain.Rentable); !ok {
		t.Error("моноколесо должно реализовывать Rentable")
	}
	if monowheel.Name() != "Моноколесо" {
		t.Errorf("наименование: получено %q", monowheel.Name())
	}
	if monowheel.DailyEnergyKWh() != 0.30 {
		t.Errorf("расход энергии: получено %v", monowheel.DailyEnergyKWh())
	}
	if domain.IsNoviceFriendly(monowheel) {
		t.Error("моноколесо с простотой 2 не должно подходить новичкам")
	}
}
