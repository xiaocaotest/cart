package handler

import (
	"cart/common"
	"cart/domain/model"
	"cart/domain/service"
	cart "cart/proto"
	"context"
)

type Cart struct {
	CartDataService service.ICartDataService
}

func (h *Cart) AddCart(ctx context.Context, request *cart.CartInfo, response *cart.ResponseAdd) (err error) {
	cart := &model.Cart{}
	common.SwapTo(request, cart)
	response.CartId, err = h.CartDataService.AddCart(cart)
	return err
}

func (h *Cart) CleanCart(ctx context.Context, request *cart.Clean, response *cart.Response) error {
	if err := h.CartDataService.CleanCart(request.UserId); err != nil {
		return err
	}
	response.Meg = "cart nil fail"
	return nil
}

func (h *Cart) Incr(ctx context.Context, request *cart.Item, response *cart.Response) error {
	if err := h.CartDataService.IncrNum(request.Id, request.ChangeNum); err != nil {
		return err
	}
	response.Meg = "cart add susess"
	return nil
}

func (h *Cart) Decr(ctx context.Context, request *cart.Item, response *cart.Response) error {
	if err := h.CartDataService.DecrNum(request.Id, request.ChangeNum); err != nil {
		return err
	}
	response.Meg = "- fail"
	return nil
}

func (h *Cart) DeleteItemByID(ctx context.Context, request *cart.CartID, response *cart.Response) error {
	if err := h.CartDataService.DeleteCart(request.Id); err != nil {
		return err
	}
	response.Meg = "add success"
	return nil
}

func (h *Cart) GetAll(ctx context.Context, request *cart.CartFindAll, response *cart.CartAll) error {
	cartAll, err := h.CartDataService.FindAllCart(request.UserId)
	if err != nil {
		return err
	}
	for _, v := range cartAll {
		cart := &cart.CartInfo{}
		if err := common.SwapTo(v, cart); err != nil {
			return err
		}
		response.CartInfo = append(response.CartInfo)
	}
	return nil
}
