package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	productProto "github.com/leonvanderhaeghen/go-grpc/pkg/product/v1"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var collection *mongo.Collection

type server struct{

	productProto.UnimplementedProductConfigFetchServiceServer
}

type Product struct {
	Id      string 			   `bson:"_id,omitempty"`
	Name 	string             `bson:"name"`
	Price  	float32            `bson:"price"`
	Tags    string             `bson:"tags"`
	Barcode	string             `bson:"barcode"`
}
// DB set up
func setupDB() *sql.DB {

	// connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://leonvanderhaeghen:EnagMsC8j00X1C8f@productscluster0.uw7owvj.mongodb.net/?retryWrites=true&w=majority"))
	fmt.Println(" mongodb connection opend")
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" passed error checksd")

	collection = client.Database("myFirstDatabase").Collection("product")
	fmt.Println(" opend collection")
}
func (pr *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(pr)
}
func (pr *Product) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(pr)
}


func (*server) CreateProduct(product *Product) (product *Product){
	fmt.Println("Create product request")
	value := req.GetValues()

	data := Product{
		Id:	value.GetId(),
		Name: value.GetName(),
		Price: value.GetPrice(),
		Tags: value.GetTags(),
		Barcode: value.GetBarcode(),
	}

	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	return &data
}


func (*server) GetProduct(Id string) (product *Product) {
	fmt.Println("Read product request")
	data := &product{}
	filter := bson.M{"_id": Id}

	res := collection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot Find product item : %v", err),
		)
	}

	return &data
}

func (*server) GetProducts(ctx context.Context, req *productProto.GetProductsRequest) (*productProto.GetProductsResponse, error) {
	fmt.Println("Get all products")
	// create empty struct
	filter := bson.D{{}}

	curs, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Printf("Error while getting records: %v", err)
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("internal error"),
		)
	}
	//Close the cursor once finished
	defer curs.Close(context.Background())
	var results []*productProto.Product
	for curs.Next(context.Background()) {
		//Create a value into which the single document can be decoded
		var elem Product
		err := curs.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, productToProto(&elem))

	}

	if err := curs.Err(); err != nil {
		log.Fatal(err)
	}

	return &productProto.GetProductsResponse{
		Values: results,
	}, nil
}


func (*server) UpdateProduct(ctx context.Context, req *productProto.UpdateProductRequest){
	fmt.Println("Update product request")
	value := req.GetValues()

	data := &Product{}
	filter := bson.M{"_id": value.Id}

	res := collection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, "product %v not found (the id should not change)", err.Error())
	}

	data.Name = value.Name
	data.Price = value.Price
	data.Barcode = value.Barcode
	data.Tags = value.Tags

	_, updateErr := collection.ReplaceOne(ctx, filter, data)
	if updateErr != nil {
		return nil, status.Errorf(codes.Internal, "Cannot update object in MongoDB: %v", updateErr)
	}

	return &productProto.UpdateProductResponse{
		Values: productToProto(data),
	}, nil
}

func (*server) DeleteProduct(ctx context.Context, req *productProto.DeleteProductRequest) (*productProto.DeleteProductResponse, error) {
	fmt.Println("Delete product request")

	filter := bson.M{"_id": req.GetId()}

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("something went wrong with deleting product %v", req.GetId()))
	}

	if res.DeletedCount == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("product %v not found", req.GetId()))
	}

	return &productProto.DeleteProductResponse{
		Id: req.GetId(),
	}, nil

}


func productToProto(data *Product) *productProto.Product {
	v := &productProto.Product{
		Id:      data.Id,
		Name:    data.Name,
		Price:   data.Price,
		Tags:    data.Tags,
		Barcode: data.Barcode,
	}

	return v
}