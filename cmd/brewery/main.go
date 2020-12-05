package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/jbowl/apibrewery"
	"github.com/jbowl/brewery/internal/pkg/obdb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type server struct {
	apibrewery.UnimplementedBreweryServiceServer
	obdb obdb.OBDB
}

func init() {
	// log as JSON not default ASCII
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	// only log warning severity or above
	//	log.SetLevel(log.WarnLevel)

	log.Printf("init()")
}

func run() error {
	log.Println("run() Server running ...")

	// TODO: use env variable for port
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}

	api := server{
		obdb: obdb.OBDB{APIUrl: "https://api.openbrewerydb.org"},
	}

	// TODO: use a selfsigned cert

	srv := grpc.NewServer()
	apibrewery.RegisterBreweryServiceServer(srv, &api)

	log.Fatalln(srv.Serve(lis))

	return nil
}

func main() {
	// not of value as a docker container
	pid := os.Getpid()
	fmt.Printf("pid for %s = %d\n", os.Args[0], pid)

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// SearchBreweries -
func (s *server) SearchBreweries(filter *apibrewery.Filter, stream apibrewery.BreweryService_SearchBreweriesServer) error {

	APIUrl := s.obdb.APIUrl + "/breweries/search?" + filter.GetBy()

	resp, err := s.obdb.RESTReq("GET", APIUrl)

	if err != nil {
		return err
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

	for _, brewery := range b {
		apib := apibrewery.Brewery{Id: brewery.ID, Name: brewery.Name, WebsiteUrl: brewery.Website}

		if err := stream.Send(&apib); err != nil {
			return err
		}
	}

	return nil
}

// ListBreweries -
func (s *server) ListBreweries(filter *apibrewery.Filter, stream apibrewery.BreweryService_ListBreweriesServer) error {

	APIUrl := s.obdb.APIUrl + "/breweries?" + filter.GetBy()

	resp, err := s.obdb.RESTReq("GET", APIUrl)

	if err != nil {
		return err
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
			Id:          brewery.ID,
			Name:        brewery.Name,
			Street:      brewery.Street,
			City:        brewery.City,
			State:       brewery.State,
			Countryprov: brewery.CountryProvince,
			Postalcode:  brewery.PostalCode,
			Country:     brewery.Country,
			Longitude:   brewery.Longitude,
			Latitude:    brewery.Latitude,
			Phone:       brewery.Phone,
			WebsiteUrl:  brewery.Website,
		}

		if err := stream.Send(&apib); err != nil {
			return err
		}
	}

	return nil
}
