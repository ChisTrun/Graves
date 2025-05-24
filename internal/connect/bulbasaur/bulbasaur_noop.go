package bulbasaur

import (
	"context"
	"fmt"
)

type Noop struct{}

func (n *Noop) IncreaseBalance(ctx context.Context, userId uint64, amount float32) error {
	return fmt.Errorf("bulbasaur noop: IncreaseBalance not implemented")
}
