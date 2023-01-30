package utils

import (
	"context"
	"fmt"
	"reflect"
)

/**
 * @Author jiang
 * @Description 布隆过滤器（redis）
 * @Date 12:00 2023/1/30
 **/
// 设置布隆过滤器容量，默认大小为100000
const (
	KeyBloomFilter = "BloomFilter"
	error_rate     = 0.01
	DefaultSize    = 100000
)

// 设置种子，保证不同哈希函数有不同的计算方式
var seeds = []uint{77, 112, 1390, 31342}

var ctx = context.Background()

// 初始化布隆过滤器
func InitBloomFilter() {
	// 先判断布隆过滤器是否已经存在，存在，则不再创建
	rnt, _ := RDB9.Exists(ctx, KeyBloomFilter).Result()
	if rnt == 1 {
		return
	}

	// 创建一个大小为capacity，错误率为error_rate的空的Bloom
	// 	BF.RESERVE {key} {error_rate} {capacity} [EXPANSION expansion] [NONSCALING]
	_, err := RDB9.Do(ctx, "BF.RESERVE", KeyBloomFilter, error_rate, DefaultSize).Result()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("BloomFilter inited ...... ")
}

// 构造哈希函数
func CreateHash(seed uint, value string) uint {
	var result uint = 0
	for i := 0; i < len(value); i++ {
		result = result*seed + uint(value[i])
	}
	//length = 2^n 时，X % length = X & (length - 1)
	return result & (DefaultSize - 1)
}

// 将值添加到布隆过滤器
func BloomFilterAdd(value string) error {
	result := [4]uint{}
	for i, seed := range seeds {
		result[i] = CreateHash(seed, value)
	}

	// 向key指定的Bloom中添加多个元素
	// BF.MADD {key} {item} [item…]
	_, err := RDB9.Do(ctx, "BF.MADD", KeyBloomFilter, result[0], result[1], result[2], result[3]).Result()
	return err
}

// 判断该值是否存在布隆过滤器
func BloomFilterCheck(value string) (bool, error) {
	result := [4]uint{}
	for i, seed := range seeds {
		result[i] = CreateHash(seed, value)
	}

	// 同时检查元素是否可能存在于key指定的Bloom中
	// BF.MEXISTS {key} {item} [item…]
	res, err := RDB9.Do(ctx, "BF.MEXISTS", KeyBloomFilter, result[0], result[1], result[2], result[3]).Result()
	if err != nil {
		return false, err
	}

	valueS := reflect.ValueOf(res)
	for i := 0; i < valueS.Len(); i++ {
		if fmt.Sprintf("%v", valueS.Index(i)) == "0" {
			return false, nil
		}
	}

	return true, nil
}

// /**
//  * @Author jiang
//  * @Description 布隆过滤器（BitSet，废弃，已经升级为redis的布隆过滤器）
//  * @Date 13:00 2023/1/17
//  **/
// //设置哈希数组默认大小为100000
// const DefaultSize = 100000

// var Filter *BloomFilter

// //设置种子，保证不同哈希函数有不同的计算方式
// var seeds = []uint{97, 112, 1390, 31342, 3237, 631}

// //布隆过滤器结构，包括二进制数组和多个哈希函数
// type BloomFilter struct {
// 	//使用第三方库
// 	set *bitset.BitSet
// 	//指定长度为6
// 	hashFuncs [6]func(seed uint, value string) uint
// }

// //构造一个布隆过滤器，包括数组和哈希函数的初始化
// func NewBloomFilter() *BloomFilter {
// 	bf := new(BloomFilter)
// 	bf.set = bitset.New(DefaultSize)

// 	for i := 0; i < len(bf.hashFuncs); i++ {
// 		bf.hashFuncs[i] = createHash()
// 	}
// 	return bf
// }

// //构造6个哈希函数，每个哈希函数有参数seed保证计算方式的不同
// func createHash() func(seed uint, value string) uint {
// 	return func(seed uint, value string) uint {
// 		var result uint = 0
// 		for i := 0; i < len(value); i++ {
// 			result = result*seed + uint(value[i])
// 		}
// 		//length = 2^n 时，X % length = X & (length - 1)
// 		return result & (DefaultSize - 1)
// 	}
// }

// //添加元素
// func (b *BloomFilter) Add(value string) {
// 	for i, f := range b.hashFuncs {
// 		//将哈希函数计算结果对应的数组位置1
// 		b.set.Set(f(seeds[i], value))
// 	}
// }

// //判断元素是否存在
// func (b *BloomFilter) Check(value string) bool {
// 	//调用每个哈希函数，并且判断数组对应位是否为1
// 	//如果不为1，直接返回false，表明一定不存在
// 	for i, f := range b.hashFuncs {
// 		//result = result && b.set.Test(f(seeds[i], value))
// 		if !b.set.Test(f(seeds[i], value)) {
// 			return false
// 		}
// 	}
// 	return true
// }

// // 初始化布隆过滤器
// func InitBloomFilter() {
// 	Filter = NewBloomFilter()
// 	fmt.Println("BloomFilter inited ...... ")
// }

// func main() {
//	filter := NewBloomFilter()
// 	filter.add("asd")
// 	fmt.Println(filter.contains("asd"))
// 	fmt.Println(filter.contains("2222"))
// 	fmt.Println(filter.contains("155343"))
// }
