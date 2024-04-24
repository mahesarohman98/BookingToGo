package customer

import (
	commonerrors "booking_to_go/internal/errors"
	"time"
)

type FamilyList struct {
	ID          string
	Relation    string
	Name        string
	DateOfBirth time.Time
}

var (
	ErrDuplicateFamilyMember = commonerrors.NewIncorrectInputError("family member already exists", "duplicate-family-member")
	ErrFamilyMemberNotFond   = commonerrors.NewIncorrectInputError("family member not found", "not-found-family-member")
)

func (c *Customer) AddFamilyMember(id string, relation string, name string, dateOfBirth time.Time) error {
	for _, fm := range c.FamilyList {
		if fm.ID == id {
			return ErrDuplicateFamilyMember
		}
	}

	c.FamilyList = append(c.FamilyList, FamilyList{
		ID:          id,
		Relation:    relation,
		Name:        name,
		DateOfBirth: dateOfBirth,
	})

	return nil
}

func (c *Customer) UpdateFamilyMember(familyListID string, relation string, name string, dateOfBirth *time.Time) error {
	index := -1
	for i, fm := range c.FamilyList {
		if fm.ID == familyListID {
			index = i
			break
		}
	}

	if index == -1 {
		return ErrFamilyMemberNotFond
	}

	if relation != "" {
		c.FamilyList[index].Relation = relation
	}

	if name != "" {
		c.FamilyList[index].Name = name
	}

	if dateOfBirth != nil {
		c.FamilyList[index].DateOfBirth = *dateOfBirth
	}

	return nil

}

func (c *Customer) RemoveFamilyMember(familyListID string) error {
	index := -1
	for i, fm := range c.FamilyList {
		if fm.ID == familyListID {
			index = i
			break
		}
	}

	if index == -1 {
		return ErrFamilyMemberNotFond
	}

	c.FamilyList = append(c.FamilyList[:index], c.FamilyList[index+1:]...)

	return nil
}
