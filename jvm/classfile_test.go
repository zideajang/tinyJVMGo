package jvm

import (
	"os"
	"path/filepath"
	"strings" // 导入 strings 包
	"testing"
)

// createTestFile 创建一个带有给定内容的临时文件。
// 它返回创建的文件的路径。
func createTestFile(t *testing.T, filename string, content []byte) string {
	dir := t.TempDir() // 为测试创建一个临时目录
	filePath := filepath.Join(dir, filename)

	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		t.Fatalf("创建测试文件 %s 失败: %v", filePath, err)
	}
	return filePath
}

func TestReadClassMagic(t *testing.T) {
	tests := []struct {
		name          string
		magicBytes    []byte
		expectedMagic uint32
		expectedErr   bool
		errSubstring  string // 我们期望在错误信息中出现的子字符串
		osSpecificErr bool   // 如果错误信息可能因操作系统而异，则为 true
	}{
		{
			name:          "有效的魔数",
			magicBytes:    []byte{0xCA, 0xFE, 0xBA, 0xBE},
			expectedMagic: MagicNumber,
			expectedErr:   false,
		},
		{
			name:          "无效的魔数",
			magicBytes:    []byte{0xDE, 0xAD, 0xBE, 0xEF},
			expectedMagic: 0xDEADBEEF,
			expectedErr:   false,
		},
		{
			name:          "文件过小",
			magicBytes:    []byte{0xCA, 0xFE, 0xBA}, // 只有 3 个字节
			expectedMagic: 0,
			expectedErr:   true,
			errSubstring:  "too small", // 这部分在不同操作系统上是一致的
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := createTestFile(t, "test.class", tt.magicBytes)

			magic, err := ReadClassMagic(filePath)

			if tt.expectedErr {
				if err == nil {
					t.Errorf("期望出现包含 '%s' 的错误，但未出现", tt.errSubstring)
				} else if tt.errSubstring != "" && !strings.Contains(err.Error(), tt.errSubstring) { // 使用 strings.Contains
					t.Errorf("期望错误包含 '%s'，但得到 '%v'", tt.errSubstring, err)
				}
			} else {
				if err != nil {
					t.Errorf("不期望出现错误，但得到: %v", err)
				}
				if magic != tt.expectedMagic {
					t.Errorf("期望魔数为 0x%X，但得到 0x%X", tt.expectedMagic, magic)
				}
			}
		})
	}

	t.Run("不存在的文件", func(t *testing.T) {
		_, err := ReadClassMagic("non_existent_file.class")
		if err == nil {
			t.Errorf("期望不存在的文件出现错误，但未出现")
		}

		// 检查操作系统特定的错误信息
		var expectedErrSubstring string
		if os.IsNotExist(err) { // 这是检查文件不存在的最健壮方法
			// 对于文件不存在的错误，通常错误信息中会包含 "file" 这个词
			// 或者在 Windows 上会包含 "The system cannot find the file specified."
			// 我们可以选择一个通用的部分，或者同时检查两种情况
			if strings.Contains(err.Error(), "The system cannot find the file specified.") { // Windows 特定错误
				expectedErrSubstring = "The system cannot find the file specified."
			} else {
				// 适用于 Linux/macOS 和其他通用情况
				expectedErrSubstring = "no such file or directory"
			}
		}

		// 如果 os.IsNotExist 为 false (即不是文件不存在错误，而是其他打开错误)，
		// 那么我们可能需要更通用的错误检查，这里为了简化，我们只关注文件不存在的情况
		if expectedErrSubstring != "" && !strings.Contains(err.Error(), expectedErrSubstring) {
			t.Errorf("期望错误包含 '%s'，但得到: %v", expectedErrSubstring, err)
		}
	})
}

func TestVerifyClassMagic(t *testing.T) {
	tests := []struct {
		name         string
		magicBytes   []byte
		expectedOk   bool
		expectedErr  bool
		errSubstring string
	}{
		{
			name:        "有效的魔数",
			magicBytes:  []byte{0xCA, 0xFE, 0xBA, 0xBE},
			expectedOk:  true,
			expectedErr: false,
		},
		{
			name:        "无效的魔数",
			magicBytes:  []byte{0xDE, 0xAD, 0xBE, 0xEF},
			expectedOk:  false,
			expectedErr: false,
		},
		{
			name:        "文件过小",
			magicBytes:  []byte{0xCA, 0xFE, 0xBA},
			expectedOk:  false,
			expectedErr: true,
			errSubstring: "too small",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := createTestFile(t, "test.class", tt.magicBytes)

			ok, err := VerifyClassMagic(filePath)

			if tt.expectedErr {
				if err == nil {
					t.Errorf("期望出现错误，但未出现")
				} else if tt.errSubstring != "" && !strings.Contains(err.Error(), tt.errSubstring) { // 使用 strings.Contains
					t.Errorf("期望错误包含 '%s'，但得到 '%v'", tt.errSubstring, err)
				}
			} else {
				if err != nil {
					t.Errorf("不期望出现错误，但得到: %v", err)
				}
				if ok != tt.expectedOk {
					t.Errorf("期望 ok 为 %t，但得到 %t", tt.expectedOk, ok)
				}
			}
		})
	}
}