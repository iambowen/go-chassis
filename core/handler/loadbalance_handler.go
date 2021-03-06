package handler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cenkalti/backoff"
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/archaius"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/loadbalancer"

	"github.com/go-chassis/go-chassis/pkg/util"

	backoffUtil "github.com/go-chassis/go-chassis/pkg/backoff"

	"github.com/go-chassis/go-chassis/session"
)

// LBHandler loadbalancer handler struct
type LBHandler struct{}

func (lb *LBHandler) getEndpoint(i *invocation.Invocation, lbConfig control.LoadBalancingConfig) (string, error) {
	var strategyFun func() loadbalancer.Strategy
	var err error
	if i.Strategy == "" {
		i.Strategy = lbConfig.Strategy
		strategyFun, err = loadbalancer.GetStrategyPlugin(i.Strategy)
		if err != nil {
			lager.Logger.Errorf("lb error [%s] because of [%s]", loadbalancer.LBError{
				Message: "Get strategy [" + i.Strategy + "] failed."}.Error(), err.Error())
		}
	} else {
		strategyFun, err = loadbalancer.GetStrategyPlugin(i.Strategy)
		if err != nil {
			lager.Logger.Errorf("lb error [%s] because of [%s]", loadbalancer.LBError{
				Message: "Get strategy [" + i.Strategy + "] failed."}.Error(), err.Error())
		}
	}
	if len(i.Filters) == 0 {
		i.Filters = lbConfig.Filters
	}

	var sessionID string
	if i.Strategy == loadbalancer.StrategySessionStickiness {
		sessionID = getSessionID(i)
	}

	s, err := loadbalancer.BuildStrategy(i.SourceServiceID, i.MicroServiceName, i.Protocol,
		sessionID, i.Filters, strategyFun(), i.RouteTags)
	if err != nil {
		return "", err
	}

	ins, err := s.Pick()
	if err != nil {
		lbErr := loadbalancer.LBError{Message: err.Error()}
		return "", lbErr
	}

	var ep string
	if i.Protocol == "" {
		i.Protocol = archaius.GetString("cse.references."+i.MicroServiceName+".transport", ins.DefaultProtocol)
	}
	if i.Protocol == "" {
		for k := range ins.EndpointsMap {
			i.Protocol = k
			break
		}
	}
	ep, ok := ins.EndpointsMap[util.GenProtoEndPoint(i.Protocol, i.Port)]
	if !ok {
		errStr := fmt.Sprintf("No available instance support ["+i.Protocol+"] protocol,"+
			" msName: "+i.MicroServiceName+" %v", ins.EndpointsMap)
		lbErr := loadbalancer.LBError{Message: errStr}
		lager.Logger.Errorf(lbErr.Error())
		return "", lbErr
	}
	return ep, nil
}

// Handle to handle the load balancing
func (lb *LBHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	lbConfig := control.DefaultPanel.GetLoadBalancing(*i)
	if !lbConfig.RetryEnabled {
		lb.handleWithNoRetry(chain, i, lbConfig, cb)
	} else {
		lb.handleWithRetry(chain, i, lbConfig, cb)
	}
}

func (lb *LBHandler) handleWithNoRetry(chain *Chain, i *invocation.Invocation, lbConfig control.LoadBalancingConfig, cb invocation.ResponseCallBack) {
	ep, err := lb.getEndpoint(i, lbConfig)
	if err != nil {
		writeErr(err, cb)
		return
	}

	i.Endpoint = ep
	chain.Next(i, cb)
}

func (lb *LBHandler) handleWithRetry(chain *Chain, i *invocation.Invocation, lbConfig control.LoadBalancingConfig, cb invocation.ResponseCallBack) {
	retryOnSame := lbConfig.RetryOnSame
	retryOnNext := lbConfig.RetryOnNext
	handlerIndex := chain.HandlerIndex
	var invResp *invocation.Response
	for j := 0; j < retryOnNext+1; j++ {
		// exchange and retry on the next server
		ep, err := lb.getEndpoint(i, lbConfig)
		if err != nil {
			// if get endpoint failed, no need to retry
			writeErr(err, cb)
			return
		}
		// retry on the same server
		lbBackoff := backoffUtil.GetBackOff(lbConfig.BackOffKind, lbConfig.BackOffMin, lbConfig.BackOffMax)
		callTimes := 0
		operation := func() error {
			if callTimes == retryOnSame+1 {
				return backoff.Permanent(errors.New("retry times expires"))
			}
			callTimes++
			i.Endpoint = ep
			var respErr error
			chain.HandlerIndex = handlerIndex
			chain.Next(i, func(r *invocation.Response) error {
				if r != nil {
					invResp = r
					respErr = invResp.Err
					return invResp.Err
				}
				return nil
			})
			return respErr
		}
		if err = backoff.Retry(operation, lbBackoff); err == nil {
			break
		}
	}
	if invResp == nil {
		invResp = &invocation.Response{}
	}
	cb(invResp)
}

// Name returns loadbalancer string
func (lb *LBHandler) Name() string {
	return "loadbalancer"
}

func newLBHandler() Handler {
	return &LBHandler{}
}

func getSessionID(i *invocation.Invocation) string {
	var metadata interface{}

	switch i.Args.(type) {
	case *rest.Request:
		req := i.Args.(*rest.Request)
		value := req.GetCookie(common.LBSessionID)
		if value != "" {
			metadata = value
		}
	default:
		value := session.GetContextMetadata(i.Ctx, common.LBSessionID)
		if value != "" {
			cookieKey := strings.Split(string(value), "=")
			if len(cookieKey) > 1 {
				metadata = cookieKey[1]
			}
		}
	}

	if metadata == nil {
		metadata = ""
	}

	return metadata.(string)
}
