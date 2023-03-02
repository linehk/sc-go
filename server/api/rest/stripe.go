package rest

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/go-chi/render"
	"github.com/stablecog/sc-go/database/ent"
	"github.com/stablecog/sc-go/server/requests"
	"github.com/stablecog/sc-go/server/responses"
	"github.com/stablecog/sc-go/utils"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/webhook"
	"golang.org/x/exp/slices"
)

var PriceIDs = map[int]string{
	// ultimate
	3: utils.GetEnv("STRIPE_ULTIMATE_PRICE_ID", "price_1Mf591ATa0ehBYTA6ggpEEkA"),
	// pro
	2: utils.GetEnv("STRIPE_PRO_PRICE_ID", "price_1Mf50bATa0ehBYTAPOcfnOjG"),
	// starter
	1: utils.GetEnv("STRIPE_STARTER_PRICE_ID", "price_1Mf56NATa0ehBYTAHkCUablG"),
}

var SinglePurchasePriceIDs = []string{
	"price_1MfRaaATa0ehBYTAVRW3LPdR",
}

// For creating customer portal session
func (c *RestAPI) HandleCreatePortalSession(w http.ResponseWriter, r *http.Request) {
	var user *ent.User
	if user = c.GetUserIfAuthenticated(w, r); user == nil {
		return
	}

	// Parse request body
	reqBody, _ := io.ReadAll(r.Body)
	var stripeReq requests.StripePortalRequest
	err := json.Unmarshal(reqBody, &stripeReq)
	if err != nil {
		responses.ErrUnableToParseJson(w, r)
		return
	}

	// Create portal session
	session, err := c.StripeClient.BillingPortalSessions.New(&stripe.BillingPortalSessionParams{
		Customer:  stripe.String(user.StripeCustomerID),
		ReturnURL: stripe.String(stripeReq.ReturnUrl),
	})

	if err != nil {
		log.Error("Error creating portal session", "err", err)
		responses.ErrInternalServerError(w, r, "An unknown error has occured")
		return
	}

	sessionResponse := responses.StripeSessionResponse{
		CustomerPortalURL: session.URL,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, sessionResponse)
}

// For creating a new subscription or upgrading one
// Rejects, if they have a subscription that is at a higher level than the target priceID
func (c *RestAPI) HandleCreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	var user *ent.User
	if user = c.GetUserIfAuthenticated(w, r); user == nil {
		return
	}
	// Parse request body
	reqBody, _ := io.ReadAll(r.Body)
	var stripeReq requests.StripeCheckoutRequest
	err := json.Unmarshal(reqBody, &stripeReq)
	if err != nil {
		responses.ErrUnableToParseJson(w, r)
		return
	}

	// Make sure price ID exists in map
	var targetPriceID string
	var targetPriceLevel int
	adhocPrice := false
	for level, priceID := range PriceIDs {
		if priceID == stripeReq.TargetPriceID {
			targetPriceID = priceID
			targetPriceLevel = level
			break
		}
	}
	if targetPriceID == "" {
		// Check if it's a single purchase price
		for _, priceID := range SinglePurchasePriceIDs {
			if priceID == stripeReq.TargetPriceID {
				targetPriceID = priceID
				adhocPrice = true
				break
			}
		}
	}
	if targetPriceID == "" {
		responses.ErrBadRequest(w, r, "invalid_price_id")
		return
	}

	// Validate currency
	if !slices.Contains([]string{"usd", "eur"}, stripeReq.Currency) {
		responses.ErrBadRequest(w, r, "invalid_currency")
		return
	}

	// Get subscription
	customer, err := c.StripeClient.Customers.Get(user.StripeCustomerID, &stripe.CustomerParams{
		Params: stripe.Params{
			Expand: []*string{
				stripe.String("subscriptions"),
			},
		},
	})

	if err != nil {
		log.Error("Error getting customer", "err", err)
		responses.ErrInternalServerError(w, r, "An unknown error has occured")
		return
	}

	var currentPriceID string
	if customer.Subscriptions != nil {
		for _, sub := range customer.Subscriptions.Data {
			if sub.Status == stripe.SubscriptionStatusActive && sub.CancelAt == 0 {
				for _, item := range sub.Items.Data {
					if item.Price.ID == targetPriceID {
						responses.ErrBadRequest(w, r, "already_subscribed")
						return
					}
					// If price ID is in map it's valid
					for _, priceID := range PriceIDs {
						if item.Price.ID == priceID {
							currentPriceID = item.Price.ID
							break
						}
					}
					// See if this is starter euro price id
					euroPriceId := utils.GetEnv("STRIPE_STARTER_EURO_PRICE_ID", "price_1Mf56NATa0ehBYTAHkCUablG")
					if item.Price.ID == euroPriceId {
						currentPriceID = item.Price.ID
					}
				}
				break
			}
		}
	}

	// If they don't have one, cannot buy adhoc
	if currentPriceID == "" && adhocPrice {
		responses.ErrBadRequest(w, r, "no_subscription")
		return
	}

	// If they have a current one, make sure they are upgrading
	if currentPriceID != "" && !adhocPrice {
		var currentPriceLevel int
		for level, priceID := range PriceIDs {
			if priceID == currentPriceID {
				currentPriceLevel = level
				break
			}
		}
		// Check euro
		euroPriceId := utils.GetEnv("STRIPE_STARTER_EURO_PRICE_ID", "price_1Mf56NATa0ehBYTAHkCUablG")
		if currentPriceID == euroPriceId {
			currentPriceLevel = 1
		}

		if currentPriceLevel >= targetPriceLevel {
			responses.ErrBadRequest(w, r, "cannot_downgrade")
			return
		}
	}

	mode := stripe.CheckoutSessionModeSubscription
	if adhocPrice {
		mode = stripe.CheckoutSessionModePayment
	}
	// Create checkout session
	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(user.StripeCustomerID),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(targetPriceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(mode)),
		SuccessURL: stripe.String(stripeReq.SuccessUrl),
		CancelURL:  stripe.String(stripeReq.CancelUrl),
		Currency:   stripe.String(stripeReq.Currency),
	}

	session, err := c.StripeClient.CheckoutSessions.New(params)
	if err != nil {
		log.Error("Error creating checkout session", "err", err)
		responses.ErrInternalServerError(w, r, "An unknown error has occured")
		return
	}

	sessionResponse := responses.StripeSessionResponse{
		CheckoutURL: session.URL,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, sessionResponse)
}

