package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"github.com/oklog/ulid/v2"
	math "math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

// 密码加密参数
const (
	saltLength  = 16
	keyLength   = 32
	iterations  = 3
	memory      = 64 * 1024
	parallelism = uint8(2) // 修改为 uint8 类型
)

// HashPassword 使用 Argon2id 算法对密码进行加密
func HashPassword(password string) (string, error) {
	// 生成随机盐值
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 使用 Argon2id 算法加密密码
	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	// 将参数和结果编码为字符串
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	// 格式: $argon2id$v=19$m=65536,t=3,p=2$<salt>$<hash>
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		memory, iterations, parallelism, encodedSalt, encodedHash), nil
}

// VerifyPassword 验证密码是否匹配
func VerifyPassword(hashedPassword, password string) bool {
	// 解析哈希字符串
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return false
	}

	var memory, iterations, parallelism uint32
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	// 使用相同参数计算输入密码的哈希值
	inputHash := argon2.IDKey([]byte(password), salt, iterations, memory, uint8(parallelism), uint32(len(decodedHash)))

	// 比较两个哈希值
	return subtle.ConstantTimeCompare(decodedHash, inputHash) == 1
}

// GenUUID 生成UUID作为用户ID
func GenUUID() string {
	entropy := ulid.Monotonic(math.New(math.NewSource(time.Now().UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}
