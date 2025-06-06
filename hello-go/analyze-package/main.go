package main

import (
	"context"
	"fmt"
	"go/types"
	"log"
	"plugin"

	"golang.org/x/tools/go/packages"
)

// 加载并分析给定的Go包
func loadAndAnalyzePackage(pkgPath string) (*types.Package, error) {
	// 设置构建环境
	//buildContext := build.Default
	//buildContext.GOPATH = buildContext.GOPATH // 可以根据需要自定义GOPATH

	// 使用go/packages加载包
	cfg := &packages.Config{
		Context:    context.Background(),
		Mode:       packages.LoadAllSyntax,
		Tests:      true,
		BuildFlags: []string{"-tags", "plugin"}, // 根据需要添加构建标签
	}
	//initialPackages := []packages.Package{
	//	{PkgPath: pkgPath},
	//}
	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil {
		return nil, err
	}

	// 分析包
	pkg := pkgs[0]
	if len(pkg.Errors) > 0 {
		return nil, pkg.Errors[0]
	}
	return pkg.Types, nil
}

// 加载并初始化插件
func loadAndInitPlugin(pluginPath string) error {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return err
	}

	initFunc, err := p.Lookup("Init")
	if err != nil {
		return err
	}

	// 调用Init函数
	initFunc.(func())()
	return nil
}

func main() {
	pluginPath := "./analyze-package/analyze-package.exe"
	// 加载并分析引用的包
	pkg, err := loadAndAnalyzePackage(pluginPath)
	if err != nil {
		log.Fatalf("加载包失败: %v", err)
	}
	log.Println(pkg)
	// 初始化插件，传入分析得到的包
	err = loadAndInitPlugin(pluginPath)
	if err != nil {
		log.Fatalf("初始化插件失败: %v", err)
	}

	fmt.Println("插件初始化成功")
}
