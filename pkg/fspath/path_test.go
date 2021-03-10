package fspath

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
		{"Self_0", args{"", ""}, PathSelf},
		{"Self_1", args{".", ""}, PathSelf},
		{"Self_2", args{"./", ""}, PathSelf},
		{"Self_3", args{"./", "."}, PathSelf},
		{"Parrent_0", args{"a", ""}, PathParrent},
		{"Parrent_1", args{"./a", ""}, PathParrent},
		{"Parrent_2", args{"./a", "."}, PathParrent},
		{"Parrent_3", args{"./a", "./"}, PathParrent},
		{"Child_0", args{"", "b"}, PathChild},
		{"Sup_0", args{"a/b", ""}, PathSup},
		{"Sub_0", args{"", "a/b"}, PathSub},
		{"Irrelevant_0", args{"a", "b"}, PathIrrelevant},
		{"Irrelevant_1", args{"/a", "b"}, PathIrrelevant},
		{"Irrelevant_2", args{"a/", "b"}, PathIrrelevant},
		{"Irrelevant_3", args{"/a/", "b"}, PathIrrelevant},
		{"Sup_1", args{"a/b/c", "a"}, PathSup},
		{"Parrent_4", args{"a/b/c", "a/b"}, PathParrent},
		{"Self_4", args{"a/b/c", "a/b/c"}, PathSelf},
		{"Self_5", args{"a/b/c", "/a/b/c"}, PathSelf},
		{"Self_6", args{"a/b/c", "a/b/c/"}, PathSelf},
		{"Self_7", args{"a/b/c", "/a/b/c/"}, PathSelf},
		{"Child_1", args{"/a/b", "a/b/c/"}, PathChild},
		{"Sub_1", args{"/a/", "a/b/c"}, PathSub},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SumPathRelation(tt.args.self, tt.args.target); got != tt.want {
				t.Errorf("SumPathRelation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_SplitPrefixDir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name       string
		args       args
		wantBucket string
		wantObject string
	}{
		{"", args{"a/b"}, "a", "b"},
		{"", args{"/a/b"}, "a", "b"},
		{"", args{"a/b/"}, "a", "b"},
		{"", args{"/a/b/"}, "a", "b"},
		{"", args{"a"}, "a", ""},
		{"", args{"a/b/c"}, "a", "b/c"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBucket, gotObject := SplitPrefixDir(tt.args.path)
			if gotBucket != tt.wantBucket {
				t.Errorf("SplitPrefixDir() gotBucket = %v, want %v", gotBucket, tt.wantBucket)
			}
			if gotObject != tt.wantObject {
				t.Errorf("SplitPrefixDir() gotObject = %v, want %v", gotObject, tt.wantObject)
			}
		})
	}
}
