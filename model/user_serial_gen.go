package model

import "github.com/tinylib/msgp/msgp"

// UnmarshalMsg implements msgp.Unmarshaler
func (z *User) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 32 {
		err = msgp.ArrayError{Wanted: 32, Got: zb0001}
		return
	}
	z.Id, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Id")
		return
	}
	z.CreateAt, bts, err = msgp.ReadInt64Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "CreateAt")
		return
	}
	z.UpdateAt, bts, err = msgp.ReadInt64Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "UpdateAt")
		return
	}
	z.DeleteAt, bts, err = msgp.ReadInt64Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "DeleteAt")
		return
	}
	z.Username, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Username")
		return
	}
	z.Password, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Password")
		return
	}
	if msgp.IsNil(bts) {
		bts, err = msgp.ReadNilBytes(bts)
		if err != nil {
			return
		}
		z.AuthData = nil
	} else {
		if z.AuthData == nil {
			z.AuthData = new(string)
		}
		*z.AuthData, bts, err = msgp.ReadStringBytes(bts)
		if err != nil {
			err = msgp.WrapError(err, "AuthData")
			return
		}
	}
	z.AuthService, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "AuthService")
		return
	}
	z.Email, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Email")
		return
	}
	z.EmailVerified, bts, err = msgp.ReadBoolBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "EmailVerified")
		return
	}
	z.Nickname, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Nickname")
		return
	}
	z.FirstName, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "FirstName")
		return
	}
	z.LastName, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "LastName")
		return
	}
	// z.Position, bts, err = msgp.ReadStringBytes(bts)
	// if err != nil {
	// 	err = msgp.WrapError(err, "Position")
	// 	return
	// }
	z.Roles, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Roles")
		return
	}
	// z.AllowMarketing, bts, err = msgp.ReadBoolBytes(bts)
	// if err != nil {
	// 	err = msgp.WrapError(err, "AllowMarketing")
	// 	return
	// }
	bts, err = z.Props.UnmarshalMsg(bts)
	if err != nil {
		err = msgp.WrapError(err, "Props")
		return
	}
	bts, err = z.NotifyProps.UnmarshalMsg(bts)
	if err != nil {
		err = msgp.WrapError(err, "NotifyProps")
		return
	}
	z.LastPasswordUpdate, bts, err = msgp.ReadInt64Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "LastPasswordUpdate")
		return
	}
	z.LastPictureUpdate, bts, err = msgp.ReadInt64Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "LastPictureUpdate")
		return
	}
	z.FailedAttempts, bts, err = msgp.ReadIntBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "FailedAttempts")
		return
	}
	z.Locale, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Locale")
		return
	}
	bts, err = z.Timezone.UnmarshalMsg(bts)
	if err != nil {
		err = msgp.WrapError(err, "Timezone")
		return
	}
	z.MfaActive, bts, err = msgp.ReadBoolBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "MfaActive")
		return
	}
	z.MfaSecret, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "MfaSecret")
		return
	}
	// if msgp.IsNil(bts) {
	// 	bts, err = msgp.ReadNilBytes(bts)
	// 	if err != nil {
	// 		return
	// 	}
	// 	z.RemoteId = nil
	// } else {
	// 	if z.RemoteId == nil {
	// 		z.RemoteId = new(string)
	// 	}
	// 	*z.RemoteId, bts, err = msgp.ReadStringBytes(bts)
	// 	if err != nil {
	// 		err = msgp.WrapError(err, "RemoteId")
	// 		return
	// 	}
	// }
	z.LastActivityAt, bts, err = msgp.ReadInt64Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "LastActivityAt")
		return
	}
	z.TermsOfServiceId, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "TermsOfServiceId")
		return
	}
	z.TermsOfServiceCreateAt, bts, err = msgp.ReadInt64Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "TermsOfServiceCreateAt")
		return
	}
	o = bts
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

// UnmarshalMsg implements msgp.Unmarshaler
func (z *UserMap) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0003 uint32
	zb0003, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if (*z) == nil {
		(*z) = make(UserMap, zb0003)
	} else if len((*z)) > 0 {
		for key := range *z {
			delete((*z), key)
		}
	}
	for zb0003 > 0 {
		var zb0001 string
		var zb0002 *User
		zb0003--
		zb0001, bts, err = msgp.ReadStringBytes(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		if msgp.IsNil(bts) {
			bts, err = msgp.ReadNilBytes(bts)
			if err != nil {
				return
			}
			zb0002 = nil
		} else {
			if zb0002 == nil {
				zb0002 = new(User)
			}
			bts, err = zb0002.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, zb0001)
				return
			}
		}
		(*z)[zb0001] = zb0002
	}
	o = bts
	return
}
