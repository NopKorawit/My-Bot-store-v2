package line

import (
	"fmt"
	"log"
	"store/pkg/line/keyword"
	"store/pkg/product"
	"strconv"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

func (b *bot) replyMessage(event *linebot.Event) {
	fmt.Printf("reply roken is: %v", event.ReplyToken)
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		b.replyTextMessage(event, message)
	default:
		b.replyUnknownMessage(event)
	}
}

func (b *bot) replyTextMessage(event *linebot.Event, message *linebot.TextMessage) {
	if message.Text == "Test" {
		b.replytest()
		return
	}

	if message.Text == keyword.Flavor {
		b.replyMenu(event, message)
		return
	}

	if message.Text == keyword.Bank {
		b.replyBankAccount(event, message)
		return
	}

	if message.Text == keyword.Promptpay || message.Text == "พร้อมเพย์" {
		b.replyBankAccount(event, message)
		return
	}
	
	if message.Text == keyword.Story {
		b.replyStory(event, message)
		return
	}

	if keyword.IsMenu(message.Text) {
		b.replyProducts(event, message)
		return
	}
	rows := strings.Split(message.Text, "\n")
	//ขายสินค้า
	if rows[0] == "ขาย" || rows[0] == "sell" || rows[0] == "เอา" {
		b.replySell(event, rows)
	}
	//คืนสินค้า
	if rows[0] == "rollback" || rows[0] == "addback" || rows[0] == "ย้อนกลับ" || rows[0] == "ขายผิด" {
		b.replyAddBack(event, rows)
	}
	// //สร้าง qr
	// if rows[0] == "pp" || rows[0] == "qr" || rows[0] == "พร้อมเพย์" {
	// 	b.replyPromptpay(event, rows)
	// }
	// b.replyOthers(event, message)

}

func (b *bot) replytest() {
	update := []product.Product{
		{Code: "E1", Quantity: 1},
		{Code: "E5", Quantity: 3},
	}
	b.sheetService.Sell(update)
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

	// fmt.Println(head)

	if _, err := b.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(head)).Do(); err != nil {
		log.Print(err)
	}
}

func (b *bot) replySell(event *linebot.Event, rows []string) {
	rows = rows[1:]
	var productsList []product.Product
	for _, row := range rows {
		split := strings.Split(row, " ")
		fmt.Println(split[1])
		amount, _ := strconv.Atoi(split[1])
		product := product.Product{
			Code:     split[0],
			Quantity: amount,
		}
		productsList = append(productsList, product)
	}
	err := b.sheetService.Sell(productsList)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		if err == keyword.ErrProductNotEnough {
			if _, err = b.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("สินค้ามีจำนวนไม่เพียงพอ")).Do(); err != nil {
				log.Print(err)
				return
			}
			return
		} else if err == keyword.ErrCodenotFound {
			if _, err = b.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ค้นหาสินค้าไม่เจอ")).Do(); err != nil {
				log.Print(err)
				return
			}
			return
		} else {
			if _, err = b.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ระบบผิดพลาด")).Do(); err != nil {
				log.Print(err)
				return
			}
			return
		}
	}
	text := "รายการ\n"
	for _, Product := range productsList {
		list := fmt.Sprintf("%v | %v\n", Product.Code, Product.Quantity)
		text = text + list
	}
	text = text + "ถูกซื้อเรียบร้อยแล้ว"
	if _, err = b.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do(); err != nil {
		log.Print(err)
		return
	}
}

func (b *bot) replyAddBack(event *linebot.Event, rows []string) {
	rows = rows[1:]
	var productsList []product.Product
	for _, row := range rows {
		split := strings.Split(row, " ")
		fmt.Println(split[1])
		amount, _ := strconv.Atoi(split[1])
		product := product.Product{
			Code:     split[0],
			Quantity: amount,
		}
		productsList = append(productsList, product)
	}
	err := b.sheetService.Sell(productsList)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		if err == keyword.ErrProductNotEnough {
			if _, err = b.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("สินค้ามีจำนวนไม่เพียงพอ")).Do(); err != nil {
				log.Print(err)
				return
			}
			return
		} else if err == keyword.ErrCodenotFound {
			if _, err = b.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ค้นหาสินค้าไม่เจอ")).Do(); err != nil {
				log.Print(err)
				return
			}
			return
		} else {
			if _, err = b.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ระบบผิดพลาด")).Do(); err != nil {
				log.Print(err)
				return
			}
			return
		}
	}
	text := "รายการ\n"
	for _, Product := range productsList {
		list := fmt.Sprintf("%v | %v\n", Product.Code, Product.Quantity)
		text = text + list
	}
	text = text + "ถูกซื้อเรียบร้อยแล้ว"
	if _, err = b.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do(); err != nil {
		log.Print(err)
		return
	}
}

func (b *bot) replyOthers(event *linebot.Event, message *linebot.TextMessage) {
	// Emoji
	sorry := linebot.NewEmoji(0, "5ac1bfd5040ab15980c9b435", "024")
	// have := linebot.NewEmoji(0, "5ac21a18040ab15980c9b43e", "007")
	// out := linebot.NewEmoji(0, "5ac21a18040ab15980c9b43e", "068")
	// few := linebot.NewEmoji(0, "5ac21a18040ab15980c9b43e", "025")
	if _, err := b.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("$ ขออภัยครับ แต่ผมยังไม่เข้าใจ ท่านอยากจะทวนอีกรอบหรือรอให้นพมาตอบคำถามดีครับ").AddEmoji(sorry)).Do(); err != nil {
		log.Print(err)
	}

}

func (b *bot) replyStory(event *linebot.Event, message *linebot.TextMessage) {
	if _, err := b.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(
	`เปิดร้านค้าบ
	หัวละ 125
	3หัว+ หัวละ 120
	3หัวขึ้นไป บริเวณม.ส่งฟรี ค้าบ
	ส่งแฟลช 40 บาท`)).Do(); err != nil {
		log.Print(err)
	}
}

func (b *bot) replyBankAccount(event *linebot.Event, message *linebot.TextMessage) {
	if _, err := b.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(`
	9862461403
	กรุงไทย
	นายกรวิชญ์ วาสนารุ่งเรืองสุข
	`)).Do(); err != nil {
		log.Print(err)
	}
}

func (b *bot) replyPromptpay(event *linebot.Event, message *linebot.TextMessage) {
	if _, err := b.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(`
	9862461403
	กรุงไทย
	นายกรวิชญ์ วาสนารุ่งเรืองสุข
	`)).Do(); err != nil {
		log.Print(err)
	}
	split := strings.Split(message.Text, " ")
		fmt.Println(split[1])
		amount, _ := strconv.Atoi(split[1])
	QRcode := fmt.Sprintf("https://promptpay.io/1900101293500/%v.png",amount)
	if _, err := b.client.ReplyMessage(event.ReplyToken,linebot.NewImageMessage(QRcode,QRcode)).Do(); err != nil {
		log.Print(err)
	}
}


func (b *bot) replyUnknownMessage(event *linebot.Event) {
	if _, err := b.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Unknown")).Do(); err != nil {
		log.Println(err)
	}
}
