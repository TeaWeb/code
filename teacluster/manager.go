package teacluster

import (
	"errors"
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/logs"
	"github.com/vmihailenco/msgpack"
	"io"
	"net"
	"sync"
	"time"
)

var ClusterManager = NewManager()

// cluster communication manager
type Manager struct {
	Context chan bool

	conn    net.Conn
	encoder *msgpack.Encoder
	decoder *msgpack.Decoder

	isStarting  bool
	startLocker sync.Mutex

	prevAction  ActionInterface
	queueLocker sync.Mutex

	error string

	isActive bool
}

func NewManager() *Manager {
	return &Manager{
		Context: make(chan bool),
	}
}

// start manager
func (this *Manager) Start() error {
	this.startLocker.Lock()
	if this.isStarting {
		return nil
	}
	this.isStarting = true
	defer func() {
		this.isStarting = false
		this.startLocker.Unlock()
	}()

	node := teaconfigs.SharedNodeConfig()
	if node == nil {
		return nil
	}

	if !node.On {
		return nil
	}

	if len(node.ClusterAddr) == 0 {
		return errors.New("'clusterAddr' should not be empty")
	}

	conn, err := net.DialTimeout("tcp", node.ClusterAddr, 10*time.Second)
	if err != nil {
		this.error = err.Error()
		return err
	}

	this.isActive = true

	this.conn = conn
	this.encoder = msgpack.NewEncoder(this.conn)
	this.decoder = msgpack.NewDecoder(this.conn)

	// register
	err = this.Write(&RegisterAction{
		ClusterId:     node.ClusterId,
		ClusterSecret: node.ClusterSecret,
		NodeId:        node.Id,
		NodeName:      node.Name,
		NodeRole:      node.Role,
	})
	if err != nil {
		logs.Error(errors.New("fail to register node"))
	}

	this.Read(func(action ActionInterface) {
		if action.Name() == "success" || action.Name() == "fail" {
			if this.prevAction != nil && action.BaseAction().RequestId == this.prevAction.BaseAction().Id {
				switch action.Name() {
				case "success":
					this.error = ""
					this.prevAction.OnSuccess(action.(*SuccessAction))
				case "fail":
					this.error = action.(*FailAction).Message
					this.prevAction.OnFail(action.(*FailAction))
				}
			}
		}
		action.Execute()
	})

	return nil
}

// read action from cluster
func (this *Manager) Read(f func(action ActionInterface)) {
	for {
		typeId, _, err := this.decoder.DecodeExtHeader()
		if err != nil {
			if err == io.EOF {
				break
			}
			this.error = err.Error()
			break
		}
		instance := FindActionInstance(typeId)
		if instance == nil {
			logs.Error(errors.New("can not find action type '" + fmt.Sprintf("%d", typeId) + "'"))
			continue
		}
		err = this.decoder.Decode(instance)
		if err != nil {
			if err == io.EOF {
				break
			}
			this.error = err.Error()
			break
		}
		f(instance)
	}

	this.isActive = false
}

// write action message to cluster manager
func (this *Manager) Write(action ActionInterface) error {
	if this.conn == nil {
		return errors.New("no connection to cluster")
	}
	this.queueLocker.Lock()
	this.prevAction = action
	action.BaseAction().Id = GenerateActionId()
	err := this.encoder.Encode(action)
	this.queueLocker.Unlock()
	return err
}

// stop manager
func (this *Manager) Stop() error {
	conn := this.conn
	if conn != nil {
		err := conn.Close()
		this.conn = nil
		return err
	}
	return nil
}

// is active
func (this *Manager) IsActive() bool {
	return this.isActive
}

func (this *Manager) Error() string {
	return this.error
}

func (this *Manager) Restart() {
	this.Stop()
	this.Context <- true
}
