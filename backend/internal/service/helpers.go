package service

import (
	"errors"
	"net/http"

	"github.com/rxcheck/rxcheck/internal/constants"
	"github.com/rxcheck/rxcheck/internal/util"
)

// wrapSvcErr 统一包装 service 层错误，保留 AppError 链（handler 再次包装）。
func wrapSvcErr(err error) error {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return util.NewAppError(http.StatusInternalServerError, constants.MsgInternalError, err)
}
