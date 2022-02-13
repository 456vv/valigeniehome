package valigeniehome


//https://www.yuque.com/qw5nze/ga14hc/cmhq2c#jEAJp
//v2:
//{
//  "header": {
//    "messageId": "754eab2b-20f3-4f67-94e9-1a5778c4e1ef",
//    "name": "thing.attribute.set",
//    "namespace": "AliGenie.Iot.Device.Control",
//    "payLoadVersion": 2
//  },
//  "payload": {
//    "accessToken": "1:zMWR5Z9SG6qz0R_Jisxx7q7o1pyuqnD4jZjVmxloBWXr0L_UhlYiprnEXCt8ltRr_nNE81HEm1uo82VdsB6leS_MVhf50.rS5QdBpk9iSaKXcdc5YMQ57xhLHn.Hason",
//    "deviceIds": [
//      "4084:0"
//    ],
//    "extensions": {
//      "4084:0": {
//        "aiIcon": "https://aligenie-bw.iot.birdswo.com.cn/img/png/product/1.png",
//        "productKey": "a1G42dxqmVP"
//      }
//    },
//    "params": {
//      "powerstate": 1
//    }
//  }
//}

type Header2 struct{
	MessageId			string	`json:"messageId"`		//1,用于跟踪请求，是不重复的消息id

	//thing.attribute.set	功能下发(设置)
	//thing.attribute.adjust功能下发(调节)
	//thing.attribute.get	功能上报
	Name				string	`json:"name"`			//1,操作类型名称
	
	//AliGenie.Iot.Device.Discovery	设备发现
	//AliGenie.Iot.Device.Control	设备控制
	//AliGenie.Iot.Device.Query		设备属性查询
	Namespace			string	`json:"namespace"`		//1,消息命名空间
	PayLoadVersion		int		`json:"payLoadVersion"`	//1,payload 的版本
}
type RequestPayload2 struct{
	AccessToken			string								`json:"accessToken"`	//1,token
	DeviceIds			[]string							`json:"deviceIds"`		//1,设备id列表
	//"powerstate": 1	//关闭：0，打开：1 
	//"brightness": 20 //20说明是调高操作；-20说明是调低操作
	Params				map[string]int						`json:"params"`			//1,控制参数
	Extensions			map[string]interface{}				`json:"extensions"`		//0,扩展，产品扩展属性,为空返回null或者不返回该字段
}
type Request2 struct{
	Header		Header2				`json:"header"`
	Payload		RequestPayload2		`json:"payload"`
}







