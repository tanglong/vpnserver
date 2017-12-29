package common

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const (
	HEART_BEAT       = 10
	TIMER_QUEUE_SIZE = 1024
	TIME_OUT_FACTOR  = 3
)

// flags + team + pos +  slot
// uint8  uint8  uint8  uint8
type UFightRecKey uint32

func NewUFightRecKey(flags, pos, slot uint8) UFightRecKey {
	return UFightRecKey(flags)<<16 | UFightRecKey(pos)<<8 | UFightRecKey(slot)
}

func (f UFightRecKey) Flags() uint8 {
	return uint8(f >> 16)
}

//func (f UFightRecKey) Team() uint8 {
//	return uint8(f >> 16)
//}

func (f UFightRecKey) Pos() uint8 {
	return uint8(f >> 8)
}

func (f UFightRecKey) Slot() uint8 {
	return uint8(f)
}

type UPhotoRecKey uint64

func NewUPhotoRecKey(rType, id uint32) UPhotoRecKey {
	return UPhotoRecKey(rType)<<32 | UPhotoRecKey(id)
}

func (f UPhotoRecKey) Types() uint32 {
	return uint32(f >> 32)
}

func (f UPhotoRecKey) Ids() uint32 {
	return uint32(f)
}

var identity uint32 = 0

func InitIdentity(id uint32) {
	identity = id
}

func NewID() uint32 {
	return atomic.AddUint32(&identity, 1)
}

func AddrToUint32(strHost string) uint32 {
	var dwIp uint32
	slcIpPort := strings.Split(strHost, ":")
	ips := strings.Split(slcIpPort[0], ".")
	if len(ips) != 4 {
		return uint32(dwIp)
	}

	b0, _ := strconv.Atoi(ips[3])
	b1, _ := strconv.Atoi(ips[2])
	b2, _ := strconv.Atoi(ips[1])
	b3, _ := strconv.Atoi(ips[0])
	dwIp += uint32(b0) << 24
	dwIp += uint32(b1) << 16
	dwIp += uint32(b2) << 8
	dwIp += uint32(b3)
	return uint32(dwIp)
}

func AddrPort(strHost string) uint32 {
	var dwPort uint32

	slcIpPort := strings.Split(strHost, ":")
	if len(slcIpPort) != 2 {
		return dwPort
	}
	nPort, _ := strconv.Atoi(slcIpPort[1])
	dwPort = uint32(nPort)
	return dwPort
}

func AddrToUint64(strHost string) uint64 {
	var dwIp uint32
	var dwPort uint32

	slcIpPort := strings.Split(strHost, ":")
	if len(slcIpPort) != 2 {
		return ((uint64(dwIp)) << 32) + uint64(dwPort)
	}

	nPort, _ := strconv.Atoi(slcIpPort[1])
	dwPort = uint32(nPort)

	ips := strings.Split(slcIpPort[0], ".")
	if len(ips) != 4 {
		return ((uint64(dwIp)) << 32) + uint64(dwPort)
	}

	b0, _ := strconv.Atoi(ips[3])
	b1, _ := strconv.Atoi(ips[2])
	b2, _ := strconv.Atoi(ips[1])
	b3, _ := strconv.Atoi(ips[0])
	dwIp += uint32(b0) << 24
	dwIp += uint32(b1) << 16
	dwIp += uint32(b2) << 8
	dwIp += uint32(b3)
	return ((uint64(dwIp)) << 32) + uint64(dwPort)
}

func Uint64ToAddr(ldwAddr uint64) string {
	return fmt.Sprintf("%d.%d.%d.%d:%d", (ldwAddr>>32)&0xff, (ldwAddr>>40)&0xff, (ldwAddr>>48)&0xff, (ldwAddr>>56)&0xff, ldwAddr&0xffffffff)
}

//获取几分之的几率
func SelectByOdds(upNum, downNum uint32) bool {
	if downNum < 1 {
		return false
	}
	if upNum < 1 {
		return false
	}
	if upNum > downNum-1 {
		return true
	}
	return (1 + uint32((float64(rand.Int63())/(1<<63))*float64(downNum))) <= upNum
}

//获取百分之的几率
func SelectByPercent(percent uint32) bool {
	return SelectByOdds(percent, 100)
}

//获取千分之的几率
func SelectByThousand(th uint32) bool {
	return SelectByOdds(th, 1000)
}

//获取万分之的几率
func SelectByTenTh(tenth uint32) bool {
	return SelectByOdds(tenth, 10000)
}

//获取十万分之的几率
func SelectByLakh(lakh uint32) bool {
	return SelectByOdds(lakh, 100000)
}

