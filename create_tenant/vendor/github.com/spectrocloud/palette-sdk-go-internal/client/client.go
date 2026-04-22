package client

import (
	"context"
	"crypto/tls"
	"net/http"

	openapiclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/spectrocloud/hapi/apiutil/transport"
	authC "github.com/spectrocloud/hapi/auth/client/v1"
	mgmtC "github.com/spectrocloud/hapi/mgmt/client/v1"
	"github.com/spectrocloud/hapi/models"
	clusterC "github.com/spectrocloud/hapi/spectrocluster/client/v1"
	systemC "github.com/spectrocloud/hapi/system/client/v1"
	userC "github.com/spectrocloud/hapi/user/client/v1"
)

type V1Client struct {
	AuthC    authC.ClientService
	ClusterC clusterC.ClientService
	MgmtC    mgmtC.ClientService
	UserC    userC.ClientService
	SystemC  systemC.ClientService

	ctx                context.Context
	apikey             string
	jwt                string
	username           string
	password           string
	systemScoped       bool
	hubbleUri          string
	projectUid         string
	schemes            []string
	insecureSkipVerify bool
	transportDebug     bool
	retryAttempts      int
}

func New(options ...func(*V1Client)) *V1Client {
	client := &V1Client{
		ctx:           context.Background(),
		retryAttempts: 0,
		schemes:       []string{"https"},
	}
	for _, o := range options {
		o(client)
	}
	client.AuthC = authC.New(client.getTransport(), strfmt.Default)
	client.ClusterC = clusterC.New(client.getTransport(), strfmt.Default)
	client.MgmtC = mgmtC.New(client.getTransport(), strfmt.Default)
	client.UserC = userC.New(client.getTransport(), strfmt.Default)
	client.SystemC = systemC.New(client.getTransport(), strfmt.Default)
	return client
}

func WithAPIKey(apiKey string) func(*V1Client) {
	return func(v *V1Client) {
		v.apikey = apiKey
	}
}

func WithJWT(jwt string) func(*V1Client) {
	return func(v *V1Client) {
		v.jwt = jwt
	}
}

func WithUsername(username string) func(*V1Client) {
	return func(v *V1Client) {
		v.username = username
	}
}

func WithPassword(password string) func(*V1Client) {
	return func(v *V1Client) {
		v.password = password
	}
}

func WithHubbleURI(hubbleUri string) func(*V1Client) {
	return func(v *V1Client) {
		v.hubbleUri = hubbleUri
	}
}

func WithInsecureSkipVerify(insecureSkipVerify bool) func(*V1Client) {
	return func(v *V1Client) {
		v.insecureSkipVerify = insecureSkipVerify
	}
}

func WithScopeProject(projectUid string) func(*V1Client) {
	return func(v *V1Client) {
		v.ctx = context.WithValue(v.ctx, transport.CUSTOM_HEADERS, transport.Values{
			HeaderMap: map[string]string{
				"ProjectUid": projectUid,
			}},
		)
	}
}

func WithScopeTenant() func(*V1Client) {
	return func(v *V1Client) {
		v.ctx = context.Background()
	}
}

func WithScopeSystem(username, password string) func(*V1Client) {
	return func(v *V1Client) {
		v.username = username
		v.password = password
		v.systemScoped = true
	}
}

func WithRetries(retries int) func(*V1Client) {
	return func(v *V1Client) {
		v.retryAttempts = retries
	}
}

func WithSchemes(schemes []string) func(*V1Client) {
	return func(v *V1Client) {
		v.schemes = schemes
	}
}

func WithTransportDebug() func(*V1Client) {
	return func(v *V1Client) {
		v.transportDebug = true
	}
}

func (h *V1Client) Clone() *V1Client {
	opts := []func(*V1Client){
		WithHubbleURI(h.hubbleUri),
		WithInsecureSkipVerify(h.insecureSkipVerify),
		WithRetries(h.retryAttempts),
		WithSchemes(h.schemes),
		WithScopeTenant(),
		WithScopeSystem(h.username, h.password),
	}
	if h.apikey != "" {
		opts = append(opts, WithAPIKey(h.apikey))
	}
	if h.jwt != "" {
		opts = append(opts, WithJWT(h.jwt))
	}
	if h.username != "" && h.password != "" {
		if h.systemScoped {
			opts = append(opts, WithScopeSystem(h.username, h.password))
		} else {
			opts = append(opts, WithUsername(h.username), WithPassword(h.password))
		}
	}
	if h.projectUid != "" {
		opts = append(opts, WithScopeProject(h.projectUid))
	}
	if h.transportDebug {
		opts = append(opts, WithTransportDebug())
	}
	return New(opts...)
}

func (h *V1Client) getTransport() (t *transport.Runtime) {
	if h.username != "" && h.password != "" {
		if err := h.authenticate(); err != nil {
			return nil
		}
	}
	if h.apikey != "" {
		t = h.apiKeyTransport()
	} else if h.jwt != "" {
		t = h.jwtTransport()
	} else {
		t = h.baseTransport()
	}
	return
}

func (h *V1Client) apiKeyTransport() *transport.Runtime {
	httpTransport := h.baseTransport()
	httpTransport.DefaultAuthentication = openapiclient.APIKeyAuth(authApiKey, authTokenInput, h.apikey)
	return httpTransport
}

func (h *V1Client) jwtTransport() *transport.Runtime {
	httpTransport := h.baseTransport()
	httpTransport.DefaultAuthentication = openapiclient.APIKeyAuth(authJwt, authTokenInput, h.jwt)
	return httpTransport
}

func (h *V1Client) authenticate() error {
	httpTransport := h.baseTransport()
	authClient := authC.New(httpTransport, strfmt.Default)

	if h.systemScoped {
		params := &authC.V1SysLoginParams{
			Body: &models.V1SysLogin{
				Username: h.username,
				Password: strfmt.Password(h.password),
			},
		}
		resp, err := authClient.V1SysLogin(params)
		if err != nil {
			return err
		}
		h.jwt = resp.Payload.Authorization
	} else {
		params := &authC.V1AuthenticateParams{
			Body: &models.V1AuthLogin{
				EmailID:  h.username,
				Password: strfmt.Password(h.password),
			},
		}
		resp, err := authClient.V1Authenticate(params)
		if err != nil {
			return err
		}
		h.jwt = resp.Payload.Authorization
	}

	return nil
}

func (h *V1Client) baseTransport() *transport.Runtime {
	httpTransport := transport.NewWithClient(h.hubbleUri, "", h.schemes, h.httpClient())
	httpTransport.RetryAttempts = h.retryAttempts
	httpTransport.Debug = h.transportDebug
	return httpTransport
}

func (h *V1Client) httpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{},
				InsecureSkipVerify: h.insecureSkipVerify, // #nosec G402
			},
		},
	}
}

func (h *V1Client) Validate() error {
	// API key can only be validated by making an API call
	if h.apikey != "" {
		_, err := h.UserC.V1ProjectsList(nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *V1Client) ValidateTenantAdmin() error {
	// API key can only be validated by making an API call
	if h.apikey != "" {
		_, err := h.UserC.V1UsersList(nil)
		if err != nil {
			return err
		}
	}
	return nil
}
