package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Grifonhard/Practicum-s5_6/internal/lib/validate"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Grifonhard/Practicum-s5_6/internal/accrual/model"
	errs "github.com/Grifonhard/Practicum-s5_6/internal/lib/errors"
)

type orderRegistrationRequest struct {
	Order string       `json:"order"`
	Goods []model.Good `json:"goods"`
}

func (h *Handler) OrderRegistrationHandler(c *gin.Context) {
	var req orderRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequestResponse(c, err)
		return
	}

	ctx := c.Request.Context()

	if ok := validate.CheckLuhn(req.Order); !ok {
		err := errors.New("invalid order")
		badRequestResponse(c, err)
		return
	}

	number, err := strconv.ParseUint(req.Order, 10, 64)
	if err != nil {
		slog.ErrorContext(ctx, "order registration handle", "err", err)
		badRequestResponse(c, err)
		return
	}

	err = h.OrderService.RegisterOrder(ctx, number, req.Goods)

	if err != nil {
		slog.ErrorContext(ctx, "order registration handle", "err", err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == errs.ErrPostgresUniqueViolation {
			c.JSON(http.StatusConflict, gin.H{
				"error": errs.NewConflict(err),
			})
			return
		}

		c.JSON(errs.Status(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Printf("order num registed: %d\n", number)

	c.JSON(http.StatusAccepted, gin.H{
		"number": req.Order,
	})
}
