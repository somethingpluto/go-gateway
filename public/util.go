package public

import (
	"crypto/sha256"
	"fmt"
)

// GenSaltPassword
// @Description: 密码加密
// @param salt 盐
// @param password 密码
// @return string 加密后的密码字符串
//
func GenSaltPassword(salt string, password string) string {
	s1 := sha256.New()
	s1.Write([]byte(password))
	str1 := fmt.Sprintf("%x", s1.Sum(nil))
	s2 := sha256.New()
	s2.Write([]byte(str1 + salt))
	return fmt.Sprintf("%x", s2.Sum(nil))
}
