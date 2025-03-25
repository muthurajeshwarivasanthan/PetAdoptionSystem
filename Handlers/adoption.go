package Handlers

import (
	"context"
	"errors"
	"log"
	"pet/Config"
	"pet/Models"
	pb "pet/pet/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Adoption struct {
	pb.UnimplementedAdoptionserviceServer
}

func (a *Adoption) CreateAdoption(ctx context.Context, req *pb.CreateAdoptionRequest) (*pb.CreateAdoptionResponse, error) {
	db := Config.GetDB()

	var pet Models.Pet
	if err := db.First(&pet, req.PetId).Error; err != nil {
		log.Printf("Pet ID %d not found", req.PetId)
		return nil, errors.New("pet ID not found")
	}
	if pet.Status == "Adopted" {
		log.Printf("Pet ID %d is already adopted", req.PetId)
		return nil, errors.New("pet is already adopted, cannot create adoption")
	}

	var buyer Models.Buyer
	if err := db.First(&buyer, req.BuyerId).Error; err != nil {
		log.Printf("Buyer ID %d not found", req.BuyerId)
		return nil, errors.New("buyer ID not found")
	}

	adoption := Models.Adoption{
		PetID:        uint(req.PetId),
		BuyerID:      uint(req.BuyerId),
		AdoptionDate: req.AdoptionDate.AsTime(),
		Status:       req.Status,
	}

	data := db.Create(&adoption)
	if data.Error != nil {
		log.Printf("Error in Adoption Creation: %v", data.Error)
		return nil, data.Error
	}
	pet.Status = "Adopted"
	if err := db.Save(&pet).Error; err != nil {
		log.Printf("Failed to update pet status: %v", err)
		return nil, err
	}
	return &pb.CreateAdoptionResponse{
		Message: "Adoption Created Successfully!",
	}, nil
}

func (a *Adoption) GetAdoption(ctx context.Context, req *pb.GetAdoptionRequest) (*pb.GetAdoptionResponse, error) {
	db := Config.GetDB()

	var adoption Models.Adoption
	data := db.First(&adoption, req.AdoptionId)

	if data.Error != nil {
		if errors.Is(data.Error, gorm.ErrRecordNotFound) {
			log.Printf("Adoption ID %d not found", req.AdoptionId)
			return nil, errors.New("adoption ID not found")
		}
		return nil, data.Error
	}

	adoptionDate := timestamppb.New(adoption.AdoptionDate)

	return &pb.GetAdoptionResponse{
		AdoptionId:   int32(adoption.AdoptionID),
		PetId:        int32(adoption.PetID),
		BuyerId:      int32(adoption.BuyerID),
		AdoptionDate: adoptionDate,
		Status:       adoption.Status,
	}, nil
}

func (a *Adoption) ListAdoptions(ctx context.Context, req *pb.Empty) (*pb.ListAdoptionResponse, error) {
	db := Config.GetDB()

	var adoptions []Models.Adoption
	data := db.Find(&adoptions)
	if data.Error != nil {
		log.Printf("Error retrieving Adoptions: %v", data.Error)
		return nil, data.Error
	}

	var adoptionResponses []*pb.GetAdoptionResponse
	for _, adoption := range adoptions {
		adoptionResponses = append(adoptionResponses, &pb.GetAdoptionResponse{
			AdoptionId:   int32(adoption.AdoptionID),
			PetId:        int32(adoption.PetID),
			BuyerId:      int32(adoption.BuyerID),
			AdoptionDate: timestamppb.New(adoption.AdoptionDate),
			Status:       adoption.Status,
		})
	}

	return &pb.ListAdoptionResponse{Adoptions: adoptionResponses}, nil
}

func (a *Adoption) UpdateAdoption(ctx context.Context, req *pb.UpdateAdoptionRequest) (*pb.UpdateAdoptionResponse, error) {
	db := Config.GetDB()

	var adoption Models.Adoption
	data := db.First(&adoption, req.AdoptionId)
	if data.Error != nil {
		if errors.Is(data.Error, gorm.ErrRecordNotFound) {
			log.Printf("Adoption ID %d not found", req.AdoptionId)
			return &pb.UpdateAdoptionResponse{
				Message: "Adoption ID not found, cannot update!",
			}, errors.New("adoption ID not found")
		}
		return nil, data.Error
	}

	adoption.PetID = uint(req.PetId)
	adoption.BuyerID = uint(req.BuyerId)
	adoption.AdoptionDate = req.AdoptionDate.AsTime()
	adoption.Status = req.Status

	data1 := db.Save(&adoption)
	if data1.Error != nil {
		log.Printf("Error in Adoption Update: %v", data1.Error)
		return nil, data1.Error
	}

	return &pb.UpdateAdoptionResponse{
		Message: "Adoption Updated Successfully!",
	}, nil
}

func (a *Adoption) DeleteAdoption(ctx context.Context, req *pb.DeleteAdoptionRequest) (*pb.DeleteAdoptionResponse, error) {
	db := Config.GetDB()

	var adoption Models.Adoption
	data := db.First(&adoption, req.AdoptionId)

	if data.Error != nil {
		if errors.Is(data.Error, gorm.ErrRecordNotFound) {
			log.Printf("Adoption ID %d not found", req.AdoptionId)
			return &pb.DeleteAdoptionResponse{
				Message: "Adoption ID not found, cannot delete!",
			}, errors.New("adoption ID not found")
		}
		return nil, data.Error
	}

	data1 := db.Delete(&adoption)
	if data1.Error != nil {
		log.Printf("Error in Adoption Deletion: %v", data1.Error)
		return nil, data1.Error
	}

	return &pb.DeleteAdoptionResponse{
		Message: "Adoption Deleted Successfully!",
	}, nil
}
