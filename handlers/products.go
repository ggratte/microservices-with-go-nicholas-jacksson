// Package classification of Product API
//
// Documentation for Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta

package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ggratte/go-microservices-nicjacksson/data"
	"github.com/gorilla/mux"
)

type Products struct{
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products{
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request){
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle POST Products")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"]) 

	p.l.Println("Handle PUT Product", id)
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err := data.UpdateProduct(id, &prod)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	
	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler{
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request){
		prod := data.Product{}

		if err := prod.FromJSON(r.Body); err != nil{
			http.Error(rw, "Unable to unmarshall json", http.StatusBadRequest)
			return
		}

		err := prod.Validate()
		if err != nil {
			p.l.Println("[Error] validating product", err)
			http.Error(rw, fmt.Sprintf("Error validating product: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}