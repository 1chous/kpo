package di_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/1chous/kpo/hw1/internal/di"
)

type greeter interface {
	Greet() string
}

type ruGreeter struct{ id int }

func (g *ruGreeter) Greet() string { return "привет" }

func TestSingletonReturnsSameInstance(t *testing.T) {
	c := di.New()
	created := 0

	di.Register[greeter](c, di.Singleton, func(*di.Container) (greeter, error) {
		created++
		return &ruGreeter{id: created}, nil
	})

	first, err := di.Resolve[greeter](c)
	if err != nil {
		t.Fatalf("первый Resolve: %v", err)
	}
	second, err := di.Resolve[greeter](c)
	if err != nil {
		t.Fatalf("второй Resolve: %v", err)
	}

	if first != second {
		t.Error("singleton вернул разные объекты")
	}
	if created != 1 {
		t.Errorf("фабрика вызвана %d раз, ожидался 1", created)
	}
}

func TestTransientReturnsNewInstance(t *testing.T) {
	c := di.New()
	created := 0

	di.Register[greeter](c, di.Transient, func(*di.Container) (greeter, error) {
		created++
		return &ruGreeter{id: created}, nil
	})

	first, err := di.Resolve[greeter](c)
	if err != nil {
		t.Fatalf("первый Resolve: %v", err)
	}
	second, err := di.Resolve[greeter](c)
	if err != nil {
		t.Fatalf("второй Resolve: %v", err)
	}

	if first == second {
		t.Error("transient вернул один и тот же объект")
	}
	if created != 2 {
		t.Errorf("фабрика вызвана %d раз, ожидалось 2", created)
	}
}

func TestResolveUnregisteredTypeFails(t *testing.T) {
	c := di.New()

	_, err := di.Resolve[greeter](c)
	if err == nil {
		t.Fatal("ожидалась ошибка для незарегистрированного типа")
	}
	if !strings.Contains(err.Error(), "не зарегистрирован") {
		t.Errorf("непонятное сообщение об ошибке: %v", err)
	}
}

func TestResolvePropagatesFactoryError(t *testing.T) {
	c := di.New()
	errBoom := errors.New("конфигурация не прочитана")

	di.Register[greeter](c, di.Singleton, func(*di.Container) (greeter, error) {
		return nil, errBoom
	})

	_, err := di.Resolve[greeter](c)
	if !errors.Is(err, errBoom) {
		t.Fatalf("получено %v, ожидалась исходная ошибка фабрики", err)
	}
}

func TestResolveBuildsDependencyGraph(t *testing.T) {
	c := di.New()

	type service struct{ g greeter }

	di.Register[greeter](c, di.Singleton, func(*di.Container) (greeter, error) {
		return &ruGreeter{}, nil
	})
	di.Register[*service](c, di.Singleton, func(inner *di.Container) (*service, error) {
		g, err := di.Resolve[greeter](inner)
		if err != nil {
			return nil, err
		}
		return &service{g: g}, nil
	})

	got, err := di.Resolve[*service](c)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if got.g == nil || got.g.Greet() != "привет" {
		t.Fatal("зависимость не внедрена")
	}
}
