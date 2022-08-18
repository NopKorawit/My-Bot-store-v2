package sheet

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"store/pkg/config"
	"store/pkg/line/keyword"
	"store/pkg/product"
	"strconv"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Service interface {
	GetProducts() ([][]product.Product, error)
	GetProductsByType(productType string) ([]product.Product, error)
	AddBack(add []product.ProductUpdate) error
	Sell(sell []product.ProductUpdate) error
}

type service struct {
	sheet *sheets.Service
	cfg   *config.AppConfig
}

func NewService(cfg *config.AppConfig) (Service, error) {
	b, err := ioutil.ReadFile(cfg.Sheet.GoogleCredentialsPath)
	if err != nil {
		return nil, err
	}

	srv, err := sheets.NewService(context.Background(), option.WithCredentialsJSON(b))
	if err != nil {
		return nil, err
	}

	return &service{sheet: srv, cfg: cfg}, nil
}

// func (s *service) GetProductslist() (*sheets.ValueRange, error) {
// 	resp, err := s.sheet.Spreadsheets.Values.Get(s.cfg.Sheet.SpreadSheetId, "Summary!B8:G").Do()
// 	if err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}
// 	return resp,nil
// }

func (s *service) GetAllProducts() ([]product.Product, error) {
	products := []product.Product{}
	resp, err := s.sheet.Spreadsheets.Values.Get(s.cfg.Sheet.SpreadSheetId, "Summary!B8:G").Do()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// fmt.Println(resp.Values)
	for i, p := range resp.Values {
		buy_int, _ := strconv.Atoi(p[2].(string))
		sell_int, _ := strconv.Atoi(p[3].(string))
		qty_int, _ := strconv.Atoi(p[4].(string))

		products = append(products, product.Product{
			Code:     p[0].(string),
			Name:     p[1].(string),
			Buy:      buy_int,
			Sell:     sell_int,
			Quantity: qty_int,
			Row:      i + 8,
		})
	}
	// fmt.Println(products)
	return products, nil
}

func (s *service) GetProducts() ([][]product.Product, error) {
	var productslist [][]product.Product
	resp, err := s.sheet.Spreadsheets.Values.Get(s.cfg.Sheet.SpreadSheetId, "Summary!B8:G").Do()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	typeGroup := []string{"A", "B", "C", "D", "E"}
	for _, typecode := range typeGroup {
		products := []product.Product{}
		for i, p := range resp.Values {
			if p[0].(string)[0:1] == typecode {
				qty, _ := strconv.Atoi(p[4].(string))
				products = append(products, product.Product{
					Code:     p[0].(string),
					Name:     p[1].(string),
					Quantity: qty,
					Row:      i + 8,
				})
			}
		}
		productslist = append(productslist, products)
	}

	return productslist, nil
}

func (s *service) GetProductsByType(productType string) ([]product.Product, error) {
	products := []product.Product{}
	resp, err := s.sheet.Spreadsheets.Values.Get(s.cfg.Sheet.SpreadSheetId, "Summary!B8:G").Do()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// fmt.Println(resp.Values)
	typecode := keyword.ConvertType(productType)
	for i, p := range resp.Values {
		if p[0].(string)[0:1] == typecode {
			qty, _ := strconv.Atoi(p[4].(string))
			products = append(products, product.Product{
				Code:     p[0].(string),
				Name:     p[1].(string),
				Quantity: qty,
				Row:      i + 8,
			})
		}
	}
	// fmt.Println(products)
	return products, nil
}

func (s *service) AddBack(add []product.ProductUpdate) error {
	list, err := s.GetAllProducts()
	if err != nil {
		log.Println(err)
		return err
	}
	
	var changes = make(map[string]int)
	// changes["A1"] = 2
	for i, input := range add {
		fmt.Println(input)
		fmt.Println("go in list")
		for _, product := range list {
			fmt.Println(product)
			if product.Code == input.Code {
				cell := fmt.Sprintf("E%v", product.Row)
				changes[cell] = product.Sell - input.Quantity
			}

		}
		fmt.Printf("change = %v, i = %v\n", len(changes), i+1)
		if len(changes) != i+1 {
			return fmt.Errorf("not Found Code %v ", input.Code)
		}

	}
	return s.updates(changes)
}

func (s *service) Sell(sell []product.ProductUpdate) error {
	fmt.Println("IN SELL")
	list, err := s.GetAllProducts()
	if err != nil {
		log.Println(err)
		return err
	}
	var changes = make(map[string]int)
	// changes["A1"] = 2
	for i, input := range sell {
		for _, product := range list {
			if product.Code == input.Code {
				if product.Quantity-input.Quantity < 0 {
					return keyword.ErrProductNotEnough
				}

				cell := fmt.Sprintf("E%v", product.Row)
				changes[cell] = product.Sell + input.Quantity

			}
		}
		fmt.Printf("change = %v, i = %v\n", len(changes), i+1)
		if len(changes) != i+1 {
			return keyword.ErrCodenotFound
		}
	}
	return s.updates(changes)
}

// func (s *service) updates(change map[string]int) error {
// 	products := []product.Product{}
// 	resp, err := s.sheet.Spreadsheets.Values.Get(s.cfg.Sheet.SpreadSheetId, "Summary!B8:G").Do()
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	// fmt.Println(resp.Values)
// 	last := (len(resp.Values) + 7)
// 	typecode := keyword.ConvertType(productType)
// 	for i, p := range resp.Values {
// 		if p[0].(string)[0:1] == typecode {
// 			qty, _ := strconv.Atoi(p[4].(string))
// 			products = append(products, product.Product{
// 				Code:     p[0].(string),
// 				Name:     p[1].(string),
// 				Quantity: qty,
// 				Row:      i + 8,
// 			})
// 		}
// 	}
// 	fmt.Println(products)
// 	return nil
// }

func (s *service) updates(change map[string]int) error {
	var vr sheets.ValueRange
	for cell, amount := range change {
		updateVal := []interface{}{amount}
		vr.Values = [][]interface{}{updateVal}
		_, err := s.sheet.Spreadsheets.Values.Update(s.cfg.Sheet.SpreadSheetId, cell, &vr).ValueInputOption("RAW").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve data from sheet. %v", err)
		}
		log.Printf("%v updated with value = %v", cell, vr.Values)

	}
	return nil
}
