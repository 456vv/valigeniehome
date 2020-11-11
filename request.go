package valigeniehome
	


//{
//	"header": {
//		"messageId": "93d03349-7f1b-45a9-8723-9dcc375a8c2d",
//		"name": "DiscoveryDevices",
//		"namespace": "AliGenie.Iot.Device.Discovery",
//		"payLoadVersion": 1
//	},
//	"payload": {
//		"accessToken": "1:9e78277ab62b3c51996fb39913ddd2abbde393cb8f55b572884ea4ff7b9c55f7eeb47cf586ada1f0812dba7e070bf564ce39baea81daf522dcecd45818508c0a"
//	}
//}

//{
//	"header": {
//		"messageId": "e0b9e2a9-8fb4-40e6-af6e-c5142077ef19",
//		"name": "TurnOn",
//		"namespace": "AliGenie.Iot.Device.Control",
//		"payLoadVersion": 1
//	},
//	"payload": {
//		"accessToken": "1:1aff43324fb8f05b0823e7804b817c7a43313d2eedf66d1561d624b5bbe47498fe9ab5cf1d079d40ec23ce3a8a5d8d19a6ca421a9326dd55ea10beb0f4248bfa",
//		"attribute": "powerstate",
//		"deviceId": "73:45264c",
//		"deviceType": "switch",
//		"value": "on"
//	}
//}

type Header struct{
	MessageId			string	`json:"messageId"`		//1,用于跟踪请求，是不重复的消息id
	
	//设备发现类（与AliGenie.Iot.Device.Discovery对应）
	//DiscoveryDevices		设备发现（获取设备列表）
	
	//操作类（与AliGenie.Iot.Device.Control对应）
	//TurnOn				打开
	//TurnOff				关闭
	//SelectChannel			频道切换
	//AdjustUpChannel		频道增加
	//AdjustDownChannel		频道减少
	//AdjustUpVolume		声音按照步长调大
	//AdjustDownVolume		声音按照步长调小
	//SetVolume				声音调到某个值
	//SetMute				设置静音
	//CancelMute			取消静音
	//Play					播放
	//Pause					暂停
	//Continue				继续
	//Next					下一首或下一台
	//Previous				上一首或上一台
	//SetBrightness			设置亮度
	//AdjustUpBrightness	调大亮度
	//AdjustDownBrightness	调小亮度
	//SetTemperature		设置温度
	//AdjustUpTemperature	调高温度
	//AdjustDownTemperature	调低温度
	//SetWindSpeed			设置风速
	//AdjustUpWindSpeed		调大风速
	//AdjustDownWindSpeed	调小风速
	//SetMode				模式的切换
	//SetColor				设置颜色
	//OpenFunction			打开功能
	//CloseFunction			关闭功能
	//Cancel				取消
	//CancelMode			取消模式(退出模式)
	
	//查询类（与AliGenie.Iot.Device.Query对应）
	//支持的查询属性方法	操作方法说明		返回值说明
	//Query					查询所有标准属性	详情见各个属性
	//QueryColor			查询颜色			Red、Yellow、Blue、White、Black等值（AliGenie以这些值为准，厂家适配）
	//QueryPowerState		查询电源开关		on(打开)、off(关闭)
	//QueryTemperature		查询温度			返回数值(AliGenie默认的单位为摄氏度，厂家适配该单位)
	//QueryHumidity			查询湿度			返回数值
	//QueryWindSpeed		查询风速			返回值参考 设备控制 中章节 1.8.1 风速值对应表
	//QueryBrightness		查询亮度			返回数值
	//QueryFog				查询雾量			返回数值
	//QueryMode				查询模式			返回值枚举参考例子模式切换中的例子
	//QueryPM25				查询pm2.5 含量		返回数值
	//QueryDirection		查询方向			返回 left,right,forward,back,up,down
	//QueryAngle			查询角度			返回数值，单位度
	Name				string	`json:"name"`			//1,操作类型名称
	
	//AliGenie.Iot.Device.Discovery	设备发现
	//AliGenie.Iot.Device.Control	设备控制
	//AliGenie.Iot.Device.Query		设备属性查询
	Namespace			string	`json:"namespace"`		//1,消息命名空间
	PayLoadVersion		int		`json:"payLoadVersion"`	//1,payload 的版本,目前版本为 1
}
type RequestPayload struct{
	AccessToken			string								`json:"accessToken"`	//1,token
	
	DeviceId			string								`json:"deviceId"`		//1,设备id
	DeviceType			string								`json:"deviceType"`		//1,设备类型,http://doc-bot.tmall.com/docs/doc.htm?spm=0.0.0.0.yEvk7c&treeId=393&articleId=108271&docType=1	Attribute				
	Attribute			string								`json:"attribute"`		//1,属性，http://doc-bot.tmall.com/docs/doc.htm?spm=0.0.0.0.wzijJu&treeId=393&articleId=108268&docType=1
	
	//最值	说明
	//max	对应最大值
	//min	对应最小值
	Value				string								`json:"value"`			//1,值，开关的值
	Extensions			map[string]string					`json:"extensions"`		//0,扩展，产品扩展属性,为空返回null或者不返回该字段
}
type Request struct{
	Header		Header				`json:"header"`
	Payload		RequestPayload		`json:"payload"`
}








