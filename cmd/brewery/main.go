//https://gist.github.com/fntlnz/cf14feb5a46b2eda428e000157447309  create CA cert

//export GO_PATH=~/go
//export PATH=$PATH:/$GO_PATH/bin
//protoc --go_out=. --go-grpc_out=. service.proto
// protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/brewery/brewery.proto
package main

import (
	"encoding/json"
	"fmt"
	"github.com/jbowl/brewery/internal/pkg/obdb"
	"io/ioutil"

	"github.com/jbowl/apibrewery"

	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	apibrewery.UnimplementedBreweryServiceServer
	obdb obdb.OBDB


}

func main() {
	log.Println("Server running ...")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}

	api := server {
		obdb: obdb.OBDB{ APIUrl:  "https://api.openbrewerydb.org"},
	}

	srv := grpc.NewServer()
	apibrewery.RegisterBreweryServiceServer(srv, &api)

	log.Fatalln(srv.Serve(lis))
}

type BreweryResults struct {

	BreweryResult [] apibrewery.BreweryResult;
}


// ListFeatures lists all features contained within the given bounding Rectangle.
func (s *server) SearchBreweries(filter *apibrewery.Filter, stream apibrewery.BreweryService_SearchBreweriesServer) error {

	APIUrl := s.obdb.APIUrl + "/breweries/search?" + filter.GetBy()

	resp, err := s.obdb.RequestImplResponse("GET", APIUrl)

	if err != nil {
		return  err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var b []apibrewery.BreweryResult


	if err := json.Unmarshal(contents, &b); err != nil {
		return err
	}

	fmt.Printf(string(contents))

	for _, brewery := range b {

		apib := apibrewery.Brewery{Id: brewery.ID, Name: brewery.Name, WebsiteUrl: brewery.Website}

		if err := stream.Send(&apib); err != nil {
			return err
		}
	}

	return nil
}



// ListFeatures lists all features contained within the given bounding Rectangle.
func (s *server) ListBreweries(filter *apibrewery.Filter, stream apibrewery.BreweryService_ListBreweriesServer) error {


	APIUrl := s.obdb.APIUrl + "/breweries?" + filter.GetBy()

	resp, err := s.obdb.RequestImplResponse("GET", APIUrl)

	if err != nil {
		return  err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

    var b []apibrewery.BreweryResult


	if err := json.Unmarshal(contents, &b); err != nil {
		return err
	}

	fmt.Printf(string(contents))

     for _, brewery := range b {

		 apib := apibrewery.Brewery{
		 	Id: brewery.ID,
		 	Name: brewery.Name,
		 	Street: brewery.Street,
		 	City: brewery.City,
		 	State: brewery.State,
		 	Countryprov: brewery.CountryProvince,
		 	Postalcode: brewery.PostalCode,
		 	Country: brewery.Country,
		 	Longitude: brewery.Longitude,
		 	Latitude: brewery.Latitude,
		 	Phone: brewery.Phone,
		 	WebsiteUrl: brewery.Website,
		 }

		 			if err := stream.Send(&apib); err != nil {
		 				return err
		 			}
	 }

	return nil
}

