package utils

import (
	"fmt"
	"simple_tiktok/logger"
	"sync"
	"time"

	"github.com/spf13/viper"
)

/**
 * @Author jiang
 * @Description 熔断器
 * @Date 21:00 2023/1/31
 **/
// api的信息
type apiSnapShop struct {
	isPaused   bool  // api是否熔断
	errCount   int64 // api在周期内失败次数
	totalCount int64 // api在周期内总次数

	accessLast int64 // api最近一次访问时间
	roundLast  int64 // 熔断器周期时间
}

// 熔断器实现体
type CircuitBreakerImp struct {
	lock            sync.RWMutex
	apiMap          map[string]*apiSnapShop // api全局map，key为API标志
	minCheck        int64                   // 接口熔断开启下限次数
	cbkErrRate      float64                 // 接口熔断开启比值
	recoverInterval int64                   // 熔断恢复区间
	roundInterval   int64                   // 计数重复区间
}

// 访问更新
func (c *CircuitBreakerImp) accessed(api *apiSnapShop) {
	/*
	* 判断是否大于周期时间
	* - 是：重置计数
	* - 否：更新计数
	 */
	now := time.Now().Unix()
	if (now - api.roundLast) > int64(c.roundInterval) {
		if api.roundLast != 0 {
			// 首次不打印日志
			logger.SugarLogger.Info("# Trigger 熔断器窗口关闭，重置API计数")
		}
		api.errCount = 0
		api.totalCount = 0
		api.roundLast = now
	}
	api.totalCount++
	api.accessLast = now
}

/*
 * Succeed 记录成功
 * 只更新api列表已有的,
 * 记录访问, 并判断是否熔断:
 * - 是, 取消熔断状态
 */
func (c *CircuitBreakerImp) Succeed(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if api, ok := c.apiMap[key]; ok {
		c.accessed(api)
		if api.isPaused {
			logger.SugarLogger.Info("# Trigger API: \"", key, "\"请求成功，关闭熔断状态.")
			api.isPaused = false
		}
	}
}

/*
 * Failed 记录失败访问
 * api列表查找,
 *	- 已有:
 *		- 记录访问/错误次数
 *		- 是否失败占比到达阈值? 是, 则标记置为熔断
 *	- 未找到:
 *		更新至api列表: 记录访问/错误次数
 */
func (c *CircuitBreakerImp) Failed(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if api, ok := c.apiMap[key]; ok {
		c.accessed(api)
		api.errCount++

		errRate := float64(api.errCount) / float64(api.totalCount)
		// 请求数量达到阈值 && 错误率高于熔断界限
		if api.totalCount > c.minCheck && errRate > c.cbkErrRate {
			logger.SugarLogger.Info("# Trigger 达到错误率, 开启熔断！API：\"", key, "\", total: ", api.totalCount,
				"   errRate: ", errRate)
			api.isPaused = true
		}
		return
	}

	api := &apiSnapShop{}
	c.accessed(api)
	api.errCount++
	// 写入全局map
	c.apiMap[key] = api
}

/*
* CanAccess 判断api是否可访问
 */
func (c *CircuitBreakerImp) CanAccess(key string) bool {
	/*
	 * 判断当前api的isPaused状态
	 *	- 未熔断, 返回true
	 *	- 已熔断, 当前时间与恢复期比较
	 *		- 大于恢复期, 返回true
	 *		- 小于恢复期, 返回false
	 */
	c.lock.RLock()
	defer c.lock.RUnlock()
	// 从api全局map查找
	if api, ok := c.apiMap[key]; ok {
		if api.isPaused {
			// 判断是否进入恢复期
			latency := time.Now().Unix() - api.accessLast
			if latency < int64(c.recoverInterval) {
				// 在恢复期之内, 快速失败，保持熔断
				return false
			}
			// 度过恢复期
			logger.SugarLogger.Info("# Trigger: 熔断器度过恢复期: ", c.recoverInterval, "s , API: \"", key, "\"!")
		}
	}
	// 给予临时恢复
	return true
}

var CB CircuitBreakerImp

// 初始化熔断器
func InitCircuitBreaker() {
	CB = CircuitBreakerImp{}
	CB.apiMap = make(map[string]*apiSnapShop)
	// 控制时间窗口，15秒一轮, 重置api错误率
	CB.roundInterval = viper.GetInt64("CircuitBreaker.roundInterval")
	// 熔断之后，5秒不出现错误再恢复
	CB.recoverInterval = viper.GetInt64("CircuitBreaker.recoverInterval")
	// 请求到达5次以上才进行熔断检测
	CB.minCheck = 5
	// 错误率到达 50% 开启熔断
	CB.cbkErrRate = 0.5

	fmt.Println("CircuitBreaker inited ...... ")
}
