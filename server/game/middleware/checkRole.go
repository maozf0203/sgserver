package middleware

import (
	"sgserver/constant"
	"sgserver/net"
)

func CheckRole() net.MiddlewareFunc {
	return func(next net.Handlerfunc) net.Handlerfunc {
		return func(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
			_, err := req.Conn.GetProperty("role")
			if err != nil {
				rsp.Body.Code = constant.SessionInvalid
				return
			}
			next(req, rsp)
		}

	}

}
