package app

import (
	cliflag "shop-v2/pkg/common/cli/flag"
)

// 定义从命令行读取参数的配置选项
type CliOptions interface {
	// 向指定的 FlagSet 对象添加标志（flags)
	// AddFlags(fs *pflag.FlagSet)
	// 返回 FlagSet 集合
	Flags() (fss cliflag.NamedFlagSets)
	// 验证配置选项，返回一个错误列表，若无错误则为空列表
	Validate() []error
}

// 定义了从配置文件读取参数的选项
type ConfigurableOptions interface {
	// 用于将命令行或配置文件中的参数应用到选项实例，并返回可能的错误列表
	ApplyFlags() []error
}

// 它用于表示可完成设置的选项，并能反馈是否成功完成
type CompletableOptions interface {
	Complete() error
}

// 可打印的配置选项
type PrintableOptions interface {
	String() string
}