// HTTP Post - handle stripe subscription downgrade
// Rejects if they don't have a subscription, or if they are not downgrading
func (c *RestAPI) HandleSubscriptionDowngrade(w http.ResponseWriter, r *http.Request) {
	var user *ent.User
	if user = c.GetUserIfAuthenticated(w, r); user == nil {
		return
	}

	// Parse request body
	reqBody, _ := io.ReadAll(r.Body)
	var stripeReq requests.StripeDowngradeRequest
	err := json.Unmarshal(reqBody, &stripeReq)
	if err != nil {
		responses.ErrUnableToParseJson(w, r)
		return
	}

	// Make sure price ID exists in map
	var targetPriceID string
	var targetPriceLevel int
	for level, priceID := range PriceIDs {
		if priceID == stripeReq.TargetPriceID {
			targetPriceID = priceID
			targetPriceLevel = level
			break
		}
	}
	if targetPriceID == "" {
		responses.ErrBadRequest(w, r, "invalid_price_id")
		return
	}

	// Get subscription
	customer, err := c.StripeClient.Customers.Get(user.StripeCustomerID, &stripe.CustomerParams{
		Params: stripe.Params{
			Expand: []*string{
				stripe.String("subscriptions"),
			},
		},
	})

	if err != nil {
		log.Error("Error getting customer", "err", err)
		responses.ErrInternalServerError(w, r, "An unknown error has occured")
		return
	}

	if customer.Subscriptions == nil || len(customer.Subscriptions.Data) == 0 || customer.Subscriptions.TotalCount == 0 {
		responses.ErrBadRequest(w, r, "no_active_subscription")
		return
	}

	var currentPriceID string
	var currentSubId string
	var currentItemId string
	for _, sub := range customer.Subscriptions.Data {
		if sub.Status == stripe.SubscriptionStatusActive && sub.CancelAt == 0 {
			for _, item := range sub.Items.Data {
				// If price ID is in map it's valid
				for _, priceID := range PriceIDs {
					if item.Price.ID == priceID {
						currentPriceID = item.Price.ID
						currentSubId = sub.ID
						currentItemId = item.ID
						break
					}
				}
				// Check if euro price ID
				euroPriceId := utils.GetEnv("STRIPE_STARTER_EURO_PRICE_ID", "price_1Mf56NATa0ehBYTAHkCUablG")
				if item.Price.ID == euroPriceId {
					currentPriceID = item.Price.ID
					currentSubId = sub.ID
					currentItemId = item.ID
				}
				break
			}
		}
	}

	if currentPriceID == "" {
		responses.ErrBadRequest(w, r, "no_active_subscription")
		return
	}

	if currentPriceID == targetPriceID {
		responses.ErrBadRequest(w, r, "not_lower")
		return
	}

	// Make sure this is a downgrade
	for level, priceID := range PriceIDs {
		if priceID == currentPriceID {
			if level <= targetPriceLevel {
				responses.ErrBadRequest(w, r, "not_lower")
				return
			}
			break
		}
	}

	// Execute subscription update
	_, err = c.StripeClient.Subscriptions.Update(currentSubId, &stripe.SubscriptionParams{
		ProrationBehavior: stripe.String("none"),
		Items: []*stripe.SubscriptionItemsParams{
			{
				ID:    stripe.String(currentItemId),
				Price: stripe.String(targetPriceID),
			},
		},
	})

	if err != nil {
		log.Error("Error updating subscription", "err", err)
		responses.ErrInternalServerError(w, r, "An unknown error has occured")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"success": true,
	})
}

