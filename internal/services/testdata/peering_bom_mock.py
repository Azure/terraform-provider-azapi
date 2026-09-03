"""mitmproxy addon that mocks the Microsoft.Peering/peerings ARM API.

It reproduces https://github.com/Azure/terraform-provider-azapi/issues/1098: the
real Azure Peering API prepends a UTF-8 byte-order mark (EF BB BF) to its JSON
response body. Go's encoding/json cannot parse a leading BOM, so azapi fails to
read the resource with:

    unmarshalling type *interface {}: invalid character 'ï' looking for beginning of value

A real peering cannot be provisioned in a test subscription (it needs a
registered ASN and a signed peering agreement), so this addon fully mocks the
peering resource lifecycle instead. It is stateful so the whole create/read/
import/delete flow works once the provider learns to strip the BOM:

    * PUT / PATCH -> 200 with a clean body, and marks the peering as "created"
    * GET         -> 404 until created, then 200 with a BOM-prefixed body (the bug)
    * DELETE      -> 200, and marks the peering as "deleted"

Only requests to the peering resource itself are intercepted; every other
request (AAD auth, resource provider registration, ...) is passed through to the
real endpoint.

Usage:
    mitmdump -p 7070 -s peering_bom_mock.py
"""

import json
import re

from mitmproxy import http

# UTF-8 byte-order mark that the Azure Peering API prepends to its JSON body.
UTF8_BOM = b"\xef\xbb\xbf"

PEERING_RE = re.compile(
    r"^/subscriptions/[^/]+/resourceGroups/[^/]+"
    r"/providers/Microsoft\.Peering/peerings/[^/?]+$",
    re.IGNORECASE,
)

_created = set()


def _peering_body(resource_id, name):
    """A realistic exchange-peering body, modelled on the 2025-05-01 API spec."""
    return {
        "name": name,
        "type": "Microsoft.Peering/peerings",
        "id": resource_id,
        "kind": "Exchange",
        "location": "eastus",
        "properties": {
            "peeringLocation": "peeringLocation0",
            "provisioningState": "Succeeded",
            "exchange": {
                "peerAsn": {
                    "id": "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Peering/peerAsns/myAsn1"
                },
                "connections": [
                    {
                        "connectionIdentifier": "CE495334-0E94-4E51-8164-8116D6CD284D",
                        "peeringDBFacilityId": 99999,
                        "connectionState": "Active",
                        "bgpSession": {
                            "maxPrefixesAdvertisedV4": 1000,
                            "maxPrefixesAdvertisedV6": 100,
                            "peerSessionIPv4Address": "192.168.2.1",
                            "peerSessionIPv6Address": "fd00::1",
                        },
                    }
                ],
            },
        },
        "sku": {
            "name": "Basic_Exchange_Free",
            "family": "Exchange",
            "size": "Free",
            "tier": "Basic",
        },
    }


def _json_response(status, body_bytes):
    return http.Response.make(
        status, body_bytes, {"Content-Type": "application/json"}
    )


def request(flow: http.HTTPFlow) -> None:
    host = flow.request.pretty_host
    path = flow.request.path.split("?", 1)[0]

    if "management.azure.com" not in host or not PEERING_RE.match(path):
        return 

    name = path.rstrip("/").rsplit("/", 1)[-1]
    body = json.dumps(_peering_body(path, name)).encode("utf-8")
    method = flow.request.method.upper()

    if method == "GET":
        if path in _created:
            flow.response = _json_response(200, UTF8_BOM + body)
        else:
            not_found = json.dumps(
                {
                    "error": {
                        "code": "NotFound",
                        "message": "The peering resource was not found.",
                    }
                }
            ).encode("utf-8")
            flow.response = _json_response(404, not_found)
    elif method in ("PUT", "PATCH"):
        _created.add(path)
        flow.response = _json_response(200, body)
    elif method == "DELETE":
        _created.discard(path)
        flow.response = _json_response(
            200, json.dumps({"status": "Succeeded"}).encode("utf-8")
        )
    else:
        flow.response = _json_response(200, body)