func DailyZero() uint32 {
	now := time.Now()
	hour, minute, second := now.Hour(), now.Minute(), now.Second()
	return uint32(now.Unix() - int64((hour*3600)+(minute*60)+second))
}

func DailyHour() uint32 {
	now := time.Now()
	minute, second := now.Minute(), now.Second()
	return uint32(now.Unix() - int64((minute*60)+second))
}

// Ascii numbers 0-9
const (
	ascii_0 = 48
	ascii_9 = 57
)

func ParseUint64(d []byte) (uint64, bool) {
	var n uint64
	d_len := len(d)
	if d_len == 0 {
		return 0, false
	}
	for i := 0; i < d_len; i++ {
		j := d[i]
		if j < ascii_0 || j > ascii_9 {
			return 0, false
		}
		n = n*10 + (uint64(j - ascii_0))
	}
	return n, true
}

//获取一年有多少天
func GetYearDays(year int) int {
	days := 0
	if ((year%4) == 0 && (year%100) != 0) || (year%400) == 0 {
		days = 366
	} else {
		days = 365
	}
	return days
}

//将 a:b|c:d的格式转换成map
func StringToMap(str string) (map[uint32]uint32, bool) {
	if str == "" {
		return nil, true
	}
	result_map := make(map[uint32]uint32)
	str1 := strings.Split(str, "|")
	st1_len := len(str1)
	if st1_len > 0 {
		for _, value := range str1 {
			str2 := strings.Split(value, ":")
			str2_len := len(str2)
			if str2_len == 2 {
				k := uint64(0)
				v := uint64(0)
				var e error
				if k, e = strconv.ParseUint(str2[0], 10, 64); e != nil {
					return nil, false
				}
				if v, e = strconv.ParseUint(str2[1], 10, 64); e != nil {
					return nil, false
				}
				result_map[uint32(k)] = uint32(v)
			} else {
				return nil, false
			}
		}
	}
	return result_map, true
}

func StringToMap64(str string) (map[uint32]uint64, bool) {
	if str == "" {
		return nil, true
	}
	result_map := make(map[uint32]uint64)
	str1 := strings.Split(str, "|")
	st1_len := len(str1)
	if st1_len > 0 {
		for _, value := range str1 {
			str2 := strings.Split(value, ":")
			str2_len := len(str2)
			if str2_len == 2 {
				k := uint64(0)
				v := uint64(0)
				var e error
				if k, e = strconv.ParseUint(str2[0], 10, 64); e != nil {
					return nil, false
				}
				if v, e = strconv.ParseUint(str2[1], 10, 64); e != nil {
					return nil, false
				}
				result_map[uint32(k)] = v
			} else {
				return nil, false
			}
		}
	}
	return result_map, true
}

func StringToMap_Uint64(str string) (map[uint64]uint64, bool) {
	if str == "" {
		return nil, true
	}
	result_map := make(map[uint64]uint64)
	str1 := strings.Split(str, "|")
	st1_len := len(str1)
	if st1_len > 0 {
		for _, value := range str1 {
			str2 := strings.Split(value, ":")
			str2_len := len(str2)
			if str2_len == 2 {
				k := uint64(0)
				v := uint64(0)
				var e error
				if k, e = strconv.ParseUint(str2[0], 10, 64); e != nil {
					return nil, false
				}
				if v, e = strconv.ParseUint(str2[1], 10, 64); e != nil {
					return nil, false
				}
				result_map[k] = v
			} else {
				return nil, false
			}
		}
	}
	return result_map, true
}

func StringToSlice(str string) []uint32 {
	if str == "" {
		return nil
	}
	str1 := strings.Split(str, "|")
	var slice []uint32
	if len(str1) > 0 {
		for _, v := range str1 {
			if value, e := strconv.ParseUint(v, 10, 64); e != nil {
				return nil
			} else {
				slice = append(slice, uint32(value))
			}
		}
		return slice
	}
	return nil
}

func SliceToString(slice []uint32) string {
	//	var str []byte
	str := make([]byte, 0, 512)
	if len(slice) == 0 {
		return string(str)
	}

	for _, v := range slice {
		//		if str != []byte("") {
		if len(str) > 0 {
			//		str += "|" + strconv.FormatUint(uint64(v), 10)
			str = append(str, []byte("|")...)
			str = append(str, []byte(strconv.FormatUint(uint64(v), 10))...)
		} else {
			//		str += strconv.FormatUint(uint64(v), 10)
			str = append(str, []byte(strconv.FormatUint(uint64(v), 10))...)
		}
	}
	return string(str)
}

