package util

import "sync"

var (
	available   bool
	availRWLock *sync.RWMutex
)

// 初始化状态变量，并标识为可用
func InitState() {
	availRWLock = &sync.RWMutex{}
	SetAvailable()
}

// 用于保护一个请求，只能用在中间件中
func UseServerStart() bool {
	if !available {
		return false
	}

	availRWLock.RLock()

	// 等待期间有可能服务器不可用了
	return available
}

// 必须和 UseServerStart 对偶使用
func UseServerEnd() {
	availRWLock.RUnlock()
}

func ChangeServerStart() {
	availRWLock.Lock()
}

func ChangeServerEnd() {
	availRWLock.Unlock()
}

func IsAvailabale() bool {
	return available
}

func SetUnavailable() {
	available = false
}

func SetAvailable() {
	available = true
}
