package line

import (
	"fmt"
	"log"
	"store/pkg/line/keyword"

	"github.com/line/line-bot-sdk-go/linebot"
)

func (b *bot) replyMessage(event *linebot.Event) {
	fmt.Println(event.ReplyToken)
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		b.replyTextMessage(event, message)
	default:
		b.replyUnknownMessage(event)
	}
}

func (b *bot) replyTextMessage(event *linebot.Event, message *linebot.TextMessage) {
	if message.Text == keyword.Flavor {
		b.replyMenu(event, message)
		return
	}

	if keyword.IsMenu(message.Text) {
		b.replyProducts(event, message)
		return
	}

	b.replyOthers(event, message)

}

func (b *bot) replyMenu(event *linebot.Event, message *linebot.TextMessage) {
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(MenuFlex))
	if err != nil {
		log.Println(err)
	}

	flexMessage := linebot.NewFlexMessage(message.Text, flexContainer)
	_, err = b.client.ReplyMessage(event.ReplyToken, flexMessage).Do()
	if err != nil {
		log.Print(err)
	}
}

func (b *bot) replyProducts(event *linebot.Event, message *linebot.TextMessage) {
	if message.Text == keyword.TypeAll {
		productslist, err := b.sheetService.GetProducts()
		if err != nil {
			log.Println(err)
		}
		TypeGroup := []string{keyword.TypeNameA, keyword.TypeNameB, keyword.TypeNameC, keyword.TypeNameD, keyword.TypeNameE}
		text := "รายการทั้งหมดมีดังนี้\n\n"
		for i, products := range productslist {
			head := fmt.Sprintf("%v มีตามนี้ค้าบ\n", TypeGroup[i])
			for _, product := range products {
				prefix := product.GetQtySymbol()
				text := fmt.Sprintf("%v | %v | %v\n", prefix, product.Code, product.Name)
				head = head + text
			}
			text = text + head + "\n"
		}
		text = text + "❌ หมดแล้วค้าบ \n⚠️ เหลือ 1-2 อัน \n✅ มีมากกว่า 2 อัน"
		// fmt.Println(text)
		if _, err := b.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do(); err != nil {
			log.Print(err)
		}
	}

	products, err := b.sheetService.GetProductsByType(message.Text)
	if err != nil {
		log.Println(err)
	}

	head := fmt.Sprintf("รายการ %v \n", message.Text)
	for _, product := range products {
		prefix := product.GetQtySymbol()
		text := fmt.Sprintf("%v | %v | %v\n", prefix, product.Code, product.Name)
		head = head + text
	}

	fmt.Println(head)

	if _, err := b.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(head)).Do(); err != nil {
		log.Print(err)
	}
}

func (b *bot) replyOthers(event *linebot.Event, message *linebot.TextMessage) {
	// rows := strings.Split(message.Text, "\n")
	// //อัพเดตสินค้า
	// if rows[0] == "อัพเดต" || rows[0] == "ตั้งค่า" {
	// 	rows := rows[1:]
	// 	var productsList []model.MultiProduct
	// 	for _, row := range rows {
	// 		split := strings.Split(row, " ")
	// 		fmt.Println(split[1])
	// 		amount, _ := strconv.Atoi(split[1])
	// 		product := model.MultiProduct{
	// 			Code:     split[0],
	// 			Quantity: amount,
	// 		}
	// 		productsList = append(productsList, product)
	// 	}
	// 	sell, err := h.productService.UpdateMultiProducts(productsList)
	// 	fmt.Println(err)
	// 	if err != nil {
	// 		if err == model.ErrProductNotEnough {
	// 			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("สินค้ามีจำนวนไม่เพียงพอ")).Do(); err != nil {
	// 				log.Print(err)
	// 				return
	// 			}
	// 			return
	// 		} else if err == model.ErrCodenotFound {
	// 			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ค้นหาสินค้าไม่เจอ")).Do(); err != nil {
	// 				log.Print(err)
	// 				return
	// 			}
	// 			return
	// 		} else {
	// 			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ระบบผิดพลาด")).Do(); err != nil {
	// 				log.Print(err)
	// 				return
	// 			}
	// 			return
	// 		}
	// 	}
	// 	text := "รายการ\n"
	// 	for _, Product := range sell {
	// 		list := fmt.Sprintf("%v หัว| %v %v\n", Product.Quantity, Product.Type, Product.Name)
	// 		text = text + list
	// 	}
	// 	text = text + "ถูกซื้อเรียบร้อยแล้ว"
	// 	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do(); err != nil {
	// 		log.Print(err)
	// 		return
	// 	}
	// }
	// //ขายสินค้า
	// if rows[0] == "ขาย" || rows[0] == "เอา" || rows[0] == "buy" || rows[0] == "order" {
	// 	rows := rows[1:]
	// 	var productsList []model.MultiProduct
	// 	for _, row := range rows {
	// 		split := strings.Split(row, " ")
	// 		fmt.Println(split[1])
	// 		amount, _ := strconv.Atoi(split[1])
	// 		product := model.MultiProduct{
	// 			Code:     split[0],
	// 			Quantity: amount,
	// 		}
	// 		productsList = append(productsList, product)
	// 	}
	// 	sell, err := h.productService.SellMultiProduct(productsList)
	// 	fmt.Println(err)
	// 	if err != nil {
	// 		if err == model.ErrProductNotEnough {
	// 			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("สินค้ามีจำนวนไม่เพียงพอ")).Do(); err != nil {
	// 				log.Print(err)
	// 				return
	// 			}
	// 			return
	// 		} else if err == model.ErrCodenotFound {
	// 			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ค้นหาสินค้าไม่เจอ")).Do(); err != nil {
	// 				log.Print(err)
	// 				return
	// 			}
	// 			return
	// 		} else {
	// 			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ระบบผิดพลาด")).Do(); err != nil {
	// 				log.Print(err)
	// 				return
	// 			}
	// 			return
	// 		}
	// 	}
	// 	text := "รายการ\n"
	// 	for _, Product := range sell {
	// 		list := fmt.Sprintf("%v หัว| %v %v\n", Product.Quantity, Product.Type, Product.Name)
	// 		text = text + list
	// 	}
	// 	text = text + "ถูกซื้อเรียบร้อยแล้ว"
	// 	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do(); err != nil {
	// 		log.Print(err)
	// 		return
	// 	}
	// } else {
	// 	// Emoji
	// 	sorry := linebot.NewEmoji(0, "5ac1bfd5040ab15980c9b435", "024")
	// 	// have := linebot.NewEmoji(0, "5ac21a18040ab15980c9b43e", "007")
	// 	// out := linebot.NewEmoji(0, "5ac21a18040ab15980c9b43e", "068")
	// 	// few := linebot.NewEmoji(0, "5ac21a18040ab15980c9b43e", "025")
	// 	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("$ ขออภัยครับ แต่ผมยังไม่เข้าใจ ท่านอยากจะทวนอีกรอบหรือรอให้นพมาตอบคำถามดีครับ").AddEmoji(sorry)).Do(); err != nil {
	// 		log.Print(err)
	// 	}
	// }
}

func (b *bot) replyUnknownMessage(event *linebot.Event) {
	if _, err := b.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Unknown")).Do(); err != nil {
		log.Println(err)
	}
}
