package fsutil

import "testing"

func TestParsePathRelation(t *testing.T) {
	type args struct {
		self   string
		target string
	}
	tests := []struct {
		name string
		args args
		want pathRelation
	}{
		{"0", args{"", ""}, PathSelf},
		{"1", args{"a", ""}, PathParrent},
		{"2", args{"", "b"}, PathChild},
		{"3", args{"a/b", ""}, PathSup},
		{"4", args{"", "a/b"}, PathSub},
		{"5", args{"a", "b"}, PathIrrelevant},
		{"6", args{"/a", "b"}, PathIrrelevant},
		{"7", args{"a/", "b"}, PathIrrelevant},
		{"8", args{"/a/", "b"}, PathIrrelevant},
		{"9", args{"a/b/c", "a"}, PathSup},
		{"10", args{"a/b/c", "a/b"}, PathParrent},
		{"11", args{"a/b/c", "a/b/c"}, PathSelf},
		{"12", args{"a/b/c", "/a/b/c"}, PathSelf},
		{"13", args{"a/b/c", "a/b/c/"}, PathSelf},
		{"14", args{"a/b/c", "/a/b/c/"}, PathSelf},
		{"15", args{"/a/b", "a/b/c/"}, PathChild},
		{"16", args{"/a/", "a/b/c"}, PathSub},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SumPathRelation(tt.args.self, tt.args.target); got != tt.want {
				t.Errorf("SumPathRelation() = %v, want %v", got, tt.want)
			}
		})
	}
}
