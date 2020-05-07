package main



import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
        "unsafe"
)

type regParseInfo struct {
	src       string
	value     string
	isKeyOnly bool
	key       string
	param     string
	valtype   uint32
	mainKey   uintptr
}

const (
	_REG_OPTION_NON_VOLATILE = 0

	_REG_CREATED_NEW_KEY     = 1
	_REG_OPENED_EXISTING_KEY = 2
)

const (
	// Registry value types.
	NONE                       = 0
	SZ                         = 1
	EXPAND_SZ                  = 2
	BINARY                     = 3
	DWORD                      = 4
	DWORD_BIG_ENDIAN           = 5
	LINK                       = 6
	MULTI_SZ                   = 7
	RESOURCE_LIST              = 8
	FULL_RESOURCE_DESCRIPTOR   = 9
	RESOURCE_REQUIREMENTS_LIST = 10
	QWORD                      = 11
)

var regParamFormats = map[string]uint32{
	"SZ":    SZ,
	"DWORD": DWORD,
	"QWORD": QWORD}

func operationCopy(options []string) {
	n := len(options)
	if n == 0 || n&1 != 0 {
		fmt.Printf("Copy requires even number of arguments")
		return
	}
	for i := 0; i < n; i += 2 {
		src := options[i]
		dst := options[i+1]
		buf, err := ioutil.ReadFile(src)
		if err != nil {
			fmt.Printf("Error reading %s: %s", src, err.Error())
			continue
		}
		err = ioutil.WriteFile(dst, buf, 0644)
		if err != nil {
			fmt.Printf("Error writing %s: %s", dst, err.Error())
		} else {
			fmt.Printf("Copy to %s successful", dst)
		}
	}
}

func operationRegParse(src, value string) (regParseInfo, error) {
	//example: HKEY_LOCAL_MACHINE\SYSTEM\ControlSet001\services\Tcpip\Parameters~SZ~SearchList  dat1.com,dat2.com
	r := regParseInfo{src: src, value: value, isKeyOnly: true, key: src}
	for f, ft := range regParamFormats {
		p := strings.Index(src, "~"+f+"~")
		if p >= 0 {
			r.key = src[:p]
			r.param = src[p+2+len(f):]
			r.isKeyOnly = false
			r.valtype = ft
			break
		}
	}
	r.mainKey = syscall.HKEY_LOCAL_MACHINE
	r.key = strings.Replace(r.key, "/", "\\", -1)
	p1 := strings.Index(r.key, "\\")
	if p1 <= 0 {
		return r, errors.New("The first parameter in key must be HKEY: " + r.key)
	}
	mainKeyName := r.key[:p1]
	r.key = r.key[p1+1:]
	if len(mainKeyName) > 4 && mainKeyName[:5] == "HKEY_" {
		mainKeyName = mainKeyName[5:]
	}
	switch mainKeyName {
	case "LOCAL_MACHINE":
	case "CLASSES_ROOT":
		r.mainKey = syscall.HKEY_CLASSES_ROOT
	case "CURRENT_USER":
		r.mainKey = syscall.HKEY_CURRENT_USER
	case "USERS":
		r.mainKey = syscall.HKEY_USERS
	case "CURRENT_CONFIG":
		r.mainKey = syscall.HKEY_CURRENT_CONFIG
	case "PERFORMANCE_DATA":
		r.mainKey = syscall.HKEY_PERFORMANCE_DATA
	default:
		return r, errors.New("Unknown HKEY: " + mainKeyName + " for key " + r.key)
	}
	return r, nil
}

type Key syscall.Handle

func CreateKey(k uintptr, path string, access uint32) (newk Key, openedExisting bool, err error) {
		var h syscall.Handle
		var d uint32
		err = regCreateKeyEx(syscall.Handle(k), syscall.StringToUTF16Ptr(path),
			0, nil,_REG_OPTION_NON_VOLATILE, access, nil, &h, &d)
		if err != nil {
			return 0, false, err
		}
		return Key(h), d == _REG_OPENED_EXISTING_KEY, nil
}

func OpenKey(k uintptr, path string, options uint32, desiredAccess uint32) (Key, error) {
	var h syscall.Handle
	err := syscall.RegOpenKeyEx(syscall.Handle(k), syscall.StringToUTF16Ptr(path),
		options, desiredAccess, &h)
	if err != nil {
		return 0, err
	}
	return Key(h), nil
}

func (k Key) setStringValue(name string, valtype uint32, value string) error {
	v, err := syscall.UTF16FromString(value)
	if err != nil {
		return err
	}
	buf := (*[1 << 29]byte)(unsafe.Pointer(&v[0]))[:len(v)*2]
	return k.setValue(name, valtype, buf)
}

func (k Key) setValue(name string, valtype uint32, data []byte) error {
	p, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		regSetValueEx(syscall.Handle(k), p, 0, valtype, nil, 0)
	}
	return regSetValueEx(syscall.Handle(k), p, 0, valtype, &data[0], uint32(len(data)))
}

