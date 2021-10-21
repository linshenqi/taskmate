package task

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/linshenqi/sptty"
	"github.com/linshenqi/taskmate/src/base"
)

func (s *Service) routeGetPythonVersions(ctx iris.Context) {
	versions, err := base.ListAvailablePythonVersions()
	if err != nil {
		_ = sptty.SimpleResponse(ctx, iris.StatusBadRequest, sptty.RequestError{
			Code: base.CurrentFuncName(),
			Msg:  err.Error(),
		})

		return
	}

	_ = sptty.SimpleResponse(ctx, iris.StatusOK, versions)
}

func (s *Service) routeGetEnvs(ctx iris.Context) {
	envs, err := base.ListEnvs()
	if err != nil {
		_ = sptty.SimpleResponse(ctx, iris.StatusBadRequest, sptty.RequestError{
			Code: base.CurrentFuncName(),
			Msg:  err.Error(),
		})

		return
	}

	_ = sptty.SimpleResponse(ctx, iris.StatusOK, envs)
}

func (s *Service) routePostEnvs(ctx iris.Context) {
	req := base.PythonEnv{}

	if err := ctx.ReadJSON(&req); err != nil {
		_ = sptty.SimpleResponse(ctx, iris.StatusBadRequest, sptty.RequestError{
			Code: base.CurrentFuncName(),
			Msg:  err.Error(),
		})

		return
	}

	if err := base.CreateEnv(&req); err != nil {
		_ = sptty.SimpleResponse(ctx, iris.StatusBadRequest, sptty.RequestError{
			Code: base.CurrentFuncName(),
			Msg:  err.Error(),
		})

		return
	}

	_ = sptty.SimpleResponse(ctx, iris.StatusCreated, req)
}

func (s *Service) routePostEnvsScript(ctx iris.Context) {
	req := base.PythonEnv{}

	if err := ctx.ReadJSON(&req); err != nil {
		_ = sptty.SimpleResponse(ctx, iris.StatusBadRequest, sptty.RequestError{
			Code: base.CurrentFuncName(),
			Msg:  err.Error(),
		})

		return
	}

	if err := base.ExecScript(&req); err != nil {
		_ = sptty.SimpleResponse(ctx, iris.StatusBadRequest, sptty.RequestError{
			Code: base.CurrentFuncName(),
			Msg:  err.Error(),
		})

		return
	}

	_ = sptty.SimpleResponse(ctx, iris.StatusCreated, req)
}

func (s *Service) routePutEnvsRemoval(ctx iris.Context) {
	req := []string{}

	if err := ctx.ReadJSON(&req); err != nil {
		_ = sptty.SimpleResponse(ctx, iris.StatusBadRequest, sptty.RequestError{
			Code: base.CurrentFuncName(),
			Msg:  err.Error(),
		})

		return
	}

	if err := base.RemoveEnvs(req); err != nil {
		_ = sptty.SimpleResponse(ctx, iris.StatusBadRequest, sptty.RequestError{
			Code: base.CurrentFuncName(),
			Msg:  err.Error(),
		})

		return
	}

	_ = sptty.SimpleResponse(ctx, iris.StatusOK, req)
}

func (s *Service) routeGetTasks(ctx iris.Context) {
	query := base.CreateQueryFromContext(&base.QueryTask{}, s.db, ctx).(*base.QueryTask)
	tasks, total, err := s.ListTasks(query)
	if err != nil {
		_ = sptty.SimpleResponse(ctx, iris.StatusBadRequest, sptty.RequestError{
			Code: base.CurrentFuncName(),
			Msg:  err.Error(),
		})

		return
	}

	_ = sptty.SimpleResponse(ctx, iris.StatusOK, tasks,
		map[string]string{
			"Access-Control-Expose-Headers": base.HeaderTotal,
			base.HeaderTotal:                fmt.Sprintf("%d", total),
		})
}

func (s *Service) routePostTasks(ctx iris.Context) {
	req := base.Task{}
	if err := ctx.ReadJSON(&req); err != nil {
		_ = sptty.SimpleResponse(ctx, iris.StatusBadRequest, sptty.RequestError{
			Code: base.CurrentFuncName(),
			Msg:  err.Error(),
		})

		return
	}

	if err := s.CreateTask(&req); err != nil {
		_ = sptty.SimpleResponse(ctx, iris.StatusBadRequest, sptty.RequestError{
			Code: base.CurrentFuncName(),
			Msg:  err.Error(),
		})

		return
	}

	_ = sptty.SimpleResponse(ctx, iris.StatusCreated, req)
}

func (s *Service) routePutTasksRemoval(ctx iris.Context) {
	req := []string{}
	if err := ctx.ReadJSON(&req); err != nil {
		_ = sptty.SimpleResponse(ctx, iris.StatusBadRequest, sptty.RequestError{
			Code: base.CurrentFuncName(),
			Msg:  err.Error(),
		})

		return
	}

	if err := s.dbDeleteTasksByIDs(req); err != nil {
		_ = sptty.SimpleResponse(ctx, iris.StatusBadRequest, sptty.RequestError{
			Code: base.CurrentFuncName(),
			Msg:  err.Error(),
		})

		return
	}

	_ = sptty.SimpleResponse(ctx, iris.StatusOK, nil)
}
