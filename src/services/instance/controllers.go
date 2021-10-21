package instance

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/linshenqi/sptty"
	"github.com/linshenqi/taskmate/src/base"
)

func (s *Service) routePostInstances(ctx iris.Context) {
	req := base.Instance{}

	if err := ctx.ReadJSON(&req); err != nil {
		_ = sptty.SimpleResponse(ctx, iris.StatusBadRequest, sptty.RequestError{
			Code: base.CurrentFuncName(),
			Msg:  err.Error(),
		})

		return
	}

	if err := s.CreateInstance(&req); err != nil {
		_ = sptty.SimpleResponse(ctx, iris.StatusBadRequest, sptty.RequestError{
			Code: base.CurrentFuncName(),
			Msg:  err.Error(),
		})

		return
	}

	_ = sptty.SimpleResponse(ctx, iris.StatusCreated, req)
}

func (s *Service) routeGetInstances(ctx iris.Context) {
	query := base.CreateQueryFromContext(&base.QueryInstance{}, s.db, ctx).(*base.QueryInstance)
	instances, total, err := s.ListInstances(query)
	if err != nil {
		_ = sptty.SimpleResponse(ctx, iris.StatusBadRequest, sptty.RequestError{
			Code: base.CurrentFuncName(),
			Msg:  err.Error(),
		})

		return
	}

	_ = sptty.SimpleResponse(ctx, iris.StatusOK, instances,
		map[string]string{
			"Access-Control-Expose-Headers": base.HeaderTotal,
			base.HeaderTotal:                fmt.Sprintf("%d", total),
		})
}

func (s *Service) routeGetInstancesByID(ctx iris.Context) {
	instance, err := s.GetInstanceByID(ctx.Params().Get("id"))
	if err != nil {
		_ = sptty.SimpleResponse(ctx, iris.StatusBadRequest, sptty.RequestError{
			Code: base.CurrentFuncName(),
			Msg:  err.Error(),
		})

		return
	}

	_ = sptty.SimpleResponse(ctx, iris.StatusOK, instance)
}

func (s *Service) routePutInstancesRemoval(ctx iris.Context) {
	ids := []string{}
	if err := ctx.ReadJSON(&ids); err != nil {
		_ = sptty.SimpleResponse(ctx, iris.StatusBadRequest, sptty.RequestError{
			Code: base.CurrentFuncName(),
			Msg:  err.Error(),
		})

		return
	}

	if err := s.DeleteInstances(ids); err != nil {
		_ = sptty.SimpleResponse(ctx, iris.StatusBadRequest, sptty.RequestError{
			Code: base.CurrentFuncName(),
			Msg:  err.Error(),
		})

		return
	}

	_ = sptty.SimpleResponse(ctx, iris.StatusOK, nil)
}
