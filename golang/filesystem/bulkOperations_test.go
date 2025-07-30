package filesystem

import (
	"reflect"
	"testing"
)

// mockFolder records calls to AddTagToFile
type mockFolder struct {
	calls []tagCall
}
type tagCall struct {
	path string
	tag  string
}

func (m *mockFolder) AddTagToFile(path, tag string) bool {
	m.calls = append(m.calls, tagCall{path: path, tag: tag})
	return true // simulate success
}

func TestBulkAddTags(t *testing.T) {
	tests := []struct {
		name      string
		bulkList  []TagsStruct
		wantCalls []tagCall
		wantPanic bool
	}{
		{
			name:      "empty bulkList",
			bulkList:  nil,
			wantCalls: nil,
		},
		{
			name: "single file, single tag",
			bulkList: []TagsStruct{
				{FilePath: "/foo.txt", Tags: []string{"a"}},
			},
			wantCalls: []tagCall{
				{path: "/foo.txt", tag: "a"},
			},
		},
		{
			name: "single file, multiple tags",
			bulkList: []TagsStruct{
				{FilePath: "/foo.txt", Tags: []string{"a", "b", "c"}},
			},
			wantCalls: []tagCall{
				{path: "/foo.txt", tag: "a"},
				{path: "/foo.txt", tag: "b"},
				{path: "/foo.txt", tag: "c"},
			},
		},
		{
			name: "multiple files, multiple tags",
			bulkList: []TagsStruct{
				{FilePath: "/a", Tags: []string{"x", "y"}},
				{FilePath: "/b", Tags: []string{"1", "2", "3"}},
			},
			wantCalls: []tagCall{
				{path: "/a", tag: "x"},
				{path: "/a", tag: "y"},
				{path: "/b", tag: "1"},
				{path: "/b", tag: "2"},
				{path: "/b", tag: "3"},
			},
		},
		{
			name:      "nil item panics",
			bulkList:  []TagsStruct{{FilePath: "p", Tags: []string{"t"}}},
			wantPanic: true,
		},
	}

	for _, tc := range tests {
		tc := tc // capture
		t.Run(tc.name, func(t *testing.T) {
			var mf *mockFolder
			if !tc.wantPanic {
				mf = &mockFolder{}
			}
			// wrap for panic detection
			didPanic := false
			func() {
				defer func() {
					if r := recover(); r != nil {
						didPanic = true
					}
				}()
				// call under test
				_ = BulkAddTags(mf, tc.bulkList)
				if mf == nil {
					t.Fatal("expected non-nil item for non-panic case")
				}
			}()
			if tc.wantPanic {
				if !didPanic {
					t.Errorf("expected panic but none occurred")
				}
				return
			}
			if didPanic {
				t.Errorf("unexpected panic")
				return
			}
			// verify calls
			if !reflect.DeepEqual(mf.calls, tc.wantCalls) {
				t.Errorf("calls = %#v; want %#v", mf.calls, tc.wantCalls)
			}
		})
	}
}
