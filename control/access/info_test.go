package access

import (
	"testing"
)

// 测试函数
func TestGetResult(t *testing.T) {
	const name = "Ubuntu22"

	info, err := GetAccessInfo(name)
	if err != nil {
		t.Fatalf("Failed to get access info: %v", err)
	}

	// 这里可以添加更多的断言来验证返回的数据
	t.Logf("Access Info: %+v", info)
}

// TestGetInfoById 测试通过 ID 查找信息的函数
func TestGetInfoById(t *testing.T) {
	tests := []struct {
		id       string
		expected Info
		hasError bool
	}{
		{
			id: "yyy1",
			expected: Info{
				Name: "Ubuntu25",
				Host: "192.168.235.128",
				Port: 1080,
			},
			hasError: false,
		},
		{
			id: "yyy2",
			expected: Info{
				Name: "Ubuntu22",
				Host: "192.168.235.133",
				Port: 1080,
			},
			hasError: false,
		},
		{
			id:       "invalid_id",
			expected: Info{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			info, err := GetInfoById(tt.id)

			if tt.hasError {
				if err == nil {
					t.Fatalf("Expected error for ID %s, but got none", tt.id)
				}
				return
			}

			if err != nil {
				t.Fatalf("Failed to get info for ID %s: %v", tt.id, err)
			}

			// 断言返回的值是否与预期一致
			if info != tt.expected {
				t.Errorf("Expected %+v, but got %+v", tt.expected, info)
			}

			t.Logf("Info for ID %s: %+v", tt.id, info)
		})
	}
}
