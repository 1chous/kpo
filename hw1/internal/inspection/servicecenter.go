package inspection

import (
	"errors"
	"fmt"

	"github.com/1chous/kpo/hw1/internal/domain"
)

var (
	ErrNoInventoryNumber = errors.New("у устройства нет инвентарного номера")
	ErrEnergyTooHigh     = errors.New("суточное энергопотребление выше нормы парка")
)

type StandardServiceCenter struct {
	maxDailyEnergyKWh float64
}

func NewStandardServiceCenter(maxDailyEnergyKWh float64) *StandardServiceCenter {
	return &StandardServiceCenter{maxDailyEnergyKWh: maxDailyEnergyKWh}
}

func (s *StandardServiceCenter) Inspect(item domain.Inventoriable) error {
	if item.InventoryNumber() == "" {
		return ErrNoInventoryNumber
	}

	consumer, ok := item.(domain.EnergyConsumer)
	if !ok {
		return nil
	}

	if consumer.DailyEnergyKWh() > s.maxDailyEnergyKWh {
		return fmt.Errorf("%w: %.2f кВт·ч при норме %.2f",
			ErrEnergyTooHigh, consumer.DailyEnergyKWh(), s.maxDailyEnergyKWh)
	}

	return nil
}
