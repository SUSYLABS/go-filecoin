package series

import (
	"context"
)

var CKMiningOnce struct{} = struct{}{}

type MiningOnceFunc func()

// MiningOnceFromCtx extracts the MiningOnceFunc through key CKMiningOnce
func MiningOnceFromCtx(ctx context.Context) {
	MiningOnce, ok := ctx.Value(CKMiningOnce).(MiningOnceFunc)
	if ok {
		MiningOnce()
	}
}
