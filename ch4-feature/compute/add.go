package compute

type AddOperator interface {

	/**
	 * 算术增加
	 */
	Add() interface{}

}

type IntParams struct{
	P1 int  // 加数
	P2 int  // 被加数
}

func (params *IntParams) Add() interface{} {
	return params.P1 + params.P2
}


