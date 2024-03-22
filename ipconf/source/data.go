package source

import (
	"blue/common/config"
	"blue/common/discovery"
	"context"
	"github.com/bytedance/gopkg/util/logger"
)

func Init() {
	eventChan = make(chan *Event, 5)
	ctx := context.Background()
	go DataHandler(&ctx)
	//if config.IsDebug() {
	//	ctx := context.Background()
	//	testServiceRegister(&ctx, "7896", "node1")
	//	testServiceRegister(&ctx, "7897", "node2")
	//	testServiceRegister(&ctx, "7898", "node3")
	//}
}

// DataHandler 服务发现处理
func DataHandler(ctx *context.Context) {
	dis := discovery.NewServiceDiscovery(ctx)
	defer dis.Close()
	setFunc := func(key, value string) {
		if ed, err := discovery.UnMarshal([]byte(value)); err == nil {
			if event := NewEvent(ed); event != nil {
				event.Type = AddNodeEvent
				eventChan <- event
			}
		} else {
			logger.CtxErrorf(*ctx, "DataHandler.setFunc.err :%s", err.Error())
		}
	}
	delFunc := func(key, value string) {
		if ed, err := discovery.UnMarshal([]byte(value)); err == nil {
			if event := NewEvent(ed); ed != nil {
				event.Type = DelNodeEvent
				eventChan <- event
			}
		} else {
			logger.CtxErrorf(*ctx, "DataHandler.delFunc.err :%s", err.Error())
		}
	}
	err := dis.WatchService(config.GetServicePathForIPConf(), setFunc, delFunc)
	if err != nil {
		panic(err)
	}
}
