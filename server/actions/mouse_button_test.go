package actions

import "testing"

func TestMouseButtonString(t *testing.T) {
	tests := []struct {
		button MouseButton
		want   string
	}{
		{LeftButton, "left"},
		{MiddleButton, "middle"},
		{RightButton, "right"},
		{MouseButton(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.button.String()
			if got != tt.want {
				t.Errorf("MouseButton(%d).String() = %q, want %q", tt.button, got, tt.want)
			}
		})
	}
}
