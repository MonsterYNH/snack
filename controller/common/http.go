package common

import "log"

const (
	SUCCESS = 0
	// 参数错误 1XXX
	PARAMETER_ERR = 1001
	// 响应错误 2XXX
	PLEASE_LOGIN       = 2001
	USER_NOT_EXIST     = 2002
	USER_EXIST         = 2003
	PLEASE_AGREEMENT   = 2004
	DIFFERENT_PASSWORD = 2005
	ID_NOT_EXIST       = 2006
	// 服务器内部错误
	SERVER_ERROR   = 3001
	MARSHAML_ERROR = 3002
	// 权限
	NO_AUTH       = 4001
	TOKEN_EXPIRED = 4002
	TOKEN_INVALID = 4003
)

var codeMessageMap map[int]string

func init() {
	codeMessageMap = make(map[int]string)
	codeMessageMap[SUCCESS] = "success"
	// 参数错误 1XXX
	codeMessageMap[PARAMETER_ERR] = "parameter err"
	// 响应错误 2XXX
	codeMessageMap[PLEASE_LOGIN] = "please login"
	codeMessageMap[USER_NOT_EXIST] = "user is not exist"
	codeMessageMap[USER_EXIST] = "user exist"
	codeMessageMap[PLEASE_AGREEMENT] = "please agreement"
	codeMessageMap[DIFFERENT_PASSWORD] = "different password"
	codeMessageMap[ID_NOT_EXIST] = "id not exist"
	// 服务器内部错误
	codeMessageMap[SERVER_ERROR] = "server err"
	codeMessageMap[MARSHAML_ERROR] = "marshal err"
	// 权限
	codeMessageMap[NO_AUTH] = "no auth"
	codeMessageMap[TOKEN_EXPIRED] = "Token is expired"
	codeMessageMap[TOKEN_INVALID] = "Token is invalid"

}

type CommonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Token   string      `json:"token,omitempty"`
}

func ResponseSuccess(data, token interface{}) CommonResponse {
	response := CommonResponse{
		Code:    SUCCESS,
		Message: "success",
		Data:    data,
	}
	if token != nil {
		response.Token = token.(string)
	}
	return response
}

func ResponseError(code int) CommonResponse {
	message, exist := codeMessageMap[code]
	if !exist {
		message = "Unknow message"
		log.Println("Error: unknow message code: ", code)
	}
	return CommonResponse{
		Code:    code,
		Message: message,
	}
}
