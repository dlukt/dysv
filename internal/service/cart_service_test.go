/*
Copyright © 2025 Darko Luketic <info@icod.de>
*/
package service_test

import (
	"github.com/deicod/dysv/internal/model"
	"github.com/deicod/dysv/internal/service"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cart Service", func() {
	Describe("Plans", func() {
		It("should have static-micro plan", func() {
			plan, ok := service.Plans["static-micro"]
			Expect(ok).To(BeTrue())
			Expect(plan.ID).To(Equal("static-micro"))
			Expect(plan.Name).To(Equal("Static Micro"))
			Expect(plan.MonthlyPrice).To(BeNumerically("==", 3.90))
		})

		It("should have node-starter plan", func() {
			plan, ok := service.Plans["node-starter"]
			Expect(ok).To(BeTrue())
			Expect(plan.ID).To(Equal("node-starter"))
			Expect(plan.Name).To(Equal("Node Starter"))
			Expect(plan.MonthlyPrice).To(BeNumerically("==", 9.90))
		})

		It("should have node-pro plan", func() {
			plan, ok := service.Plans["node-pro"]
			Expect(ok).To(BeTrue())
			Expect(plan.ID).To(Equal("node-pro"))
			Expect(plan.Name).To(Equal("Node Pro"))
			Expect(plan.MonthlyPrice).To(BeNumerically("==", 39.90))
		})

		It("should have 3 plans total", func() {
			Expect(service.Plans).To(HaveLen(3))
		})
	})

	Describe("Addons", func() {
		It("should have de-domain addon", func() {
			addon, ok := service.Addons["de-domain"]
			Expect(ok).To(BeTrue())
			Expect(addon.ID).To(Equal("de-domain"))
			Expect(addon.Name).To(Equal(".de Domain"))
			Expect(addon.MonthlyPrice).To(BeNumerically("==", 1.00))
		})
	})

	Describe("YearlyDiscountMonths", func() {
		It("should be 2 months discount", func() {
			Expect(service.YearlyDiscountMonths).To(Equal(2))
		})
	})

	Describe("GetCartTotal", func() {
		var cartService *service.CartService

		BeforeEach(func() {
			// CartService with nil repo for total calculation tests
			cartService = service.NewCartService(nil)
		})

		It("should calculate monthly total correctly for single plan", func() {
			cart := &model.Cart{
				Items: []model.LineItem{
					{ItemID: "static-micro", ItemType: "plan", Name: "Static Micro", Price: 3.90, Quantity: 1},
				},
				BillingCycle: model.BillingMonthly,
			}

			monthly, yearly := cartService.GetCartTotal(cart)

			Expect(monthly).To(BeNumerically("~", 3.90, 0.01))
			Expect(yearly).To(BeNumerically("~", 39.00, 0.01)) // 3.90 * 10 months
		})

		It("should calculate monthly total correctly for plan with addon", func() {
			cart := &model.Cart{
				Items: []model.LineItem{
					{ItemID: "node-starter", ItemType: "plan", Name: "Node Starter", Price: 9.90, Quantity: 1},
					{ItemID: "de-domain", ItemType: "addon", Name: ".de Domain", Price: 1.00, Quantity: 1},
				},
				BillingCycle: model.BillingMonthly,
			}

			monthly, yearly := cartService.GetCartTotal(cart)

			Expect(monthly).To(BeNumerically("~", 10.90, 0.01)) // 9.90 + 1.00
			Expect(yearly).To(BeNumerically("~", 111.00, 0.01)) // plan 9.90 * 10 months + addon 1.00 * 12 months
		})

		It("should calculate yearly total with 2 months discount", func() {
			cart := &model.Cart{
				Items: []model.LineItem{
					{ItemID: "node-pro", ItemType: "plan", Name: "Node Pro", Price: 39.90, Quantity: 1},
				},
				BillingCycle: model.BillingYearly,
			}

			monthly, yearly := cartService.GetCartTotal(cart)

			Expect(monthly).To(BeNumerically("~", 39.90, 0.01))
			Expect(yearly).To(BeNumerically("~", 399.00, 0.01)) // 39.90 * 10 months (2 free)
		})

		It("should return zero for empty cart", func() {
			cart := &model.Cart{
				Items:        []model.LineItem{},
				BillingCycle: model.BillingMonthly,
			}

			monthly, yearly := cartService.GetCartTotal(cart)

			Expect(monthly).To(Equal(0.0))
			Expect(yearly).To(Equal(0.0))
		})

		It("should handle multiple quantity items", func() {
			cart := &model.Cart{
				Items: []model.LineItem{
					{ItemID: "de-domain", ItemType: "addon", Name: ".de Domain", Price: 1.00, Quantity: 3},
				},
				BillingCycle: model.BillingMonthly,
			}

			monthly, yearly := cartService.GetCartTotal(cart)

			Expect(monthly).To(BeNumerically("~", 3.00, 0.01))
			Expect(yearly).To(BeNumerically("~", 36.00, 0.01)) // addons billed for 12 months
		})
	})
})

var _ = Describe("Service Errors", func() {
	Describe("ErrInvalidPlan", func() {
		It("should have correct error message", func() {
			Expect(service.ErrInvalidPlan.Error()).To(Equal("invalid plan ID"))
		})
	})

	Describe("ErrInvalidAddon", func() {
		It("should have correct error message", func() {
			Expect(service.ErrInvalidAddon.Error()).To(Equal("invalid addon ID"))
		})
	})

	Describe("ErrEmptyCart", func() {
		It("should have correct error message", func() {
			Expect(service.ErrEmptyCart.Error()).To(Equal("cart is empty"))
		})
	})
})

var _ = Describe("Model Types", func() {
	Describe("BillingCycle", func() {
		It("should have monthly constant", func() {
			Expect(model.BillingMonthly).To(Equal(model.BillingCycle("monthly")))
		})

		It("should have yearly constant", func() {
			Expect(model.BillingYearly).To(Equal(model.BillingCycle("yearly")))
		})
	})

	Describe("LineItem", func() {
		It("should have correct JSON tags", func() {
			item := model.LineItem{
				ItemID:   "test-id",
				ItemType: "plan",
				Name:     "Test Plan",
				Price:    9.99,
				Quantity: 1,
			}

			Expect(item.ItemID).To(Equal("test-id"))
			Expect(item.ItemType).To(Equal("plan"))
			Expect(item.Name).To(Equal("Test Plan"))
			Expect(item.Price).To(BeNumerically("==", 9.99))
			Expect(item.Quantity).To(Equal(1))
		})
	})

	Describe("Cart", func() {
		It("should initialize with default values", func() {
			cart := model.Cart{
				SessionID:    "session-123",
				Items:        []model.LineItem{},
				BillingCycle: model.BillingMonthly,
			}

			Expect(cart.SessionID).To(Equal("session-123"))
			Expect(cart.Items).To(BeEmpty())
			Expect(cart.BillingCycle).To(Equal(model.BillingMonthly))
		})
	})

	Describe("Order", func() {
		It("should track order status", func() {
			order := model.Order{
				Status: "pending",
			}

			Expect(order.Status).To(Equal("pending"))

			order.Status = "paid"
			Expect(order.Status).To(Equal("paid"))
		})
	})
})
