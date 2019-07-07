// +build !go1.12

package httpclient

//Close Client release connection resource
func (c *Client) Close() {}
