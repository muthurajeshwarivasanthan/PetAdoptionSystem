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

type Buyer struct {
	pb.UnimplementedBuyerserviceServer
}

func (b *Buyer) CreateBuyer(ctx context.Context, req *pb.CreateBuyerRequest) (*pb.CreateBuyerResponse, error) {
	db := Config.GetDB()

	var existingBuyer Models.Buyer

	data := db.First(&existingBuyer, req.Id)
	if data.RowsAffected > 0 {
		log.Printf("Buyer ID %d already exists", req.Id)
		return &pb.CreateBuyerResponse{
			Message: "Buyer ID already present, cannot create duplicate!",
		}, errors.New("buyer ID already exists")
	}

	buyer := Models.Buyer{
		ID:          uint(req.Id),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Age:         int(req.Age),
	}

	data = db.Create(&buyer)
	if data.Error != nil {
		log.Printf("Error in Buyer Creation: %v", data.Error)
		return nil, data.Error
	}

	return &pb.CreateBuyerResponse{
		Message: "Buyer Created Successfully!",
	}, nil
}

func (b *Buyer) GetBuyer(ctx context.Context, req *pb.GetBuyerIDRequest) (*pb.GetBuyerIDResponse, error) {
	db := Config.GetDB()

	var buyer Models.Buyer
	data := db.First(&buyer, req.Id)

	if data.Error != nil {
		if errors.Is(data.Error, gorm.ErrRecordNotFound) {
			log.Printf("Buyer ID %d not found", req.Id)
			return nil, errors.New("buyer ID not found")
		}
		return nil, data.Error
	}

	return &pb.GetBuyerIDResponse{
		Id:          int32(buyer.ID),
		FirstName:   buyer.FirstName,
		LastName:    buyer.LastName,
		PhoneNumber: buyer.PhoneNumber,
		Address:     buyer.Address,
		Age:         int32(buyer.Age),
	}, nil
}

func (b *Buyer) ListBuyer(ctx context.Context, req *pb.UserEmpty) (*pb.ListBuyerResponse, error) {
	db := Config.GetDB()

	var buyers []Models.Buyer
	data := db.Find(&buyers)
	if data.Error != nil {
		log.Printf("Error retrieving Buyers: %v", data.Error)
		return nil, data.Error
	}

	var buyerResponses []*pb.GetBuyerIDResponse
	for _, buyer := range buyers {
		buyerResponses = append(buyerResponses, &pb.GetBuyerIDResponse{
			Id:          int32(buyer.ID),
			FirstName:   buyer.FirstName,
			LastName:    buyer.LastName,
			PhoneNumber: buyer.PhoneNumber,
			Address:     buyer.Address,
			Age:         int32(buyer.Age),
		})
	}

	return &pb.ListBuyerResponse{Buyers: buyerResponses}, nil
}

func (b *Buyer) UpdateBuyer(ctx context.Context, req *pb.UpdateBuyerRequest) (*pb.UpdateBuyerResponse, error) {
	db := Config.GetDB()

	var buyer Models.Buyer
	data := db.First(&buyer, req.Id)

	if data.Error != nil {
		if errors.Is(data.Error, gorm.ErrRecordNotFound) {
			log.Printf("Buyer ID %d not found", req.Id)
			return &pb.UpdateBuyerResponse{
				Message: "Buyer ID not found, cannot update!",
			}, errors.New("buyer ID not found")
		}
		return nil, data.Error
	}

	buyer.FirstName = req.FirstName
	buyer.LastName = req.LastName
	buyer.PhoneNumber = req.PhoneNumber
	buyer.Address = req.Address
	buyer.Age = int(req.Age)

	data1 := db.Save(&buyer)
	if data1.Error != nil {
		log.Printf("Error in Buyer Update: %v", data1.Error)
		return nil, data1.Error
	}

	return &pb.UpdateBuyerResponse{
		Message: "Buyer Updated Successfully!",
	}, nil
}

func (b *Buyer) DeleteBuyer(ctx context.Context, req *pb.DeleteBuyerRequest) (*pb.DeleteBuyerResponse, error) {
	db := Config.GetDB()

	var buyer Models.Buyer
	data := db.First(&buyer, req.Id)

	if data.Error != nil {
		if errors.Is(data.Error, gorm.ErrRecordNotFound) {
			log.Printf("Buyer ID %d not found", req.Id)
			return &pb.DeleteBuyerResponse{
				Message: "Buyer ID not found, cannot delete!",
			}, errors.New("buyer ID not found")
		}
		return nil, data.Error
	}

	data1 := db.Delete(&buyer)
	if data1.Error != nil {
		log.Printf("Error in Buyer Deletion: %v", data1.Error)
		return nil, data1.Error
	}

	return &pb.DeleteBuyerResponse{
		Message: "Buyer Deleted Successfully!",
	}, nil
}
