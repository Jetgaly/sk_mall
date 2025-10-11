package main

import (
    "fmt"
)

func main() {
    fmt.Println("调用 RPC 前")
    
    callWithInterceptor(realRPC)

    fmt.Println("这行不会执行") // ❌ 永远不会到这里
}

func callWithInterceptor(next func() error) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("拦截器捕获 panic，直接返回响应")
        }
    }()

    err := next()
    if err != nil {
        panic("RPC 内部错误") // ❗模拟 go-zero 拦截器
    }
}

func realRPC() error {
    fmt.Println("RPC 内部执行")
    return fmt.Errorf("下游错误")
}
