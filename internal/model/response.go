package model

import "strconv"

type Code int

const (
	OK  Code = 0
	ERR Code = -1
)

func (c Code) String() string {
	switch c {
	case OK:
		return "OK"
	case ERR:
		return "ERR"
	default:
		return "Code(" + strconv.FormatInt(int64(c), 10) + ")"
	}
}

type Response struct {
	Code    Code        `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseOK(data interface{}) Response {
	return Response{
		OK,
		OK.String(),
		data,
	}
}

func ResponseERR(msg string, data interface{}) Response {
	return Response{
		ERR,
		ERR.String() + ": " + msg,
		data,
	}
}
