package valigeniehome
	
import (
	"io"
)

type ResponseWriter interface{
	Error(errStr, code string)
	Discovery(product ResponsePayloadDevice, devices []ResponsePayloadDevice)
	Control()
	Query(properties []ResponseProperties)
	WriteTo(w io.Writer) (n int64, err error) 
}