func StringToSlice64(str string) []uint64 {
	if str == "" {
		return nil
	}
	str1 := strings.Split(str, "|")
	var slice []uint64
	if len(str1) > 0 {
		for _, v := range str1 {
			if value, e := strconv.ParseUint(v, 10, 64); e != nil {
				return nil
			} else {
				slice = append(slice, value)
			}
		}
		return slice
	}
	return nil
}

func Slice64ToString(slice []uint64) string {
	str := make([]byte, 0, 512)
	if len(slice) == 0 {
		return string(str)
	}

	for _, v := range slice {
		if len(str) > 0 {
			str = append(str, []byte("|")...)
			str = append(str, []byte(strconv.FormatUint(uint64(v), 10))...)
		} else {
			str = append(str, []byte(strconv.FormatUint(uint64(v), 10))...)
		}
	}
	return string(str)
}

//将map 转换成a:b|c:d
func MapToString(result_map map[uint32]uint32) string {
	//	str := []byte("")
	//var str []byte
	str := make([]byte, 0, 512)
	if result_map == nil || len(result_map) <= 0 {
		return string(str)
	}
	for k, v := range result_map {
		//		if str != []byte"" {
		if len(str) > 0 {
			//str = str + "|" + strconv.FormatUint(uint64(k), 10) + ":" + strconv.FormatUint(uint64(v), 10)
			str = append(str, []byte("|")...)
			str = append(str, []byte(strconv.FormatUint(uint64(k), 10))...)
			str = append(str, []byte(":")...)
			str = append(str, []byte(strconv.FormatUint(uint64(v), 10))...)
		} else {
			//			str = strconv.FormatUint(uint64(k), 10) + ":" + strconv.FormatUint(uint64(v), 10)
			str = append(str, []byte(strconv.FormatUint(uint64(k), 10))...)
			str = append(str, []byte(":")...)
			str = append(str, []byte(strconv.FormatUint(uint64(v), 10))...)

		}
	}
	return string(str)
}

func Map64ToString(result_map map[uint32]uint64) string {
	//	str := ""
	//var str []byte
	str := make([]byte, 0, 512)
	if result_map == nil || len(result_map) <= 0 {
		return string(str)
	}
	for k, v := range result_map {
		//		if str != "" {
		if len(str) > 0 {
			//			str = str + "|" + strconv.FormatUint(uint64(k), 10) + ":" + strconv.FormatUint(v, 10)
			str = append(str, []byte("|")...)
			str = append(str, []byte(strconv.FormatUint(uint64(k), 10))...)
			str = append(str, []byte(":")...)
			str = append(str, []byte(strconv.FormatUint(uint64(v), 10))...)

		} else {
			//			str = strconv.FormatUint(uint64(k), 10) + ":" + strconv.FormatUint(v, 10)
			str = append(str, []byte(strconv.FormatUint(uint64(k), 10))...)
			str = append(str, []byte(":")...)
			str = append(str, []byte(strconv.FormatUint(uint64(v), 10))...)
		}
	}
	return string(str)

}

func MapToString_Uint64(result_map map[uint64]uint64) string {
	//	str := ""
	//var str []byte
	str := make([]byte, 0, 512)
	if result_map == nil || len(result_map) <= 0 {
		return string(str)
	}
	for k, v := range result_map {
		//		if str != "" {
		if len(str) > 0 {
			//		str = str + "|" + strconv.FormatUint(k, 10) + ":" + strconv.FormatUint(v, 10)
			str = append(str, []byte("|")...)
			str = append(str, []byte(strconv.FormatUint(uint64(k), 10))...)
			str = append(str, []byte(":")...)
			str = append(str, []byte(strconv.FormatUint(uint64(v), 10))...)
		} else {
			//		str = strconv.FormatUint(k, 10) + ":" + strconv.FormatUint(v, 10)
			str = append(str, []byte(strconv.FormatUint(uint64(k), 10))...)
			str = append(str, []byte(":")...)
			str = append(str, []byte(strconv.FormatUint(uint64(v), 10))...)
		}
	}
	return string(str)

}

