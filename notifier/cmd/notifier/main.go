package main

import (
	"github.com/vctrl/currency-service/notifier/internal/consumer"
	"github.com/vctrl/currency-service/notifier/internal/notifier"
)

func main() {
	c, err := consumer.New()

	n, err := notifier.New()

	for pb := range c.Consume() {
		// трансформируем []byte в notifier.Notification
		n.Send(ctx, notification)
	}
}

/*
crate table notifiactions (
id string
telegram string
currency_id string -> service Currency
expr string ("current_value > 10"), ("old_value > current_value")
)


1. Получаем новую цену валюты
2. Получить текущую цену валюты
3. Меняем в базе + записывапем что надо отправить событие
4. Публикум событие что цена поменялась


*/
