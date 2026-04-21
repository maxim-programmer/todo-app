// Package core_http_types содержит типы для HTTP-слоя, расширяющие доменные типы.
package core_http_types

import (
	"encoding/json"

	"github.com/nilchan-social/golang-todoapp/internal/core/domain"
)

// Nullable[T] — HTTP-версия доменного Nullable[T] с реализацией UnmarshalJSON.
//
// Встраивание domain.Nullable[T] даёт поля Value и Set.
// Добавляем UnmarshalJSON, чтобы json.Decoder корректно различал три случая:
//
//   - Поле отсутствует в JSON → UnmarshalJSON не вызывается         → Set==false (zero value)
//   - Поле = null             → UnmarshalJSON вызывается с b=`null` → Set=true, Value=*nil
//   - Поле = конкретное значение                                    → Set=true, Value=&value
//
// После десериализации вызываем ToDomain() для передачи в сервисный слой.
type Nullable[T any] struct {
	domain.Nullable[T]
}

// UnmarshalJSON реализует encoding/json.Unmarshaler.
// Вызывается json.Decoder когда поле присутствует в теле запроса.
// Факт вызова (Set=true) означает, что HTTP клиент намеренно передал это поле.
func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	n.Set = true

	if string(b) == "null" {
		n.Value = nil

		return nil
	}

	var value T
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	n.Value = &value

	return nil
}

// ToDomain конвертирует HTTP Nullable в доменный Nullable для передачи в сервис.
func (n *Nullable[T]) ToDomain() domain.Nullable[T] {
	return domain.Nullable[T]{
		Value: n.Value,
		Set:   n.Set,
	}
}

/*
-------------------
JSON: {}
Nullable:
	- Value: *nil
	- Set: false


-------------------
JSON: {
	"phone_number": "+79998887766"
}
Nullable:
	- Value: *"+79998887766"
	- Set: true


-------------------
JSON: {
	"phone_number": null
}
Nullable:
	- Value: *nil
	- Set: true
*/
