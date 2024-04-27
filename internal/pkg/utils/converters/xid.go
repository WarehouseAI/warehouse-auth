package wh_converters

import "github.com/rs/xid"

func FastConvertToXid(s string) xid.ID {
	id, _ := xid.FromString(s)
	return id
}