// Handle stripe webhooks in the following ways:
// invoice.payment_succeeded
//   - Apply credits to user depending on type (subscription, adhoc)
//   - For subscriptions, set active_product_id
//
// customer.subscription.deleted"
//   - For an immediate cancellation, we set active_product_id to nil if this is a cancellation
//     of the product ID we currently have set for them. (In case they upgraded, it won't unset their upgrade)
//
// customer.subscription.created
//   - For a subscription upgrade, we cancel all old subscriptions
func (c *RestAPI) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("Unable reading stripe webhook body", "err", err)
		responses.ErrBadRequest(w, r, "invalid stripe webhook body")
		return
	}

	// Verify signature
	endpointSecret := utils.GetEnv("STRIPE_ENDPOINT_SECRET", "")

	event, err := webhook.ConstructEvent(reqBody, r.Header.Get("Stripe-Signature"), endpointSecret)
	if err != nil {
		log.Error("Unable verifying stripe webhook signature", "err", err)
		responses.ErrBadRequest(w, r, "invalid stripe webhook signature")
		return
	}

	switch event.Type {
	// For subscription upgrades, we want to cancel all old subscriptions
	case "customer.subscription.created":
		newSub, err := stripeObjectMapToSubscriptionObject(event.Data.Object)
		var newProduct string
		var oldProduct string
		if err != nil || newSub == nil {
			log.Error("Unable parsing stripe subscription object", "err", err)
			responses.ErrInternalServerError(w, r, err.Error())
			return
		}
		if newSub.Items != nil && len(newSub.Items.Data) > 0 && newSub.Items.Data[0].Price != nil && newSub.Items.Data[0].Price.Product != nil {
			newProduct = newSub.Items.Data[0].Price.Product.ID
		}
		// We need to see if they have more than one subscription
		subIter := c.StripeClient.Subscriptions.List(&stripe.SubscriptionListParams{
			Customer: stripe.String(newSub.Customer.ID),
		})
		for subIter.Next() {
			sub := subIter.Subscription()
			if sub.ID != newSub.ID {
				if sub.Items != nil && len(sub.Items.Data) > 0 && sub.Items.Data[0].Price != nil && sub.Items.Data[0].Price.Product != nil {
					oldProduct = newSub.Items.Data[0].Price.Product.ID
				}
				// We need to cancel this subscription
				_, err := c.StripeClient.Subscriptions.Cancel(sub.ID, &stripe.SubscriptionCancelParams{
					Prorate: stripe.Bool(false),
				})
				if err != nil {
					log.Error("Unable canceling stripe subscription", "err", err)
					responses.ErrInternalServerError(w, r, err.Error())
					return
				}
			}
		}
		// Analytics
		if newProduct != "" && oldProduct != "" {
			go func() {
				user, err := c.Repo.GetUserByStripeCustomerId(newSub.Customer.ID)
				if err != nil {
					log.Error("Unable getting user from stripe customer id in upgrade subscription event", "err", err)
					return
				}
				go c.Track.SubscriptionUpgraded(user, oldProduct, newProduct)
			}()
		}
	case "customer.subscription.deleted":
		sub, err := stripeObjectMapToCustomSubscriptionObject(event.Data.Object)
		if err != nil || sub == nil {
			log.Error("Unable parsing stripe subscription object", "err", err)
			responses.ErrInternalServerError(w, r, err.Error())
			return
		}
		user, err := c.Repo.GetUserByStripeCustomerId(sub.Customer)
		if err != nil {
			log.Error("Unable getting user from stripe customer id", "err", err)
			responses.ErrInternalServerError(w, r, err.Error())
			return
		} else if user == nil {
			log.Error("User does not exist with stripe customer id: %s", sub.Customer)
			responses.ErrInternalServerError(w, r, "User does not exist with stripe customer id")
			return
		}
		// Get product Id from subscription
		if sub.Items != nil && len(sub.Items.Data) > 0 && sub.Items.Data[0].Price != nil {
			affected, err := c.Repo.UnsetActiveProductID(user.ID, sub.Items.Data[0].Price.Product, nil)
			if err != nil {
				log.Error("Unable unsetting stripe product id", "err", err)
				responses.ErrInternalServerError(w, r, err.Error())
				return
			}
			if affected > 0 {
				// Subscription cancelled
				go c.Track.SubscriptionCancelled(user, sub.Items.Data[0].Price.Product)
			}
		}
	case "invoice.payment_succeeded":
		// We can parse the object as an invoice since that's the only thing we care about
		invoice, err := stripeObjectMapToInvoiceObject(event.Data.Object)
		if err != nil || invoice == nil {
			log.Error("Unable parsing stripe invoice object", "err", err)
			responses.ErrInternalServerError(w, r, err.Error())
			return
		}

		// We only care about renewal (cycle), create, and manual
		if invoice.BillingReason != InvoiceBillingReasonSubscriptionCycle && invoice.BillingReason != InvoiceBillingReasonSubscriptionCreate && invoice.BillingReason != InvoiceBillingReasonManual {
			render.Status(r, http.StatusOK)
			render.PlainText(w, r, "OK")
			return
		}

		if invoice.Lines == nil {
			log.Error("Stripe invoice lines is nil %s", invoice.ID)
			responses.ErrInternalServerError(w, r, "Stripe invoice lines is nil")
			return
		}

		for _, line := range invoice.Lines.Data {
			var product string
			if line.Plan == nil && invoice.BillingReason != InvoiceBillingReasonManual {
				log.Error("Stripe plan is nil in line item %s", line.ID)
				responses.ErrInternalServerError(w, r, "Stripe plan is nil in line item")
				return
			}

			if line.Price == nil && invoice.BillingReason == InvoiceBillingReasonManual {
				log.Error("Stripe price is nil in line item %s", line.ID)
				responses.ErrInternalServerError(w, r, "Stripe price is nil in line item")
				return
			}

			if invoice.BillingReason == InvoiceBillingReasonManual {
				product = line.Price.Product
			} else {
				product = line.Plan.Product
			}

			if product == "" {
				log.Error("Stripe product is nil in line item %s", line.ID)
				responses.ErrInternalServerError(w, r, "Stripe product is nil in line item")
				return
			}

			// Get user from customer ID
			user, err := c.Repo.GetUserByStripeCustomerId(invoice.Customer)
			if err != nil {
				log.Error("Unable getting user from stripe customer id", "err", err)
				responses.ErrInternalServerError(w, r, err.Error())
				return
			} else if user == nil {
				log.Error("User does not exist with stripe customer id: %s", invoice.Customer)
				responses.ErrInternalServerError(w, r, "User does not exist with stripe customer id")
				return
			}

			// Get the credit type for this plan
			creditType, err := c.Repo.GetCreditTypeByStripeProductID(product)
			if err != nil {
				log.Error("Unable getting credit type from stripe product id", "err", err)
				responses.ErrInternalServerError(w, r, err.Error())
				return
			} else if creditType == nil {
				log.Error("Credit type does not exist with stripe product id: %s", line.Plan.Product)
				responses.ErrInternalServerError(w, r, "Credit type does not exist with stripe product id")
				return
			}

			if invoice.BillingReason == InvoiceBillingReasonManual {
				// Ad-hoc credit add
				_, err = c.Repo.AddAdhocCreditsIfEligible(creditType, user.ID, line.ID)
				if err != nil {
					log.Error("Unable adding credits to user %s: %v", user.ID.String(), err)
					responses.ErrInternalServerError(w, r, err.Error())
					return
				}
				go c.Track.CreditPurchase(user, product, int(creditType.Amount))
			} else {
				expiresAt := utils.SecondsSinceEpochToTime(line.Period.End)
				// Update user credit
				if err := c.Repo.WithTx(func(tx *ent.Tx) error {
					client := tx.Client()
					_, err = c.Repo.AddCreditsIfEligible(creditType, user.ID, expiresAt, line.ID, client)
					if err != nil {
						log.Error("Unable adding credits to user %s: %v", user.ID.String(), err)
						responses.ErrInternalServerError(w, r, err.Error())
						return err
					}
					if user.ActiveProductID == nil {
						// New subscriber
						go c.Track.Subscription(user, product)
					} else {
						// Renewal
						go c.Track.SubscriptionRenewal(user, product)
					}
					err = c.Repo.SetActiveProductID(user.ID, product, client)
					if err != nil {
						log.Error("Unable setting stripe product id for user %s: %v", user.ID.String(), err)
						responses.ErrInternalServerError(w, r, err.Error())
						return err
					}
					return nil
				}); err != nil {
					log.Error("Unable adding credits to user %s: %v", user.ID.String(), err)
					return
				}
			}
		}
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "OK")
}

