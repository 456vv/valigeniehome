package valigeniehome

type Requester interface{
	Version() int
	V1() *Request1
	V2() *Request2
}

type reqSynthesis struct{
	v1 *Request1
	v2 *Request2
}

func (T *reqSynthesis) Version() int {
	if T.v1 != nil {
		return T.v1.Header.PayLoadVersion
	}
	if T.v2 != nil {
		return T.v2.Header.PayLoadVersion
	}
	return 0
}
func (T *reqSynthesis) V1() *Request1 {
	return T.v1
}
func (T *reqSynthesis) V2() *Request2 {
	return T.v2
}
