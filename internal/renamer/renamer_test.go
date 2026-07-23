package renamer

import "testing"

func TestName_PadsToFileCountWidth(t *testing.T) {
	r := New(150)

	tests := []struct {
		index int
		want  string
	}{
		{1, "001.png"},
		{10, "010.png"},
		{150, "150.png"},
	}

	for _, tt := range tests {
		if got := r.Name(tt.index, ".png"); got != tt.want {
			t.Errorf("Name(%d, .png) = %q, want %q", tt.index, got, tt.want)
		}
	}
}

func TestName_ZeroOrNegativeFileCountDefaultsTo999Width(t *testing.T) {
	for _, count := range []int{0, -1} {
		r := New(count)
		if got := r.Name(5, ".jpg"); got != "005.jpg" {
			t.Errorf("New(%d).Name(5, .jpg) = %q, want %q", count, got, "005.jpg")
		}
	}
}
