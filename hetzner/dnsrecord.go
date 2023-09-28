package hetzner

type DnsRecordsResponse struct {
	Records []DnsRecord
}

type DnsRecordResponse struct {
	Record DnsRecord
}

type DnsRecord struct {
	Id    string
	Name  string
	Value string
	Type  string
}
