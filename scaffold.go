package main

import (
	"bytes"
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var tmplFS embed.FS

type Opts struct {
	Project string
	Module  string
	Port    int
	With    map[string]bool // api/admin
}

type TplData struct {
	Project string
	Module  string
	Port    int
}

func main() {
	if len(os.Args) < 2 || os.Args[1] != "create" {
		fmt.Println("Usage: go run scaffold.go create <projectName> --module=<module> --port=8080 --with=api,admin")
		os.Exit(1)
	}

	var (
		project string
		module  string
		port    int
		with    string
	)
	if len(os.Args) >= 3 {
		project = os.Args[2]
	}
	flagset := flag.NewFlagSet("create", flag.ExitOnError)
	flagset.StringVar(&module, "module", "", "go module path (e.g. github.com/acme/backend-example)")
	flagset.IntVar(&port, "port", 8080, "http port")
	flagset.StringVar(&with, "with", "api,admin,task", "components to include, csv: api,admin,task")
	_ = flagset.Parse(os.Args[3:])

	if project == "" || module == "" {
		fmt.Println("project and --module are required")
		os.Exit(1)
	}

	opt := Opts{
		Project: project,
		Module:  module,
		Port:    port,
		With:    parseWith(with),
	}
	if err := run(opt); err != nil {
		fmt.Println("scaffold error:", err)
		os.Exit(2)
	}
	fmt.Println("✅ Done. Next:")
	fmt.Printf("cd %s && go mod tidy && go run ./cmd/api\n", project)
}

func parseWith(s string) map[string]bool {
	m := map[string]bool{}
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			m[p] = true
		}
	}
	return m
}

func run(o Opts) error {
	dstRoot := filepath.Join(".", o.Project)
	if err := os.MkdirAll(dstRoot, 0o755); err != nil {
		return err
	}

	data := TplData{Project: o.Project, Module: o.Module, Port: o.Port}

	return fs.WalkDir(tmplFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "templates" {
			return nil
		}

		rel := strings.TrimPrefix(path, "templates/")
		dst := filepath.Join(dstRoot, rel)
		if d.IsDir() {
			// 组件过滤：若选择只要 api，不生成 admin 的目录
			if strings.HasPrefix(rel, "cmd/api") && !o.With["api"] {
				return fs.SkipDir
			}
			if strings.HasPrefix(rel, "cmd/admin") && !o.With["admin"] {
				return fs.SkipDir
			}

			if strings.HasPrefix(rel, "cmd/task") && !o.With["task"] {
				return fs.SkipDir
			}
			return os.MkdirAll(dst, 0o755)
		}

		// 组件过滤（文件级）
		if strings.HasPrefix(rel, "cmd/api") && !o.With["api"] {
			return nil
		}
		if strings.HasPrefix(rel, "cmd/admin") && !o.With["admin"] {
			return nil
		}
		if strings.HasPrefix(rel, "cmd/task") && !o.With["task"] {
			return nil
		}

		b, readErr := fs.ReadFile(tmplFS, path)
		if readErr != nil {
			return readErr
		}

		// 渲染 .tmpl；否则原样复制
		if strings.HasSuffix(rel, ".tmpl") {
			dst = strings.TrimSuffix(dst, ".tmpl")
			return renderToFile(b, dst, data)
		}
		return writeFile(dst, b)
	})
}

func renderToFile(tpl []byte, dst string, data TplData) error {
	t, err := template.New(filepath.Base(dst)).Parse(string(tpl))
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return err
	}
	return writeFile(dst, buf.Bytes())
}

func writeFile(dst string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	// 如果文件已存在，避免覆盖（可按需改成覆盖）
	if _, err := os.Stat(dst); err == nil {
		return errors.New("file exists: " + dst)
	}
	return os.WriteFile(dst, b, 0o644)
}
