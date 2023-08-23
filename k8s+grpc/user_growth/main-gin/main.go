package main

import (
	"context"
	"log"
	"net/http"
	"user_growth/conf"
	"user_growth/dbhelper"
	"user_growth/pb"
	"user_growth/ugserver"

	"google.golang.org/grpc/metadata"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func initDb() {
	//time.Local=time.UTC
	conf.LoadConfigs()
	dbhelper.InitDb()
}

// 允许跨域的白名单
var AllowOrigin = map[string]bool{
	"http://a.site.com": true,
	"http://b.site.com": true,
	"http://web.com":    true,
}

func mainGateway() {
	initDb()
	// 初始化服务
	s := grpc.NewServer()
	// 注册服务
	pb.RegisterUserCoinServer(s, &ugserver.UgCoinServer{})
	pb.RegisterUserGradeServer(s, &ugserver.UgGradeServer{})

	// grpc-gateway注册服务
	// 设置支持跨域
	serverMuxOpt := []runtime.ServeMuxOption{
		runtime.WithOutgoingHeaderMatcher(func(s string) (string, bool) {
			return s, true
		}),
		// 生成metadata
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			origin := request.Header.Get("Origin")
			if AllowOrigin[origin] {
				md := metadata.New(map[string]string{
					"Access-Control-Allow-Origin":      origin,
					"Access-Control-Allow-Methods":     "GET,POST,PUT,DELETE,OPTION",
					"Access-Control-Allow-Headers":     "*",
					"Access-Control-Allow-Credentials": "true",
				})
				grpc.SetHeader(ctx, md)
			}
			return nil
		}),
	}
	mux := runtime.NewServeMux(serverMuxOpt...)
	ctx := context.Background()
	if err := pb.RegisterUserCoinHandlerServer(ctx, mux, &ugserver.UgCoinServer{}); err != nil {
		log.Printf("Faile to RegisterUserCoinHandlerServer error=%v", err)
	}
	if err := pb.RegisterUserGradeHandlerServer(ctx, mux, &ugserver.UgGradeServer{}); err != nil {
		log.Printf("Faile to RegisterUserGradeHandlerServer error=%v", err)
	}
	httpMux := http.NewServeMux()
	httpMux.Handle("/v1/UserGrowth", mux)
	// 配置http服务
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("http.HandlerFunc url=%s", r.URL)
			mux.ServeHTTP(w, r)
		}),
	}
	// 启动http服务
	log.Printf("server.ListenAndServer(%s)", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("ListenAndServe error=%v", err)
	}
}

func mainGin() {
	// 连接到grpc服务的客户端
	conn, err := grpc.Dial("localhost:80", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect:%v", err)
	}
	defer conn.Close()

	clientCoin := pb.NewUserCoinClient(conn)
	clientGrade := pb.NewUserGradeClient(conn)

	router := gin.New()
	router.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello")
	})

	// 用户积分服务的方法，定义服务路由组
	// 在根路径设置允许跨域请求
	v1Group := router.Group("/v1", func(ctx *gin.Context) {
		// 支持跨域,只允许白名单的域名访问
		origin := ctx.GetHeader("Origin")
		if AllowOrigin[origin] {
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTION")
			ctx.Header("Access-Control-Allow-Headers", "*")
			ctx.Header("Access-Control-Allow-Credentials", "true")
		}
		ctx.Next()
	})
	gUserCoin := v1Group.Group("/UserGrowth.UserCoin")
	gUserCoin.GET("/ListTasks", func(ctx *gin.Context) {
		out, err := clientCoin.ListTasks(ctx, &pb.ListTasksRequest{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    2,
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, out)
	})
	gUserCoin.POST("/UserCoinChange", func(ctx *gin.Context) {
		body := &pb.UserCoinChangeRequest{}
		err := ctx.BindJSON(body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    2,
				"message": err.Error(),
			})
			return
		}
		out, err := clientCoin.UserCoinChange(ctx, body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    2,
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, out)
	})

	// 用户等级服务的方法
	gUserGrade := v1Group.Group("/UserGrowth.UserGrade")
	gUserGrade.GET("/ListGrades", func(ctx *gin.Context) {
		out, err := clientGrade.ListGrades(ctx, &pb.ListGradesRequest{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    2,
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, out)
	})

	// 为http2配置参数
	h2Handler := h2c.NewHandler(router, &http2.Server{})
	// 配置http服务
	server := &http.Server{
		Addr:    ":8080",
		Handler: h2Handler,
	}
	// 启动http服务
	server.ListenAndServe()
}

func main() {
	//mainGin()
	mainGateway()
}
