package payos

import (
	"graves/pkg/config"
	"strconv"

	payosapi "github.com/payOSHQ/payos-lib-golang"

	"time"
)

func GenerateNumber() int64 {
	millis := time.Now().UnixNano() / int64(time.Millisecond)
	millisStr := strconv.FormatInt(millis, 10)
	number, _ := strconv.Atoi(millisStr[len(millisStr)-6:])
	return int64(number)
}

type Payos interface {
	GetTimeOut() int32
	CreatePaymentLink(body payosapi.CheckoutRequestType) (*payosapi.CheckoutResponseDataType, error)
	GetPaymentLinkInfo(orderId string) (*payosapi.PaymentLinkDataType, error)
	CancelPaymentLink(orderId string) (*payosapi.PaymentLinkDataType, error)
	VerifyPaymentWebhookData(body payosapi.WebhookType) (*payosapi.WebhookDataType, error)
}

type payos struct {
	cancleUrl string
	returnUrl string
	timeout   int32
}

var payosInstance Payos

func GetInstance() Payos {
	cfg, err := config.GetInstance()
	if err != nil || cfg.Payos == nil {
		return &Noop{}
	}
	if payosInstance != nil {
		return payosInstance
	}
	payosapi.Key(cfg.Payos.ClientId, cfg.Payos.ApiKey, cfg.Payos.Checksum)
	payosInstance = &payos{
		cancleUrl: cfg.Payos.CancelUrl,
		returnUrl: cfg.Payos.SuccessUrl,
		timeout:   int32(cfg.Payos.TimeOut),
	}
	return payosInstance
}

func (p *payos) GetTimeOut() int32 {
	return p.timeout
}

func (p *payos) CreatePaymentLink(body payosapi.CheckoutRequestType) (*payosapi.CheckoutResponseDataType, error) {
	expAt := int(time.Now().Add(time.Duration(p.timeout) * time.Second).Unix())
	body.OrderCode = GenerateNumber()
	body.CancelUrl = p.cancleUrl
	body.ReturnUrl = p.returnUrl
	body.ExpiredAt = &expAt
	return payosapi.CreatePaymentLink(body)
}

func (*payos) GetPaymentLinkInfo(orderId string) (*payosapi.PaymentLinkDataType, error) {
	return payosapi.GetPaymentLinkInformation(orderId)
}

func (*payos) CancelPaymentLink(orderId string) (*payosapi.PaymentLinkDataType, error) {
	return payosapi.CancelPaymentLink(orderId, nil)
}

func (*payos) VerifyPaymentWebhookData(webhookDataReq payosapi.WebhookType) (*payosapi.WebhookDataType, error) {
	return payosapi.VerifyPaymentWebhookData(webhookDataReq)
}
