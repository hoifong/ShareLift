package Lift

import (
	"errors"
	"strconv"
)

//	错误
const (
	ErrorLift = "lift append a error"
	ErrorLiftOverweight = ErrorLift + ":overweight"
	ErrorLiftEmpty = ErrorLift + ":empty"
	ErrorLiftPersonNotFound = ErrorLift + ":person not found"
)

//	电梯状态
const (
	//	关机
	StatusLiftShutdown = iota
	//	静止
	StatusLiftStatic
	//	向下
	StatusLiftDown
	//	向上
	StatusLiftUp
)

type Lift struct {
	//	顶层
	top int
	//	底层
	bottom int
	//	当前层
	Level int
	//	速度：移动一层需要的时间,单位（秒）
	speed float32
	//	容量 单位（个人）
	capacity int
	//	电梯内所有人，键值为person id
	persons map[int]*Person
	//	状态
	status int
	//	电梯内按键状态
	levelPress map[int]bool
	//	电梯外向下键按键状态
	downPress map[int]bool
	//	电梯外向上键按键状态
	upPress map[int]bool

	//	电梯运行程序channel
	pressLevel chan int
	pressDown chan int
	pressUp chan int
	stop chan bool
}

func NewLift(capacity int, level int, top int, bottom int) *Lift {
	return &Lift{
		top,
		bottom,
		level,
		2.0,
		capacity,
		make(map[int]*Person),
		StatusLiftShutdown,
		make(map[int]bool),
		make(map[int]bool),
		make(map[int]bool),
		make(chan int, 1),
		make(chan int, 1),
		make(chan int, 1),
		make(chan bool, 1),
	}
}

//	电梯上人
func (lift *Lift) AddPerson(person *Person) error {
	if len(lift.persons) < lift.capacity {
		lift.persons[person.id] = person
		return nil
	}
	return errors.New(ErrorLiftOverweight)
}

//	电梯下人
func (lift *Lift) RemovePersonById(id int) error {
	if len(lift.persons) == 0 {
		return errors.New(ErrorLiftEmpty)
	}
	_, ok := lift.persons[id]
	if ok {
		delete(lift.persons, id)
		return nil
	}
	return errors.New(ErrorLiftPersonNotFound+" by id " + strconv.Itoa(id))
}

//	按电梯内的数字键
func (lift *Lift) PressDown(level int) {
	if lift.status == StatusLiftShutdown {
		return
	}
	lift.levelPress[level] = true
}

//	按电梯外的向上键
func (lift *Lift) PressUp(level int) {
	if level == lift.top || lift.status == StatusLiftShutdown {
		return
	}
	lift.upPress[level] = true
}


//	按电梯外的向下键
func (lift *Lift) PressLevel(level int) {
	if level == lift.bottom || lift.status == StatusLiftShutdown {
		return
	}
	lift.downPress[level] = true
}

//	运行
func (lift *Lift) Run() error {
	go func() {
		for {
			select {
				case <-lift.pressLevel:
				//	按电梯内的数字键
				case <-lift.pressUp:
				//	按电梯外的向上键
				case <-lift.pressDown:
				//	按电梯外的向下键
				case <-lift.stop:
					return
				//	停止
			}
		}
	}()
	return nil
}

func (lift *Lift) Stop() error {
	lift.stop <- true
	return nil
}