// 夹具：子包由父包行（`sandbox`）覆盖，且带两条「必须绿」的正腿——
// 同仓内部 import 允许（kit/logging），标准库 os/signal 的段不算交易信号。
package runner

import (
	"os/signal"
	"syscall"
	"time"

	"github.com/ruoxizhnya/tellurion-kit/logging"
)

type Limits struct {
	TimeoutMS int64
	MemoryMB  int64
}

type Runner struct {
	log    logging.Logger
	notify chan syscall.Signal
}

func New(l logging.Logger, timeout time.Duration) *Runner {
	return &Runner{log: l, notify: make(chan syscall.Signal, 1)}
}

func (r *Runner) Run(name string, args []string) error {
	signal.Notify(r.notify, syscall.SIGKILL)
	return nil
}
