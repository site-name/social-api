package model

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Session) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Id":
			z.Id, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Id")
				return
			}
		case "Token":
			z.Token, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Token")
				return
			}
		case "CreateAt":
			z.CreateAt, err = dc.ReadInt64()
			if err != nil {
				err = msgp.WrapError(err, "CreateAt")
				return
			}
		case "ExpiresAt":
			z.ExpiresAt, err = dc.ReadInt64()
			if err != nil {
				err = msgp.WrapError(err, "ExpiresAt")
				return
			}
		case "LastActivityAt":
			z.LastActivityAt, err = dc.ReadInt64()
			if err != nil {
				err = msgp.WrapError(err, "LastActivityAt")
				return
			}
		case "UserId":
			z.UserId, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "UserId")
				return
			}
		case "DeviceId":
			z.DeviceId, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "DeviceId")
				return
			}
		case "Roles":
			z.Roles, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Roles")
				return
			}
		case "IsOAuth":
			z.IsOAuth, err = dc.ReadBool()
			if err != nil {
				err = msgp.WrapError(err, "IsOAuth")
				return
			}
		case "ExpiredNotify":
			z.ExpiredNotify, err = dc.ReadBool()
			if err != nil {
				err = msgp.WrapError(err, "ExpiredNotify")
				return
			}
		case "Props":
			var zb0002 uint32
			zb0002, err = dc.ReadMapHeader()
			if err != nil {
				err = msgp.WrapError(err, "Props")
				return
			}
			if z.Props == nil {
				z.Props = make(StringMap, zb0002)
			} else if len(z.Props) > 0 {
				for key := range z.Props {
					delete(z.Props, key)
				}
			}
			for zb0002 > 0 {
				zb0002--
				var za0001 string
				var za0002 string
				za0001, err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "Props")
					return
				}
				za0002, err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "Props", za0001)
					return
				}
				z.Props[za0001] = za0002
			}
		case "Local":
			z.Local, err = dc.ReadBool()
			if err != nil {
				err = msgp.WrapError(err, "Local")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Session) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 12
	// write "Id"
	err = en.Append(0x8c, 0xa2, 0x49, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.Id)
	if err != nil {
		err = msgp.WrapError(err, "Id")
		return
	}
	// write "Token"
	err = en.Append(0xa5, 0x54, 0x6f, 0x6b, 0x65, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteString(z.Token)
	if err != nil {
		err = msgp.WrapError(err, "Token")
		return
	}
	// write "CreateAt"
	err = en.Append(0xa8, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x74)
	if err != nil {
		return
	}
	err = en.WriteInt64(z.CreateAt)
	if err != nil {
		err = msgp.WrapError(err, "CreateAt")
		return
	}
	// write "ExpiresAt"
	err = en.Append(0xa9, 0x45, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x41, 0x74)
	if err != nil {
		return
	}
	err = en.WriteInt64(z.ExpiresAt)
	if err != nil {
		err = msgp.WrapError(err, "ExpiresAt")
		return
	}
	// write "LastActivityAt"
	err = en.Append(0xae, 0x4c, 0x61, 0x73, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x41, 0x74)
	if err != nil {
		return
	}
	err = en.WriteInt64(z.LastActivityAt)
	if err != nil {
		err = msgp.WrapError(err, "LastActivityAt")
		return
	}
	// write "UserId"
	err = en.Append(0xa6, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.UserId)
	if err != nil {
		err = msgp.WrapError(err, "UserId")
		return
	}
	// write "DeviceId"
	err = en.Append(0xa8, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.DeviceId)
	if err != nil {
		err = msgp.WrapError(err, "DeviceId")
		return
	}
	// write "Roles"
	err = en.Append(0xa5, 0x52, 0x6f, 0x6c, 0x65, 0x73)
	if err != nil {
		return
	}
	err = en.WriteString(z.Roles)
	if err != nil {
		err = msgp.WrapError(err, "Roles")
		return
	}
	// write "IsOAuth"
	err = en.Append(0xa7, 0x49, 0x73, 0x4f, 0x41, 0x75, 0x74, 0x68)
	if err != nil {
		return
	}
	err = en.WriteBool(z.IsOAuth)
	if err != nil {
		err = msgp.WrapError(err, "IsOAuth")
		return
	}
	// write "ExpiredNotify"
	err = en.Append(0xad, 0x45, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79)
	if err != nil {
		return
	}
	err = en.WriteBool(z.ExpiredNotify)
	if err != nil {
		err = msgp.WrapError(err, "ExpiredNotify")
		return
	}
	// write "Props"
	err = en.Append(0xa5, 0x50, 0x72, 0x6f, 0x70, 0x73)
	if err != nil {
		return
	}
	err = en.WriteMapHeader(uint32(len(z.Props)))
	if err != nil {
		err = msgp.WrapError(err, "Props")
		return
	}
	for za0001, za0002 := range z.Props {
		err = en.WriteString(za0001)
		if err != nil {
			err = msgp.WrapError(err, "Props")
			return
		}
		err = en.WriteString(za0002)
		if err != nil {
			err = msgp.WrapError(err, "Props", za0001)
			return
		}
	}
	// write "Local"
	err = en.Append(0xa5, 0x4c, 0x6f, 0x63, 0x61, 0x6c)
	if err != nil {
		return
	}
	err = en.WriteBool(z.Local)
	if err != nil {
		err = msgp.WrapError(err, "Local")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Session) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 12
	// string "Id"
	o = append(o, 0x8c, 0xa2, 0x49, 0x64)
	o = msgp.AppendString(o, z.Id)
	// string "Token"
	o = append(o, 0xa5, 0x54, 0x6f, 0x6b, 0x65, 0x6e)
	o = msgp.AppendString(o, z.Token)
	// string "CreateAt"
	o = append(o, 0xa8, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x74)
	o = msgp.AppendInt64(o, z.CreateAt)
	// string "ExpiresAt"
	o = append(o, 0xa9, 0x45, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x41, 0x74)
	o = msgp.AppendInt64(o, z.ExpiresAt)
	// string "LastActivityAt"
	o = append(o, 0xae, 0x4c, 0x61, 0x73, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x41, 0x74)
	o = msgp.AppendInt64(o, z.LastActivityAt)
	// string "UserId"
	o = append(o, 0xa6, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64)
	o = msgp.AppendString(o, z.UserId)
	// string "DeviceId"
	o = append(o, 0xa8, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64)
	o = msgp.AppendString(o, z.DeviceId)
	// string "Roles"
	o = append(o, 0xa5, 0x52, 0x6f, 0x6c, 0x65, 0x73)
	o = msgp.AppendString(o, z.Roles)
	// string "IsOAuth"
	o = append(o, 0xa7, 0x49, 0x73, 0x4f, 0x41, 0x75, 0x74, 0x68)
	o = msgp.AppendBool(o, z.IsOAuth)
	// string "ExpiredNotify"
	o = append(o, 0xad, 0x45, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79)
	o = msgp.AppendBool(o, z.ExpiredNotify)
	// string "Props"
	o = append(o, 0xa5, 0x50, 0x72, 0x6f, 0x70, 0x73)
	o = msgp.AppendMapHeader(o, uint32(len(z.Props)))
	for za0001, za0002 := range z.Props {
		o = msgp.AppendString(o, za0001)
		o = msgp.AppendString(o, za0002)
	}
	// string "Local"
	o = append(o, 0xa5, 0x4c, 0x6f, 0x63, 0x61, 0x6c)
	o = msgp.AppendBool(o, z.Local)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Session) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Id":
			z.Id, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Id")
				return
			}
		case "Token":
			z.Token, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Token")
				return
			}
		case "CreateAt":
			z.CreateAt, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "CreateAt")
				return
			}
		case "ExpiresAt":
			z.ExpiresAt, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ExpiresAt")
				return
			}
		case "LastActivityAt":
			z.LastActivityAt, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "LastActivityAt")
				return
			}
		case "UserId":
			z.UserId, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "UserId")
				return
			}
		case "DeviceId":
			z.DeviceId, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "DeviceId")
				return
			}
		case "Roles":
			z.Roles, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Roles")
				return
			}
		case "IsOAuth":
			z.IsOAuth, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "IsOAuth")
				return
			}
		case "ExpiredNotify":
			z.ExpiredNotify, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ExpiredNotify")
				return
			}
		case "Props":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Props")
				return
			}
			if z.Props == nil {
				z.Props = make(StringMap, zb0002)
			} else if len(z.Props) > 0 {
				for key := range z.Props {
					delete(z.Props, key)
				}
			}
			for zb0002 > 0 {
				var za0001 string
				var za0002 string
				zb0002--
				za0001, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Props")
					return
				}
				za0002, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Props", za0001)
					return
				}
				z.Props[za0001] = za0002
			}
		case "Local":
			z.Local, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Local")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Session) Msgsize() (s int) {
	s = 1 + 3 + msgp.StringPrefixSize + len(z.Id) + 6 + msgp.StringPrefixSize + len(z.Token) + 9 + msgp.Int64Size + 10 + msgp.Int64Size + 15 + msgp.Int64Size + 7 + msgp.StringPrefixSize + len(z.UserId) + 9 + msgp.StringPrefixSize + len(z.DeviceId) + 6 + msgp.StringPrefixSize + len(z.Roles) + 8 + msgp.BoolSize + 14 + msgp.BoolSize + 6 + msgp.MapHeaderSize
	if z.Props != nil {
		for za0001, za0002 := range z.Props {
			_ = za0002
			s += msgp.StringPrefixSize + len(za0001) + msgp.StringPrefixSize + len(za0002)
		}
	}
	s += 6 + msgp.BoolSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *StringMap) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0003 uint32
	zb0003, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if (*z) == nil {
		(*z) = make(StringMap, zb0003)
	} else if len((*z)) > 0 {
		for key := range *z {
			delete((*z), key)
		}
	}
	for zb0003 > 0 {
		zb0003--
		var zb0001 string
		var zb0002 string
		zb0001, err = dc.ReadString()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		zb0002, err = dc.ReadString()
		if err != nil {
			err = msgp.WrapError(err, zb0001)
			return
		}
		(*z)[zb0001] = zb0002
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z StringMap) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteMapHeader(uint32(len(z)))
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0004, zb0005 := range z {
		err = en.WriteString(zb0004)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		err = en.WriteString(zb0005)
		if err != nil {
			err = msgp.WrapError(err, zb0004)
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z StringMap) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendMapHeader(o, uint32(len(z)))
	for zb0004, zb0005 := range z {
		o = msgp.AppendString(o, zb0004)
		o = msgp.AppendString(o, zb0005)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *StringMap) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0003 uint32
	zb0003, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if (*z) == nil {
		(*z) = make(StringMap, zb0003)
	} else if len((*z)) > 0 {
		for key := range *z {
			delete((*z), key)
		}
	}
	for zb0003 > 0 {
		var zb0001 string
		var zb0002 string
		zb0003--
		zb0001, bts, err = msgp.ReadStringBytes(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		zb0002, bts, err = msgp.ReadStringBytes(bts)
		if err != nil {
			err = msgp.WrapError(err, zb0001)
			return
		}
		(*z)[zb0001] = zb0002
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z StringMap) Msgsize() (s int) {
	s = msgp.MapHeaderSize
	if z != nil {
		for zb0004, zb0005 := range z {
			_ = zb0005
			s += msgp.StringPrefixSize + len(zb0004) + msgp.StringPrefixSize + len(zb0005)
		}
	}
	return
}