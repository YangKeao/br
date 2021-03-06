package backup

import (
	"context"

	"github.com/pingcap/errors"
	pd "github.com/pingcap/pd/client"

	"github.com/pingcap/br/pkg/utils"
)

// GetGCSafePoint returns the current gc safe point.
// TODO: Some cluster may not enable distributed GC.
func GetGCSafePoint(ctx context.Context, pdClient pd.Client) (utils.Timestamp, error) {
	safePoint, err := pdClient.UpdateGCSafePoint(ctx, 0)
	if err != nil {
		return utils.Timestamp{}, errors.Trace(err)
	}
	return utils.DecodeTs(safePoint), nil
}

// CheckGCSafepoint checks whether the ts is older than GC safepoint.
func CheckGCSafepoint(ctx context.Context, pdClient pd.Client, ts uint64) error {
	// TODO: use PDClient.GetGCSafePoint instead once PD client exports it.
	safePoint, err := GetGCSafePoint(ctx, pdClient)
	if err != nil {
		return err
	}
	safePointTS := utils.EncodeTs(safePoint)
	if ts <= safePointTS {
		return errors.Errorf("GC safepoint %d exceed TS %d", safePointTS, ts)
	}
	return nil
}
