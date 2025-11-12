package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/eino-examples/adk/common/model"
	"github.com/cloudwego/eino-examples/adk/common/prints"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/compose"

	"week05-homework-go/agents"
	"week05-homework-go/tools"
)

func main() {
	ctx := context.Background()
	if os.Getenv("OPENAI_API_KEY") == "" && os.Getenv("ARK_API_KEY") == "" {
		fmt.Println("错误：请在环境变量中设置 OPENAI_API_KEY 或 ARK_API_KEY。")
		return
	}
	cm := model.NewChatModel()
	rd := bufio.NewReader(os.Stdin)
	fmt.Print("请输入您想写的文章主题 (或按回车使用默认主题): ")
	t, _ := rd.ReadString('\n')
	topic := strings.TrimRight(t, "\r\n")
	if strings.TrimSpace(topic) == "" {
		topic = "帮我写一篇关于AI Agent的文章"
	}

	style := "通俗易懂"
	length := 1000

	cmd, args := tools.DefaultMCPCommand()
	researchInstr, err := tools.FetchPromptUsingCommand(ctx, cmd, args, "research", "", 0)
	must(err)
	searchTool, err := tools.NewMCPSearchTool(ctx, cmd, args)
	must(err)
	researchAgent, err := agents.NewResearchAgent(ctx, cm, researchInstr, searchTool)
	must(err)
	writingInstr, err := tools.FetchPromptUsingCommand(ctx, cmd, args, "write", style, length)
	must(err)
	writingAgent, err := agents.NewWritingAgent(ctx, cm, writingInstr)
	must(err)
	reviewInstr, err := tools.FetchPromptUsingCommand(ctx, cmd, args, "review", "", 0)
	must(err)
	reviewer, err := agents.NewReviewerAgent(ctx, cm, reviewInstr)
	must(err)
	polishInstr, err := tools.FetchPromptUsingCommand(ctx, cmd, args, "polish", "", 0)
	must(err)
	polisher, err := agents.NewPolisherAgent(ctx, cm, polishInstr)
	must(err)

	seq, err := adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
		Name:        "WriterWorkflow",
		Description: "研究→撰写→审核→润色的顺序编排",
		SubAgents:   []adk.Agent{researchAgent, writingAgent, reviewer, polisher},
	})
	must(err)

	runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: seq, EnableStreaming: false, CheckPointStore: newInMemoryStore()})
	query := fmt.Sprintf("主题：%s\n风格：%s\n长度：%d", topic, style, length)
	iter := runner.Query(ctx, query, adk.WithCheckPointID("pipeline"))
	finalText, researchText, draft, suggestions := collectOutputs(iter)

	ts := time.Now().Format("20060102_150405")
	fn := fmt.Sprintf("article_output_%s.md", ts)
	content := fmt.Sprintf("# 最终文章：%s\n\n%s\n\n---\n\n## 研究报告\n\n%s\n\n## 文章初稿\n\n%s\n\n## 审核建议\n\n%s\n", topic, finalText, researchText, draft, suggestions)
	os.WriteFile(fn, []byte(content), 0644)
	fmt.Printf("🎉 成功生成输出文件：%s\n", fn)
}

func collectOutputs(iter *adk.AsyncIterator[*adk.AgentEvent]) (string, string, string, string) {
	ft := ""
	rr := ""
	dr := ""
	sg := ""
	for {
		e, ok := iter.Next()
		if !ok {
			break
		}
		prints.Event(e)
		if e.Output != nil && e.Output.MessageOutput != nil && e.Output.MessageOutput.Message != nil {
			c := e.Output.MessageOutput.Message.Content
			n := e.AgentName
			if n == "ResearchAgent" {
				rr = c
			}
			if n == "WritingAgent" {
				dr = c
			}
			if n == "ReviewerAgent" {
				sg = c
			}
			if n == "PolisherAgent" {
				ft = c
			}
		}
	}
	return ft, rr, dr, sg
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

type inMemoryStore struct{ mem map[string][]byte }

func newInMemoryStore() compose.CheckPointStore { return &inMemoryStore{mem: map[string][]byte{}} }
func (s *inMemoryStore) Set(ctx context.Context, key string, value []byte) error {
	s.mem[key] = value
	return nil
}
func (s *inMemoryStore) Get(ctx context.Context, key string) ([]byte, bool, error) {
	v, ok := s.mem[key]
	return v, ok, nil
}
