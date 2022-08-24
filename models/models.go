package models


type Product struct {
	Id      string  `json:"_id,omitempty" bson:"_id,omitempty"`
	Name    string  `json:"name" bson:"name"`
	Price   float32 `json:"price" bson:"price"`
	Tags    string  `json:"tags" bson:"tags"`
	Barcode string  `json:"barcode" bson:"barcode"`
}
