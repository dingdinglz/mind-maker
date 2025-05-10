package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/tool/duckduckgo"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()

	fmt.Println("mind-maker by dinglz")
	fmt.Println("https://github.com/dingdinglz/mind-maker")

	LoadConfig()

	if FileExist("mind.html") {
		os.Remove("mind.html")
	}

	fmt.Print("请输入要生成思维导图的知识点:")
	reader := bufio.NewReader(os.Stdin)
	question, e := reader.ReadString('\n')
	if e != nil {
		panic(e)
	}
	question = strings.ReplaceAll(question, " ", "")
	question = strings.ReplaceAll(question, "\n", "")
	fmt.Println("开始对知识点", question, "绘制思维导图")

	chatModel, e := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: ActivateConfig.Model.BaseURL,
		APIKey:  ActivateConfig.Model.ApiKey,
		Model:   ActivateConfig.Model.Model,
	})
	if e != nil {
		panic(e)
	}

	tools := GenerateTools()
	if ActivateConfig.Search {
		fmt.Println("search mode enabled")
		ducktool, e := duckduckgo.NewTool(ctx, &duckduckgo.Config{})
		if e != nil {
			panic(e)
		}
		tools = append(tools, ducktool)
	}
	ragent, e := react.NewAgent(ctx, &react.AgentConfig{
		Model: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: tools,
		},
	})
	if e != nil {
		panic(e)
	}

	startTime := time.Now()
	sr, e := ragent.Stream(ctx, []*schema.Message{
		{
			Role:    schema.System,
			Content: "下面，你需要对用户给出的知识点进行思维导图的绘制，然后将结果保存在mind.html，要求思维导图结构清晰，逻辑正确，覆盖面广，生成的内容的语言应当使用中文\n\n制作前你需要进行搜索收集相关信息，如果收集到的信息不是中文，请翻译成中文，结果一定要是中文，记住，一定要将结果保存到mind.html，生成的内容的语言应当使用中文",
		},
		{
			Role:    schema.User,
			Content: question,
		},
	}, agent.WithComposeOptions(compose.WithCallbacks(&LoggerCallback{})))
	if e != nil {
		panic(e)
	}

	defer sr.Close()
	for {
		msg, e := sr.Recv()
		if e != nil {
			if errors.Is(e, io.EOF) {
				break
			}
			panic(e)
		}
		fmt.Print(msg.Content)
	}
	fmt.Println()
	fmt.Println("生成思维导图用时：", time.Since(startTime))
	if FileExist("mind.html") {
		fmt.Println("生成思维导图已保存到：mind.html")
	} else {
		fmt.Println("生成失败！请重试")
	}
}
