package valigeniehome


type Responser interface{
	Version() int
	V1() *Response1
	V2() *Response2
}

type resSynthesis struct{
	version int
	v1 *Response1
	v2 *Response2
}

func (T *resSynthesis) Version() int {
	return T.version
}

func (T *resSynthesis) V1() *Response1 {
	return T.v1
}

func (T *resSynthesis) V2() *Response2 {
	return T.v2
}

type SCode int
const (
	FAIL						SCode	= iota	//未知
	SUCCESS											//成功
	INVALIDATE_PARAMS								//请求参数有误
	ACCESS_TOKEN_INVALIDATE							//access_token 无效
	DEVICE_IS_NOT_EXIST								//设备未找到
	INVALIDATE_CONTROL_ORDER						//控制指令不正确
	SERVICE_ERROR									//服务异常
	DEVICE_NOT_SUPPORT_FUNCTION						//设备不支持此功能
	IOT_DEVICE_OFFLINE								//设备已离线
)
var sCodeName = map[SCode]string{
	FAIL						: "FAIL",
  	SUCCESS						: "SUCCESS",
  	INVALIDATE_PARAMS			: "INVALIDATE_PARAMS",
  	ACCESS_TOKEN_INVALIDATE		: "ACCESS_TOKEN_INVALIDATE",
  	DEVICE_IS_NOT_EXIST			: "DEVICE_IS_NOT_EXIST",
  	INVALIDATE_CONTROL_ORDER	: "INVALIDATE_CONTROL_ORDER",
  	SERVICE_ERROR				: "SERVICE_ERROR",
  	DEVICE_NOT_SUPPORT_FUNCTION	: "DEVICE_NOT_SUPPORT_FUNCTION",
  	IOT_DEVICE_OFFLINE			: "IOT_DEVICE_OFFLINE",
}
func (T SCode) Name() string {
	return sCodeName[T]
}
var sCodeMessage = map[SCode]string{
	FAIL						: "unknown error",
  	SUCCESS						: "SUCCESS",
  	INVALIDATE_PARAMS			: "invalidate params",
  	ACCESS_TOKEN_INVALIDATE		: "access_token is invalidate",
  	DEVICE_IS_NOT_EXIST			: "device is not exist",
  	INVALIDATE_CONTROL_ORDER	: "invalidate control order",
  	SERVICE_ERROR				: "Service error",
  	DEVICE_NOT_SUPPORT_FUNCTION	: "device not support",
  	IOT_DEVICE_OFFLINE			: "device is offline",
}
func (T SCode) String() string {
	return sCodeMessage[T]
}

