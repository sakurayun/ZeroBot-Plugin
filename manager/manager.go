package manager

import (
	"bot/manager/utils"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	zero.RegisterPlugin(manager{}) // 注册插件
}

type manager struct{} // pixivSearch 搜索P站插图

func (_ manager) GetPluginInfo() zero.PluginInfo { // 返回插件信息
	return zero.PluginInfo{
		Author:     "kanri",
		PluginName: "Manager",
		Version:    "0.0.1",
		Details:    "群管",
	}
}

func (_ manager) Start() { // 插件主体
	// TODO 菜单
	zero.OnFullMatch("群管系统").SetBlock(true).SetPriority(40).
		Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			zero.Send(event, `====群管====
- 禁言@QQ 1
- 解除禁言 @QQ
- 我要自闭 1
- 开启全员禁言
- 解除全员禁言
- 升为管理@QQ
- 取消管理@QQ
- 修改名片@QQ XXX
- 修改头衔@QQ XXX
- 申请头衔 XXX
- 群聊转发 1234 XXX
- 私聊转发 0000 XXX`)
			return zero.SuccessResponse
		})
	// TODO 升为管理
	zero.OnRegex(`^升为管理.*?(\d+)`, zero.OnlyGroup, zero.SuperUserPermission).SetBlock(true).SetPriority(40).
		Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			zero.SetGroupAdmin(
				event.GroupID,
				utils.Str2Int(state["regex_matched"].([]string)[1]), // 被升为管理的人的qq
				true,
			)
			nickname := zero.GetGroupMemberInfo( // 被升为管理的人的昵称
				event.GroupID,
				utils.Str2Int(state["regex_matched"].([]string)[1]), // 被升为管理的人的qq
				false,
			).Get("nickname").Str
			zero.Send(
				event,
				nickname+" 升为了管理~",
			)
			return zero.SuccessResponse
		})
	// TODO 取消管理
	zero.OnRegex(`^取消管理.*?(\d+)`, zero.OnlyGroup, zero.SuperUserPermission).SetBlock(true).SetPriority(40).
		Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			zero.SetGroupAdmin(
				event.GroupID,
				utils.Str2Int(state["regex_matched"].([]string)[1]), // 被取消管理的人的qq
				false,
			)
			nickname := zero.GetGroupMemberInfo( // 被取消管理的人的昵称
				event.GroupID,
				utils.Str2Int(state["regex_matched"].([]string)[1]), // 被取消管理的人的qq
				false,
			).Get("nickname").Str
			zero.Send(
				event,
				"残念~ "+nickname+" 暂时失去了管理员的资格",
			)
			return zero.SuccessResponse
		})
	// TODO 踢出群聊
	zero.OnRegex(`^踢出群聊.*?(\d+)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).SetPriority(40).
		Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			zero.SetGroupKick(
				event.GroupID,
				utils.Str2Int(state["regex_matched"].([]string)[1]), // 被踢出群聊的人的qq
				false,
			)
			nickname := zero.GetGroupMemberInfo( // 被踢出群聊的人的昵称
				event.GroupID,
				utils.Str2Int(state["regex_matched"].([]string)[1]), // 被踢出群聊的人的qq
				false,
			).Get("nickname").Str
			zero.Send(
				event,
				"残念~ "+nickname+" 被放逐",
			)
			return zero.SuccessResponse
		})
	// TODO 退出群聊
	zero.OnRegex(`^退出群聊.*?(\d+)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).SetPriority(40).
		Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			zero.SetGroupLeave(
				utils.Str2Int(state["regex_matched"].([]string)[1]), // 要退出的群的群号
				true,
			)
			return zero.SuccessResponse
		})
	// TODO 开启全体禁言
	zero.OnRegex(`^开启全员禁言$`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).SetPriority(40).
		Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			zero.SetGroupWholeBan(
				event.GroupID,
				true,
			)
			zero.Send(event, "全员自闭开始~")
			return zero.SuccessResponse
		})
	// TODO 解除全体禁言
	zero.OnRegex(`^解除全体禁言$`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).SetPriority(40).
		Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			zero.SetGroupWholeBan(
				event.GroupID,
				false,
			)
			zero.Send(event, "全员自闭结束~")
			return zero.SuccessResponse
		})
	// TODO 禁言
	zero.OnRegex(`^禁言.*?(\d+).*?\s(\d+)(.*)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).SetPriority(40).
		Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			duration := utils.Str2Int(state["regex_matched"].([]string)[2])
			switch state["regex_matched"].([]string)[3] {
			case "分钟":
				//
			case "小时":
				duration = duration * 60
			case "天":
				duration = duration * 60 * 24
			default:
				//
			}
			if duration >= 43200 {
				duration = 43199 // qq禁言最大时长为一个月
			}
			zero.SetGroupBan(
				event.GroupID,
				utils.Str2Int(state["regex_matched"].([]string)[1]), // 要禁言的人的qq
				duration*60, // 要禁言的时间（分钟）
			)
			zero.Send(event, "小黑屋收留成功~")
			return zero.SuccessResponse
		})
	// TODO 解除禁言
	zero.OnRegex(`^解除禁言.*?(\d+)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).SetPriority(40).
		Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			zero.SetGroupBan(
				event.GroupID,
				utils.Str2Int(state["regex_matched"].([]string)[1]), // 要解除禁言的人的qq
				0,
			)
			zero.Send(event, "小黑屋释放成功~")
			return zero.SuccessResponse
		})
	// TODO 自闭禁言
	zero.OnRegex(`^我要自闭.*?(\d+)(.*)`, zero.OnlyGroup).SetBlock(true).SetPriority(40).
		Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			duration := utils.Str2Int(state["regex_matched"].([]string)[2])
			switch state["regex_matched"].([]string)[2] {
			case "分钟":
				//
			case "小时":
				duration = duration * 60
			case "天":
				duration = duration * 60 * 24
			default:
				//
			}
			if duration >= 43200 {
				duration = 43199 // qq禁言最大时长为一个月
			}
			zero.SetGroupBan(
				event.GroupID,
				event.UserID,
				utils.Str2Int(state["regex_matched"].([]string)[1])*60, // 要自闭的时间（分钟）
			)
			zero.Send(event, "那我就不手下留情了~")
			return zero.SuccessResponse
		})
	// TODO 修改名片
	zero.OnRegex(`^修改名片.*?(\d+).*?\s(.*)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).SetPriority(40).
		Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			zero.SetGroupCard(
				event.GroupID,
				utils.Str2Int(state["regex_matched"].([]string)[1]), // 被修改群名片的人
				state["regex_matched"].([]string)[2],                // 修改成的群名片
			)
			zero.Send(
				event,
				"嗯！已经修改了",
			)
			return zero.SuccessResponse
		})
	// TODO 修改头衔
	zero.OnRegex(`^修改头衔.*?(\d+).*?\s(.*)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).SetPriority(40).
		Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			zero.SetGroupSpecialTitle(
				event.GroupID,
				utils.Str2Int(state["regex_matched"].([]string)[1]), // 被修改群头衔的人
				state["regex_matched"].([]string)[2],                // 修改成的群头衔
			)
			zero.Send(
				event,
				"嗯！已经修改了",
			)
			return zero.SuccessResponse
		})
	// TODO 申请头衔
	zero.OnRegex(`^申请头衔(.*)`, zero.OnlyGroup).SetBlock(true).SetPriority(40).
		Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			zero.SetGroupSpecialTitle(
				event.GroupID,
				utils.Str2Int(state["regex_matched"].([]string)[1]), // 被修改群头衔的人
				state["regex_matched"].([]string)[2],                // 修改成的群头衔
			)
			zero.Send(
				event,
				"嗯！不错的头衔呢~",
			)
			return zero.SuccessResponse
		})
	// TODO 群聊转发
	zero.OnRegex(`^群聊转发.*?(\d+)\s(.*)`, zero.SuperUserPermission).SetBlock(true).SetPriority(40).
		Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			zero.SendGroupMessage(
				utils.Str2Int(state["regex_matched"].([]string)[1]), // 需要发送的群
				state["regex_matched"].([]string)[1],                // 需要发送的信息
			)
			zero.Send(
				event,
				"📧 --> "+state["regex_matched"].([]string)[1],
			)
			return zero.SuccessResponse
		})
	// TODO 私聊转发
	zero.OnRegex(`^私聊转发.*?(\d+)\s(.*)`, zero.SuperUserPermission).SetBlock(true).SetPriority(40).
		Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			zero.SendPrivateMessage(
				utils.Str2Int(state["regex_matched"].([]string)[1]), // 需要发送的人的qq
				state["regex_matched"].([]string)[1],                // 需要发送的信息
			)
			zero.Send(
				event,
				"📧 --> "+state["regex_matched"].([]string)[1],
			)
			return zero.SuccessResponse
		})
	// TODO 戳一戳
	zero.OnNotice().SetBlock(false).SetPriority(40).
		Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			if event.NoticeType == "notify" && event.SubType == "poke" && event.RawEvent.Get("target_id").Int() == utils.Str2Int(zero.BotConfig.SelfID) {
				time.Sleep(time.Second * 1)
				zero.Send(event, "请不要戳我 >_<")
			}
			return zero.SuccessResponse
		})
	// TODO 入群欢迎
	zero.OnNotice().SetBlock(false).SetPriority(40).
		Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			if event.NoticeType == "group_increase" {
				zero.Send(event, "欢迎~")
			}
			return zero.SuccessResponse
		})
	// TODO 退群提醒
	zero.OnNotice().SetBlock(false).SetPriority(40).
		Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			if event.NoticeType == "group_decrease" {
				zero.Send(event, "有人跑路了~")
			}
			return zero.SuccessResponse
		})
}
