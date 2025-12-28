package errutil

// 此文件为占位文件
// 错误处理需要手动拆分：
// 
// 1. stack.go → errutil/stack.go
//    移动整个文件内容
//
// 2. errbox.go 中的 panicError → errutil/panic.go
//    - 将 panicError 改为 PanicError（导出）
//    - 提供 NewPanicError() 构造函数
//    - 更新 errbox.go 中的导入
//
// 请参考 MIGRATION_GUIDE.md 进行拆分
