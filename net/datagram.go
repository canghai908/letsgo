/*=============================================================================
#     FileName: datagram.go
#         Desc: Datagram pack/unpack
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-06 14:35:37
#      History:
=============================================================================*/
package net

import (
    "encoding/binary"
)

//type DatagramInterface interface {
//    Fetch(data []byte) int,[]byte
//    Pack(int,[]byte) []byte
//}

const (
    mask1 = 0x59
    mask2 = 0x7a

    DATAPACKET_TYPE_GENERAL = 0
    DATAPACKET_TYPE_DELAY = 1
    DATAPACKET_TYPE_BOARDCAST = 3
)

type Datagram struct { }

//对数据进行拆包
func (d *Datagram) Fetch(c *Transport) (n int, dps []*DataPacket) {
    dps = []*DataPacket{}

    ilen := len(c.Buff)
    if ilen == 0 {
        return
    }

    pos := 0
    dataType := 0
    dpSize := 0

    for {
        //拆包
        if c.DPSize > 0 {
            if ilen-pos < c.DPSize {
                //如果缓存去数据长度不够就退出接着等后续数据
                return
            }
        } else {
            if ilen-pos < 7 {
                return
            }

            if c.Buff[pos] == mask1 && c.Buff[pos+1] == mask2 {
                dataType = int(c.Buff[pos+2])

                dpSize = int(binary.BigEndian.Uint32(c.Buff[pos+3 : pos+7]))

                if ilen < dpSize+7 {
                    //如果缓存去数据长度不够就退出接着等后续数据
                    c.Buff = c.Buff[7:]

                    c.DPSize = dpSize
                    c.DataType = dataType

                    return
                }

                pos += 7

            } else {
                //如果错位则将缓存数据抛弃
                c.InitBuff()
                return
            }
        }

        dp := &DataPacket{Type: dataType, Data: c.Buff[pos : pos+dpSize]}
        dps = append(dps,dp)
        n += 1

        c.DPSize = 0
        c.DataType = 0

        //send to channel for consume
        //c.ProcessDP(dp)

        if ilen > dpSize+7 {
            //c.Buff = c.Buff[5+dpSize:]
            pos += dpSize
            continue

        } else {
            c.InitBuff()
            return
        }
    }
    return
}

//对数据进行封包
func (d *Datagram) Pack(dp *DataPacket) []byte {
    ilen := len(dp.Data)
    buf := make([]byte, ilen+7)

    buf[0] = byte(mask1)
    buf[1] = byte(mask2)
    buf[2] = byte(dp.Type)
    binary.BigEndian.PutUint32(buf[3:], uint32(ilen))

    copy(buf[7:], dp.Data)
    return buf
}