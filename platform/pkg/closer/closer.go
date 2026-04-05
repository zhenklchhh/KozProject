package closer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/zhenklchhh/KozProject/platform/pkg/logger"
	"go.uber.org/zap"
)

const (
	shutdownTimeout = 5 * time.Second
)

var (
	globalCloser = NewWithLogger(&logger.NoopLogger{})
)

type Logger interface {
	Info(context.Context, string, ...zap.Field)
	Error(context.Context, string, ...zap.Field)
}

type Closer struct {
	mu     sync.Mutex
	once   sync.Once
	done   chan struct{}
	funcs  []func(context.Context) error
	logger Logger
}

func New(signals ...os.Signal) *Closer {
	return NewWithLogger(logger.Logger(), signals...)
}

func NewWithLogger(logger Logger, signals ...os.Signal) *Closer {
	closer := &Closer{
		done:   make(chan struct{}),
		logger: logger,
	}
	if len(signals) > 0 {
		go closer.handleSignals(signals...)
	}
	return closer
}

func (c *Closer) Add(f ...func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f...)
}

func Add(f ...func(context.Context) error) {
	globalCloser.Add(f...)
}

func (c *Closer) AddNamed(name string, f func(context.Context) error) {
	c.Add(func(ctx context.Context) error {
		start := time.Now()
		c.logger.Info(ctx, fmt.Sprintf("Закрываем %s...", name))
		err := f(ctx)
		duration := time.Since(start)
		if err != nil {
			c.logger.Error(ctx, fmt.Sprintf("❌ Ошибка при закрытии %s: %v (заняло %s)", name, err, duration))
		} else {
			c.logger.Info(ctx, fmt.Sprintf("%s успешно закрыт за %s", name, duration))
		}
		return err
	})
}

func AddNamed(name string, f func(context.Context) error) {
	globalCloser.AddNamed(name, f)
}

func (c *Closer) CloseAll(ctx context.Context) error {
	var result error
	c.once.Do(func() {
		defer close(c.done)
		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		if len(funcs) == 0 {
			c.logger.Info(ctx, "Нет функций для закрытия")
			return
		}

		errCh := make(chan error, len(funcs))
		var wg sync.WaitGroup
		for i := len(funcs) - 1; i >= 0; i-- {
			f := funcs[i]
			wg.Add(1)
			go func(f func(context.Context) error) {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						errCh <- errors.New("panic recover in closer")
						c.logger.Error(ctx, "Паника в функции закрытия: %w", zap.Any("error", r))
					}
				}()
			}(f)
		}
		go func() {
			wg.Wait()
			close(errCh)
		}()
		for {
			select {
			case <-ctx.Done():
				c.logger.Info(ctx, "Контекст отменён во время закрытия", zap.Error(ctx.Err()))
				if result == nil {
					result = ctx.Err()
				}
				return
			case err, ok := <-errCh:
				if !ok {
					c.logger.Info(ctx, " Все ресурсы успешно закрыты")
					return
				}
				c.logger.Error(ctx, "Ошибка при закрытии", zap.Error(err))
				if result == nil {
					result = err
				}
			}
		}
	})
	return result
}

func CloseAll(ctx context.Context) error {
	return globalCloser.CloseAll(ctx)
}

func SetLogger(logger Logger) {
	globalCloser.SetLogger(logger)
}

func Configure(signals ...os.Signal) {
	go globalCloser.handleSignals(signals...)
}

func (c *Closer) handleSignals(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	defer signal.Stop(ch)
	select {
	case <-ch:
		c.logger.Info(context.Background(), "Получен системный сигнал, graceful shutdown...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer shutdownCancel()
		if err := c.CloseAll(shutdownCtx); err != nil {
			c.logger.Error(context.Background(), "Ошибка при закрытии ресурсов: %w", zap.Error(err))
		}
	case <-c.done:

	}
}

func (c *Closer) SetLogger(logger Logger) {
	c.logger = logger
}
