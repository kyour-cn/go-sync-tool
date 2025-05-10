package task

import "golang.org/x/exp/slog"

func NewOrderOutbound() *OrderOutbound {
    return &OrderOutbound{}
}

// OrderOutbound 同步ERP订单出库到商城
type OrderOutbound struct{}

func (g OrderOutbound) GetName() string {
    return "orderOutbound"
}

func (g OrderOutbound) Run(t *Task) error {

    // TODO: 待实现

    slog.Debug("同步同步ERP订单出库到商城")

    return nil
}
