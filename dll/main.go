package main

import "C"
import (
	"encoding/json"
	"golang.org/x/text/encoding/simplifiedchinese"
	"syscall"
	"unsafe"
)

var h = syscall.NewLazyDLL("termb.dll")
var funcCVRInitComm = h.NewProc("CVR_InitComm")

var funcCVRAuthenticate = h.NewProc("CVR_Authenticate")
var funcCVRReadContent = h.NewProc("CVR_Read_Content")

//var funcCVRCloseComm = h.NewProc("CVR_CloseComm")

var funcGetPeopleName = h.NewProc("GetPeopleName")
var funcGetPeopleSex = h.NewProc("GetPeopleSex")
var funcGetPeopleIDCode = h.NewProc("GetPeopleIDCode")
var funcGetbase64BMPData = h.NewProc("Getbase64BMPData")

var funcGetPeopleNation = h.NewProc("GetPeopleNation")
var funcGetPeopleBirthday = h.NewProc("GetPeopleBirthday")
var funcGetPeopleAddress = h.NewProc("GetPeopleAddress")
var funcGetDepartment = h.NewProc("GetDepartment")
var funcGetStartDate = h.NewProc("GetStartDate")
var funcGetEndDate = h.NewProc("GetEndDate")

//export open
func open() int {
	s1, _, _ := funcCVRInitComm.Call(uintptr(1001))
	return int(s1)
}

//export read
func read() int {
	s3, _, _ := funcCVRAuthenticate.Call()
	return int(s3)
}

//export get
func get() *C.char {
	funcCVRReadContent.Call(uintptr(4))
	img := parse(funcGetbase64BMPData)
	name := parse(funcGetPeopleName)
	idCode := parse(funcGetPeopleIDCode)
	sex := parse(funcGetPeopleSex)
	nation := parse(funcGetPeopleNation)
	birthday := parse(funcGetPeopleBirthday)
	address := parse(funcGetPeopleAddress)
	department := parse(funcGetDepartment)
	startDate := parse(funcGetStartDate)
	endDate := parse(funcGetEndDate)

	idCard := IdCard{
		img,
		name,
		idCode,
		sex,
		nation,
		birthday,
		address,
		department,
		startDate,
		endDate,
	}

	idCardJson, _ := json.Marshal(idCard)
	return C.CString(string(idCardJson))
}

type IdCard struct {
	Img        string `json:"img"`
	Name       string `json:"name"`
	IdCode     string `json:"idCode"`
	Sex        string `json:"sex"`
	Nation     string `json:"nation"`
	Birthday   string `json:"birthday"`
	Address    string `json:"address"`
	Department string `json:"department"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
}

func parse(proc *syscall.LazyProc) string {
	var content = make([]byte, 1024*100)
	var length int
	proc.Call(uintptr(unsafe.Pointer(&content[0])), uintptr(unsafe.Pointer(&length)))
	reply, _ := simplifiedchinese.GBK.NewDecoder().Bytes(content[:length])
	return string(reply)
}

func main() {
}
