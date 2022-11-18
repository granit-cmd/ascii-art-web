package ascii_art

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func Draw(text string, banner string) (string, error) {
	// Записанные ключи шрифтов
	checkSumArr := []string{
		"a51f800619146db0c42d26db3114c99f",
		"c5b7947f1ac707b55c9b4d4e12994165",
		"93e20c2510dfd28993db87352915826a",
	}
	banner += ".txt"

	checkSumIsRight := false

	// Проверка чексуммы файлов шрифтов
	for i := range checkSumArr {
		temp := CheckSum(banner)
		if checkSumArr[i] == temp {
			checkSumIsRight = true
			break
		}
	}
	if !checkSumIsRight {
		return "", fmt.Errorf("invalid banner")
	}

	fontText := FontReader(banner)

	arrText, err := WriteText(text, fontText)
	if err != nil {
		return "", err
	}
	result := ""
	for i := 0; i < len(arrText); i++ {
		result += arrText[i] + "\n"
	}

	return result, nil
}

func AppendInResult(resultArr []string, textPrint [8]string) ([]string, [8]string) {
	if !isArrEmpty(textPrint) {
		for j := range textPrint {
			resultArr = append(resultArr, textPrint[j])
			textPrint[j] = ""
		}
	} else {
		resultArr = append(resultArr, " \n")
	}

	return resultArr, textPrint
}

func WriteText(text string, fontStringArr []string) ([]string, error) {
	textPrint := [8]string{}
	var resultArr []string

	lastIndex := len(text) - 1
	for i := range text {
		// Переходы на следующую строку
		if (i != lastIndex && text[i] == '\\' && text[i+1] == 'n') ||
			text[i] == '\n' {
			resultArr, textPrint = AppendInResult(resultArr, textPrint)

			if (i == lastIndex-1 && i != 0 && text[i] == '\\' &&
				text[i+1] == 'n') ||
				(i == lastIndex && i != 0 && text[i] == '\n') {
				resultArr = append(resultArr, "\n")
			}
		} else if i == 0 || (text[i-1] != '\\' || text[i] != 'n') {
			for j := range textPrint {
				if text[i] == 13 {
					resultArr = append(resultArr, " \n")
					break
				}
				if (text[i] > 126 || text[i] < 32) && text[i] != '\n' {
					return []string{}, fmt.Errorf("invalid character in text")
				}

				textPrint[j] += fontStringArr[j+1+(int(text[i])-32)*9]
			}
		}
	}
	resultArr, _ = AppendInResult(resultArr, textPrint)

	return resultArr, nil
}

func CheckSum(text string) string {
	h := md5.New()
	f, err := os.Open("ascii-art/" + text)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = io.Copy(h, f)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func isArrEmpty(text [8]string) bool {
	for i := range text {
		if text[i] != "" {
			return false
		}
	}
	return true
}

func FontReader(banner string) []string {
	// Открывает файл шрифта соответствующего названия
	file, _ := os.Open("ascii-art/" + banner)
	fileScanner := bufio.NewScanner(file)

	fontStringArr := []string{}

	// Запись текста из шрифта в массив построчно
	for fileScanner.Scan() {
		fontStringArr = append(fontStringArr, fileScanner.Text())
	}

	file.Close()
	return fontStringArr
}
