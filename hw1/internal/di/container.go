package di

import (
	"fmt"
	"reflect"
)

type Lifetime int

const (
	Transient Lifetime = iota
	Singleton
)

type registration struct {
	lifetime Lifetime
	factory  func(*Container) (any, error)
	instance any
	resolved bool
}

type Container struct {
	registrations map[reflect.Type]*registration
}

func New() *Container {
	return &Container{registrations: make(map[reflect.Type]*registration)}
}

func Register[T any](c *Container, lifetime Lifetime, factory func(*Container) (T, error)) {
	c.registrations[reflect.TypeFor[T]()] = &registration{
		lifetime: lifetime,
		factory: func(inner *Container) (any, error) {
			return factory(inner)
		},
	}
}

func Resolve[T any](c *Container) (T, error) {
	var zero T

	key := reflect.TypeFor[T]()
	reg, ok := c.registrations[key]
	if !ok {
		return zero, fmt.Errorf("di: тип %s не зарегистрирован", key)
	}

	if reg.lifetime == Singleton && reg.resolved {
		return reg.instance.(T), nil
	}

	value, err := reg.factory(c)
	if err != nil {
		return zero, fmt.Errorf("di: не удалось создать %s: %w", key, err)
	}

	if reg.lifetime == Singleton {
		reg.instance = value
		reg.resolved = true
	}

	return value.(T), nil
}
