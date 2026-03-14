package binding

import (
	"skillmanager/internal/model"
	"skillmanager/internal/service"
)

type InventoryBinding struct {
	service *service.InventoryService
}

func NewInventoryBinding(svc *service.InventoryService) *InventoryBinding {
	return &InventoryBinding{service: svc}
}

func (b *InventoryBinding) BuildReport() (*model.InventoryReport, error) {
	return b.service.BuildReport()
}
