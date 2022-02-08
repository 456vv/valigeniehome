package valigeniehome

import (
	"encoding/json"
	"log"
	"io"
	"bytes"
)

type ResponseProperties struct{
	Name				string					`json:"name"`		//1,名称
	Value				string					`json:"value"`		//1,状态
}
type ResponsePayloadDevice struct{
	DeviceId			string								`json:"deviceId"`			//1,设备id
	DeviceName			string								`json:"deviceName"`			//1,设备名称,https://open.bot.tmall.com/oauth/api/aliaslist
	DeviceType			string								`json:"deviceType"`			//1,设备类型,具体参考AliGenie支持的品类列表
	Zone				string								`json:"zone,omitempty"`		//0,位置, https://open.bot.tmall.com/oauth/api/placelist
    Brand				string								`json:"brand"`				//1,品牌
	Model				string								`json:"model"`				//1,型号
	Icon				string								`json:"icon"`				//1,产品icon(https协议的url链接),像素最好160*160 以免在app显示模糊

	//[{"name": "powerstate","value": "off"}]
	Properties			[]ResponseProperties				`json:"properties"`				//1,性能，返回当前设备支持的属性状态列表，产品支持的属性列表参考 设备控制与设备状态查询页

	//TurnOn和TurnOf
	Actions				[]string							`json:"actions"`				//1，开关值，产品支持的操作

	//parentId这个是key,存放着父设备ID
	Extensions			map[string]string					`json:"extensions,omitempty"`	//0，扩展，产品扩展属性,为空返回null或者不返回该字段
}

//产品信息
//	deviceName			设备名称
//	brand, model, icon ,zone string	品牌，型号，图标，位置
func (T *ResponsePayloadDevice) Info(deviceName, brand, model, icon, zone string){
	T.DeviceName = deviceName
	T.Brand	= brand
	T.Model	= model
	T.Icon	= icon
	T.Zone	= zone
}

//开关
//	deviceId			设备id
//	powerState string	电源状态 on(打开)，off(关闭)
func (T *ResponsePayloadDevice) Switch(deviceId, powerState string){
	T.DeviceId	 = deviceId
	T.DeviceType = "switch"
	//当前设备状态
	T.Properties = []ResponseProperties{
		{Name:"powerstate", Value:powerState},
	}
	//设备支持行为
	if T.Actions == nil {
		T.Actions = []string{
			"TurnOn",
			"TurnOff",
			"Query",
		}
	}
}

//插座
//	deviceId			设备id
//	powerState string	电源状态 on(打开)，off(关闭)
func (T *ResponsePayloadDevice) Outlet(deviceId, powerState string){
	T.DeviceId	 = deviceId
	T.DeviceType = "outlet"
	//当前设备状态
	T.Properties = []ResponseProperties{
		{Name:"powerstate", Value:powerState},
	}
	//设备支持行为
	if T.Actions == nil {
		T.Actions = []string{
			"TurnOn",
			"TurnOff",
			"Query",
		}
	}
}

//灯
//	deviceId			设备id
//	powerState			电源状态 on(打开)，off(关闭)
//	color				颜色 参考颜色对应表 http://doc-bot.tmall.com/docs/doc.htm?spm=0.0.0.0.wzijJu&treeId=393&articleId=108268&docType=1
//	brightness			亮度 数值
//	colorTemperature	色温 数值
func (T *ResponsePayloadDevice) Light(deviceId, powerState, color, brightness, colorTemperature string){
	T.DeviceId	 = deviceId
	T.DeviceType = "light"
	//当前设备状态
	T.Properties = []ResponseProperties{
		{Name:"powerstate", Value:powerState},
		{Name:"color", Value:color},
		{Name:"brightness", Value:brightness},
		{Name:"colorTemperature", Value:colorTemperature},
	}
	//设备支持行为
	if T.Actions == nil {
		T.Actions = []string{
			"TurnOn",
			"TurnOff",
			"SetColor",
			"SetBrightness",
			"AdjustUpBrightness",
			"AdjustDownBrightness",
			"SetTemperature",
			"Query",
		}
	}
}


