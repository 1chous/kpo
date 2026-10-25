package inspection_test

import (
	"errors"
	"testing"

	"github.com/1chous/kpo/hw1/internal/domain"
	"github.com/1chous/kpo/hw1/internal/inspection"
)

func TestStandardServiceCenterInspect(t *testing.T) {
	const norm = 1.5

	scooterWithinNorm := build(t, func() (domain.Inventoriable, error) {
		return domain.NewScooter("SC-1", 8, 1.40)
	})
	scooterOverNorm := build(t, func() (domain.Inventoriable, error) {
		return domain.NewScooter("SC-2", 8, 1.60)
	})
	scooterAtNorm := build(t, func() (domain.Inventoriable, error) {
		return domain.NewScooter("SC-3", 8, norm)
	})
	bicycle := build(t, func() (domain.Inventoriable, error) {
		return domain.NewBicycle("BC-1", 9)
	})
	helmet := build(t, func() (domain.Inventoriable, error) {
		return domain.NewHelmet("HM-1")
	})

	tests := []struct {
		name    string
		item    domain.Inventoriable
		wantErr error
	}{
		{"расход в пределах нормы — принято", scooterWithinNorm, nil},
		{"расход ровно по норме — принято", scooterAtNorm, nil},
		{"расход выше нормы — отклонено", scooterOverNorm, inspection.ErrEnergyTooHigh},
		{"велосипед не потребляет энергию — принято", bicycle, nil},
		{"вещь без энергопотребления — принято", helmet, nil},
	}

	center := inspection.NewStandardServiceCenter(norm)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := center.Inspect(tt.item)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("ожидалась приёмка, получено %v", err)
				}
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("получено %v, ожидалось %v", err, tt.wantErr)
			}
		})
	}
}

type emptyNumberItem struct{}

func (emptyNumberItem) InventoryNumber() string { return "" }
func (emptyNumberItem) Name() string            { return "Неизвестное устройство" }

func TestStandardServiceCenterRejectsItemWithoutInventoryNumber(t *testing.T) {
	center := inspection.NewStandardServiceCenter(1.5)

	err := center.Inspect(emptyNumberItem{})
	if !errors.Is(err, inspection.ErrNoInventoryNumber) {
		t.Fatalf("получено %v, ожидалось %v", err, inspection.ErrNoInventoryNumber)
	}
}

func build(t *testing.T, make func() (domain.Inventoriable, error)) domain.Inventoriable {
	t.Helper()
	item, err := make()
	if err != nil {
		t.Fatalf("не удалось подготовить объект: %v", err)
	}
	return item
}
