package onfleet

import "context"

type AdminsService service

type Admin struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AdminServiceInterface interface {
	List(ctx context.Context) ([]Admin, error)
}

func (s *AdminsService) List(ctx context.Context) ([]Admin, error) {
	req, err := s.client.NewRequest("GET", "admins", nil)
	if err != nil {
		return nil, err
	}

	var admins []Admin
	err = s.client.Do(ctx, req, &admins)
	if err != nil {
		return nil, err
	}

	return admins, nil
}
