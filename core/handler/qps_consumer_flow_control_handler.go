package handler

import (
	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/qpslimiter"
)

// ConsumerRateLimiterHandler consumer rate limiter handler
type ConsumerRateLimiterHandler struct{}

// Handle is handles the consumer rate limiter APIs
func (rl *ConsumerRateLimiterHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	rlc := control.DefaultPanel.GetRateLimiting(*i, common.Consumer)
	if !rlc.Enabled {
		chain.Next(i, cb)

		return
	}
	//get operation meta info ms.schema, ms.schema.operation, ms
	rl.GetOrCreate(rlc)
	chain.Next(i, cb)
}

func newConsumerRateLimiterHandler() Handler {
	return &ConsumerRateLimiterHandler{}
}

// Name returns consumerratelimiter string
func (rl *ConsumerRateLimiterHandler) Name() string {
	return "consumerratelimiter"
}

// GetOrCreate is for getting or creating qps limiter meta data
func (rl *ConsumerRateLimiterHandler) GetOrCreate(rlc control.RateLimitingConfig) {
	qpslimiter.GetQPSTrafficLimiter().ProcessQPSTokenReq(rlc.Key, rlc.Rate)
	return
}
