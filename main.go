package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
)

var (
	cotx   = context.Background()
	db     *redis.Client
	dblock = &sync.RWMutex{}
)

func generateToken(data string) string {
	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%s=%d", data, time.Now().Unix())))
	key := hex.EncodeToString(hasher.Sum(nil))
	go savetodb(key, data)
	return key
}

func decidebot(data, month string, getdes bool) string {
	var res strings.Builder
	if getdes {
		res.WriteString("Subscription to ")
	}
	if data == "1" {
		res.WriteString("Premium channel")
	}
	if getdes {
		res.WriteString(" for ")
		res.WriteString(month)
		res.WriteString(" month(s)")
	}
	return res.String()
}

func culatemoney(month string) int {
	mont, _ := strconv.Atoi(month)
	return 349 * mont
}

func savetodb(key, data string) {
	dblock.Lock()
	defer dblock.Unlock()
	if data != "" {
		db.Set(cotx, key, data, time.Hour*2)
	} else {
		db.Set(cotx, key, 1, time.Hour*2)
	}
}

func delfromdb(key string) {
	dblock.Lock()
	defer dblock.Unlock()
	db.Del(cotx, key)
}

func getdfromdb(key string) string {
	dblock.RLock()
	defer dblock.RUnlock()
	return db.Get(cotx, key).Val()
}

func buysub(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery
	msg := query.Message
	data := strings.Split(query.Data, "_")
	month, whichbot := data[1], data[2]
	if data[3] == "0" {
		_, _, _ = msg.EditText(b, fmt.Sprintf("You need to pay ₹%d on the below mentioned payment method and send the screenshot with this token in the caption: <code>%s</code>\nNote: Payment token will expire in next 2 hours!\n> UPI: <code>hlinfo@axl</code>", culatemoney(month), generateToken(fmt.Sprintf("<code>%d</code>_-100%s_%s", query.From.Id, whichbot, month))),
			&gotgbot.EditMessageTextOpts{
				ParseMode: "html",
			})
	} else {
		_, _, _ = msg.EditText(b, "Something went wrong!\nContact @Annihilatorrrr!", nil)
	}
	return ext.EndGroups
}

func selmt(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery
	msg := query.Message
	data := strings.Split(query.Data, "_")
	whichbot := data[1]
	_, _, _ = msg.EditText(b, "Choose the number of month(s) you want the subscription for: "+decidebot(whichbot, "", false),
		&gotgbot.EditMessageTextOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{
					Text:         "1 Month",
					CallbackData: fmt.Sprintf("b_1_%s_%s", whichbot, data[2]),
				}},
				{{
					Text:         "2 Months",
					CallbackData: fmt.Sprintf("b_2_%s_%s", whichbot, data[2]),
				}},
				{{
					Text:         "3 Months",
					CallbackData: fmt.Sprintf("b_3_%s_%s", whichbot, data[2]),
				}},
				{{
					Text:         "4 Months",
					CallbackData: fmt.Sprintf("b_4_%s_%s", whichbot, data[2]),
				}},
				{{
					Text:         "5 Months",
					CallbackData: fmt.Sprintf("b_5_%s_%s", whichbot, data[2]),
				}},
				{{
					Text:         "6 Months",
					CallbackData: fmt.Sprintf("b_6_%s_%s", whichbot, data[2]),
				}},
				{{
					Text:         "7 Months",
					CallbackData: fmt.Sprintf("b_7_%s_%s", whichbot, data[2]),
				}},
				{{
					Text:         "8 Months",
					CallbackData: fmt.Sprintf("b_8_%s_%s", whichbot, data[2]),
				}},
				{{
					Text:         "9 Months",
					CallbackData: fmt.Sprintf("b_9_%s_%s", whichbot, data[2]),
				}},
				{{
					Text:         "10 Months",
					CallbackData: fmt.Sprintf("b_10_%s_%s", whichbot, data[2]),
				}},
				{{
					Text:         "11 Months",
					CallbackData: fmt.Sprintf("b_11_%s_%s", whichbot, data[2]),
				}},
				{{
					Text:         "12 Months",
					CallbackData: fmt.Sprintf("b_12_%s_%s", whichbot, data[2]),
				}},
			}},
		})
	return ext.EndGroups
}

