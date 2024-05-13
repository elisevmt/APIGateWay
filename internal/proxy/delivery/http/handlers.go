package proxy_delivery_http

import (
	GL "APIGateWay/pkg/guzzle_logger"
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
	//"strings"
	"time"

	"APIGateWay/internal/proxy"

	"github.com/gofiber/fiber/v2"
	fiber_proxy "github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/valyala/fasthttp"
)

type proxyHTTPHadnler struct {
	proxyUC proxy.UC
	cli     *fasthttp.Client
	gl      GL.API
}

func NewProxyHTTPHandlers(proxyUC proxy.UC, gl GL.API) proxy.HTTPHandlers {
	cli := fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return fasthttp.DialTimeout(addr, time.Second*60)
		},
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		NoDefaultUserAgentHeader: true,
	}
	return &proxyHTTPHadnler{
		proxyUC: proxyUC,
		cli:     &cli,
		gl:      gl,
	}
}

func (p *proxyHTTPHadnler) Proxy() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var finalUrl string
		subURL := c.Params("*")
		queryString := string(c.Request().URI().QueryString())
		serviceId, err := strconv.Atoi(c.Params("service_id"))
		if err != nil {
			return err
		}
		serviceID := int64(serviceId)
		instance, err := p.proxyUC.GetProxyInstance(&serviceID)
		if err != nil {
			return err
		}
		defer p.proxyUC.DecreaseLoad(instance.Url)
		//body := c.Body()
		if len(queryString) != 0 {
			finalUrl = *instance.Url + "/" + subURL + "?" + queryString
		} else {
			finalUrl = *instance.Url + "/" + subURL
		}
		err = fiber_proxy.Do(c, finalUrl, p.cli)
		fmt.Println("url", finalUrl)
		if err != nil {
			if serviceID == 1{
				go p.gl.SendLog(GL.LEVEL_ERROR, err.Error(), fmt.Sprintf("TRON FATAL ERROR %s\nproxy url: %s", subURL, finalUrl))
			} 
			return err
		}
		//response := c.Response()
		return nil
	}
}
