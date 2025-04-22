package main

import (
	"context"
	"fmt"
	"time"

	einomcp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func GenerateTools() []tool.BaseTool {
	startTime := time.Now()
	ctx := context.Background()
	var res []tool.BaseTool
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "mind-maker-client",
		Version: "1.0.0",
	}
	for _, item := range ActivateConfig.Mcps {
		cli, e := client.NewStdioMCPClient(item.Command, item.Env, item.Args...)
		if e != nil {
			panic(e)
		}
		_, e = cli.Initialize(ctx, initRequest)
		if e != nil {
			panic(e)
		}
		tools, e := einomcp.GetTools(ctx, &einomcp.Config{
			Cli: cli,
		})
		if e != nil {
			panic(e)
		}
		res = append(res, tools...)
	}
	timeLast := time.Since(startTime)
	fmt.Println("生成mcp工具列表用时:", timeLast)
	return res
}
