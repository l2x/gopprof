package util

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// GetBinFile return current running process file
func GetBinFile() (string, []byte, error) {
	bf, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", nil, err
	}
	b, err := ioutil.ReadFile(bf)
	if err != nil {
		return bf, nil, err
	}
	return bf, b, nil
}

// GetNetInterfaceIP return net interface ip
func GetNetInterfaceIP() ([]string, []string, error) {
	var internalIP, externalIP []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, nil, err
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			ipv4 := ipnet.IP.To4()
			if ipv4 == nil {
				continue
			}
			if isInternalIP(ipv4.String()) {
				internalIP = append(internalIP, ipv4.String())
			} else {
				externalIP = append(externalIP, ipv4.String())
			}
		}
	}
	return internalIP, externalIP, nil
}

func isInternalIP(ipStr string) bool {
	if strings.HasPrefix(ipStr, "10.") || strings.HasPrefix(ipStr, "192.168.") {
		return true
	}
	if strings.HasPrefix(ipStr, "172.") {
		// 172.16.0.0-172.31.255.255
		arr := strings.Split(ipStr, ".")
		if len(arr) != 4 {
			return false
		}
		second, err := strconv.ParseInt(arr[1], 10, 64)
		if err != nil {
			return false
		}
		if second >= 16 && second <= 31 {
			return true
		}
	}
	return false
}

// InStringSlice return true if s in arr
func InStringSlice(s string, arr []string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

// PostFile post a file
func PostFile(uri, file string, params map[string]string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, f)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	contentType := writer.FormDataContentType()
	writer.Close()

	resp, err := http.Post(uri, contentType, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

// PostData post data
func PostData(uri, fname string, data []byte, params map[string]string) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fname)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	contentType := writer.FormDataContentType()
	writer.Close()

	resp, err := http.Post(uri, contentType, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}
