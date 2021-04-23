package model

const (
	HEADER_REQUEST_ID          = "X-Request-ID"
	HEADER_VERSION_ID          = "X-Version-ID"
	HEADER_CLUSTER_ID          = "X-Cluster-ID"
	HEADER_ETAG_SERVER         = "ETag"
	HEADER_ETAG_CLIENT         = "If-None-Match"
	HEADER_FORWARDED           = "X-Forwarded-For"
	HEADER_REAL_IP             = "X-Real-IP"
	HEADER_FORWARDED_PROTO     = "X-Forwarded-Proto"
	HEADER_TOKEN               = "token"
	HEADER_CSRF_TOKEN          = "X-CSRF-Token"
	HEADER_BEARER              = "BEARER"
	HEADER_AUTH                = "Authorization"
	HEADER_CLOUD_TOKEN         = "X-Cloud-Token"
	HEADER_REMOTECLUSTER_TOKEN = "X-RemoteCluster-Token"
	HEADER_REMOTECLUSTER_ID    = "X-RemoteCluster-Id"
	HEADER_REQUESTED_WITH      = "X-Requested-With"
	HEADER_REQUESTED_WITH_XML  = "XMLHttpRequest"
	HEADER_RANGE               = "Range"
	STATUS                     = "status"
	STATUS_OK                  = "OK"
	STATUS_FAIL                = "FAIL"
	STATUS_UNHEALTHY           = "UNHEALTHY"
	STATUS_REMOVE              = "REMOVE"

	CLIENT_DIR = "client"
)
