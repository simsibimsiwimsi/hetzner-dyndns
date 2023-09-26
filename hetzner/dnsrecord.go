package hetzner

type DnsRecordResponse struct {
	Records []DnsRecord
}

type DnsRecord struct {
	Id    string
	Name  string
	Value string
	Type  string
}