// Parse generic object into stripe invoice struct
func stripeObjectMapToInvoiceObject(obj map[string]interface{}) (*Invoice, error) {
	marshalled, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	var invoice Invoice
	err = json.Unmarshal(marshalled, &invoice)
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

// Parse generic object into stripe subscription struct
func stripeObjectMapToSubscriptionObject(obj map[string]interface{}) (*stripe.Subscription, error) {
	marshalled, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	var subscription stripe.Subscription
	err = json.Unmarshal(marshalled, &subscription)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// Parse generic object into custom stripe subscription struct with correct types
func stripeObjectMapToCustomSubscriptionObject(obj map[string]interface{}) (*Subscription, error) {
	marshalled, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	var subscription Subscription
	err = json.Unmarshal(marshalled, &subscription)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// ! Stripe types are busted so we modify the ones included in their lib
// InvoiceBillingReason is the reason why a given invoice was created
type InvoiceBillingReason string

// List of values that InvoiceBillingReason can take.
const (
	InvoiceBillingReasonManual                InvoiceBillingReason = "manual"
	InvoiceBillingReasonSubscription          InvoiceBillingReason = "subscription"
	InvoiceBillingReasonSubscriptionCreate    InvoiceBillingReason = "subscription_create"
	InvoiceBillingReasonSubscriptionCycle     InvoiceBillingReason = "subscription_cycle"
	InvoiceBillingReasonSubscriptionThreshold InvoiceBillingReason = "subscription_threshold"
	InvoiceBillingReasonSubscriptionUpdate    InvoiceBillingReason = "subscription_update"
	InvoiceBillingReasonUpcoming              InvoiceBillingReason = "upcoming"
)

// ListMeta is the structure that contains the common properties
// of List iterators. The Count property is only populated if the
// total_count include option is passed in (see tests for example).
type ListMeta struct {
	HasMore    bool   `json:"has_more"`
	TotalCount uint32 `json:"total_count"`
	URL        string `json:"url"`
}

// Period is a structure representing a start and end dates.
type Period struct {
	End   int64 `json:"end"`
	Start int64 `json:"start"`
}

type Plan struct {
	Product string `json:"product"`
}

type Price struct {
	Product string `json:"product"`
}

// InvoiceLine is the resource representing a Stripe invoice line item.
// For more details see https://stripe.com/docs/api#invoice_line_item_object.
type InvoiceLine struct {
	ID     string  `json:"id"`
	Period *Period `json:"period"`
	Plan   *Plan   `json:"plan"`
	Price  *Price  `json:"price"`
}

type InvoiceLineList struct {
	ListMeta
	Data []*InvoiceLine `json:"data"`
}

type Invoice struct {
	ID            string               `json:"id"`
	BillingReason InvoiceBillingReason `json:"billing_reason"`
	Lines         *InvoiceLineList     `json:"lines"`
	Customer      string               `json:"customer"`
}

// Subscription object is also pbroken in stripe
type SubscriptionItem struct {
	Price *Price `json:"price"`
}
type SubscriptionItemList struct {
	Data []*SubscriptionItem `json:"data"`
}
type Subscription struct {
	Items    *SubscriptionItemList `json:"items"`
	Customer string                `json:"customer"`
}
