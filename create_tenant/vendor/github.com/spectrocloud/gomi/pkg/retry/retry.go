package retry

import (
	"fmt"
	"github.com/spectrocloud/gomi/util"
	"math/big"
	"time"
)

type RetryOption struct {
	retryMsg  string
	attempts  int
	sleep     time.Duration
	retryFlag bool
}

func NewRetryOption(retryMsg string, attempts int, sleep time.Duration, retryFlag bool) *RetryOption {
	return &RetryOption{retryMsg: retryMsg, attempts: attempts, sleep: sleep, retryFlag: retryFlag}
}

func (retryOption *RetryOption) Retry(f func() error) error {
	return retryOp(retryOption.retryMsg, retryOption.attempts, retryOption.sleep, f, retryOption.retryFlag)
}

func Retry(retryMsg string, attempts int, sleep time.Duration, f func() error) error {
	return retryOp(retryMsg, attempts, sleep, f, true)
}

func RetryWithErrRetryCond(retryMsg string, attempts int, sleep time.Duration, f func() error) error {
	return retryOp(retryMsg, attempts, sleep, f, false)
}

func retryOp(retryMsg string, attempts int, sleep time.Duration, f func() error, retryFlag bool) error {
	var err error
	t := attempts
	for ; attempts >= 0; attempts-- {
		if t > attempts {
			fmt.Printf("retrying (%d/%d): %s ", t-attempts, t, retryMsg)
		}
		err = f()
		if err != nil {
			if sleep > 0 {
				time.Sleep(sleep)
				jitter := time.Duration(util.GenerateRandomInt64(big.NewInt(int64(sleep))))
				sleep = (2 * sleep) + jitter/2 //exponential sleep with jitter
			}
		} else {
			return nil
		}
	}
	return err
}
