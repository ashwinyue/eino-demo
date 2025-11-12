package common

import (
    "context"
    "fmt"

    "github.com/cloudwego/eino/callbacks"
    "github.com/cloudwego/eino/schema"
)

type LoggerCallbacks struct{}

func (l *LoggerCallbacks) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
    fmt.Println("start", info.Name, info.Type, info.Component)
    return ctx
}

func (l *LoggerCallbacks) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
    fmt.Println("end", info.Name, info.Type, info.Component)
    return ctx
}

func (l *LoggerCallbacks) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
    fmt.Println("error", info.Name, info.Type, info.Component, err)
    return ctx
}

func (l *LoggerCallbacks) OnStartWithStreamInput(ctx context.Context, info *callbacks.RunInfo, input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
    return ctx
}

func (l *LoggerCallbacks) OnEndWithStreamOutput(ctx context.Context, info *callbacks.RunInfo, output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {
    return ctx
}

