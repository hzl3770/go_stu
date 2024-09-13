package rwlatch

import "sync"

type ReaderWriterLatch03 struct {
	mu          *sync.Mutex
	readerCount uint32
	readerCond  *sync.Cond
	writerCond  *sync.Cond

	isWriting     bool
	writerWaiting bool
}

func NewReaderWriterLatch03() *ReaderWriterLatch03 {
	l := &ReaderWriterLatch03{}
	mu := &sync.Mutex{}

	l.mu = mu
	l.writerCond = sync.NewCond(mu)
	l.readerCond = sync.NewCond(mu)
	return l
}

func (l *ReaderWriterLatch03) RLock() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for l.isWriting || l.writerWaiting {
		l.readerCond.Wait()
	}

	l.readerCount++
}

func (l *ReaderWriterLatch03) RUnLock() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.readerCount--
	if l.readerCount == 0 && l.writerWaiting {
		l.writerCond.Signal()
	}
}

func (l *ReaderWriterLatch03) WLock() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writerWaiting = true

	for l.readerCount > 0 || l.isWriting {
		l.writerCond.Wait()
	}

	l.isWriting = true
	l.writerWaiting = true
}

func (l *ReaderWriterLatch03) WUnLock() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.isWriting = false
	if l.writerWaiting {
		l.writerCond.Signal()
	} else {
		l.readerCond.Broadcast()
	}
}
