package Handlers

import (
	"context"
	"fmt"
	"pet/Config"
	Models "pet/Models"
	pb "pet/pet/pb"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type PetHealth struct {
	pb.UnimplementedPetHealthServiceServer
}

// Converts *timestamppb.Timestamp to *time.Time
func parseTimestamp(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

// Converts *time.Time to *timestamppb.Timestamp
func formatTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

// CreatePetHealth - Adds a new pet health record
func (s *PetHealth) CreatePetHealth(ctx context.Context, req *pb.CreatePetHealthRequest) (*pb.CreatePetHealthResponse, error) {
	db := Config.GetDB()

	var pet Models.Pet
	if err := db.First(&pet, req.PetId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("pet with ID %d not found", req.PetId)
		}
		return nil, fmt.Errorf("error finding pet: %v", err)
	}

	petHealth := Models.PetHealth{
		PetID:            uint(req.PetId),
		Vaccinated:       req.Vaccinated,
		VaccinationDate:  parseTimestamp(req.VaccinationDate),
		Allergies:        req.Allergies,
		LastVetVisitDate: parseTimestamp(req.LastVetVisitDate),
		HealthRemarks:    req.HealthRemarks,
	}

	// Save record to the database
	if err := db.Create(&petHealth).Error; err != nil {
		return nil, fmt.Errorf("failed to create pet health record: %v", err)
	}

	return &pb.CreatePetHealthResponse{
		Message: "Pet health record created successfully",
	}, nil
}

func (s *PetHealth) GetPetHealth(ctx context.Context, req *pb.GetPetHealthIDRequest) (*pb.GetPetHealthIDResponse, error) {
	db := Config.GetDB()

	var petHealth Models.PetHealth
	if err := db.First(&petHealth, req.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("pet health record with ID %d not found", req.Id)
		}
		return nil, fmt.Errorf("error retrieving pet health record: %v", err)
	}

	return &pb.GetPetHealthIDResponse{
		HealthId:         int32(petHealth.HealthID),
		PetId:            int32(petHealth.PetID),
		Vaccinated:       petHealth.Vaccinated,
		VaccinationDate:  formatTimestamp(petHealth.VaccinationDate),
		Allergies:        petHealth.Allergies,
		LastVetVisitDate: formatTimestamp(petHealth.LastVetVisitDate),
		HealthRemarks:    petHealth.HealthRemarks,
	}, nil
}

func (s *PetHealth) ListPetHealth(ctx context.Context, req *emptypb.Empty) (*pb.ListPetHealthResponse, error) {
	db := Config.GetDB()

	var petHealths []Models.PetHealth
	result := db.Find(&petHealths)
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving pet health records: %v", result.Error)
	}

	var petHealthResponses []*pb.GetPetHealthIDResponse
	for _, record := range petHealths {
		petHealthResponses = append(petHealthResponses, &pb.GetPetHealthIDResponse{
			HealthId:         int32(record.HealthID),
			PetId:            int32(record.PetID),
			Vaccinated:       record.Vaccinated,
			VaccinationDate:  formatTimestamp(record.VaccinationDate),
			Allergies:        record.Allergies,
			LastVetVisitDate: formatTimestamp(record.LastVetVisitDate),
			HealthRemarks:    record.HealthRemarks,
		})
	}

	return &pb.ListPetHealthResponse{PetHealths: petHealthResponses}, nil
}

func (s *PetHealth) UpdatePetHealth(ctx context.Context, req *pb.UpdatePetHealthRequest) (*pb.UpdatePetHealthResponse, error) {
	db := Config.GetDB()

	var petHealth Models.PetHealth
	if err := db.First(&petHealth, req.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("pet health record with ID %d not found", req.Id)
		}
		return nil, fmt.Errorf("error finding pet health record: %v", err)
	}

	petHealth.PetID = uint(req.PetId)
	petHealth.Vaccinated = req.Vaccinated
	petHealth.VaccinationDate = parseTimestamp(req.VaccinationDate)
	petHealth.Allergies = req.Allergies
	petHealth.LastVetVisitDate = parseTimestamp(req.LastVetVisitDate)
	petHealth.HealthRemarks = req.HealthRemarks

	if err := db.Save(&petHealth).Error; err != nil {
		return nil, fmt.Errorf("failed to update pet health record: %v", err)
	}

	return &pb.UpdatePetHealthResponse{
		Message: "Pet health record updated successfully",
	}, nil
}

func (s *PetHealth) DeletePetHealth(ctx context.Context, req *pb.DeletePetHealthRequest) (*pb.DeletePetHealthResponse, error) {
	db := Config.GetDB()

	if err := db.Delete(&Models.PetHealth{}, req.Id).Error; err != nil {
		return nil, fmt.Errorf("failed to delete pet health record with ID %d: %v", req.Id, err)
	}

	return &pb.DeletePetHealthResponse{
		Message: fmt.Sprintf("Pet health record with ID %d deleted successfully", req.Id),
	}, nil
}
