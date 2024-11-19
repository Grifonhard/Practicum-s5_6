package accrual

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/model"
)

type Manager struct {
	uri string
}

func New(uri string) (*Manager, error) {
	var m Manager
	m.uri = uri
	return &m, nil
}

func (m *Manager) AccrualReq(orderID int) (*model.OrderAccrual, error) {
	url := fmt.Sprintf("%s/api/orders/%s", m.uri, strconv.Itoa(orderID))

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(1)
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// 200 OK - успешный запрос
		var result model.OrderAccrual
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}
		return &result, nil

	case http.StatusNoContent:
		// 204 - заказ не зарегистрирован
		return nil, ErrStatusNoContent

	case http.StatusTooManyRequests:
		// 429 - превышено количество запросов
		retryAfter := resp.Header.Get("Retry-After")
		return nil, fmt.Errorf("%w, retry after %s seconds", ErrTooManyRequests, retryAfter)

	case http.StatusInternalServerError:
		// 500 - внутренняя ошибка сервера
		return nil, ErrStatusInternalServerError

	default:
		// Неизвестный статус ответа
		return nil, fmt.Errorf("%w: %d", ErrUnexpected, resp.StatusCode)
	}

}
