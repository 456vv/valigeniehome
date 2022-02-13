package valigeniehome

import (
	"encoding/json"
	"io"
	"bytes"
)

type ResponseProperties1 struct{
	Name				string					`json:"name"`		//1,名称
	Value				string					`json:"value"`		//1,状态
}
type ResponsePayloadDevice1 struct{
	DeviceId			string								`json:"deviceId"`			//1,设备id
	DeviceName			string								`json:"deviceName"`			//1,设备名称,https://open.bot.tmall.com/oauth/api/aliaslist
	DeviceType			string								`json:"deviceType"`			//1,设备类型,具体参考AliGenie支持的品类列表
	Zone				string								`json:"zone,omitempty"`		//0,位置, https://open.bot.tmall.com/oauth/api/placelist
    Brand				string								`json:"brand"`				//1,品牌
	Model				string								`json:"model"`				//1,型号
	Icon				string								`json:"icon"`				//1,产品icon(https协议的url链接),像素最好160*160 以免在app显示模糊

	//[{"name": "powerstate","value": "off"}]
	Properties			[]ResponseProperties1				`json:"properties"`				//1,性能，返回当前设备支持的属性状态列表，产品支持的属性列表参考 设备控制与设备状态查询页

	//TurnOn和TurnOf
	Actions				[]string							`json:"actions"`				//1，开关值，产品支持的操作

	//parentId这个是key,存放着父设备ID
	Extensions			map[string]interface{}				`json:"extensions,omitempty"`	//0，扩展，产品扩展属性,为空返回null或者不返回该字段
	
}

//产品信息
//	deviceName			设备名称
//	brand, model ,zone, icon string	品牌，型号，图标，位置
func (T *ResponsePayloadDevice1) Info(deviceName, brand, model, zone, icon string){
	T.DeviceName = deviceName
	T.Brand	= brand
	T.Model	= model
	T.Zone	= zone
	T.Icon	= icon
}

//开关
//	deviceId			设备id
//	powerState string	电源状态 on(打开)，off(关闭)
func (T *ResponsePayloadDevice1) Switch(deviceId, powerState string){
	T.DeviceId	 = deviceId
	T.DeviceType = "switch"
	//当前设备状态
	T.Properties = []ResponseProperties1{
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
func (T *ResponsePayloadDevice1) Outlet(deviceId, powerState string){
	T.DeviceId	 = deviceId
	T.DeviceType = "outlet"
	//当前设备状态
	T.Properties = []ResponseProperties1{
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
func (T *ResponsePayloadDevice1) Light(deviceId, powerState, color, brightness, colorTemperature string){
	T.DeviceId	 = deviceId
	T.DeviceType = "light"
	//当前设备状态
	T.Properties = []ResponseProperties1{
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


type ResponsePayload1 struct{
	//产品与设备列表
	Devices				[]*ResponsePayloadDevice1	`json:"devices,omitempty"`		//1，用户设备列表
	
	//控制设备
	DeviceId			string						`json:"deviceId,omitempty"`		//1,设备ID
	//错误
    ErrorCode			string						`json:"errorCode,omitempty"`	//0,错误代码
    Message				string						`json:"message,omitempty"`		//0,描述
}
type Response1 struct{
	aligenie	*Aligenie
	req			*Request1
	Properties	[]*ResponseProperties1	`json:"properties,omitempty"`
	Header		Header1					`json:"header"`
	Payload		ResponsePayload1		`json:"payload"`

}

//错误
//errCode SCode		文本，错误代码
func (T *Response1) Error(errCode SCode) {
	T.Header		= T.req.Header
	T.Header.Name	= "ErrorResponse"
	
	T.Payload.DeviceId 	= T.req.Payload.DeviceId
	T.Payload.ErrorCode	= errCode.Name()
	T.Payload.Message	= errCode.String()
}

func (T *Response1) header(){
	if T.Header.Name == "" {
		T.Header		= T.req.Header
		T.Header.Name	= 	T.Header.Name+"Response"
	}
}

//设备发现
//	product *ResponsePayloadDevice1		分类
//	devices *ResponsePayloadDevice1...	产品
func (T *Response1) Discovery(product *ResponsePayloadDevice1, devices ...*ResponsePayloadDevice1) {
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
				d.Extensions = make(map[string]interface{})
			}
			d.Extensions["parentId"] = product.DeviceId
		}
		T.Payload.Devices	= append(T.Payload.Devices, append([]*ResponsePayloadDevice1{product}, devices...)...)
	}else{
		//其它
		T.Payload.Devices	= append(T.Payload.Devices, devices...)
	}
}

//设备控制
func (T *Response1) Control(){
	T.Payload.DeviceId 	= T.req.Payload.DeviceId
}

//设备状态查询
//	properties []ResponseProperties1	属性值
func (T *Response1) Query(properties []*ResponseProperties1){
	T.Properties = properties
	T.Payload.DeviceId 	= T.req.Payload.DeviceId
}

//写入到w
//w io.Writer	写入
func (T *Response1) WriteTo(w io.Writer) (n int64, err error) {
	T.header()
	
	buf := bytes.NewBuffer(nil)
	buf.Grow(1024)
	err = json.NewEncoder(buf).Encode(T)
	if err != nil {
		w.Write([]byte(`{"header":{"namespace":"?","name":"ErrorResponse","messageId":"?","payLoadVersion":1}, "payload":{"deviceId":"?","errorCode":"SERVICE_ERROR","message":"Fatal mistake"}}`))
		return 0, err
	}
	
	if T.aligenie.Debug {
		lr := &io.LimitedReader{R:buf, N:int64(buf.Len())}
		r := io.TeeReader(lr, buf)
		b, _ := io.ReadAll(r)
		T.aligenie.Logf("valigeniehome: 响应>%s\n\n", b)
	}
	return buf.WriteTo(w)
}



