package main

import (
	"fmt"
	"github.com/atsushinee/golang-win-read-id-card/idll"
	"golang.org/x/text/encoding/simplifiedchinese"
	"syscall"
	"time"
	"unsafe"
)

var index = 0
var progress = []string{`\`, `|`, `-`, `/`, `-`}

func getProgress() string {
	index++
	if index == 5 {
		index = 0
	}
	return progress[index]

}

func receiveIdCard() {
	idll.RestoreAsset(".", "UnPack.dll")
	idll.RestoreAsset(".", "WltRS.dll")
	idll.RestoreAsset(".", "sdtapi.dll")
	idll.RestoreAsset(".", "termb.dll")

	h := syscall.NewLazyDLL("termb.dll")

	funcCVRInitComm := h.NewProc("CVR_InitComm")
	s1, _, _ := funcCVRInitComm.Call(uintptr(1001))
	fmt.Println("初始化读卡器:", s1)
	funcCVRAuthenticate := h.NewProc("CVR_Authenticate")
	funcCVRReadContent := h.NewProc("CVR_Read_Content")
	// funcCVRCloseComm := h.NewProc("CVR_CloseComm")

	funcGetPeopleName := h.NewProc("GetPeopleName")
	funcGetPeopleSex := h.NewProc("GetPeopleSex")
	funcGetPeopleIDCode := h.NewProc("GetPeopleIDCode")
	funcGetbase64BMPData := h.NewProc("Getbase64BMPData")

	funcGetPeopleNation := h.NewProc("GetPeopleNation")
	funcGetPeopleBirthday := h.NewProc("GetPeopleBirthday")
	funcGetPeopleAddress := h.NewProc("GetPeopleAddress")
	funcGetDepartment := h.NewProc("GetDepartment")
	funcGetStartDate := h.NewProc("GetStartDate")
	funcGetEndDate := h.NewProc("GetEndDate")

	var isReading = false
	var auth int
	for {
		auth = 0
		for auth != 1 {
			s3, _, _ := funcCVRAuthenticate.Call()
			auth = int(s3)
			if auth != 1 {
				isReading = false
			}
			if isReading {
				auth = 0
			}
			time.Sleep(time.Millisecond * 300)
			fmt.Printf("读卡中:%s\r", getProgress())
		}
		if auth == 1 {
			isReading = true
			fmt.Println("读卡成功!")
			funcCVRReadContent.Call(uintptr(4))

			img := parse(funcGetbase64BMPData)
			fmt.Println(img)

			name := parse(funcGetPeopleName)
			fmt.Println(name)

			idCode := parse(funcGetPeopleIDCode)
			fmt.Println(idCode)

			sex := parse(funcGetPeopleSex)
			fmt.Println(sex)

			nation := parse(funcGetPeopleNation)
			fmt.Println(nation)

			birthday := parse(funcGetPeopleBirthday)
			fmt.Println(birthday)

			address := parse(funcGetPeopleAddress)
			fmt.Println(address)

			department := parse(funcGetDepartment)
			fmt.Println(department)

			startDate := parse(funcGetStartDate)
			fmt.Println(startDate)

			endDate := parse(funcGetEndDate)
			fmt.Println(endDate)
		}
	}
}

func parse(proc *syscall.LazyProc) string {
	var content = make([]byte, 1024*100)
	var length int
	proc.Call(uintptr(unsafe.Pointer(&content[0])), uintptr(unsafe.Pointer(&length)))
	reply, _ := simplifiedchinese.GBK.NewDecoder().Bytes(content[:length])
	return string(reply)
}

const AuthorTag = `+---------------------+
|    AUTHOR:LICHUN    |
+---------------------+
`

func main() {
	fmt.Print(AuthorTag)
	receiveIdCard()
}
