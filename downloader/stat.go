package main

import (
	"fmt"
	"sync"
	"sort"
	"time"
	"sync/atomic"
	"github.com/piotrnar/gocoin/lib/btc"
	"github.com/piotrnar/gocoin/lib/others/sys"
)

var (
	_CNT map[string] uint = make(map[string] uint)
	cnt_mut sync.Mutex
	EmptyInProgressCnt uint64
	LastBlockAsked uint32
)


func COUNTER(s string) {
	cnt_mut.Lock()
	_CNT[s]++
	cnt_mut.Unlock()
}


func print_counters() {
	var s string
	cnt_mut.Lock()
	ss := make([]string, len(_CNT))
	i := 0
	for k, v := range _CNT {
		ss[i] = fmt.Sprintf("%s=%d", k, v)
		i++
	}
	cnt_mut.Unlock()
	sort.Strings(ss)
	for i = range ss {
		s += "  "+ss[i]
	}
	fmt.Println(s)
	return
}

func print_stats() {
	sec := float64(time.Now().Sub(DlStartTime)) / 1e6
	BlocksMutex.Lock()
	var max_block_height_in_progress uint32
	for _, v := range BlocksInProgress {
		if v.Height>max_block_height_in_progress {
			max_block_height_in_progress = v.Height
		}
	}
	aloc, _ := sys.MemUsed()
	s := fmt.Sprintf("Blck:%d/%d/%d/%d  Pend:%d  Que:%d (%dMB)  Dl:%d  "+
		"Cach:%d (%dMB)  BLen:%dkB  Net:%d  [%.0f => %.0f KBps]  Mem:%dMB  EC:%d  %.1fmin",
		atomic.LoadUint32(&LastStoredBlock), BlocksComplete, max_block_height_in_progress,
		LastBlockHeight, len(BlocksToGet), len(BlockQueue),
		atomic.LoadUint64(&BlocksQueuedSize)>>20, len(BlocksInProgress),
		len(BlocksCached), BlocksCachedSize>>20, avg_block_size()/1000,
		open_connection_count(), float64(atomic.LoadUint64(&DlBytesDownloaded))/sec,
		float64(atomic.LoadUint64(&DlBytesProcessed))/sec,
		aloc>>20, btc.EcdsaVerifyCnt, time.Now().Sub(StartTime).Minutes())
	BlocksMutex.Unlock()
	fmt.Println(s)
}