func callmanualorauto(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery
	msg := query.Message
	mode := strings.Split(query.Data, "_")[1]
	_, _, _ = msg.EditText(b, "Select the service to which you would like to buy a subscription:",
		&gotgbot.EditMessageTextOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{
					Text:         "Premium channel",
					CallbackData: fmt.Sprintf("m_2079919643_%s", mode),
				}},
			}},
		})
	return ext.EndGroups
}

func getmanualorauto(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	_, _ = msg.Reply(b, "Choose the mode of approval for your payment you would like to have:",
		&gotgbot.SendMessageOpts{
			ReplyMarkup: &gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{
					Text:         "Manual Approval",
					CallbackData: "callm_0",
				}},
			}},
		})
	return ext.EndGroups
}

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	_, _ = msg.Reply(b, "<b>Hello! I'm alive!\n\nI'm here to take subscriptions for services.\n\nCommands:\n> /buy: Buy a subscription for a service.\n> /status chat_id: Check the status of your subscription of a chat.\nSupport: @wowxyz\n\nMade By: @Annihilatorrrr using GotgBot!</b>",
		&gotgbot.SendMessageOpts{
			DisableWebPagePreview: true,
			ParseMode:             "html",
			ReplyMarkup: &gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{
					Text: "Source Code!",
					Url:  "https://publicearn.in/GetSourceCode",
				}},
			}},
		})
	return ext.EndGroups
}

func callmediachk(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery
	msg := query.Message
	mode := strings.Split(query.Data, "_")
	isok := mode[2]
	ruid, _ := strconv.Atoi(mode[1])
	uid := int64(ruid)
	if isok == "0" {
		parsem := strings.Split(msg.Caption, "_")
		key := parsem[0] + "_" + parsem[1]
		db.SAdd(cotx, "paid", key)
		dd, _ := strconv.Atoi(parsem[2])
		db.Set(cotx, key, 1, time.Duration(dd*2419200))
		chatid, _ := strconv.ParseInt(parsem[2], 10, 64)
		if iv, err := b.CreateChatInviteLink(chatid, &gotgbot.CreateChatInviteLinkOpts{
			MemberLimit: 1,
		}); err != nil {
			_, _ = msg.Reply(b, "Something went wrong while creating invite link for the user to join!", nil)
		} else {
			_, _ = b.SendMessage(uid, "You payment is approved and subscription has been added!\nThanks for the payment!\nChat ID: "+parsem[2]+"\nInvite Link: "+iv.InviteLink, nil)
			_, _, _ = msg.EditCaption(b, &gotgbot.EditMessageCaptionOpts{
				Caption: msg.Caption + "\nApproved!",
			})
		}
	} else {
		_, _, _ = msg.EditCaption(b, &gotgbot.EditMessageCaptionOpts{
			Caption: msg.Caption + "\nDeclined!",
		})
		_, _ = b.SendMessage(uid, "You payment is rejected!\nMake sure you send a valid proof with exact amount.", nil)
	}
	return ext.EndGroups
}

func checkmedia(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	if msg.Photo == nil && msg.Document == nil {
		return ext.EndGroups
	}
	if msg.From.Id == 748511091 {
		return ext.EndGroups
	}
	if data := getdfromdb(msg.Caption); data != "" {
		go delfromdb(msg.Caption)
		_, _ = msg.Copy(b, 748511091, &gotgbot.CopyMessageOpts{
			Caption:   &data,
			ParseMode: "html",
			ReplyMarkup: &gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{
					Text:         "Approve!",
					CallbackData: fmt.Sprintf("chkm_%d_0", msg.From.Id),
				}},
				{{
					Text:         "Decline!",
					CallbackData: fmt.Sprintf("chkm_%d_1", msg.From.Id),
				}},
			}},
		})
		_, _ = msg.Reply(b, "Posted for verification!\nYou will be notified about it once approved or rejected!", nil)
	}
	return ext.EndGroups
}

