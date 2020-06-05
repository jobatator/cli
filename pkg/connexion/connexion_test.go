package connexion

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestUrlParsing(t *testing.T) {
	url := "jobatator-instance.organization.org"
	o := ParseURL(url)
	assert.Equal(t, o.URL, url)
	assert.Equal(t, o.Port, "8962")
	assert.Equal(t, o.Host, "jobatator-instance.organization.org")
	assert.Equal(t, o.Username, "")
	assert.Equal(t, o.Password, "")
	assert.Equal(t, o.Group, "")

	o = ParseURL("jobatator-instance.organization.org:1212")
	assert.Equal(t, o.Port, "1212")
	assert.Equal(t, o.Host, "jobatator-instance.organization.org")
	assert.Equal(t, o.Username, "")
	assert.Equal(t, o.Password, "")
	assert.Equal(t, o.Group, "")

	o = ParseURL("host.com/group")
	assert.Equal(t, o.Port, "8962")
	assert.Equal(t, o.Host, "host.com")
	assert.Equal(t, o.Username, "")
	assert.Equal(t, o.Password, "")
	assert.Equal(t, o.Group, "group")

	o = ParseURL("root@localhost")
	assert.Equal(t, o.Port, "8962")
	assert.Equal(t, o.Host, "localhost")
	assert.Equal(t, o.Username, "root")
	assert.Equal(t, o.Password, "")
	assert.Equal(t, o.Group, "")

	o = ParseURL("th_ah-is-0username:and7-_pass@ho-st.com:12/gro-_up12")
	assert.Equal(t, o.Port, "12")
	assert.Equal(t, o.Host, "ho-st.com")
	assert.Equal(t, o.Username, "th_ah-is-0username")
	assert.Equal(t, o.Password, "and7-_pass")
	assert.Equal(t, o.Group, "gro-_up12")

	o = ParseURL("username1234:verysecurepassword_verysecure_yeah@sub.hell.yay.com/customgroup")
	assert.Equal(t, o.Port, "8962")
	assert.Equal(t, o.Host, "sub.hell.yay.com")
	assert.Equal(t, o.Username, "username1234")
	assert.Equal(t, o.Password, "verysecurepassword_verysecure_yeah")
	assert.Equal(t, o.Group, "customgroup")

	o = ParseURL("auser@do-you-really.tech/customgroup")
	assert.Equal(t, o.Port, "8962")
	assert.Equal(t, o.Host, "do-you-really.tech")
	assert.Equal(t, o.Username, "auser")
	assert.Equal(t, o.Password, "")
	assert.Equal(t, o.Group, "customgroup")

	o = ParseURL("practical_usage@node03.cluster.someshittycompany.fr/my_really_great_group")
	assert.Equal(t, o.Port, "8962")
	assert.Equal(t, o.Host, "node03.cluster.someshittycompany.fr")
	assert.Equal(t, o.Username, "practical_usage")
	assert.Equal(t, o.Password, "")
	assert.Equal(t, o.Group, "my_really_great_group")
}
