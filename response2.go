package valigeniehome

import (
	"encoding/json"
	"io"
	"bytes"
)


//v2:
// ｛
//  "header":{
//      "namespace":"AliGenie.Iot.Device.Control",
//      "name":"CorrectResponse",
//      "messageId":"1bd5d003-31b9-476f-ad03-71d471922820",
//      "payLoadVersion":2
//   },
//   "payload":{
//     	"deviceResponseList":[
//        {
//         "deviceId":"devId-34234",
//         "errorCode":"SUCCESS",
//         "message":"SUCCESS"
//     	 },
//       {
//         "deviceId":"devId-abc",
//         "errorCode":"SUCCESS",
//         "message":"SUCCESS"
//     	 },
//        {
//         "deviceId":"devId-1234",
//          //说明devId-1234对应的设备不支持此功能
//         "errorCode":"DEVICE_NOT_SUPPORT_FUNCTION",
//         "message":"device not support"
//     	 },
//       {
//         "deviceId":"devId-5678",
//         //说明devId-5678对应的设备已离线
//         "errorCode":"IOT_DEVICE_OFFLINE",
//         "message":"device is offline"
//     	 }
//      ]
//    }
//  ｝

type ResponseDeviceResponseList2 struct{
	DeviceId			string					`json:"deviceId"`						//1,设备id
	//SUCCESS						成功
	//DEVICE_NOT_SUPPORT_FUNCTION	设备不支持此功能
	//IOT_DEVICE_OFFLINE			设备已离线
	ErrorCode			string					`json:"errorCode"`						//1,错误代码
	Message				string					`json:"message"`						//1,信息
}
type ResponsePayloadDevices2 struct{
	DeviceId			string								`json:"deviceId"`			//1,设备id
	DeviceName			string								`json:"deviceName"`			//1,设备名称,https://open.bot.tmall.com/oauth/api/aliaslist
	//必须是猫精平台支持的品牌名，注意！！！
	DeviceType			string								`json:"deviceType"`			//1,设备类型,https://www.aligenie.com/doc/357554/eq19cg
	//必须是猫精平台支持的品牌名，注意！！！
    Brand				string								`json:"brand"`				//1,品牌
	Model				string								`json:"model"`				//1,型号
	Zone				string								`json:"zone,omitempty"`		//0,位置, https://open.bot.tmall.com/oauth/api/placelist
	
	//"powerstate":1
	//"brightness":30
	Status				map[string]int						`json:"status"`				//1,查询的状态
	
	Extensions			map[string]interface{}				`json:"extensions,omitempty"`				//0，扩展，产品扩展属性,为空返回null或者不返回该字段
}

//产品信息
//	deviceType			设备类型
//	deviceName			设备名称
//	brand, model ,zone string 类型，品牌，型号，位置
func (T *ResponsePayloadDevices2) Info(deviceType, deviceName, brand, model, zone string){
	T.DeviceType = deviceType
	T.DeviceName = deviceName
	T.Brand	= brand
	T.Model	= model
	T.Zone	= zone
}

//设备状态
//	deviceId	string			设备id
//	status 		map[string]int	状态
func (T *ResponsePayloadDevices2) Result(deviceId string, status map[string]int){
	T.DeviceId	= deviceId
	T.Status	= status
}

type ResponsePayload2 struct{
	Devices				[]*ResponsePayloadDevices2			`json:"devices"`							//1,发现设备列表
	DeviceResponseList	[]*ResponseDeviceResponseList2		`json:"deviceResponseList,omitempty"`		//1,控制设备的结果
	//错误
    ErrorCode			string								`json:"errorCode,omitempty"`	//0,错误代码
    Message				string								`json:"message,omitempty"`		//0,描述
}
type Response2 struct{
	aligenie	*Aligenie
	req			*Request2
	Header		Header2					`json:"header"`
	Payload		ResponsePayload2		`json:"payload"`
}

//错误
//errCode SCode		文本，错误代码
func (T *Response2) Error(errCode SCode) {
	T.Header		= T.req.Header
	T.Header.Name	= "ErrorResponse"
	
	T.Payload.ErrorCode	= errCode.Name()
	T.Payload.Message	= errCode.String()
}

//设备发现
//	product *ResponsePayloadDevices2		分类
//	devices *ResponsePayloadDevices2...	产品
func (T *Response2) Discovery(product *ResponsePayloadDevices2, devices ...*ResponsePayloadDevices2) {
	T.Header		= T.req.Header
	T.Header.Name	= 	T.Header.Name+"Response"
	
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
			if d.Extensions == nil {
				d.Extensions = make(map[string]interface{})
			}
			d.Extensions["parentId"] = product.DeviceId
		}
		T.Payload.Devices	= append(T.Payload.Devices, append([]*ResponsePayloadDevices2{product}, devices...)...)
	}else{
		//其它
		T.Payload.Devices	= append(T.Payload.Devices, devices...)
	}
}

//控制
//	deviceId string	设备id
//	scode SCode		状态代码
func (T *Response2) Control(deviceId string, scode SCode) {
	if T.Header.Name == "" {
		T.Header		= T.req.Header
		T.Header.Name	= "CorrectResponse"
	}
	
	for _, dr := range T.Payload.DeviceResponseList {
		if dr.DeviceId == deviceId {
			return
		}
	}
	deviceResponseList := &ResponseDeviceResponseList2{
		DeviceId:deviceId,
		ErrorCode:scode.Name(),
		Message:scode.String(),
	}
	T.Payload.DeviceResponseList = append(T.Payload.DeviceResponseList, deviceResponseList)
}

//写入到w
//w io.Writer	写入
func (T *Response2) WriteTo(w io.Writer) (n int64, err error) {
	buf := bytes.NewBuffer(nil)
	buf.Grow(1024)
	err = json.NewEncoder(buf).Encode(T)
	if err != nil {
		w.Write([]byte(`{"header":{"namespace":"?","name":"ErrorResponse","messageId":"?","payLoadVersion":2}, "payload":{"devices":[]}}`))
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









