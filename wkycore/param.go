package wkycore

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

/*
md5LowerString sum md5 hash code with lower case
*/
func md5LowerString(s string) string {
	b := md5.Sum([]byte(s))
	str := hex.EncodeToString(b[:])

	return strings.ToLower(str)
}

/*
getDevID generate device id
*/
func GetDevID(phone string) string {
	s := md5LowerString(phone)

	//Convert to upper case
	return strings.ToLower(s[0:14])
}

/*
getPWD Get getPWD string
*/
func GetPWD(text string) string {
	s := md5LowerString(text)

	str := s[0:2] + string(s[8]) + s[3:8] + string(s[2]) + s[9:17] +
		string(s[27]) + s[18:27] + string(s[17]) + s[28:]

	return md5LowerString(str)
}

/*
getSign calculate the sign via some config
*/
func GetSign(isGet bool, body map[string]string, session string) string {
	var list []string

	//Generate list
	for k, v := range body {
		list = append(list, k+"="+v)
	}

	//Sort
	sort.Strings(list)

	//Join
	s := strings.Join(list, "&") + "&key=" + session

	return md5LowerString(s)
}

/*
getIMEI generate IMEI via phone number, it's not a real imem number
*/
func GetIMEI(phone string) string {
	s := md5LowerString(phone)

	//Convert to upper case
	return strings.ToLower(s[0:16])
}

func GenIpaddr() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}
