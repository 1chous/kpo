package main

import (
	"fmt"
	"os"

	"github.com/1chous/kpo/hw1/internal/di"
	"github.com/1chous/kpo/hw1/internal/domain"
	"github.com/1chous/kpo/hw1/internal/inspection"
	"github.com/1chous/kpo/hw1/internal/park"
)

const maxDailyEnergyKWh = 1.5

func main() {
	container := buildContainer()

	fleet, err := di.Resolve[*park.Park](container)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ошибка сборки приложения:", err)
		os.Exit(1)
	}

	runDemo(fleet)
}

func buildContainer() *di.Container {
	c := di.New()

	di.Register[park.ServiceCenter](c, di.Singleton, func(*di.Container) (park.ServiceCenter, error) {
		return inspection.NewStandardServiceCenter(maxDailyEnergyKWh), nil
	})

	di.Register[*park.Park](c, di.Singleton, func(c *di.Container) (*park.Park, error) {
		center, err := di.Resolve[park.ServiceCenter](c)
		if err != nil {
			return nil, err
		}
		return park.New(center), nil
	})

	return c
}

func runDemo(fleet *park.Park) {
	section("1. Приём транспорта: техосмотр пройден")
	accept(fleet, must(domain.NewScooter("SC-001", 8, 0.45)))
	accept(fleet, must(domain.NewEBike("EB-001", 6, 0.90)))
	accept(fleet, must(domain.NewBicycle("BC-001", 9)))
	accept(fleet, must(domain.NewScooter("SC-002", 4, 0.50)))
	accept(fleet, must(domain.NewMonowheel("MW-001", 2, 0.30)))

	section("2. Приём вещей на баланс")
	accept(fleet, must(domain.NewHelmet("HM-001")))
	accept(fleet, must(domain.NewDockStation("DS-001")))
	accept(fleet, must(domain.NewChargingCabinet("CH-001", 1.20)))

	section("3. Техосмотр отклоняет: расход выше нормы парка")
	accept(fleet, must(domain.NewScooter("SC-999", 7, 3.00)))

	section("4. Инварианты домена: некорректный объект не создаётся")
	if _, err := domain.NewScooter("", 5, 0.30); err != nil {
		fmt.Println("   отклонено конструктором:", err)
	}
	if _, err := domain.NewScooter("SC-003", 42, 0.30); err != nil {
		fmt.Println("   отклонено конструктором:", err)
	}
	if _, err := domain.NewChargingCabinet("CH-002", -1); err != nil {
		fmt.Println("   отклонено конструктором:", err)
	}

	section("5. Отчёт: число единиц на балансе")
	fmt.Printf("   единиц: %d\n", fleet.Count())

	section("6. Отчёт: суммарное суточное энергопотребление парка")
	fmt.Printf("   %.2f кВт·ч\n", fleet.TotalDailyEnergyKWh())

	section("7. Прокатная техника для новичков (простота >= 6)")
	for _, item := range fleet.NoviceFriendly() {
		fmt.Printf("   %-18s %-8s простота %d\n",
			item.Name(), item.InventoryNumber(), item.NoviceFriendliness())
	}

	section("8. Всё на балансе: наименование и инвентарный номер")
	for _, item := range fleet.Items() {
		fmt.Printf("   %-18s %s\n", item.Name(), item.InventoryNumber())
	}
}

func accept(fleet *park.Park, item domain.Inventoriable) {
	if err := fleet.Accept(item); err != nil {
		fmt.Printf("   ОТКЛОНЕНО %-8s %-18s %v\n", item.InventoryNumber(), item.Name(), err)
		return
	}
	fmt.Printf("   принято   %-8s %s\n", item.InventoryNumber(), item.Name())
}

func must(item domain.Inventoriable, err error) domain.Inventoriable {
	if err != nil {
		panic(fmt.Sprintf("некорректные демонстрационные данные: %v", err))
	}
	return item
}

func section(title string) {
	fmt.Printf("\n=== %s ===\n", title)
}
