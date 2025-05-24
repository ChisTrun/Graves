package bulbasaur

import (
	"context"
	"flag"
	"fmt"
	"graves/pkg/config"
	"graves/pkg/logger/pkg/logging"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials/insecure"
)

var instaccnce Bulbasaur

type Bulbasaur interface {
	IncreaseBalance(ctx context.Context, userId uint64, amount float32) error
}

type bulbasaur struct {
	client VenusaurClient
}

func init() {
	cfg, err := config.GetInstance()
	if err != nil {
		logging.Logger(context.Background()).Error(fmt.Sprintf("failed to get config instance: %v", err.Error()))
		instaccnce = &Noop{}
		return
	}

	ctx := context.Background()
	if cfg.Bulbasaur == nil {
		logging.Logger(ctx).Error("Bulbasaur configuration is not provided")
		instaccnce = &Noop{}
		return
	}

	serverAddr := flag.String("addr", fmt.Sprintf("%v:%v", cfg.Bulbasaur.Host, cfg.Bulbasaur.Port), "The server address in the format of host:port")
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(*serverAddr, opts...)
	if err != nil {
		logging.Logger(ctx).Error(fmt.Sprintf("failed to connect: %v", err.Error()))
		instaccnce = &Noop{}
		return
	}

	client := NewVenusaurClient(conn)

	instaccnce = &bulbasaur{client: client}
}

func GetInstance() Bulbasaur {
	return instaccnce
}

func (b *bulbasaur) IncreaseBalance(ctx context.Context, userId uint64, amount float32) error {
	if b.client == nil {
		return fmt.Errorf("bulbasaur client is not initialized")
	}

	req := &IncreaseBalanceRequest{
		UserId: userId,
		Amount: amount,
	}

	_, err := b.client.IncreaseBalance(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
