package rwlatch

import "sync"

type ReaderWriterLatch02 struct {
	mu          *sync.Mutex
	readerCount uint32
	readerCond  *sync.Cond
	writerCond  *sync.Cond

	isWriting bool
}

func NewReaderWriterLatch02() *ReaderWriterLatch02 {
	l := &ReaderWriterLatch02{}
	mu := &sync.Mutex{}

	l.mu = mu
	l.writerCond = sync.NewCond(mu)
	l.readerCond = sync.NewCond(mu)
	return l
}

func (l *ReaderWriterLatch02) RLock() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for l.isWriting {
		l.readerCond.Wait()
	}

	l.readerCount++
}

func (l *ReaderWriterLatch02) RUnLock() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.readerCount--
	if l.readerCount == 0 {
		l.writerCond.Signal()
	}
}

func (l *ReaderWriterLatch02) WLock() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for l.readerCount > 0 || l.isWriting {
		l.writerCond.Wait()
	}

	l.isWriting = true
}

func (l *ReaderWriterLatch02) WUnLock() {
	l.mu.Lock()
	l.isWriting = false
	l.mu.Unlock()
	l.readerCond.Broadcast()
}
