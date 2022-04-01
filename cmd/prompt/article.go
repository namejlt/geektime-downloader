package prompt

import (
	"errors"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/namejlt/geektime-downloader/internal/geektime"
)

// SelectArticles show select promt to choose an article
func SelectArticles(articles []geektime.ArticleSummary) int {
	items := []geektime.ArticleSummary{
		{
			AID:   -1,
			Title: "返回上一级",
		},
	}
	items = append(items, articles...)

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U00002714 {{ .Title | red }}",
		Inactive: "{{if eq .AID -1}} {{ .Title | green }} {{else}} {{ .Title | cyan }} {{end}}",
	}

	prompt := promptui.Select{
		Label:        "请选择文章: ",
		Items:        items,
		Templates:    templates,
		Size:         len(items),
		HideSelected: true,
		CursorPos:    0,
	}

	index, _, err := prompt.Run()

	if err != nil {
		if !errors.Is(err, promptui.ErrInterrupt) {
			fmt.Printf("Prompt failed %v\n", err)
		}
		os.Exit(1)
	}

	return index
}
