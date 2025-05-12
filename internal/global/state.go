package global

var State = struct {
	Status     int  // 0: 初始化中 1: 待启动 2: 启动中 3: 运行中 4: 停止中
	HideWindow bool // 是否隐藏窗口
}{
	Status:     0,
	HideWindow: false,
}
