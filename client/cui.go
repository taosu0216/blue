package client

import (
	"blue/client/sdk"
	"errors"
	"fmt"
	"github.com/gookit/color"
	"github.com/rocket049/gocui"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	verbose = true
	step = 1
}

var (
	buf     string
	chat    *sdk.Chat
	step    int
	pos     int
	verbose bool
)

type VOT struct {
	Name, Msg, Sep string
}

func RunMain() {
	// step1 创建chat的核心对象
	chat = sdk.NewChat(net.ParseIP("127.0.0.1"), 8900, "taosu", "12312321", "2131")

	// step2 创建 GUI 图层对象并进行参与与回调函数的配置
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	g.Cursor = true
	g.Mouse = false
	g.ASCII = true
	// 设置编排函数
	g.SetManagerFunc(layout)

	// 注册回调事件
	if err := g.SetKeybinding("main", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("main", gocui.KeyEnter, gocui.ModNone, viewUpdate); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("main", gocui.KeyPgup, gocui.ModNone, viewUpScroll); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("main", gocui.KeyPgdn, gocui.ModNone, viewDownScroll); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, pasteDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone, pasteUP); err != nil {
		log.Panicln(err)
	}
	//go func() {
	//	time.Sleep(10 * time.Second)
	//	chat.ReConn()
	//}()
	// 启动消费函数
	go doRecv(g)
	if err := g.MainLoop(); err != nil {
		log.Println(err)
	}
	_ = ioutil.WriteFile("chat.log", []byte(buf), 0644)
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if err := viewHead(g, 1, 1, maxX-1, 3); err != nil {
		return err
	}
	if err := viewOutput(g, 1, 4, maxX-1, maxY-4); err != nil {
		return err
	}
	if err := viewInput(g, 1, maxY-3, maxX-1, maxY-1); err != nil {
		return err
	}
	return nil
}

func viewHead(g *gocui.Gui, x0, y0, x1, y1 int) error {
	if v, err := g.SetView("head", x0, y0, x1, y1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Wrap = false
		v.Overwrite = true
		msg := "blue: im系统聊天对话框"
		setHeadText(g, msg)
	}
	return nil
}

func viewOutput(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView("out", x0, y0, x1, y1)
	if err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Wrap = true
		v.Overwrite = false
		v.Autoscroll = true
		v.SelBgColor = gocui.ColorWhite
		v.Title = "Messages"
	}
	return nil
}

func viewInput(g *gocui.Gui, x0, y0, x1, y1 int) error {
	if v, err := g.SetView("main", x0, y0, x1, y1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		//当 err == gocui.ErrUnknownView 时运行
		v.Editable = true
		v.Wrap = true
		v.Overwrite = false
		v.Title = "Input"
		if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	return nil
}

func setHeadText(g *gocui.Gui, msg string) {
	v, err := g.View("head")
	if err == nil {
		v.Clear()
		_, _ = fmt.Fprint(v, color.FgCyan.Text(msg))
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	chat.Close()
	ov, _ := g.View("out")
	buf = ov.Buffer()
	g.Close()
	return gocui.ErrQuit
}

func viewUpdate(g *gocui.Gui, cv *gocui.View) error {
	doSay(g, cv)
	l := len(cv.Buffer())
	cv.MoveCursor(0-l, 0, true)
	cv.Clear()
	return nil
}

func viewUpScroll(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("out")
	v.Autoscroll = false
	ox, oy := v.Origin()
	if err == nil {
		_ = v.SetOrigin(ox, oy-1)
	}
	return nil
}

func viewDownScroll(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("out")
	_, y := v.Size()
	ox, oy := v.Origin()
	lnum := len(v.BufferLines())
	if err == nil {
		if oy > lnum-y-1 {
			v.Autoscroll = true
		} else {
			_ = v.SetOrigin(ox, oy+1)
		}
	}
	return nil
}

func pasteUP(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("out")
	if err != nil {
		_, _ = fmt.Fprintf(cv, "error:%s", err)
		return nil
	}
	bls := v.BufferLines()
	lnum := len(bls)
	if pos < lnum-1 {
		pos++
	}
	cv.Clear()
	_, _ = fmt.Fprintf(cv, "%s", bls[lnum-pos-1])
	return nil
}

func pasteDown(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("out")
	if err != nil {
		_, _ = fmt.Fprintf(cv, "error:%s", err)
		return nil
	}
	if pos > 0 {
		pos--
	}
	bls := v.BufferLines()
	lnum := len(bls)
	cv.Clear()
	_, _ = fmt.Fprintf(cv, "%s", bls[lnum-pos-1])
	return nil
}

func doRecv(g *gocui.Gui) {
	recvChannel := chat.Recv()
	for msg := range recvChannel {
		switch msg.Type {
		case sdk.MsgTypeText:
			viewPrint(g, msg.Name, msg.Content, false)
		}
	}
	g.Close()
}

func viewPrint(g *gocui.Gui, name, msg string, newline bool) {
	var out VOT
	out.Name = name
	out.Msg = msg
	if newline {
		out.Sep = "\n"
	} else {
		out.Sep = " "
	}
	g.Update(out.Show)
}

func (sl VOT) Show(g *gocui.Gui) error {
	v, err := g.View("out")
	if err != nil {
		log.Println("No output view")
		return nil
	}
	_, _ = fmt.Fprintf(v, "%v:%v%v\n", color.FgGreen.Text(sl.Name), sl.Sep,
		color.FgYellow.Text(sl.Msg))
	return nil
}

func doSay(g *gocui.Gui, cv *gocui.View) {
	v, err := g.View("out")
	if cv != nil && err == nil {
		p := cv.ReadEditor()
		if p != nil {
			msg := &sdk.Message{
				Type:       sdk.MsgTypeText,
				Name:       "taosu",
				FromUserID: "123213",
				ToUserID:   "222222",
				Content:    string(p)}
			// 先把自己说的话显示到消息流中 fmt.Sprintf("%d", chat.GetCurClientID())
			//idKey := `1`
			viewPrint(g, "me", msg.Content, false)
			chat.Send(msg)
		}
		v.Autoscroll = true
	}
}