//将 a:b#c#d|a1:b1#c1#d1的格式转换成map
func StringToMapSlice(str string) (map[uint32][]uint32, bool) {
	if str == "" {
		return nil, true
	}
	result_map := make(map[uint32][]uint32)
	str1 := strings.Split(str, "|")
	st1_len := len(str1)
	if st1_len > 0 {
		for _, value := range str1 {
			str2 := strings.Split(value, ":")
			str2_len := len(str2)
			if str2_len > 0 {
				k := uint64(0)
				var e error
				if k, e = strconv.ParseUint(str2[0], 10, 64); e != nil {
					return nil, false
				}
				str3 := strings.Split(str2[1], "#")
				str3_len := len(str3)
				if str3_len > 0 {
					for _, value2 := range str3 {
						value3 := uint64(0)
						if value3, e = strconv.ParseUint(value2, 10, 64); e != nil {
							return nil, false
						}
						if _, exists := result_map[uint32(k)]; exists {
							result_map[uint32(k)] = append(result_map[uint32(k)], uint32(value3))
						} else {
							result_map[uint32(k)] = make([]uint32, 0)
							result_map[uint32(k)] = append(result_map[uint32(k)], uint32(value3))
						}
					}
				}
			}
		}
	}
	return result_map, true
}

func MapSliceToString(result_map map[uint32][]uint32) string {
	//	str := ""
	//var str []byte
	str := make([]byte, 0, 512)
	if result_map == nil || len(result_map) <= 0 {
		return string(str)
	}
	for k, v := range result_map {
		//str2 := ""
		var str2 []byte
		for _, v2 := range v {
			//if str2 != "" {
			if len(str2) > 0 {
				//	str2 = str2 + "#" + strconv.FormatUint(uint64(v2), 10)
				str2 = append(str2, []byte("#")...)
				str2 = append(str2, []byte(strconv.FormatUint(uint64(v2), 10))...)
			} else {
				//	str2 = strconv.FormatUint(uint64(v2), 10)
				str2 = append(str2, []byte(strconv.FormatUint(uint64(v2), 10))...)
			}
		}
		//if str != "" {
		if len(str) > 0 {
			//			str = str + "|" + strconv.FormatUint(uint64(k), 10) + ":" + str2
			str = append(str, []byte("|")...)
			str = append(str, []byte(strconv.FormatUint(uint64(k), 10))...)
			str = append(str, []byte(":")...)
			str = append(str, str2...)
		} else {
			//			str = strconv.FormatUint(uint64(k), 10) + ":" + str2
			str = append(str, []byte(strconv.FormatUint(uint64(k), 10))...)
			str = append(str, []byte(":")...)
			str = append(str, str2...)
		}
	}
	return string(str)

}

//从map 里面随机一个key 出来 根据value的相对概率
func MapRand(result_map map[uint32]uint32) uint32 {
	total_rate := uint32(0)
	for _, v := range result_map {
		total_rate += v
	}
	if total_rate > 0 {
		rand_rate := uint32(rand.Intn(int(total_rate)) + 1)
		now_rate := uint32(0)
		for k, v := range result_map {
			now_rate += v
			if rand_rate <= now_rate {
				return k
			}
		}
	}
	return 0
}

//取范围随机数[min, max)
func Random(min, max int) int {
	//	rand.Seed(time.Now().Unix())
	if min >= max {
		return min
	}
	return rand.Intn(max-min) + min
}

func FromUnixTimeToTodayFormat(t *time.Time) int32 {
	return int32(t.Year()*10000 + int(t.Month())*100 + t.Day())
}

func FromTodayFormatToUnixTime(v int32) *time.Time {
	year := int(v / 10000)
	month := time.Month((v % 10000) / 100)
	day := int(v % 100)
	time := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return &time
}

//从slice 里随机出rand_num个数据 返回的是一个slice
func RandFromInt32Slice(data []uint32, rand_num int) []uint32 {
	result := make([]uint32, 0, rand_num)
	lg := len(data)
	for i := 0; i < rand_num; i++ {
		if lg <= 0 {
			break
		}
		index := rand.Intn(lg)
		result = append(result, data[index])
		if index != lg-1 {
			data[index] = data[lg-1]
		}
		lg--
	}
	return result
}

//从slice 里随机出1个数据
func Uint32SliceRand(data []uint32) uint32 {
	lg := len(data)
	if lg > 0 {
		index := rand.Intn(lg)
		return data[index]
	}
	return uint32(0)

}

func RandFromInt64Slice(data []uint64, rand_num int) []uint64 {
	result := make([]uint64, 0, rand_num)
	lg := len(data)
	for i := 0; i < rand_num; i++ {
		if lg <= 0 {
			break
		}
		index := rand.Intn(lg)
		result = append(result, data[index])
		if index != lg-1 {
			data[index] = data[lg-1]
		}
		lg--
	}
	return result
}

