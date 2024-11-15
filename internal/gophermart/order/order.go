package order

import (
	"sort"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/drivers/accrual"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/drivers/psql"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/storage"
	"github.com/Grifonhard/Practicum-s5_6/internal/model"
)

// TODO запись в логи при возникновении ошибок

type Manager struct {
	s *storage.Storage
	a *accrual.Manager
	p *psql.DB
}

func New(p *psql.DB, stor *storage.Storage, acm *accrual.Manager) (*Manager, error) {
	var m Manager
	m.s = stor
	m.a = acm
	m.p = p 
	return &m, nil
}

func (m *Manager) AddOrder(username string, orderID int) error {
	err := checkLuhn(orderID)
	if err != nil {
		return err
	}
	err = m.s.NewOrder(username, orderID)
	if err != nil {
		return err
	}		
	return nil
}

func (m *Manager) ListOrders(username string) ([]model.OrderFront, error) {
	orders, err := m.s.GetOrders(username)
	if err != nil {
		return nil, err
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].Created.Before(orders[j].Created)
	})

	// часть получения инфы о заказах, по которым ещё нет данных
	// собираем недостающую инфу
	var ordersFront []model.OrderFront
	for _, o := range orders {
		order, err := m.checkAndConverOrder(&o)
		if err != nil {
			// TODO может просто писать в логи, а поломанные игнорить или менять статус?
			return nil, err
		}
		ordersFront = append(ordersFront, *order)
	}

	return ordersFront, err
}

func (m *Manager) Balance(username string) (int, error) {

}

func (m *Manager) Withdraw(username, order string, sum int) error {

}

func (m *Manager) Withdrawls(username string) error {
	
}

func (m *Manager) checkAndConverOrder(o *model.Order) (*model.OrderFront, error) {
	var orderFront model.OrderFront
	var accrual int
	var err error
	var newO *model.Order
	if o.Status != model.PROCESSED || o.Status != model.INVALID {
		newO, accrual, err = m.updateOrderInfo(o)
		if err != nil {
			return nil, err
		}
		o = newO
	}
	err = orderFront.ConvertOrder(o, accrual)
	if err != nil {
		return nil, err
	}
	return &orderFront, nil
}

func (m *Manager) updateOrderInfo(o *model.Order) (*model.Order, int, error) {
	var newOrder model.Order
	var accrual, status int
	var isUpdate bool

	info, err := m.a.AccrualReq(o.Id)
	if err != nil {
		return nil, 0, err 
	}

	if info.Status != o.Status {
		accrual, status, err = newOrder.ConvertAccrual(info)
		if err != nil {
			return nil, 0, err
		}
		err = m.p.UpdateOrderStatus(o.Id, status)
		if err != nil {
			return nil, 0, err
		}
		isUpdate = true
	}

	if info.Status == model.PROCESSED && isUpdate {
		err = m.p.InsertBalanceTransaction(o.UserId, o.Id, accrual)
		if err != nil {
			return nil, 0, err
		}
	}

	return &newOrder, accrual, nil
}

func checkLuhn(orderId int) error {
	var sum int
	shouldDouble := false

	for orderId > 0 {
		digit := orderId % 10


		if shouldDouble {
			digit <<= 1
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		shouldDouble = !shouldDouble
		orderId /= 10
	}

	if sum%10 == 0 {
		return nil
	}

	return ErrLuhnFail
}