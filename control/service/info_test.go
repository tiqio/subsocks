package service

import (
	"testing"
)

// 测试函数
func TestGetResult(t *testing.T) {
	const name = "IPSB"

	info, err := GetServiceInfo(name)
	if err != nil {
		t.Fatalf("Failed to get service info: %v", err)
	}

	// 这里可以添加更多的断言来验证返回的数据
	t.Logf("Service Info: %+v", info)
}

// TestGetServiceInfoById 测试通过 ID 查找服务信息的函数
func TestGetServiceInfoById(t *testing.T) {
	tests := []struct {
		id       string
		expected Info
		hasError bool
	}{
		{
			id: "zzz1",
			expected: Info{
				Name: "IPSB",
				Host: "ip.sb",
				Port: 443,
			},
			hasError: false,
		},
		{
			id: "zzz2",
			expected: Info{
				Name: "YouTube",
				Host: "www.youtube.com",
				Port: 443,
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