//从0-lg 中随机一组不重复的数据出来
func RandFromZeroToMax(lg int, rand_num int) []int {
	result := make([]int, 0, rand_num)
	data := make([]int, lg)
	for i := 0; i < lg; i++ {
		data[i] = i
	}
	for i := 0; i < rand_num; i++ {
		if lg <= 0 {
			break
		}
		index := rand.Intn(lg)
		result = append(result, data[index])
		if index != lg-1 {
			data[index] = data[lg-1]
		}
		lg--
	}
	return result
}

func BoolToUint32(data bool) uint32 {
	if data {
		return 1
	} else {
		return 0
	}
	return 0
}

func MinInSlice(v []uint32) (m uint32) {
	if len(v) > 0 {
		m = v[0]
	}
	for i := 1; i < len(v); i++ {
		if v[i] < m {
			m = v[i]
		}
	}
	return
}

func Paginate(length int, skip int, size int) (int, int) {
	if skip > length {
		skip = length
	}

	end := skip + size
	if end > length {
		end = length
	}

	return skip, end
}

/*//float32保留N位小数点的截断方法
func CutFloat32(d float32, n uint32) float32 {
	nTmp := float32(math.Pow(10, float64(n)))
	dTmp := int(d * nTmp)
	return float32(dTmp) / nTmp
}*/

//将map 转换成a:b|c:d
func MapToString_64_32(result_map map[uint64]uint32) string {
	str := make([]byte, 0, 512)
	if result_map == nil || len(result_map) <= 0 {
		return string(str)
	}
	for k, v := range result_map {
		if len(str) > 0 {
			str = append(str, []byte("|")...)
			str = append(str, []byte(strconv.FormatUint(uint64(k), 10))...)
			str = append(str, []byte(":")...)
			str = append(str, []byte(strconv.FormatUint(uint64(v), 10))...)
		} else {
			str = append(str, []byte(strconv.FormatUint(uint64(k), 10))...)
			str = append(str, []byte(":")...)
			str = append(str, []byte(strconv.FormatUint(uint64(v), 10))...)

		}
	}
	return string(str)
}

//将 a:b|c:d的格式转换成map
func StringToMap_64_32(str string) (map[uint64]uint32, bool) {
	if str == "" {
		return nil, true
	}
	result_map := make(map[uint64]uint32)
	str1 := strings.Split(str, "|")
	st1_len := len(str1)
	if st1_len > 0 {
		for _, value := range str1 {
			str2 := strings.Split(value, ":")
			str2_len := len(str2)
			if str2_len == 2 {
				k := uint64(0)
				v := uint64(0)
				var e error
				if k, e = strconv.ParseUint(str2[0], 10, 64); e != nil {
					return nil, false
				}
				if v, e = strconv.ParseUint(str2[1], 10, 64); e != nil {
					return nil, false
				}
				result_map[k] = uint32(v)
			} else {
				return nil, false
			}
		}
	}
	return result_map, true
}

//将map 转换成a:b|c:d
func MapToString_64_str(result_map map[uint64]string) string {
	str := make([]byte, 0, 512)
	if result_map == nil || len(result_map) <= 0 {
		return string(str)
	}
	for k, v := range result_map {
		if len(str) > 0 {
			str = append(str, []byte("|")...)
			str = append(str, []byte(strconv.FormatUint(uint64(k), 10))...)
			str = append(str, []byte(",")...)
			str = append(str, []byte(v)...)
		} else {
			str = append(str, []byte(strconv.FormatUint(uint64(k), 10))...)
			str = append(str, []byte(",")...)
			str = append(str, []byte(v)...)

		}
	}
	return string(str)
}

//将 a:b|c:d的格式转换成map
func StringToMap_64_str(str string) (map[uint64]string, bool) {
	if str == "" {
		return nil, true
	}
	result_map := make(map[uint64]string)
	str1 := strings.Split(str, "|")
	st1_len := len(str1)
	if st1_len > 0 {
		for _, value := range str1 {
			str2 := strings.Split(value, ",")
			str2_len := len(str2)
			if str2_len == 2 {
				k := uint64(0)
				var e error
				if k, e = strconv.ParseUint(str2[0], 10, 64); e != nil {
					return nil, false
				}
				result_map[k] = str2[1]
			} else {
				return nil, false
			}
		}
	}
	return result_map, true
}

// s2中的元素是否全在s1中
func ContainSlice(s1, s2 []uint32) bool {
	if len(s1) < len(s2) {
		return false
	}
	for _, val2 := range s2 {
		has := false
		for _, val1 := range s1 {
			if val1 == val2 {
				has = true
				break
			}
		}
		if !has {
			return false
		}
	}
	return true
}
