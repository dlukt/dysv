package service_test

import (
	"context"
	"errors"
	"time"

	"github.com/deicod/dysv/internal/mocks"
	"github.com/deicod/dysv/internal/model"
	"github.com/deicod/dysv/internal/service"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AddressService", func() {
	var (
		mockRepo    *mocks.MockAddressRepo
		addrService *service.AddressService
		ctx         context.Context
		userID      string
	)

	BeforeEach(func() {
		mockRepo = mocks.NewMockAddressRepo()
		addrService = service.NewAddressService(mockRepo)
		ctx = context.Background()
		userID = "user_123"
	})

	Describe("CreateAddress", func() {
		It("should create a new address successfully", func() {
			addr := &model.Address{
				UserID:     userID,
				Label:      "Home",
				Line1:      "123 Street",
				City:       "City",
				PostalCode: "12345",
				Country:    "DE",
			}

			err := addrService.CreateAddress(ctx, addr)
			Expect(err).NotTo(HaveOccurred())
			Expect(addr.ID).NotTo(BeEmpty())
			Expect(addr.CreatedAt).NotTo(BeZero())
			Expect(addr.UpdatedAt).NotTo(BeZero())

			// Verify stored
			stored, err := mockRepo.Get(ctx, addr.ID, userID)
			Expect(err).NotTo(HaveOccurred())
			Expect(stored).NotTo(BeNil())
			Expect(stored.Line1).To(Equal("123 Street"))
		})

		It("should unset defaults if new address is default", func() {
			// Create existing default address
			existing := &model.Address{
				UserID:    userID,
				IsDefault: true,
				Label:     "Old Default",
			}
			err := mockRepo.Create(ctx, existing)
			Expect(err).NotTo(HaveOccurred())

			// Create new default address
			newAddr := &model.Address{
				UserID:    userID,
				IsDefault: true,
				Label:     "New Default",
			}

			err = addrService.CreateAddress(ctx, newAddr)
			Expect(err).NotTo(HaveOccurred())

			// Verify old default is unset
			updatedExisting, err := mockRepo.Get(ctx, existing.ID, userID)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedExisting.IsDefault).To(BeFalse())

			// Verify new address is default
			storedNew, err := mockRepo.Get(ctx, newAddr.ID, userID)
			Expect(err).NotTo(HaveOccurred())
			Expect(storedNew.IsDefault).To(BeTrue())
		})

		It("should fail if repo create fails", func() {
			mockRepo.CreateError = errors.New("db error")
			addr := &model.Address{UserID: userID}

			err := addrService.CreateAddress(ctx, addr)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("db error"))
		})
	})

	Describe("UpdateAddress", func() {
		var existingAddr *model.Address

		BeforeEach(func() {
			existingAddr = &model.Address{
				UserID: userID,
				Label:  "Work",
				Line1:  "Work Place",
			}
			err := mockRepo.Create(ctx, existingAddr)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should update address successfully", func() {
			existingAddr.Line1 = "New Work Place"
			err := addrService.UpdateAddress(ctx, existingAddr)
			Expect(err).NotTo(HaveOccurred())
			Expect(existingAddr.UpdatedAt).To(BeTemporally("~", time.Now(), time.Second))

			stored, err := mockRepo.Get(ctx, existingAddr.ID, userID)
			Expect(err).NotTo(HaveOccurred())
			Expect(stored.Line1).To(Equal("New Work Place"))
		})

		It("should unset defaults if updated address becomes default", func() {
			// Create another default address
			otherDefault := &model.Address{
				UserID:    userID,
				IsDefault: true,
			}
			mockRepo.Create(ctx, otherDefault)

			existingAddr.IsDefault = true
			err := addrService.UpdateAddress(ctx, existingAddr)
			Expect(err).NotTo(HaveOccurred())

			// Verify other default is unset
			storedOther, _ := mockRepo.Get(ctx, otherDefault.ID, userID)
			Expect(storedOther.IsDefault).To(BeFalse())

			// Verify current is default
			storedCurrent, _ := mockRepo.Get(ctx, existingAddr.ID, userID)
			Expect(storedCurrent.IsDefault).To(BeTrue())
		})
	})

	Describe("ListAddresses", func() {
		It("should list addresses for user", func() {
			mockRepo.Create(ctx, &model.Address{UserID: userID, Label: "A1"})
			mockRepo.Create(ctx, &model.Address{UserID: userID, Label: "A2"})
			mockRepo.Create(ctx, &model.Address{UserID: "other_user", Label: "B1"})

			list, err := addrService.ListAddresses(ctx, userID)
			Expect(err).NotTo(HaveOccurred())
			Expect(list).To(HaveLen(2))
		})
	})

	Describe("DeleteAddress", func() {
		It("should delete address", func() {
			addr := &model.Address{UserID: userID}
			mockRepo.Create(ctx, addr)

			err := addrService.DeleteAddress(ctx, addr.ID, userID)
			Expect(err).NotTo(HaveOccurred())

			stored, err := mockRepo.Get(ctx, addr.ID, userID)
			Expect(err).NotTo(HaveOccurred())
			Expect(stored).To(BeNil())
		})
	})
})
