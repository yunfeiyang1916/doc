package api

import (
	"fmt"
	"net/http"
	"shop/shop-api/user-web/forms"
	"shop/shop-api/user-web/global"
	"shop/shop-api/user-web/global/response"
	"shop/shop-api/user-web/middlewares"
	"shop/shop-api/user-web/models"
	"shop/shop-srv/user-srv/proto"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/dgrijalva/jwt-go"

	"github.com/go-playground/validator/v10"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

// 将grpc的code转成http状态码
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Code(),
				})
			}
		}
	}
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func GetUserList(c *gin.Context) {
	zap.S().Info("获取用户列表页")

	//claims, _ := c.Get("claims")
	//currentUser := claims.(*models.CustomClaims)
	//zap.S().Infof("访问用户: %d", currentUser.ID)
	pn := c.DefaultQuery("pn", "1")
	pSize := c.DefaultQuery("psize", "10")
	pnInt, _ := strconv.Atoi(pn)
	pSizeInt, _ := strconv.Atoi(pSize)
	// 调用接口
	res, err := global.UserSrvClient.GetUserList(c.Request.Context(), &proto.PageInfo{Pn: uint32(pnInt), PSize: uint32(pSizeInt)})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表】失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}
	reMap := gin.H{
		"total": res.Total,
	}
	list := make([]response.UserResponse, 0, len(res.Data))
	for _, v := range res.Data {
		user := response.UserResponse{
			Id:       v.Id,
			NickName: v.NickName,
			Birthday: response.JsonTime(time.Unix(int64(v.BirthDay), 0)),
			Gender:   v.Gender,
			Mobile:   v.Mobile,
		}
		list = append(list, user)
	}
	reMap["data"] = list
	c.JSON(http.StatusOK, reMap)
}

// 使用用户名密码登录
func PassWordLogin(c *gin.Context) {
	// 表单验证
	passwordLoginForm := forms.PassWordLoginForm{}
	if err := c.ShouldBindJSON(&passwordLoginForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	// 验证码验证
	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, false) {
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	if rsp, err := global.UserSrvClient.GetUserByMobile(c.Request.Context(), &proto.MobileRequest{Mobile: passwordLoginForm.Mobile}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{"mobile": "用户不存在"})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{"mobile": "登录失败"})
			}
			return
		}
		c.JSON(http.StatusInternalServerError, map[string]string{"mobile": "登录失败"})
		return
	} else {
		// 校验密码
		if passRsp, err := global.UserSrvClient.CheckPassword(c.Request.Context(), &proto.PasswordCheckInfo{
			Password:          passwordLoginForm.PassWord,
			EncryptedPassword: rsp.PassWord,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{"password": "登录失败"})
		} else {
			if passRsp.Success {
				// 生成token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						// 签名生效时间
						NotBefore: time.Now().Unix(),
						// 30天过期
						ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Unix(),
						// 签发者
						Issuer: "yunfeiyang",
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"msg": "生成token失败"})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"id":        rsp.Id,
					"nick_name": rsp.NickName,
					"token":     token,
					// 30天毫秒级别
					"expired_at": (time.Now().Add(30 * 24 * time.Hour).Unix()) * 1000,
				})
			} else {
				c.JSON(http.StatusBadRequest, map[string]string{"password": "密码错误"})
			}
		}
	}
}

// 用户注册
func Register(c *gin.Context) {
	registerForm := forms.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	// 从redis读取验证码
	rdb := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port)})
	value, err := rdb.Get(c.Request.Context(), registerForm.Mobile).Result()

	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码错误",
		})
		return
	} else {
		if value != registerForm.Code {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": "验证码错误",
			})
			return
		}
	}
	user, err := global.UserSrvClient.CreateUser(c.Request.Context(), &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		Password: registerForm.PassWord,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		zap.S().Errorf("[Register] 查询 【新建用户失败】失败: %s", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}
	// 生成token
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(user.Id),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			// 签名生效时间
			NotBefore: time.Now().Unix(),
			// 30天过期
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Unix(),
			// 签发者
			Issuer: "yunfeiyang",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "生成token失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":        user.Id,
		"nick_name": user.NickName,
		"token":     token,
		// 30天毫秒级别
		"expired_at": (time.Now().Add(30 * 24 * time.Hour).Unix()) * 1000,
	})
}
