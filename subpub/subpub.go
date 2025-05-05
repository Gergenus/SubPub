package subpub

import (
	"context"
	"errors"
	"sync"
)

// MessageHandler is a callback function that processes messages delivered to subscribers.
type MessageHandler func(msg interface{})

type Subscription interface {
	// Unsubscribe will remove interest in the current subject subscription is for.
	Unsubscribe()
}

type SubPub interface {
	// Subscribe creates an asynchronous queue subscriber on the given subject.
	Subscribe(subject string, cb MessageHandler) (Subscription, error)

	// Publish publishes the msg argument to the given subject.
	Publish(subject string, msg interface{}) error

	// Close will shutdown sub-pub system.
	// May be blocked by data delivery until the context is canceled.
	Close(ctx context.Context) error
}

type subscription struct {
	subject  string
	handler  MessageHandler
	msgChan  chan interface{}
	doneChan chan struct{}
	sp       *subPub
}

func (s *subscription) Unsubscribe() {
	s.sp.mu.Lock()
	defer s.sp.mu.Unlock()
	subs := s.sp.subjects[s.subject]
	for i, sub := range subs {
		if sub == s {
			s.sp.subjects[s.subject] = append(subs[:i], subs[i+1:]...)
			close(s.doneChan)
			close(s.msgChan)
			return
		}
	}
}

type subPub struct {
	mu        sync.RWMutex
	subjects  map[string][]*subscription
	closed    bool
	closeChan chan struct{}
}

func NewSubPub() SubPub {
	return &subPub{
		subjects:  make(map[string][]*subscription),
		closeChan: make(chan struct{}),
	}
}

func (sp *subPub) Subscribe(subject string, cb MessageHandler) (Subscription, error) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if sp.closed {
		return nil, errors.New("subpub closed")
	}

	sub := &subscription{
		subject:  subject,
		handler:  cb,
		msgChan:  make(chan interface{}, 100),
		doneChan: make(chan struct{}),
		sp:       sp,
	}

	sp.subjects[subject] = append(sp.subjects[subject], sub)

	go func() {
		for {
			select {
			case msg := <-sub.msgChan:
				sub.handler(msg) // пункт 3 использование каналов позволяет нам выполнить требование три о FIFO очереди
			case <-sub.doneChan:
				return
			case <-sp.closeChan:
				return
			}
		}
	}()

	return sub, nil
}

func (sp *subPub) Publish(subject string, msg interface{}) error {
	sp.mu.RLock()
	defer sp.mu.RUnlock()

	if sp.closed {
		return errors.New("subpub closed")
	}

	subs := make([]*subscription, len(sp.subjects[subject]))
	copy(subs, sp.subjects[subject])

	for _, sub := range subs {
		select {
		case sub.msgChan <- msg:
		default:
		}
	}

	return nil
}

func (sp *subPub) Close(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err() // 4 пункт проверка контекста
	}
	sp.mu.Lock()
	if sp.closed {
		sp.mu.Unlock()
		return nil
	}
	sp.closed = true
	close(sp.closeChan)
	sp.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
