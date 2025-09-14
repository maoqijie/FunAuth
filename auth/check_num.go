package auth

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	g79 "github.com/Yeah114/g79client"
	"github.com/Yeah114/unmcpk"
)

func md5Hex(parts ...string) string {
	h := md5.New()
	for _, p := range parts {
		_, _ = h.Write([]byte(p))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func extractS1S2FromSource(code string) (string, string, error) {
	re := regexp.MustCompile(`message = '([^']*)' \+ data \+ '([^']*)'`)
	m := re.FindStringSubmatch(code)
	if len(m) != 3 {
		return "", "", fmt.Errorf("pattern not found")
	}
	/*
	s1 := m[1]
	s2 := m[2]
	fmt.Println(s1)
	fmt.Println(len(s1))
	fmt.Println(s2)
	fmt.Println(len(s2))
	*/
	return m[1], m[2], nil
}

func extractS1S2FromRepairedByRegex(repaired []byte) (string, string, error) {
	startIndex := 251
	endIndex := 0
	for i := startIndex; i < len(repaired); i++ {
		if repaired[i] == 0x28 {
			endIndex = i
			break
		}
	}
	target := repaired[startIndex:endIndex]
	re := regexp.MustCompile(`(?s)s.\x00{3}`)
	matches := re.FindIndex(target)
	if matches == nil {
		return "", "", fmt.Errorf("pattern not found")
	}
	s1 := string(target[:matches[0]])
	s2 := string(target[matches[1]:])
	/*
	fmt.Println(target)
	fmt.Println(s1)
	fmt.Println(len(s1))
	fmt.Println(s2)
	fmt.Println(len(s2))
	*/
	return s1, s2, nil
}

// TransferCheckNum
func TransferCheckNum(ctx context.Context, data string) (string, error) {
	var arr []any
	if err := json.Unmarshal([]byte(data), &arr); err != nil || len(arr) < 3 {
		return "", fmt.Errorf("bad data")
	}
	val, _ := arr[1].(string)
	uniqueFloat, _ := arr[2].(float64)
	uniqueID := int64(uniqueFloat)
	mcpHex, _ := arr[0].(string)
	mcpBytes, err := hex.DecodeString(mcpHex)
	if err != nil {
		return "", fmt.Errorf("bad mcp hex")
	}

	// 使用正则的分支
	useRegex := true
	var s1, s2 string
	if useRegex {
		repaired, err := unmcpk.DecryptDynamicMCP(mcpBytes)
		if err != nil {
			return "", fmt.Errorf("decompile failed")
		}
		s1, s2, err = extractS1S2FromRepairedByRegex(repaired)
		if err != nil {
			return "", fmt.Errorf("pattern not found")
		}
	} else {
		python3Path := os.Getenv("FUNAUTH_PYTHON3")
		if python3Path == "" {
			python3Path = "python3"
		}
		source, err := unmcpk.DecompileDynamicMCP(mcpBytes, python3Path)
		if err != nil {
			return "", fmt.Errorf("decompile failed")
		}
		s1, s2, err = extractS1S2FromSource(source)
		if err != nil {
			return "", fmt.Errorf("pattern not found")
		}
	}

	engineVersion := g79.EngineVersion
	patchVersion, err := g79.GetGlobalLatestVersion()
	if err != nil {
		return "", fmt.Errorf("get latest version failed")
	}
	if engineVersion == "" {
		engineVersion = g79.EngineVersion
	}

	valm := md5Hex(s1, val+"0", s2)
	tmps := make([]string, 0, len(valm)+7)
	for _, ch := range valm {
		tmps = append(tmps, fmt.Sprintf("%d", ((int(ch))*2+5)^255))
	}
	tmps = append(tmps, engineVersion, "android", patchVersion, "android", "2", "12", fmt.Sprintf("%d", uniqueID))
	//tmps = append(tmps, engineVersion, "iOS", patchVersion, "ios", "1", "12", fmt.Sprintf("%d", uniqueID))
	tmpsStr := ""
	for _, s := range tmps {
		tmpsStr += s
	}
	tmpsNum := md5Hex(s1, tmpsStr, s2)

	raw := valm[16:] + "False[]3" + tmpsNum + valm[:16]
	sign := md5Hex(s1, raw, s2)

	valueJSON := fmt.Sprintf("[\"%s\",\"%s\",false,[],\"\",\"\",3,\"%s\"]", valm, sign, tmpsNum)
	return valueJSON, nil
}
