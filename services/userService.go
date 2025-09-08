package services

import (
	"context"
	"myblog/models"
	Req "myblog/repository"
	"myblog/tools"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 注册用户业务逻辑
func RegisterUserService(c *gin.Context, db *gorm.DB, user models.User) {
	//检查用户是否存在
	u := Req.UserRep{Db: db}
	slectErr := u.SelectUserByName(user.Username)
	if slectErr == nil {
		c.JSON(400, gin.H{
			"code": "400",
			"mgs":  "用户名已存在",
		})
		return
	}
	//检查邮件地址是否重复
	selectEmailErr := u.SelectUserByEmail(user.Email)
	if selectEmailErr == nil {
		c.JSON(400, gin.H{
			"code": "400",
			"msg":  "邮箱已被注册",
		})
		return
	}
	//加密密码,存入数据库
	hashedPassword, hashErr := tools.HashPassword(user.Password)
	if hashErr != nil {
		c.JSON(500, gin.H{
			"code":   "500",
			"msg":    "注册失败",
			"errors": hashErr,
		})
		return
	}
	//插入数据
	registerErr := u.CreateUser(&models.User{Username: user.Username, Password: hashedPassword, Email: user.Email})
	if registerErr != nil {
		c.JSON(500, gin.H{
			"code":   "500",
			"msg":    "注册失败",
			"errors": registerErr,
		})
		return
	}
	//插入成功，直接返回用户信息
	roles := []string{"admin"} // 这里使用中间件进行设置，待优化
	var msgValue = map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
		"roles":    roles,
	}
	c.JSON(200, gin.H{
		//"RequestID": c.MustGet("RequestID").(string),
		"code": 200,
		"msg":  "注册成功",
		"data": msgValue,
	})

}

// 用户登录业务逻辑
func LoginUserService(c *gin.Context, db *gorm.DB, rdb *redis.Client, user models.User) {
	ctx := context.Background() // redis上下文

	u := Req.UserRep{Db: db}
	var userResult models.User
	zap.S().Debug("检查用户登录的用户名和密码是否存在数据库中")
	selectErr := u.SelectUserByNameScanValue(&user, &userResult)
	if selectErr == nil {
		zap.S().Debug()
	}
	zap.S().Debug("查询成功：", userResult)

	if userResult.Username == "" {
		c.JSON(400, gin.H{
			"code":   "400",
			"msg":    "用户未注册",
			"errors": selectErr,
		})
		return
	}
	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(userResult.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": "400",
			"msg":  "无效的用户名或密码",
		})
		return
	}
	zap.S().Debug("用户名和密码校验通过！开始生成token")
	roles := []string{"admin"}
	//生成claims
	x, err := strconv.Atoi(os.Getenv("JWT_TIME")) //字符串转成int类型
	if err != nil {
		panic(err)
	}
	var exptime int64
	exptime = time.Now().Add(time.Duration(x) * time.Hour).Unix()
	var claims = tools.MyClaims{ //这里改成了自定义的Claims结构体，若使用jwt自带jwt.MapClaims,则无法拿到username。需要理解
		UserId:   userResult.ID,
		Username: userResult.Username,
		Roles:    roles,
		Exp:      exptime, //失效时间，8小时失效，观察效果
	}
	//判断redis是否存在缓存，如果存在则直接从redis进行获取
	currentToken, getRedisErr := rdb.Get(ctx, userResult.Username).Result()
	if getRedisErr != nil { //表示没有找到redis缓存，进行加入缓存操作
		token, generateErr := tools.GenerateToken(&claims)
		if generateErr != nil {
			c.JSON(500, gin.H{
				"code":   "500",
				"msg":    "生成token失败",
				"errors": generateErr,
			})
			return
		}
		zap.S().Debug("成功生成用户名:", user.Username, "token:", token)
		///生成token,exp过期时间并存入redis
		zap.S().Debug("启动redis缓存当前登录成功的username")
		zap.S().Debug("exptime: ", exptime, " second: ", time.Duration(x)*time.Hour)
		redisSetErr := rdb.Set(ctx, user.Username, token, time.Duration(x)*time.Hour)

		if redisSetErr.Err() != nil {
			panic(redisSetErr.Err())
		}
		zap.S().Debug("token保存redis缓存成功")
		currentToken = token
	}
	//判断token是否失效
	c.JSON(http.StatusOK, gin.H{
		"code":   200,
		"msg":    "用户登录成功",
		"userId": userResult.ID,
		"rols":   roles,
		"token":  currentToken,
		"exp":    exptime,
	})

}
