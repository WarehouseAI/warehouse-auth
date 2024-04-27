package dependencies

import (
	"auth-service/internal/adapter/mail"
	"auth-service/internal/adapter/random"
	"auth-service/internal/adapter/time"
	"auth-service/internal/adapter/user"

	"go.uber.org/zap"
)

func (d *dependencies) TimeAdapter() time.Adapter {
	if d.timeAdapter == nil {
		d.timeAdapter = time.NewAdapter(
			d.cfg.Time.Locale,
		)
	}

	return d.timeAdapter
}

func (d *dependencies) RandomAdapter() random.Adapter {
	if d.randomAdapter == nil {
		d.randomAdapter = random.NewAdapter()
	}
	return d.randomAdapter
}

func (d *dependencies) UserAdapter() user.Adapter {
	if d.userAdapter == nil {
		var err error
		if d.userAdapter, err = user.NewAdapter(d.cfg.Grpc); err != nil {
			d.log.Zap().Panic("create user grpc adapter", zap.Error(err))
		}
	}

	return d.userAdapter
}

func (d *dependencies) MailAdapter() mail.Adapter {
	if d.mailAdapter == nil {
		var err error
		if d.mailAdapter, err = mail.NewAdapter(d.cfg.Rabbit.MailQueue, d.RabbitClient()); err != nil {
			d.log.Zap().Panic("create mail broker adapter", zap.Error(err))
		}
	}

	return d.mailAdapter
}
