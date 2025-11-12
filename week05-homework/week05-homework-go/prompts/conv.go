package prompts

import "fmt"

func BuildWritingUserInput(researchReport string, style string, length int) string {
	return fmt.Sprintf("根据以下研究报告撰写文章初稿。风格：%s；长度：%d字左右。\n\n[研究报告]\n\n%s", style, length, researchReport)
}

func BuildReviewUserInput(draft string) string {
	return fmt.Sprintf("请基于以下文章初稿给出修改建议（列表形式）。\n\n[文章初稿]\n\n%s", draft)
}

func BuildPolishUserInput(draft string, suggestions string) string {
	return fmt.Sprintf("请根据初稿与审核建议进行最终润色，输出终稿。\n\n[文章初稿]\n\n%s\n\n[审核建议]\n\n%s", draft, suggestions)
}
