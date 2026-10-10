package xiaohongshu

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xpzouying/xiaohongshu-mcp/browser"
)

func TestMatchBoardByName(t *testing.T) {
	boards := []Board{
		{ID: "board_1", Name: "美食探店", NoteCount: 15},
		{ID: "board_2", Name: "AI工具与开发", NoteCount: 42},
		{ID: "board_3", Name: "Reading List", NoteCount: 8},
	}

	t.Run("精确匹配", func(t *testing.T) {
		b, err := matchBoardByName(boards, "美食探店")
		require.NoError(t, err)
		require.Equal(t, "board_1", b.ID)
		require.Equal(t, "美食探店", b.Name)
	})

	t.Run("首尾带空格匹配", func(t *testing.T) {
		b, err := matchBoardByName(boards, "  AI工具与开发  ")
		require.NoError(t, err)
		require.Equal(t, "board_2", b.ID)
	})

	t.Run("大小写不敏感匹配", func(t *testing.T) {
		b, err := matchBoardByName(boards, "reading list")
		require.NoError(t, err)
		require.Equal(t, "board_3", b.ID)
	})

	t.Run("找不到专辑时列出所有候选名称", func(t *testing.T) {
		b, err := matchBoardByName(boards, "不存在的专辑")
		require.Error(t, err)
		require.Nil(t, b)
		require.Contains(t, err.Error(), "未找到名为 \"不存在的专辑\" 的专辑")
		require.Contains(t, err.Error(), "美食探店")
		require.Contains(t, err.Error(), "AI工具与开发")
		require.Contains(t, err.Error(), "Reading List")
	})

	t.Run("空名称报错", func(t *testing.T) {
		b, err := matchBoardByName(boards, "   ")
		require.Error(t, err)
		require.Nil(t, b)
		require.Contains(t, err.Error(), "专辑名称不能为空")
	})
}

// 本地有登录态 cookies 时可打开跳过运行
func TestGetBoardFeeds(t *testing.T) {
	t.Skip("SKIP: 需要登录态")

	b := browser.NewBrowser(true)
	defer b.Close()

	page := b.NewPage()
	defer page.Close()

	action := NewBoardAction(page)
	ctx := context.Background()

	// 1. 列出专辑
	boards, err := action.ListBoards(ctx, "")
	require.NoError(t, err)
	require.NotEmpty(t, boards)

	target := boards[0]
	fmt.Printf("测试专辑: %s (ID: %s, 帖子数: %d)\n", target.Name, target.ID, target.NoteCount)

	// 2. 查前 5 条
	res5, err := action.GetBoardFeeds(ctx, target.ID, 5)
	require.NoError(t, err)
	require.NotNil(t, res5)
	fmt.Printf("获取前 5 条成功，返回 %d 条\n", res5.TotalReturned)
	if target.NoteCount > 0 {
		require.NotEmpty(t, res5.Feeds)
		require.NotEmpty(t, res5.Feeds[0].ID)
	}

	// 3. 按名称解析测试
	resolved, err := action.ResolveBoard(ctx, target.Name, "")
	require.NoError(t, err)
	require.Equal(t, target.ID, resolved.ID)
}
