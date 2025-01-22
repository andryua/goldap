package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-ldap/ldap"
)

type LDAPClient struct {
	conn *ldap.Conn
}

type cfg struct {
	Host string
	Port int
}

var DefaultAttrs = []string{"*"} //"uid", "uidNumber", "cn", "mail", "title", "ou"}

func New() (*LDAPClient, error) {
	l := &LDAPClient{}
	var conf cfg
	conf.Host = "ldap.server.com"
	conf.Port = 636
	log.Println("connecting")
	err := l.connect(conf.Host, conf.Port)
	if err != nil {
		return nil, err
	}
	log.Println("connected")
	return l, err
}

func (c *LDAPClient) buildDNFromSerial(userName string) string {
	return "ou=...,o=..."
}

func (c *LDAPClient) connect(host string, port int) error {
	var err error
	c.conn, err = ldap.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	return err
}

func (c *LDAPClient) search(dn string, userName string, attrs []string) (*ldap.SearchResult, error) {
	searchRequest := ldap.NewSearchRequest(
		dn, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(uid=%s)", strings.ToUpper(userName)), // The filter to apply
		attrs,
		nil,
	)
	return c.conn.Search(searchRequest)
}

func (c *LDAPClient) Find(userName string, attrs []string) (map[string][]string, error) {
	baseDn := c.buildDNFromSerial(userName)
	log.Println("start", baseDn)
	sr, err := c.search(baseDn, userName, attrs)
	if err != nil {
		return nil, err
	}
	log.Println("end")
	if len(sr.Entries) == 0 {
		return nil, errors.New("empty search results")
	}
	result := make(map[string][]string)
	for _, entry := range sr.Entries[0].Attributes {
		result[entry.Name] = entry.Values
	}
	return result, nil
}

func (c *LDAPClient) Close() {
	c.conn.Close()
}

func main() {
	c, err := New()
	if err != nil {
		log.Println("connect ", err)
	}
	result, err := c.Find("uid", DefaultAttrs)
	if err != nil {
		log.Println("find ", err)
	}
	for key, val := range result {
		fmt.Println(key, val)
	}
}
