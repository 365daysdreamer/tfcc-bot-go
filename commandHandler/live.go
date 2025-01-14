package commandHandler

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/Touhou-Freshman-Camp/tfcc-bot-go/bilibili"
	"github.com/Touhou-Freshman-Camp/tfcc-bot-go/perm"
	"github.com/ozgio/strutil"
)

func init() {
	register(&getLiveState{})
	register(&startLive{})
	register(&stopLive{})
	register(&changeLiveTitle{})
}

type getLiveState struct{}

func (g *getLiveState) Name() string {
	return "直播状态"
}

func (g *getLiveState) ShowTips(int64, int64) string {
	return "直播状态"
}

func (g *getLiveState) CheckAuth(int64, int64) bool {
	return true
}

func (g *getLiveState) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) != 0 {
		return
	}
	ret, err := bilibili.GetLiveStatus()
	if err != nil {
		logger.WithError(err).Errorln("获取直播状态失败")
		return
	}
	if ret.Code != 0 {
		logger.Errorf("请求直播间状态失败，错误码：%d，错误信息：%s\n", ret.Code, ret.Message)
		return
	}
	var text string
	if ret.Data.LiveStatus == 0 {
		text = "直播间状态：未开播"
	} else {
		text = fmt.Sprintf("直播间状态：开播\n直播标题：%s\n人气：%d\n直播间地址：%s", ret.Data.Title, ret.Data.Online, bilibili.GetLiveUrl())
	}
	groupMsg = message.NewSendingMessage().Append(message.NewText(text))
	return
}

type startLive struct{}

func (s *startLive) Name() string {
	return "开始直播"
}

func (s *startLive) ShowTips(int64, int64) string {
	return "开始直播"
}

func (s *startLive) CheckAuth(_ int64, senderId int64) bool {
	return perm.IsWhitelist(senderId)
}

func (s *startLive) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) != 0 {
		return
	}
	if len(content) != 0 {
		return
	}
	ret, err := bilibili.StartLive()
	if err != nil {
		logger.WithError(err).Errorln("开启直播间失败")
		return
	}
	if ret.Code != 0 {
		logger.Errorf("开启直播间失败，错误码：%d，错误信息1：%s，错误信息2：%s\n", ret.Code, ret.Message, ret.Msg)
		return
	}
	var publicText string
	if ret.Data.Change == 0 {
		publicText = fmt.Sprintf("直播间本来就是开启的，推流码已私聊\n直播间地址：%s\n快来围观吧！", bilibili.GetLiveUrl())
	} else {
		publicText = fmt.Sprintf("直播间已开启，推流码已私聊，别忘了修改直播间标题哦！\n直播间地址：%s\n快来围观吧！", bilibili.GetLiveUrl())
	}
	rtmpAddr := ret.Data.Rtmp.Addr
	rtmpCode := ret.Data.Rtmp.Code
	privateText := fmt.Sprintf("RTMP推流地址：%s\n密钥：%s", rtmpAddr, rtmpCode)
	groupMsg = message.NewSendingMessage().Append(message.NewText(publicText))
	privateMsg = message.NewSendingMessage().Append(message.NewText(privateText))
	return
}

type stopLive struct{}

func (s *stopLive) Name() string {
	return "关闭直播"
}

func (s *stopLive) ShowTips(int64, int64) string {
	return "关闭直播"
}

func (s *stopLive) CheckAuth(_ int64, senderId int64) bool {
	return perm.IsWhitelist(senderId)
}

func (s *stopLive) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) != 0 {
		return
	}
	ret, err := bilibili.StopLive()
	if err != nil {
		logger.WithError(err).Errorln("关闭直播间失败")
		return
	}
	if ret.Code != 0 {
		logger.Errorf("关闭直播间失败，错误码：%d，错误信息1：%s，错误信息2：%s\n", ret.Code, ret.Message, ret.Msg)
		return
	}
	var text string
	if ret.Data.Change == 0 {
		text = "直播间本来就是关闭的"
	} else {
		text = "直播间已关闭"
	}
	groupMsg = message.NewSendingMessage().Append(message.NewText(text))
	return
}

type changeLiveTitle struct{}

func (c *changeLiveTitle) Name() string {
	return "修改直播标题"
}

func (c *changeLiveTitle) ShowTips(int64, int64) string {
	return "修改直播标题 新标题"
}

func (c *changeLiveTitle) CheckAuth(_ int64, senderId int64) bool {
	return perm.IsWhitelist(senderId)
}

func (c *changeLiveTitle) Execute(_ *message.GroupMessage, content string) (groupMsg *message.SendingMessage, privateMsg *message.SendingMessage) {
	if len(content) == 0 {
		groupMsg = message.NewSendingMessage().Append(message.NewText("指令格式如下：\n修改直播标题 新标题"))
		return
	}
	if strutil.Len(content) > 20 {
		return
	}
	ret, err := bilibili.ChangeLiveTitle(content)
	if err != nil {
		logger.WithError(err).Errorln("修改直播间标题失败")
		return
	}
	var text string
	if ret.Code != 0 {
		logger.Errorf("修改直播间标题失败，错误码：%d，错误信息1：%s，错误信息2：%s\n", ret.Code, ret.Message, ret.Msg)
		text = "修改直播间标题失败，请联系管理员"
	} else {
		text = "直播间标题已修改为：" + content
	}
	groupMsg = message.NewSendingMessage().Append(message.NewText(text))
	return
}
