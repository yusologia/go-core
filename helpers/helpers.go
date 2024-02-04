package helpers

import (
	"math/rand"
	"os"
	"strconv"
	"time"
)

func RandomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	randomBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		randomBytes[i] = chars[rand.Intn(len(chars))]
	}

	return string(randomBytes) + strconv.FormatInt(time.Now().UnixNano(), 10)
}

func CheckAndCreateDirectory(path string) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(path, os.ModePerm)
		}
	}
}

func SetStorageDir(path ...string) string {
	storagePath := os.Getenv("STORAGE_DIR")
	if len(storagePath) == 0 {
		storagePath = "storages"
	}

	if len(path) > 0 {
		storagePath += "/" + path[0]
	}

	return storagePath
}

func SetStorageAppDir(path ...string) string {
	appDir := "app"
	if len(path) > 0 {
		appDir += "/" + path[0]
	}

	return SetStorageDir(appDir)
}

func SetStorageAppPublicDir(path ...string) string {
	publicDir := "app/public"
	if len(path) > 0 {
		publicDir += "/" + path[0]
	}

	return SetStorageDir(publicDir)
}
