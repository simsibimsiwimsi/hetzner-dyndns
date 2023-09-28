package hetzner

import (
	Bytes "bytes"
	Json "encoding/json"
	Errors "errors"
	Format "fmt"
	IO "io"
	Http "net/http"
)

type dnsZone struct {
	Id    string
	Token string
}

func NewDnsZone(Id string, Token string) dnsZone {
	return dnsZone{
		Id:    Id,
		Token: Token,
	}
}

func (self *dnsZone) doRequest(request *Http.Request, response any) (any, error) {

	// Create client
	client := &Http.Client{}

	// Apply request headers
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Auth-API-Token", self.Token)

	// Fetch Request
	httpResponse, err := client.Do(request)
	if err != nil {
		Format.Println("Failure : ", err)
		return nil, err
	}

	// Read Response Body
	httpResponseBody, _ := IO.ReadAll(httpResponse.Body)
	if httpResponse.StatusCode != 200 {
		return nil, Errors.New(string(httpResponseBody))
	}

	// Display Results
	Format.Println("response Status : ", httpResponse.Status)
	Format.Println("response Headers : ", httpResponse.Header)
	Format.Println("response Body : ", string(httpResponseBody))

	err = Json.Unmarshal(httpResponseBody, response)
	if err != nil {
		Format.Println(err)
		return nil, err
	}
	if response == nil {
		return nil, Errors.New("Unexpected response did not unmarshal correctly. Response was " + string(httpResponseBody))
	}
	return response, nil
}

func (self *dnsZone) GetRecordByName(recordName string, recordType string) (*DnsRecord, error) {

	// Get Records (GET https://dns.hetzner.com/api/v1/records?zone_id={ZoneID})
	req, err := Http.NewRequest("GET", "https://dns.hetzner.com/api/v1/records?zone_id="+self.Id, nil)
	if err != nil {
		return nil, err
	}

	response, err := self.doRequest(req, &DnsRecordsResponse{})
	if err != nil {
		Format.Println("Request failed : ", err)
		return nil, err
	}
	records := response.(*DnsRecordsResponse)

	// find an existing record, if exists
	for _, v := range records.Records {
		if v.Name == recordName && v.Type == recordType {
			return &v, nil
		}
	}

	return nil, Errors.New("Unable to find Record " + recordName + " in DNS Zone " + self.Id)
}

func (self *dnsZone) UpdateRecord(record *DnsRecord, ip string) (*DnsRecord, error) {

	// Update Record (PUT https://dns.hetzner.com/api/v1/records/{RecordID})
	body := []byte(`{"value": "` + ip + `","ttl": 60,"type": "` + record.Type + `","name": "` + record.Name + `","zone_id": "` + self.Id + `"}`)
	req, err := Http.NewRequest("PUT", "https://dns.hetzner.com/api/v1/records/"+record.Id, Bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	response, err := self.doRequest(req, &DnsRecordResponse{})
	if err != nil {
		Format.Println("Request failed : ", err)
		return nil, err
	}
	dnsRecordResponse := response.(*DnsRecordResponse)

	return &dnsRecordResponse.Record, nil
}

func (self *dnsZone) CreateRecord(record *DnsRecord) (*DnsRecord, error) {

	// Create Record (POST https://dns.hetzner.com/api/v1/records)
	body := []byte(`{"value": "` + record.Value + `","ttl": 60,"type": "` + record.Type + `","name": "` + record.Name + `","zone_id": "` + self.Id + `"}`)
	req, err := Http.NewRequest("POST", "https://dns.hetzner.com/api/v1/records", Bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	response, err := self.doRequest(req, &DnsRecordResponse{})
	if err != nil {
		Format.Println("Request failed : ", err)
		return nil, err
	}
	dnsRecordResponse := response.(*DnsRecordResponse)

	return &dnsRecordResponse.Record, nil
}

func (self *dnsZone) CreateOrUpdateRecord(name string, ip string, recordType string) (*DnsRecord, error) {

	record, err := self.GetRecordByName(name, recordType)
	if err != nil {
		record, err = self.CreateRecord(&DnsRecord{
			Name:  name,
			Value: ip,
			Type:  recordType,
		})
		if err != nil {
			return nil, err
		}
	} else {
		record, err = self.UpdateRecord(record, ip)
		if err != nil {
			return nil, err
		}
	}
	return record, nil
}

func (self *dnsZone) CreateOrUpdateIpV4andV6Records(name string, ipv4 string, ipv6 string) (*DnsRecord, *DnsRecord, error) {

	ipv4record, err := self.CreateOrUpdateRecord(name, ipv4, "A")
	if err != nil {
		return nil, nil, err
	}

	ipv6record, err := self.CreateOrUpdateRecord(name, ipv6, "AAAA")
	if err != nil {
		return ipv4record, nil, err
	}

	return ipv4record, ipv6record, nil
}
