package park

import (
	"fmt"
	"slices"

	"github.com/1chous/kpo/hw1/internal/domain"
)

type ServiceCenter interface {
	Inspect(item domain.Inventoriable) error
}

type Park struct {
	serviceCenter ServiceCenter
	items         []domain.Inventoriable
}

func New(serviceCenter ServiceCenter) *Park {
	return &Park{serviceCenter: serviceCenter}
}

func (p *Park) Accept(item domain.Inventoriable) error {
	if err := p.serviceCenter.Inspect(item); err != nil {
		return fmt.Errorf("техосмотр не пройден: %w", err)
	}
	p.items = append(p.items, item)
	return nil
}

func (p *Park) Count() int {
	return len(p.items)
}

func (p *Park) Items() []domain.Inventoriable {
	return slices.Clone(p.items)
}

func (p *Park) TotalDailyEnergyKWh() float64 {
	var total float64
	for _, item := range p.items {
		if consumer, ok := item.(domain.EnergyConsumer); ok {
			total += consumer.DailyEnergyKWh()
		}
	}
	return total
}

func (p *Park) NoviceFriendly() []domain.Rentable {
	var result []domain.Rentable
	for _, item := range p.items {
		rentable, ok := item.(domain.Rentable)
		if !ok {
			continue
		}
		if domain.IsNoviceFriendly(rentable) {
			result = append(result, rentable)
		}
	}
	return result
}
