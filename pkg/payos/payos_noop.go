package payos

import (
	"fmt"

	payosapi "github.com/payOSHQ/payos-lib-golang"
)

type Noop struct {
}

func (n *Noop) GetTimeOut() int32 {
	return 0
}

func (n *Noop) CreatePaymentLink(body payosapi.CheckoutRequestType) (*payosapi.CheckoutResponseDataType, error) {
	return nil, fmt.Errorf("create payment not implemented")
}

func (n *Noop) GetPaymentLinkInfo(orderId string) (*payosapi.PaymentLinkDataType, error) {
	return nil, fmt.Errorf("get payment link info not implemented")
}

func (n *Noop) CancelPaymentLink(orderId string) (*payosapi.PaymentLinkDataType, error) {
	return nil, fmt.Errorf("cancel payment link not implemented")
}

func (n *Noop) VerifyPaymentWebhookData(body payosapi.WebhookType) (*payosapi.WebhookDataType, error) {
	return nil, fmt.Errorf("verify payment webhook data not implemented")
}