func (k Key) Close() {
	syscall.RegCloseKey(syscall.Handle(k))
}

func (k Key) SetValueByKind(param, value string, valtype uint32) error {
	switch valtype {
	case SZ:
		return k.setStringValue(param, valtype, value)
	case DWORD:
	case QWORD:
	default:
		return errors.New("Unknown value type " + strconv.Itoa(int(valtype)) + " to set " + param + " to " + value)
	}
	return nil
}

func operationRegadd(options []string) {
	n := len(options)
	if n == 0 || n&1 != 0 {
		fmt.Printf("Regadd requires even number of arguments")
		return
	}
	for i := 0; i < n; i += 2 {
		r, err1 := operationRegParse(options[i], options[i+1])
		if err1 != nil {
			fmt.Printf("Error: %s", err1.Error())
			continue
		}
		if r.isKeyOnly {
			k, exist, err := CreateKey(r.mainKey, r.key, syscall.KEY_CREATE_SUB_KEY)
			if err != nil {
				fmt.Printf("Registry: error creating key=%s  error=%s", r.src, err.Error())
			} else {
				if exist {
					fmt.Printf("Registry: already exists key=%s", r.src)
				} else {
					fmt.Printf("Registry: successfully created key=%s", r.src)
				}
				k.Close()
			}
			continue
		}
		k, err := OpenKey(r.mainKey, r.key, 0, syscall.KEY_WRITE)
		if err != nil {
			fmt.Printf("Registry: Error opening key=%s param=%s value=%s error=%s", r.src, r.param, r.value, err.Error())
			continue
		}
		err = k.SetValueByKind(r.param, r.value, r.valtype)
		k.Close()
		if err != nil {
			fmt.Printf("Error writing to register key=%s param=%s value=%s error=%s", r.src, r.param, r.value, err.Error())
		} else {
			fmt.Printf("Reg %s param %s value %s successfully added", r.key, r.param, r.value)
		}
	}

}

func operationRegdel(options []string) {
}

func main() {
	argsLen := len(os.Args)
	if argsLen < 2 {
		fmt.Printf("DvAdmin requires args")
		return
	}
	for i := 1; i < argsLen; i++ {
		name := os.Args[i]
		p := strings.IndexByte(name, '_')
		if p <= 0 {
			fmt.Printf("Illegal operation combined name %s (no underline in the middle)", name)
			return
		}
		cntStr, err := strconv.Atoi(name[p+1:])
		if err != nil {
			fmt.Printf("Expected number at the end of the operation name: %s", name)
			return
		}
		if i+cntStr >= argsLen {
			fmt.Printf("Not sufficient number of params: expected %d but found %d in %q", i+cntStr, argsLen-1, os.Args)
			return
		}
		options := make([]string, cntStr)
		for j := 0; j < cntStr; j++ {
			options[j] = os.Args[i+j+1]
		}
		i += cntStr
		switch strings.ToLower(name[:p]) {
		case "copy":
			operationCopy(options)
		case "regadd":
			operationRegadd(options)
		case "regdel":
			operationRegdel(options)
		default:
			fmt.Printf("Unsupported operation: %s", name)
		}
	}

}

// windows api calls
//sys	regSetValueEx(key syscall.Handle, valueName *uint16, reserved uint32, vtype uint32, buf *byte, bufsize uint32) (regerrno error) = advapi32.RegSetValueExW

var (
	modadvapi32 = syscall.NewLazyDLL(syscall.SysdllAdd("advapi32.dll"))

	procRegCreateKeyExW           = modadvapi32.NewProc("RegCreateKeyExW")
	procRegDeleteKeyW             = modadvapi32.NewProc("RegDeleteKeyW")
	procRegSetValueExW            = modadvapi32.NewProc("RegSetValueExW")
)

func regSetValueEx(key syscall.Handle, valueName *uint16, reserved uint32, vtype uint32, buf *byte, bufsize uint32) (regerrno error) {
	r0, _, _ := syscall.Syscall6(procRegSetValueExW.Addr(), 6, uintptr(key), uintptr(unsafe.Pointer(valueName)), uintptr(reserved), uintptr(vtype), uintptr(unsafe.Pointer(buf)), uintptr(bufsize))
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}

func regCreateKeyEx(key syscall.Handle, subkey *uint16, reserved uint32, class *uint16, options uint32, desired uint32, sa *syscall.SecurityAttributes, result *syscall.Handle, disposition *uint32) (regerrno error) {
	r0, _, _ := syscall.Syscall9(procRegCreateKeyExW.Addr(), 9, uintptr(key), uintptr(unsafe.Pointer(subkey)), uintptr(reserved), uintptr(unsafe.Pointer(class)), uintptr(options), uintptr(desired), uintptr(unsafe.Pointer(sa)), uintptr(unsafe.Pointer(result)), uintptr(unsafe.Pointer(disposition)))
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}

