# valigeniehome [![Build Status](https://travis-ci.org/456vv/valigeniehome.svg?branch=master)](https://travis-ci.org/456vv/valigeniehome)
golang valigeniehome，天猫精灵家居版本

# **列表：**
```go
valigenieHome.go-----------------------------------------------------------------------------------------------------------------------------
type Aligenie struct {                                                  //天猫精灵
    AppIdAttr       string                                                  // 属性id
    apps            map[string]*Application                                 // app集
    m               sync.Mutex                                              // 锁
}
    func (T *Aligenie) SetApp(id string, app *Application)                  // 设置APP
    func (T *Aligenie) ServeHTTP(w http.ResponseWriter, r *http.Request)    // 服务
application.go-----------------------------------------------------------------------------------------------------------------------------

type Application struct{                                                // 程序
    HandleFunc              http.HandlerFunc                                // 处理函数
    ValidReqTimestamp       int                                             // 有效时间，秒为单位
}

home.go-----------------------------------------------------------------------------------------------------------------------------

type Home struct{                                                       // Home
    Response        *Response                                               // 响应
    Request         *Request                                                // 请求
    App             *Application                                            // app
}
    func (T *Home) ServeHTTP(w http.ResponseWriter, r *http.Request)        //服务处理

request.go-----------------------------------------------------------------------------------------------------------------------------

type Header struct{
    MessageId       string      `json:"messageId"`                          //1,用于跟踪请求，是不重复的消息id
    Name            string      `json:"name"`                               //1,操作类型名称
    Namespace       int         `json:"namespace"`                          //1,消息命名空间
    PayLoadVersion  string      `json:"payLoadVersion"`                     //1,payload 的版本,目前版本为 1
}
type RequestPayload struct{
    AccessToken         string                      `json:"accessToken"`    //1,token
    DeviceId            string                      `json:"deviceId"`       //1,设备id
    DeviceType          string                      `json:"deviceType"`     //1,设备类型,http://doc-bot.tmall.com/docs/doc.htm?spm=0.0.0.0.yEvk7c&treeId=393&articleId=108271&docType=1
    Attribute           string                      `json:"attribute"`      //1,属性
    Attribute           string                      `json:"attribute"`      //1,属性，http://doc-bot.tmall.com/docs/doc.htm?spm=0.0.0.0.wzijJu&treeId=393&articleId=108268&docType=1
    Value               string                      `json:"value"`          //1,值，开关的值
    Extensions          map[string]string           `json:"extensions"`     //0,扩展，产品扩展属性,为空返回null或者不返回该字段
}
type Request struct{
    Header      Header              `json:"header"`
    Payload     RequestPayload      `json:"payload"`
}
responseWriter.go-----------------------------------------------------------------------------------------------------------------------------
type ResponseWriter interface{
    Error(errStr, code string)
    Discovery(product ResponsePayloadDevice, devices []ResponsePayloadDevice)
    Control()
    Query(properties []ResponseProperties)
    WriteTo(w io.Writer) (n int64, err error) 
}

response.go-----------------------------------------------------------------------------------------------------------------------------

type ResponseProperties struct{
    Name                string                  `json:"name"`                           //1,名称
    Value               string                  `json:"value"`                          //1,状态
}
type ResponsePayloadDevice struct{
    DeviceId        string                  `json:"deviceId"                            //1,设备id
    DeviceName      string                  `json:"deviceName"                          //1,设备名称,https://open.bot.tmall.com/oauth/api/aliaslist
    DeviceType      string                  `json:"deviceType"                          //1,设备类型,具体参考AliGenie支持的品类列表
    Zone            string                  `json:"zone,omitempty                       //0,位置, https://open.bot.tmall.com/oauth/api/placelist
    Brand           string                  `json:"brand"`                              //1,品牌
    Model           string                  `json:"model"`                              //1,型号
    Icon            string                  `json:"icon"`                               //1,产品icon(https协议的url链接),像素最好160*160 以免在app显示模糊
    Properties      []ResponseProper        `json:"properties"                          //1,性能，返回当前设备支持的属性状态列表，产品支持的属性列表参考 设备控制与设备状态查询页
    Actions         []strin                 `json:"actions"                             //1，开关值，产品支持的操作
    Extensions      map[string]str         `json:"extensions"                           //0，扩展，产品扩展属性,为空返回null或者不返回该字段
}
    func (T *ResponsePayloadDevice) Info(deviceName, brand, model, icon string)         // 产品基本产品
    func (T *ResponsePayloadDevice) Switch(deviceId, powerState string)                 // 开关
    func (T *ResponsePayloadDevice) Outlet(deviceId, powerState string)                 // 插座
    func (T *ResponsePayloadDevice) Light(deviceId, powerState, color, brightness, colorTemperature string) //灯
type ResponsePayload struct{
    Devices             []ResponsePayloadDevice     `json:"devices"`                    //1，用户设备列表
    DeviceId            string                      `json:"deviceId,omitempty"`         //1,设备ID
    ErrorCode           string                      `json:"errorCode,omitempty"`        //0,错误代码
    Message             string                      `json:"message,omitempty"`          //0,描述
}
type Response struct{
    req         *Request
    Properties  []ResponseProperties    `json:"properties"`
    Header      Header                  `json:"header"`
    Payload     ResponsePayload         `json:"payload"`

}
    func (T *Response) Error(errStr, code string)                                       // 错误
    func (T *Response) Discovery(product ResponsePayloadDevice, devices []ResponsePayloadDevice)    //发现
    func (T *Response) Control()                                                        // 控制
    func (T *Response) Query(properties []ResponseProperties)                           // 查询
    func (T *Response) WriteTo(w io.Writer) (n int64, err error)                        // 写入
```