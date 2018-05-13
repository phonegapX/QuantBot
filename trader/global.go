package trader

import (
	//"encoding/json"
	//"fmt"
	"log"
	//"reflect"
	"sync"
	"time"

	"github.com/miaolz123/conver"
	"github.com/phonegapX/QuantBot/api"
	"github.com/phonegapX/QuantBot/constant"
	"github.com/phonegapX/QuantBot/model"
	"github.com/robertkrimen/otto"
)

type Tasks map[string][]task

// Global ...
type Global struct {
	model.Trader
	Logger  model.Logger   //利用这个对象保存日志
	ctx     *otto.Otto     //js虚拟机
	es      []api.Exchange //交易所列表
	tasks   Tasks          //任务列表
	running bool
	//statusLog string
}

//js中的一个任务,目的是可以并发工作
type task struct {
	ctx  *otto.Otto    //js虚拟机
	fn   otto.Value    //代表该任务的js函数
	args []interface{} //函数的参数
}

// Sleep ...
func (g *Global) Sleep(intervals ...interface{}) {
	interval := int64(0)
	if len(intervals) > 0 {
		interval = conver.Int64Must(intervals[0])
	}
	if interval > 0 {
		time.Sleep(time.Duration(interval * 1000000))
	} else {
		for _, e := range g.es {
			e.AutoSleep()
		}
	}
}

// Console ...
func (g *Global) Console(msgs ...interface{}) {
	log.Printf("%v %v\n", constant.INFO, msgs)
}

// Log ...
func (g *Global) Log(msgs ...interface{}) {
	g.Logger.Log(constant.INFO, "", 0.0, 0.0, msgs...)
}

// LogProfit ...
func (g *Global) LogProfit(msgs ...interface{}) {
	profit := 0.0
	if len(msgs) > 0 {
		profit = conver.Float64Must(msgs[0])
	}
	g.Logger.Log(constant.PROFIT, "", 0.0, profit, msgs[1:]...)
}

// LogStatus ...
//func (g *Global) LogStatus(msgs ...interface{}) {
//	go func() {
//		msg := ""
//		for _, m := range msgs {
//			v := reflect.ValueOf(m)
//			switch v.Kind() {
//			case reflect.Struct, reflect.Map, reflect.Slice:
//				if bs, err := json.Marshal(m); err == nil {
//					msg += string(bs)
//					continue
//				}
//			}
//			msg += fmt.Sprintf("%+v", m)
//		}
//		g.statusLog = msg
//	}()
//}

// AddTask ...
func (g *Global) AddTask(group otto.Value, fn otto.Value, args ...interface{}) bool {
	if g.running {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "AddTask(), tasks are running")
		return false
	}
	if !group.IsString() {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "AddTask(), Invalid group name")
		return false
	}
	if !fn.IsString() {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "AddTask(), Invalid function name")
		return false
	}
	if _, ok := g.tasks[group.String()]; !ok {
		g.tasks[group.String()] = []task{}
	}
	t := task{ctx: g.ctx.Copy(), fn: fn, args: args}
	t.ctx.Interrupt = make(chan func(), 1)
	g.tasks[group.String()] = append(g.tasks[group.String()], t)
	return true
}

// BindTaskParam ...
func (g *Global) BindTaskParam(group otto.Value, fn otto.Value, args ...interface{}) bool {
	if g.running {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "BindTaskParam(), tasks are running")
		return false
	}
	if !group.IsString() {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "BindTaskParam(), Invalid group name")
		return false
	}
	if !fn.IsString() {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "BindTaskParam(), Invalid function name")
		return false
	}
	if _, ok := g.tasks[group.String()]; !ok {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "BindTaskParam(), group not exist")
		return false
	}
	ts := g.tasks[group.String()]
	for i := 0; i < len(ts); i++ {
		t := &ts[i]
		if t.fn.String() == fn.String() {
			t.args = args
			return true
		}
	}
	g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "BindTaskParam(), function not exist")
	return false
}

// ExecTasks ...
func (g *Global) ExecTasks(group otto.Value) (results []interface{}) {
	if !group.IsString() {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "ExecTasks(), Invalid group name")
		return
	}
	if _, ok := g.tasks[group.String()]; !ok {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "ExecTasks(), group not exist")
		return
	}
	if g.running {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "ExecTasks(), tasks are running")
		return
	}
	g.running = true
	ts := g.tasks[group.String()]
	for range ts {
		results = append(results, false)
	}
	wg := sync.WaitGroup{}
	for i, t := range ts {
		wg.Add(1)
		go func(i int, t task) {
			if f, err := t.ctx.Get(t.fn.String()); err != nil || !f.IsFunction() {
				g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Can not get the task function")
			} else {
				result, err := f.Call(f, t.args...)
				if err != nil || result.IsUndefined() || result.IsNull() {
					results[i] = false
				} else {
					results[i] = result
				}
			}
			wg.Done()
		}(i, t)
	}
	wg.Wait()
	g.running = false
	return
}
