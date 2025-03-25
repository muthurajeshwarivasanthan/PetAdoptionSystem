package Handlers

import (
	"context"
	"errors"
	"log"
	"pet/Config"
	Models "pet/Models"
	pb "pet/pet/pb"

	"gorm.io/gorm"
)

type Seller struct {
	pb.UnimplementedSellerserviceServer
}

func (s *Seller) CreateSeller(ctx context.Context, req *pb.CreateSellerRequest) (*pb.CreateSellerResponse, error) {
	db := Config.GetDB()

	var existingSeller Models.Seller

	data := db.First(&existingSeller, req.Id)
	if data.RowsAffected > 0 {
		log.Printf("Seller ID %d already exists", req.Id)
		return &pb.CreateSellerResponse{
			Message: "Seller ID already present, cannot create duplicate!",
		}, errors.New("seller ID already exists")
	}

	seller := Models.Seller{
		ID:          uint(req.Id),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Age:         int(req.Age),
	}

	data = db.Create(&seller)
	if data.Error != nil {
		log.Printf("Error in Seller Creation: %v", data.Error)
		return nil, data.Error
	}

	return &pb.CreateSellerResponse{
		Message: "Seller Created Successfully!",
	}, nil
}

func (s *Seller) GetSeller(ctx context.Context, req *pb.GetSellerIDRequest) (*pb.GetSellerIDResponse, error) {
	db := Config.GetDB()

	var seller Models.Seller
	data := db.First(&seller, req.Id)

	if data.Error != nil {
		if errors.Is(data.Error, gorm.ErrRecordNotFound) {
			log.Printf("Seller ID %d not found", req.Id)
			return nil, errors.New("seller ID not found")
		}
		return nil, data.Error
	}

	return &pb.GetSellerIDResponse{
		Id:          int32(seller.ID),
		FirstName:   seller.FirstName,
		LastName:    seller.LastName,
		PhoneNumber: seller.PhoneNumber,
		Address:     seller.Address,
		Age:         int32(seller.Age),
	}, nil
}

func (s *Seller) ListSeller(ctx context.Context, req *pb.UserEmpty) (*pb.ListSellerResponse, error) {
	db := Config.GetDB()

	var sellers []Models.Seller
	data := db.Find(&sellers)
	if data.Error != nil {
		log.Printf("Error retrieving Sellers: %v", data.Error)
		return nil, data.Error
	}

	var sellerResponses []*pb.GetSellerIDResponse
	for _, seller := range sellers {
		sellerResponses = append(sellerResponses, &pb.GetSellerIDResponse{
			Id:          int32(seller.ID),
			FirstName:   seller.FirstName,
			LastName:    seller.LastName,
			PhoneNumber: seller.PhoneNumber,
			Address:     seller.Address,
			Age:         int32(seller.Age),
		})
	}

	return &pb.ListSellerResponse{Sellers: sellerResponses}, nil
}

func (s *Seller) UpdateUser(ctx context.Context, req *pb.UpdateSellerRequest) (*pb.UpdateSellerResponse, error) {
	db := Config.GetDB()

	var seller Models.Seller
	data := db.First(&seller, req.Id)

	if data.Error != nil {
		if errors.Is(data.Error, gorm.ErrRecordNotFound) {
			log.Printf("Seller ID %d not found", req.Id)
			return &pb.UpdateSellerResponse{
				Message: "Seller ID not found, cannot update!",
			}, errors.New("seller ID not found")
		}
		return nil, data.Error
	}

	seller.FirstName = req.FirstName
	seller.LastName = req.LastName
	seller.PhoneNumber = req.PhoneNumber
	seller.Address = req.Address
	seller.Age = int(req.Age)

	data1 := db.Save(&seller)
	if data1.Error != nil {
		log.Printf("Error in Seller Update: %v", data1.Error)
		return nil, data1.Error
	}

	return &pb.UpdateSellerResponse{
		Message: "Seller Updated Successfully!",
	}, nil
}

func (s *Seller) DeleteSeller(ctx context.Context, req *pb.DeleteSellerRequest) (*pb.DeleteSellerResponse, error) {
	db := Config.GetDB()

	var seller Models.Seller
	data := db.First(&seller, req.Id)

	if data.Error != nil {
		if errors.Is(data.Error, gorm.ErrRecordNotFound) {
			log.Printf("Seller ID %d not found", req.Id)
			return &pb.DeleteSellerResponse{
				Message: "Seller ID not found, cannot delete!",
			}, errors.New("seller ID not found")
		}
		return nil, data.Error
	}

	data1 := db.Delete(&seller)
	if data1.Error != nil {
		log.Printf("Error in Seller Deletion: %v", data1.Error)
		return nil, data1.Error
	}

	return &pb.DeleteSellerResponse{
		Message: "Seller Deleted Successfully!",
	}, nil
}
