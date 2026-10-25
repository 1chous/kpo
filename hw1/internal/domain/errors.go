package domain

import "errors"

var (
	ErrEmptyInventoryNumber = errors.New("инвентарный номер не может быть пустым")
	ErrEmptyName            = errors.New("наименование не может быть пустым")
	ErrNoviceOutOfRange     = errors.New("простота для новичка должна быть от 1 до 10")
	ErrNegativeEnergy       = errors.New("суточный расход энергии не может быть отрицательным")
)
