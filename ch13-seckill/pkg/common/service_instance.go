package common


type ServiceInstance struct {
	Host      string //  Host
	Port      int    //  Port
	Weight    int    // 权重
	CurWeight int    // 当前权重

	GrpcPort int
}

