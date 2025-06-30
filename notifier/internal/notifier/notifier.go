package notifier

import "context"

type Notifier interface {
	Send(context.Context, Notification) error
}

type Notification struct {
	To      string
	Message string
}

func New() (Notifier, error) {

}

type telegramNotifier struct{}

func (n *telegramNotifier) Send(context.Context, Notification) error {}
