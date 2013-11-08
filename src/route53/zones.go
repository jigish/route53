package route53

import (
	"encoding/xml"
	"errors"
	"fmt"
)

// XML RPC types.

type HostedZone struct {
	r53                    *Route53 `xml:"-"`
	Id                     string
	Name                   string
	CallerReference        string
	Comment                string `xml:"Config>Comment"`
	ResourceRecordSetCount int
}

type CreateHostedZoneRequest struct {
	XMLName         xml.Name `xml:"CreateHostedZoneRequest"`
	Name            string
	CallerReference string
	Comment         string `xml:"HostedZoneConfig>Comment"`
}

type CreateHostedZoneResponse struct {
	XMLName     xml.Name `xml:"CreateHostedZoneResponse"`
	HostedZone  HostedZone
	ChangeInfo  ChangeInfo
	NameServers []string `xml:"DelegationSet>NameServers>NameServer"`
}

type GetHostedZoneResponse struct {
	XMLName     xml.Name `xml:"GetHostedZoneResponse"`
	HostedZone  HostedZone
	NameServers []string `xml:"DelegationSet>NameServers>NameServer"`
}

type ListHostedZonesResponse struct {
	XMLName     xml.Name `xml:"ListHostedZonesResponse"`
	HostedZones []HostedZone
	IsTruncated bool
	Marker      string
	NextMarker  string
	MaxItems    uint
}

type DeleteHostedZoneResponse struct {
	ChangeInfo ChangeInfo
}

// Route53 API requests.

func (r53 *Route53) CreateHostedZone(name, reference, comment string) (ChangeInfo, error) {
	xmlReq := &CreateHostedZoneRequest{
		Name:            name,
		CallerReference: reference,
		Comment:         comment,
	}

	req := request{
		method: "POST",
		path:   "/2012-12-12/hostedzone",
		body:   xmlReq,
	}

	xmlRes := &CreateHostedZoneResponse{}

	if err := r53.run(req, xmlRes); err != nil {
		return ChangeInfo{}, err
	}

	return xmlRes.ChangeInfo, nil
}

func (r53 *Route53) GetHostedZone(id string) (HostedZone, error) {
	req := request{
		method: "GET",
		path:   fmt.Sprintf("/2012-12-12/hostedzone/%s", id),
	}

	xmlRes := &GetHostedZoneResponse{}

	if err := r53.run(req, xmlRes); err != nil {
		return HostedZone{}, err
	}

	xmlRes.HostedZone.r53 = r53

	return xmlRes.HostedZone, nil
}

func (r53 *Route53) ListHostedZones() ([]HostedZone, error) {
	req := request{
		method: "GET",
		path:   "/2012-12-12/hostedzone",
	}

	xmlRes := &ListHostedZonesResponse{}

	if err := r53.run(req, xmlRes); err != nil {
		return []HostedZone{}, err
	}
	if xmlRes.IsTruncated {
		return []HostedZone{}, errors.New("cannot handle truncated responses")
	}

	for _, zone := range xmlRes.HostedZones {
		zone.r53 = r53
	}

	return xmlRes.HostedZones, nil
}

func (r53 *Route53) DeleteHostedZone(id string) error {
	req := request{
		method: "DELETE",
		path:   fmt.Sprintf("/2012-12-12/hostedzone/%s", id),
	}

	xmlRes := &DeleteHostedZoneResponse{}

	if err := r53.run(req, xmlRes); err != nil {
		return err
	}

	return nil
}