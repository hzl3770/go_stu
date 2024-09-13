package rwlatch

import "sync"

const (
	maxReaders = uint32(1<<32 - 1)
)

type ReaderWriterLatch01 struct {
	readerCount uint32
	isWriting   bool
	mu          *sync.Mutex
	readerCond  *sync.Cond
	writerCond  *sync.Cond
}

func New() *ReaderWriterLatch01 {
	l := &ReaderWriterLatch01{}
	l.mu = &sync.Mutex{}
	l.readerCond = sync.NewCond(l.mu)
	l.writerCond = sync.NewCond(l.mu)
	return l
}

func (l *ReaderWriterLatch01) isReaderFull() bool {
	return l.readerCount == maxReaders
}

func (l *ReaderWriterLatch01) isReaderEmpty() bool {
	return l.readerCount == 0
}

func (l *ReaderWriterLatch01) RLock() {
	l.mu.Lock()
	for l.isWriting || l.isReaderFull() {
		l.readerCond.Wait()
	}
	l.readerCount++
	l.mu.Unlock()
}

func (l *ReaderWriterLatch01) RUnLock() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.readerCount--
	if l.isWriting {
		if l.isReaderEmpty() {
			l.writerCond.Broadcast()
		}
	} else {
		if l.readerCount == maxReaders-1 {
			l.readerCond.Signal()
		}
	}

}

func (l *ReaderWriterLatch01) WLock() {
	l.mu.Lock()
	for l.isWriting {
		l.writerCond.Wait()
	}

	l.isWriting = true
	for !l.isReaderEmpty() {
		l.readerCond.Wait()
	}

	l.mu.Unlock()
}

func (l *ReaderWriterLatch01) WULock() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.isWriting = false
	l.readerCond.Broadcast()
}
