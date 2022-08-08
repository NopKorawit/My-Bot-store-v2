package sheet

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"store/pkg/config"
	"store/pkg/product"
	"strconv"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Service interface {
	GetProducts() ([]product.Product, error)
	GetProductsByType(productType string) ([]product.Product, error)
	Add()
	Sell()
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

func (s *service) GetProducts() ([]product.Product, error) {
	products := []product.Product{}
	resp, err := s.sheet.Spreadsheets.Values.Get(s.cfg.Sheet.SpreadSheetId, "Summary!B8:G").Do()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Println(resp.Values)
	for _, p := range resp.Values {
		qty, _ := strconv.Atoi(p[4].(string))
		products = append(products, product.Product{
			Code:     p[0].(string),
			Name:     p[1].(string),
			Quantity: qty,
			// Cell:     "",
		})
	}
	fmt.Println(products)
	return products, nil
}

func (s *service) GetProductsByType(productType string) ([]product.Product, error) {
	if productType == "All Flavor" {
		return s.GetProducts()
	}
	products := []product.Product{}
	resp, err := s.sheet.Spreadsheets.Values.Get(s.cfg.Sheet.SpreadSheetId, "Summary!B8:G").Do()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Println(resp.Values)
	for _, p := range resp.Values {
		if p[0].(string)[0:1] == productType {
			qty, _ := strconv.Atoi(p[4].(string))
			products = append(products, product.Product{
				Code:     p[0].(string),
				Name:     p[1].(string),
				Quantity: qty,
			})
		}
	}
	fmt.Println(products)
	return products, nil
}

func (s *service) Add() {
	s.updates(map[string]int{})
}

func (s *service) Sell() {
	s.updates(map[string]int{})
}

func (s *service) updates(map[string]int) error {
	return nil
}
