package tools

import "github.com/mozillazg/go-pinyin"

func GenFirstPinyin(text string) (py string) {
	for _, chars := range pinyin.Pinyin(text, pinyin.NewArgs()) {
		for _, char := range chars {
			if len(char) > 0 {
				py = py + string(char[0])
			}
		}
	}
	return
}
