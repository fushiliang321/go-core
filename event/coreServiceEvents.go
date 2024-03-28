package event

// 核心服务事件
const (
	BeforeServerStart = "BeforeServerStart" //服务启动之前
	AfterServerStart  = "AfterServerStart"  //服务启动之后
	ServerEnd         = "ServerEnd"         //服务关闭

	BeforeAmqpConsumerServerStart = "BeforeAmqpConsumerServerStart" //Amqp Consumer服务启动之前
	AmqpConsumerServerStart       = "AmqpConsumerServerStart"       //Amqp Consumer服务启动
	AfterAmqpConsumerServerStart  = "AfterAmqpConsumerServerStart"  //Amqp Consumer服务启动之后

	BeforeConsulConsumerServerStart = "BeforeConsulConsumerServerStart" //Consul Consumer服务启动之前
	ConsulConsumerServerStart       = "ConsulConsumerServerStart"       //Consul Consumer服务启动
	ConsulConsumerServiceInfoChange = "ConsulConsumerServiceInfoChange" //Consul Consumer服务节点信息变化
	ConsulServiceRegister           = "ConsulConsumerServiceRegister"   //Consul服务注册
	AfterConsulConsumerServerStart  = "AfterConsulConsumerServerStart"  //Consul Consumer服务启动之后
	ConsulInitFinish                = "ConsulInitFinish"                //Consul初始化完成

	BeforeJsonRpcServerStart = "BeforeJsonRpcHttpServerStart" //JsonRpc服务启动之前
	JsonRpcServerRegister    = "JsonRpcServerRegister"        //JsonRpc服务注册
	AfterJsonRpcServerStart  = "AfterJsonRpcServerStart"      //JsonRpc服务启动之后

	BeforeGrpcServerStart = "BeforeGrpcServerStart" //Grpc服务启动之前
	GrpcServerRegister    = "GrpcServerRegister"    //Grpc服务注册
	AfterGrpcServerStart  = "AfterGrpcServerStart"  //Grpc服务启动之后

	BeforeTaskServerStart = "BeforeTaskServerStart" //Task服务启动之前
	TaskRegister          = "TaskRegister"          //Task注册
	BeforeTaskExecute     = "BeforeTaskExecute"     //Task执行之前
	AfterTaskExecute      = "AfterTaskExecute"      //Task执行之后
	AfterTaskServerStart  = "AfterTaskServerStart"  //Task服务启动之后

	BeforeRateLimitServerStart = "BeforeRateLimitServerStart" //RateLimit服务启动之前
	AfterRateLimitServerStart  = "AfterRateLimitServerStart"  //RateLimit服务启动之后

	BeforeLoggerServerStart = "BeforeLoggerServerStart" //Logger服务启动之前
	AfterLoggerServerStart  = "AfterLoggerServerStart"  //Logger服务启动之后

	HttpServerListen         = "HttpServerListenStart"      //http服务监听开始
	HttpServerListenEnd      = "HttpServerListenEnd"        //http服务监听结束
	WebsocketServerListen    = "WebsocketServerListenStart" //Websocket服务监听开始
	WebsocketServerListenEnd = "WebsocketServerListenEnd"   //Websocket服务监听结束

)
