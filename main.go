package extractEmbedFile

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

func IsExistFile(fname string) bool {
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetBufHash(buf []byte) string {
	h := sha256.New()
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func GetFileHash(fname string) string {
	f, err := os.Open(fname)

	if err != nil {
		return ""
	}

	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}

	return hex.EncodeToString(h.Sum(nil))
}

func Extract(basePath string, fname string, data []byte, isExcute bool, checkUpdate bool) error {
	fp := filepath.Join(basePath, fname)

	if !IsExistFile(fp) {
		var werr error

		if isExcute {
			werr = os.WriteFile(fp, data, 0755)
		} else {
			werr = os.WriteFile(fp, data, 0644)
		}

		if werr != nil {
			return werr
		}
	} else if checkUpdate {

		oldHash := GetFileHash(fp)
		newHash := GetBufHash(data)

		// 다른 파일
		if oldHash != newHash {
			_ = os.Remove(fp)

			var werr error

			if isExcute {
				werr = os.WriteFile(fp, data, 0755)
			} else {
				werr = os.WriteFile(fp, data, 0644)
			}

			if werr != nil {
				return werr
			}
		}
	}

	return nil
}