type ResponsePayload struct{
	//产品与设备列表
	Devices				[]*ResponsePayloadDevice	`json:"devices,omitempty"`		//1，用户设备列表
	
	//控制设备
	DeviceId			string						`json:"deviceId,omitempty"`		//1,设备ID
	//错误
    ErrorCode			string						`json:"errorCode,omitempty"`	//0,错误代码
    Message				string						`json:"message,omitempty"`		//0,描述
}
type Response struct{
	req			*Request
	Properties	[]*ResponseProperties	`json:"properties,omitempty"`
	Header		Header					`json:"header"`
	Payload		ResponsePayload			`json:"payload"`

}

//错误
//errStr, code string		文本，错误代码
func (T *Response) Error(errStr, code string) {
	//错误码 errorCode				 	错误码说明			对应message
	//INVALIDATE_PARAMS				 400 请求参数有误		invalidate params
	//ACCESS_TOKEN_INVALIDATE		 401 access_token 无效	access_token is invalidate
	//DEVICE_IS_NOT_EXIST			 404 设备未找到			device is not exist
	//INVALIDATE_CONTROL_ORDER		 417 控制指令不正确		invalidate control order
	//SERVICE_ERROR					 500 服务异常			服务错误原因（方便观察原因）
	//DEVICE_NOT_SUPPORT_FUNCTION	 501 设备不支持该操作	device not support
	//IOT_DEVICE_OFFLINE			 504 设备离线状态		device is offline
	var err string
	switch code {
	case "400":
		err	= "INVALIDATE_PARAMS"
	case "401":
		err = "ACCESS_TOKEN_INVALIDATE"
	case "404":
		err = "DEVICE_NOT_SUPPORT_FUNCTION"
	case "417":
		err = "INVALIDATE_CONTROL_ORDER"
	case "500","":
		err	= "SERVICE_ERROR"
	case "501":
		err	= "DEVICE_NOT_SUPPORT_FUNCTION"
	case "504":
		err	= "IOT_DEVICE_OFFLINE"
	default:
		err = code
	}
	
	T.Header		= T.req.Header
	T.Header.Name	= "ErrorResponse"
	
	T.Payload.DeviceId 	= T.req.Payload.DeviceId
	T.Payload.ErrorCode	= err
	T.Payload.Message	= errStr
}


func (T *Response) header(){
	T.Header		= T.req.Header
	T.Header.Name	= 	T.Header.Name+"Response"
}

//设备发现
//	product *ResponsePayloadDevice		分类
//	devices *ResponsePayloadDevice...	产品
func (T *Response) Discovery(product *ResponsePayloadDevice, devices ...*ResponsePayloadDevice) {
	T.header()
	//续承
	if product != nil {
		//适用插座
		for _ , d := range devices {
			if d.DeviceName == "" {
				d.DeviceName = product.DeviceName
			}
			if d.DeviceType == "" {
				d.DeviceType = product.DeviceType
			}
			if d.Brand == "" {
				d.Brand = product.Brand
			}
			if d.Model == "" {
				d.Model = product.Model
			}
			if d.Icon == "" {
				d.Icon = product.Icon
			}
			if d.Extensions == nil {
				d.Extensions = make(map[string]string)
			}
			d.Extensions["parentId"] = product.DeviceId
		}
		T.Payload.Devices	= append(T.Payload.Devices, append([]*ResponsePayloadDevice{product}, devices...)...)
	}else{
		//其它
		T.Payload.Devices	= append(T.Payload.Devices, devices...)
	}
}

//设备控制
func (T *Response) Control(){
	T.header()
	T.Payload.DeviceId 	= T.req.Payload.DeviceId
}

//设备状态查询
//	properties []ResponseProperties	属性值
func (T *Response) Query(properties []*ResponseProperties){
	T.Properties = properties
	T.header()
	T.Payload.DeviceId 	= T.req.Payload.DeviceId
}

//写入到w
//w io.Writer	写入
func (T *Response) WriteTo(w io.Writer) (n int64, err error) {
	if T.Header.Name == "" {
		T.Error("内部错误，这是程序员忘了调用回复接口造成的。请联系官方修复。谢谢！", "500")
	}
	buf := bytes.NewBuffer(nil)
	buf.Grow(1024)
	err = json.NewEncoder(buf).Encode(T)
	if err != nil {
		w.Write([]byte(`{"header":{"namespace":"?","name":"ErrorResponse","messageId":"?","payLoadVersion":1}, "payload":{"deviceId":"?","errorCode":"SERVICE_ERROR","message":"致命的错误"}}`))
		log.Println(err)
		return 0, err
	}
	return buf.WriteTo(w)
}















