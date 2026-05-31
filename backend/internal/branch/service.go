package branch

import (
	"context"
	"errors"
	"time"
)

type Service interface {
	GetAllBranches(ctx context.Context, role string, includeInactive bool) ([]Branch, error)
	CreateBranch(ctx context.Context, req CreateBranchRequest) error
	UpdateBranch(ctx context.Context, id int, req UpdateBranchRequest) error
	DeleteBranch(ctx context.Context, id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) GetAllBranches(ctx context.Context, role string, includeInactive bool) ([]Branch, error) {
	// Jika bukan Admin, paksa includeInactive menjadi false (Customer/Manager hanya lihat yang aktif)
	if role != "Admin" {
		includeInactive = false
	}
	
	branches, err := s.repo.FindAll(ctx, includeInactive)
	if err != nil {
		return nil, err
	}
	
	if branches == nil {
		branches = []Branch{}
	}
	
	return branches, nil
}

func (s *service) CreateBranch(ctx context.Context, req CreateBranchRequest) error {
	branch := &Branch{
		Name:      req.Name,
		Address:   req.Address,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	return s.repo.Create(ctx, branch)
}

func (s *service) UpdateBranch(ctx context.Context, id int, req UpdateBranchRequest) error {
	branch, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if branch == nil {
		return errors.New("branch not found")
	}

	if req.Name != nil {
		branch.Name = *req.Name
	}
	if req.Address != nil {
		branch.Address = *req.Address
	}
	if req.IsActive != nil {
		branch.IsActive = *req.IsActive
	}
	
	branch.UpdatedAt = time.Now()
	
	return s.repo.Update(ctx, branch)
}

func (s *service) DeleteBranch(ctx context.Context, id int) error {
	branch, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if branch == nil {
		return errors.New("branch not found")
	}

	// Soft delete
	branch.IsActive = false
	branch.UpdatedAt = time.Now()
	
	return s.repo.Update(ctx, branch)
}
