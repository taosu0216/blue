package domain

import (
	"github.com/hardcore-os/plato/ipconf/source"
	"sort"
	"sync"
)

type Dispatcher struct {
	candidateTable map[string]*Endport
	sync.RWMutex   //锁
}

var dp *Dispatcher

func init() {
	dp := &Dispatcher{}
	dp.candidateTable = make(map[string]*Endport)
	go func() {
		for event := range source.EventChan() {
			switch event.Type {
			case source.AddNodeEvent:
				dp.addNode(event)
			case source.DelNodeEvent:
				dp.delNode(event)
			}
		}
	}()
}

func Dispatch(ctx *IpConfContext) []*Endport {
	//step1: 获得候选endPort
	eds := dp.getCandidateEndPort(ctx)
	//step2: 计算得分
	for _, ed := range eds {
		ed.CalculateScore(ctx)
	}
	//step3: 全局排序
	sort.Slice(eds, func(i, j int) bool {
		//优先根据活跃分进行排序
		if eds[i].ActiveScore > eds[j].ActiveScore {
			return true
		}
		//然后再根据静态分排序
		if eds[i].ActiveScore == eds[j].ActiveScore {
			if eds[i].StaticScore > eds[j].StaticScore {
				return true
			}
			return false
		}
		return false
	})
	return eds
}

func (dp *Dispatcher) getCandidateEndPort(ctx *IpConfContext) []*Endport {
	dp.RLock()
	defer dp.RUnlock()
	candidateList := make([]*Endport, 0, len(dp.candidateTable))
	for _, ed := range dp.candidateTable {
		candidateList = append(candidateList, ed)
	}
	return candidateList
}

func (dp *Dispatcher) delNode(event *source.Event) {
	dp.Lock()
	defer dp.Unlock()
	delete(dp.candidateTable, event.Key())
}

func (dp *Dispatcher) addNode(event *source.Event) {
	dp.Lock()
	defer dp.Unlock()
	var (
		ed *Endport
		ok bool
	)
	if ed, ok = dp.candidateTable[event.Key()]; !ok { // 不存在
		ed = NewEndport(event.IP, event.Port)
		dp.candidateTable[event.Key()] = ed
	}
	ed.UpdateStat(&Stat{
		ConnectNum:   event.ConnectNum,
		MessageBytes: event.MessageBytes,
	})
}
