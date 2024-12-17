package filter

import "testing"

func TestFilter_ShouldProcess(t *testing.T) {
    tests := []struct {
        name    string
        opts    Options
        path    string
        want    bool
    }{
        {
            name: "no filters",
            opts: Options{},
            path: "any/path/file.txt",
            want: true,
        },
        {
            name: "extension include",
            opts: Options{
                Includes: []string{".go"},
            },
            path: "src/main.go",
            want: true,
        },
        {
            name: "extension exclude",
            opts: Options{
                Excludes: []string{".tmp"},
            },
            path: "src/file.tmp",
            want: false,
        },
        {
            name: "directory exclude",
            opts: Options{
                Excludes: []string{"vendor/"},
            },
            path: "vendor/module/file.go",
            want: false,
        },
        {
            name: "default ignores",
            opts: Options{
                IgnoreDefault: true,
            },
            path: "node_modules/package.json",
            want: false,
        },
        {
            name: "mixed patterns",
            opts: Options{
                Includes: []string{".go", ".md"},
                Excludes: []string{"vendor/", "test/"},
                IgnoreDefault: true,
            },
            path: "src/main.go",
            want: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            f := New(tt.opts)
            if got := f.ShouldProcess(tt.path); got != tt.want {
                t.Errorf("ShouldProcess(%q) = %v, want %v", tt.path, got, tt.want)
            }
        })
    }
}
