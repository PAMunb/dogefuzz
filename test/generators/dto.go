package generators

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/dto"
)

func NewContractDTOGen() *dto.NewContractDTO {
	return &dto.NewContractDTO{
		Name:   gofakeit.Name(),
		Source: gofakeit.LetterN(255),
	}

}

func ContractDTOGen() *dto.ContractDTO {
	return &dto.ContractDTO{
		Id:      gofakeit.LetterN(255),
		Address: SmartContractGen(),
		Source:  gofakeit.LetterN(255),
		Name:    gofakeit.Name(),
	}
}