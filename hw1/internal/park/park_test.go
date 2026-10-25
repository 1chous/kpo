package park_test

import (
	"errors"
	"math"
	"testing"

	"github.com/1chous/kpo/hw1/internal/domain"
	"github.com/1chous/kpo/hw1/internal/park"
)

type stubServiceCenter struct {
	err      error
	calls    int
	lastItem domain.Inventoriable
}

func (s *stubServiceCenter) Inspect(item domain.Inventoriable) error {
	s.calls++
	s.lastItem = item
	return s.err
}

var errRejected = errors.New("техника не прошла проверку тормозов")

func TestAcceptAddsItemWhenInspectionPasses(t *testing.T) {
	center := &stubServiceCenter{err: nil}
	fleet := park.New(center)
	scooter := newScooter(t, "SC-1", 8, 0.45)

	if err := fleet.Accept(scooter); err != nil {
		t.Fatalf("ожидалась приёмка, получено %v", err)
	}

	if fleet.Count() != 1 {
		t.Fatalf("на балансе %d единиц, ожидалась 1", fleet.Count())
	}
}

func TestAcceptRejectsItemWhenInspectionFails(t *testing.T) {
	center := &stubServiceCenter{err: errRejected}
	fleet := park.New(center)
	scooter := newScooter(t, "SC-1", 8, 0.45)

	err := fleet.Accept(scooter)

	if !errors.Is(err, errRejected) {
		t.Fatalf("получено %v, ожидалась причина отказа техосмотра", err)
	}
	if fleet.Count() != 0 {
		t.Fatalf("отклонённая техника попала на баланс: %d единиц", fleet.Count())
	}
}

func TestAcceptAlwaysAsksServiceCenter(t *testing.T) {
	center := &stubServiceCenter{}
	fleet := park.New(center)
	scooter := newScooter(t, "SC-1", 8, 0.45)

	if err := fleet.Accept(scooter); err != nil {
		t.Fatalf("ожидалась приёмка, получено %v", err)
	}

	if center.calls != 1 {
		t.Errorf("техосмотр вызван %d раз, ожидался 1", center.calls)
	}
	if center.lastItem != domain.Inventoriable(scooter) {
		t.Errorf("на техосмотр ушёл не тот объект: %v", center.lastItem)
	}
}

func TestTotalDailyEnergyKWhCountsOnlyConsumers(t *testing.T) {
	fleet := park.New(&stubServiceCenter{})

	accept(t, fleet, newScooter(t, "SC-1", 8, 0.45))
	accept(t, fleet, newEBike(t, "EB-1", 6, 0.90))
	accept(t, fleet, newBicycle(t, "BC-1", 9))
	accept(t, fleet, newHelmet(t, "HM-1"))
	accept(t, fleet, newCabinet(t, "CH-1", 1.20))
	accept(t, fleet, newMonowheel(t, "MW-1", 2, 0.30))

	const want = 0.45 + 0.90 + 1.20 + 0.30
	got := fleet.TotalDailyEnergyKWh()

	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("получено %.4f кВт·ч, ожидалось %.4f", got, want)
	}
}

func TestNoviceFriendlyFiltersByThresholdAndRentability(t *testing.T) {
	fleet := park.New(&stubServiceCenter{})

	accept(t, fleet, newScooter(t, "SC-OK", 8, 0.45))
	accept(t, fleet, newEBike(t, "EB-EDGE", 6, 0.90))
	accept(t, fleet, newScooter(t, "SC-LOW", 5, 0.45))
	accept(t, fleet, newCabinet(t, "CH-1", 1.20))
	accept(t, fleet, newHelmet(t, "HM-1"))
	accept(t, fleet, newMonowheel(t, "MW-1", 2, 0.30))

	got := fleet.NoviceFriendly()

	want := []string{"SC-OK", "EB-EDGE"}
	if len(got) != len(want) {
		t.Fatalf("получено %d единиц, ожидалось %d: %v", len(got), len(want), numbers(got))
	}
	for i, number := range want {
		if got[i].InventoryNumber() != number {
			t.Errorf("позиция %d: получено %q, ожидалось %q", i, got[i].InventoryNumber(), number)
		}
	}
}

func TestItemsReturnsCopy(t *testing.T) {
	fleet := park.New(&stubServiceCenter{})
	accept(t, fleet, newScooter(t, "SC-1", 8, 0.45))

	items := fleet.Items()
	items[0] = nil

	if fleet.Items()[0] == nil {
		t.Fatal("изменение возвращённого среза изменило состояние парка")
	}
}

func accept(t *testing.T, fleet *park.Park, item domain.Inventoriable) {
	t.Helper()
	if err := fleet.Accept(item); err != nil {
		t.Fatalf("не удалось принять %s: %v", item.InventoryNumber(), err)
	}
}

func newScooter(t *testing.T, number string, novice int, energy float64) domain.Scooter {
	t.Helper()
	item, err := domain.NewScooter(number, novice, energy)
	if err != nil {
		t.Fatalf("не удалось создать самокат: %v", err)
	}
	return item
}

func newEBike(t *testing.T, number string, novice int, energy float64) domain.EBike {
	t.Helper()
	item, err := domain.NewEBike(number, novice, energy)
	if err != nil {
		t.Fatalf("не удалось создать электровелосипед: %v", err)
	}
	return item
}

func newBicycle(t *testing.T, number string, novice int) domain.Bicycle {
	t.Helper()
	item, err := domain.NewBicycle(number, novice)
	if err != nil {
		t.Fatalf("не удалось создать велосипед: %v", err)
	}
	return item
}

func newHelmet(t *testing.T, number string) domain.Helmet {
	t.Helper()
	item, err := domain.NewHelmet(number)
	if err != nil {
		t.Fatalf("не удалось создать шлем: %v", err)
	}
	return item
}

func newCabinet(t *testing.T, number string, energy float64) domain.ChargingCabinet {
	t.Helper()
	item, err := domain.NewChargingCabinet(number, energy)
	if err != nil {
		t.Fatalf("не удалось создать зарядный шкаф: %v", err)
	}
	return item
}

func newMonowheel(t *testing.T, number string, novice int, energy float64) domain.Monowheel {
	t.Helper()
	item, err := domain.NewMonowheel(number, novice, energy)
	if err != nil {
		t.Fatalf("не удалось создать моноколесо: %v", err)
	}
	return item
}

func numbers(items []domain.Rentable) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, item.InventoryNumber())
	}
	return result
}
