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

type Pet struct {
	pb.UnimplementedPetserviceServer
}

func (p *Pet) CreatePet(ctx context.Context, req *pb.CreatePetRequest) (*pb.CreatePetResponse, error) {
	db := Config.GetDB()

	var existingPet Models.Pet

	data := db.First(&existingPet, req.Id)
	if data.RowsAffected > 0 {
		log.Printf("Pet ID %d already exists", req.Id)
		return &pb.CreatePetResponse{
			Message: "Pet ID already present, cannot create duplicate!",
		}, errors.New("pet ID already exists")
	}

	var seller Models.Seller
	if err := db.First(&seller, req.SellerId).Error; err != nil {
		log.Printf("Seller ID %d not found", req.SellerId)
		return &pb.CreatePetResponse{
			Message: "Seller ID does not exist, cannot add Pet!",
		}, errors.New("seller ID not found")
	}

	pet := Models.Pet{
		ID:       uint(req.Id),
		SellerID: uint(req.SellerId),
		PetName:  req.PetName,
		PetType:  req.PetType,
		Breed:    req.Breed,
		Age:      int(req.Age),
		Gender:   req.Gender,
		Status:   req.Status,
		PetImage: req.PetImage,
	}

	data = db.Create(&pet)
	if data.Error != nil {
		log.Printf("Error in Pet Creation: %v", data.Error)
		return nil, data.Error
	}

	return &pb.CreatePetResponse{
		Message: "Pet Created Successfully!",
	}, nil
}

func (p *Pet) GetPet(ctx context.Context, req *pb.GetPetIDRequest) (*pb.GetPetIDResponse, error) {
	db := Config.GetDB()

	var pet Models.Pet
	data := db.First(&pet, req.Id)

	if data.Error != nil {
		if errors.Is(data.Error, gorm.ErrRecordNotFound) {
			log.Printf("Pet ID %d not found", req.Id)
			return nil, errors.New("pet ID not found")
		}
		return nil, data.Error
	}

	return &pb.GetPetIDResponse{
		Id:       int32(pet.ID),
		SellerId: int32(pet.SellerID),
		PetName:  pet.PetName,
		PetType:  pet.PetType,
		Breed:    pet.Breed,
		Age:      int32(pet.Age),
		Gender:   pet.Gender,
		Status:   pet.Status,
		PetImage: pet.PetImage,
	}, nil
}

func (p *Pet) ListPet(ctx context.Context, req *pb.UserEmpty) (*pb.ListPetResponse, error) {
	db := Config.GetDB()

	var pets []Models.Pet
	data := db.Find(&pets)
	if data.Error != nil {
		log.Printf("Error retrieving Pets: %v", data.Error)
		return nil, data.Error
	}

	var petResponses []*pb.GetPetIDResponse
	for _, pet := range pets {
		petResponses = append(petResponses, &pb.GetPetIDResponse{
			Id:       int32(pet.ID),
			SellerId: int32(pet.SellerID),
			PetName:  pet.PetName,
			PetType:  pet.PetType,
			Breed:    pet.Breed,
			Age:      int32(pet.Age),
			Gender:   pet.Gender,
			Status:   pet.Status,
			PetImage: pet.PetImage,
		})
	}

	return &pb.ListPetResponse{Pets: petResponses}, nil
}

func (p *Pet) UpdatePet(ctx context.Context, req *pb.UpdatePetRequest) (*pb.UpdatePetResponse, error) {
	db := Config.GetDB()

	var pet Models.Pet
	data := db.First(&pet, req.Id)

	if data.Error != nil {
		if errors.Is(data.Error, gorm.ErrRecordNotFound) {
			log.Printf("Pet ID %d not found", req.Id)
			return &pb.UpdatePetResponse{
				Message: "Pet ID not found, cannot update!",
			}, errors.New("pet ID not found")
		}
		return nil, data.Error
	}

	pet.PetName = req.PetName
	pet.PetType = req.PetType
	pet.Breed = req.Breed
	pet.Age = int(req.Age)
	pet.Gender = req.Gender
	pet.Status = req.Status
	pet.PetImage = req.PetImage

	data1 := db.Save(&pet)
	if data1.Error != nil {
		log.Printf("Error in Pet Update: %v", data1.Error)
		return nil, data1.Error
	}

	return &pb.UpdatePetResponse{
		Message: "Pet Updated Successfully!",
	}, nil
}

// DeletePet - Delete a pet by ID
func (p *Pet) DeletePet(ctx context.Context, req *pb.DeletePetRequest) (*pb.DeletePetResponse, error) {
	db := Config.GetDB()

	var pet Models.Pet
	data := db.First(&pet, req.Id)

	if data.Error != nil {
		if errors.Is(data.Error, gorm.ErrRecordNotFound) {
			log.Printf("Pet ID %d not found", req.Id)
			return &pb.DeletePetResponse{
				Message: "Pet ID not found, cannot delete!",
			}, errors.New("pet ID not found")
		}
		return nil, data.Error
	}

	data1 := db.Delete(&pet)
	if data1.Error != nil {
		log.Printf("Error in Pet Deletion: %v", data1.Error)
		return nil, data1.Error
	}

	return &pb.DeletePetResponse{
		Message: "Pet Deleted Successfully!",
	}, nil
}