func chkandremover(b *gotgbot.Bot) {
	for {
		time.Sleep(1 * time.Hour)
		for _, key := range db.SMembers(cotx, "paid").Val() {
			if db.Get(cotx, key).Val() == "" {
				parsem := strings.Split(key, "_")
				dd, _ := strconv.ParseInt(parsem[0], 10, 64)
				chatid, _ := strconv.ParseInt(parsem[1], 10, 64)
				_, _ = b.UnbanChatMember(chatid, dd, nil)
				db.SRem(cotx, "paid", key)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func checkstatus(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	args := ctx.Args()[1:]
	if len(args) >= 1 {
		if !strings.HasPrefix(args[0], "-100") {
			_, _ = msg.Reply(b, "Chat id always start with -100!", nil)
			return ext.EndGroups
		}
		into, _ := strconv.ParseInt(args[0], 10, 64)
		if into == 0 {
			_, _ = msg.Reply(b, "Invalid Chat Id!", nil)
			return ext.EndGroups
		}
		_, _ = msg.Reply(b, fmt.Sprintf("On %d your subscription is valid upto: %.1f", into, db.TTL(cotx, fmt.Sprintf("%d_%d", msg.From.Id, into)).Val().Hours()), nil)
	} else {
		_, _ = msg.Reply(b, "Provide the chat id to check the status!", nil)
	}
	return ext.EndGroups
}

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		token = "6976118547:AAF_6_dei2UdLJkIToHkh3t_fIAk4OLjbI8"
	}
	opt, err := redis.ParseURL("redis://default:jmvNAeJe2NUJV8OU3BP2EmRRkWmyUJht@redis-12419.c246.us-east-1-4.ec2.cloud.redislabs.com:12419")
	if err != nil {
		log.Fatal(err.Error())
	}
	bot, err := gotgbot.NewBot(token, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	db = redis.NewClient(opt)
	if err = db.Ping(cotx).Err(); err != nil {
		log.Fatal(err.Error())
	}
	disp := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(_ *gotgbot.Bot, _ *ext.Context, err error) ext.DispatcherAction {
			_, _ = bot.SendMessage(1594433798, err.Error(), nil)
			return ext.DispatcherActionNoop
		},
		MaxRoutines: -1,
	})
	updater := ext.NewUpdater(disp, nil)
	disp.AddHandler(handlers.NewCommand("start", start))
	disp.AddHandler(handlers.NewCommand("buy", getmanualorauto))
	disp.AddHandler(handlers.NewCallback(callbackquery.Prefix("chkm_"), callmediachk))
	disp.AddHandler(handlers.NewCallback(callbackquery.Prefix("callm_"), callmanualorauto))
	disp.AddHandler(handlers.NewCallback(callbackquery.Prefix("m_"), selmt))
	disp.AddHandler(handlers.NewCallback(callbackquery.Prefix("b_"), buysub))
	disp.AddHandler(handlers.NewCommand("status", checkstatus))
	disp.AddHandler(handlers.NewMessage(message.Caption, checkmedia))
	if updater.StartPolling(bot, &ext.PollingOpts{
		EnableWebhookDeletion: true,
		DropPendingUpdates:    false,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			AllowedUpdates: []string{"message", "callback_query"},
		},
	}) != nil {
		log.Fatal("unable to start polling ...")
	}
	go chkandremover(bot)
	log.Println(bot.User.FirstName, "has been started!")
	updater.Idle()
	_ = updater.Stop()
	log.Println(db.Close())
	log.Println("Bye!")
}
