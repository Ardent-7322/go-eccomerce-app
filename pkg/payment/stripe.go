package payment

import (
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
)

type PaymentClient interface {
	CreatePayment(amount float64, userId uint, orderId string) (*stripe.PaymentIntent, error)
	GetPaymentStatus(pId string) (*stripe.PaymentIntent, error)
}

type payment struct {
	stripeSecretKey string
	successUrl      string
	cancelUrl       string
}

// CreatePayment implements [PaymentClient].
func (p *payment) CreatePayment(
	amount float64,
	userId uint,
	orderId string,
) (*stripe.PaymentIntent, error) {

	stripe.Key = p.stripeSecretKey

	amountInCents := int64(math.Round(amount * 100))
	if amountInCents <= 0 {
		return nil, errors.New("invalid payment amount")
	}

	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(int64(amountInCents)),
		Currency:           stripe.String(string(stripe.CurrencyUSD)),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
	}

	params.AddMetadata("order_id", fmt.Sprintf("%s", orderId))
	params.AddMetadata("user_id", fmt.Sprintf("%d", userId))

	pi, err := paymentintent.New(params)
	if err != nil {
		log.Printf("stripe session error: %v", err)
		return nil, errors.New("payment intent creation failed")
	}

	return pi, nil
}

// GetPaymentStatus implements [PaymentClient].
func (p *payment) GetPaymentStatus(pId string) (*stripe.PaymentIntent, error) {

	stripe.Key = p.stripeSecretKey

	result, err := paymentintent.Get(pId, nil)
	if err != nil {
		log.Printf("error getting payment intent: %v", err)
		return nil, errors.New("get payment intent failed")
	}

	return result, nil
}

func NewPaymentClient(stripeSecretKey, successUrl, cancelUrl string) PaymentClient {
	return &payment{
		stripeSecretKey: stripeSecretKey,
		successUrl:      successUrl,
		cancelUrl:       cancelUrl,
	}
}

// params := &stripe.CheckoutSessionParams{
// 	PaymentMethodTypes: stripe.StringSlice([]string{"card"}),

// 	PaymentMethodOptions: &stripe.CheckoutSessionPaymentMethodOptionsParams{
// 		Card: &stripe.CheckoutSessionPaymentMethodOptionsCardParams{
// 			SetupFutureUsage: stripe.String("none"), // 👈 disables Link
// 		},
// 	},

// 	LineItems: []*stripe.CheckoutSessionLineItemParams{
// 		{
// 			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
// 				UnitAmount: stripe.Int64(amountInCents),
// 				Currency:   stripe.String("usd"),
// 				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
// 					Name: stripe.String("Order Payment"),
// 				},
// 			},
// 			Quantity: stripe.Int64(1),
// 		},
// 	},
