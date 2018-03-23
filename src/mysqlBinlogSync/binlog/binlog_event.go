package binlog

import (
	"io"
	"fmt"
	"mysqlBinlogSync/util"
	"strings"
	"strconv"
	"unicode"
	log "github.com/sirupsen/logrus"
	"encoding/binary"
	"encoding/hex"
	"mysqlBinlogSync/comm"
	"time"
	"bytes"
	"github.com/shopspring/decimal"
)

type BinlogEvent struct {
	RawData []byte
	Header  *EventHeader
	Event
}

type Event interface {
	Write(w io.Writer)
	Read(data []byte) error
}

var (
	checksumVersionSplitMysql   []int = []int{5, 6, 1}
	checksumVersionProductMysql int   = (checksumVersionSplitMysql[0]*256+checksumVersionSplitMysql[1])*256 + checksumVersionSplitMysql[2]
)

// server version format X.Y.Zabc, a is not . or number
func splitServerVersion(server string) []int {
	seps := strings.Split(server, ".")
	if len(seps) < 3 {
		return []int{0, 0, 0}
	}

	x, _ := strconv.Atoi(seps[0])
	y, _ := strconv.Atoi(seps[1])

	index := 0
	for i, c := range seps[2] {
		if !unicode.IsNumber(c) {
			index = i
			break
		}
	}

	z, _ := strconv.Atoi(seps[2][0:index])

	return []int{x, y, z}
}

func calcVersionProduct(server string) int {
	versionSplit := splitServerVersion(server)
	return (versionSplit[0]*256+versionSplit[1])*256 + versionSplit[2]
}

//----------------------------------GenericEvent----------------------------------------------//
// 不解析的event统一用GenericEvent表示
type GenericEvent struct {
	Data []byte
}

func (e *GenericEvent) Write(w io.Writer) {
	fmt.Fprintf(w, "Event data: ; %s", hex.Dump(e.Data))
	fmt.Fprintln(w)
}

func (e *GenericEvent) Read(data []byte) error {
	e.Data = data
	return nil
}

//----------------------------------FormatDescriptionEvent----------------------------------------------//

// https://dev.mysql.com/doc/internals/en/format-description-event.html
type FormatDescriptionEvent struct {
	Version                uint16
	ServerVersion          []byte //len = 50
	CreateTimestamp        uint32
	EventHeaderLength      uint8
	EventTypeHeaderLengths []byte

	// 0 is off, 1 is for CRC32, 255 is undefined
	ChecksumAlgorithm byte
}

func (e *FormatDescriptionEvent) Read(data []byte) error {
	pos := 0

	util.DecodeUint16(data, pos, &e.Version)
	pos += 2

	e.ServerVersion = make([]byte, 50)
	copy(e.ServerVersion, data[pos:])
	pos += 50

	util.DecodeUint32(data, pos, &e.CreateTimestamp)
	pos += 4

	e.EventHeaderLength = data[pos]
	pos++

	if e.EventHeaderLength != byte(EventHeaderSize) {
		return fmt.Errorf("invalid event header length %d, must 19", e.EventHeaderLength)
	}

	// TODO checksum
	checksumProduct := checksumVersionProductMysql
	if calcVersionProduct(string(e.ServerVersion)) >= checksumProduct {
		// here, the last 5 bytes is 1 byte check sum alg type and 4 byte checksum if exists
		e.ChecksumAlgorithm = data[len(data)-5]
		e.EventTypeHeaderLengths = data[pos: len(data)-5]
	} else {
		e.ChecksumAlgorithm = BINLOG_CHECKSUM_ALG_UNDEF
		e.EventTypeHeaderLengths = data[pos:]
	}

	return nil
}

func (e *FormatDescriptionEvent) Write(w io.Writer) {
	fmt.Fprintf(w, "Version: %d; ", e.Version)
	fmt.Fprintf(w, "Server version: %s; ", e.ServerVersion)
	//fmt.Fprintf(w, "Create date: %s; ", time.Unix(int64(e.CreateTimestamp), 0).Format(TimeFormat))
	fmt.Fprintf(w, "Checksum algorithm: %d; ", e.ChecksumAlgorithm)
	//fmt.Fprintf(w, "Event header lengths: ; %s", hex.Dump(e.EventTypeHeaderLengths))
	fmt.Fprintln(w)
}

//----------------------------------RotateEvent----------------------------------------------//

//https://dev.mysql.com/doc/internals/en/rotate-event.html
type RotateEvent struct {
	Position    uint64
	NextLogName []byte
}

func (e *RotateEvent) Read(data []byte) error {
	util.DecodeUint64(data, 0, &e.Position)
	e.NextLogName = data[8:]
	return nil
}

func (e *RotateEvent) Write(w io.Writer) {
	fmt.Fprintf(w, "Position: %d; ", e.Position)
	fmt.Fprintf(w, "Next log name: %s; ", e.NextLogName)
	fmt.Fprintln(w)
}


//----------------------------------XIDEvent----------------------------------------------//
type XIDEvent struct {
	XID uint64
}

func (e *XIDEvent) Read(data []byte) error {
	e.XID = binary.LittleEndian.Uint64(data)
	return nil
}

func (e *XIDEvent) Write(w io.Writer) {
	fmt.Fprintf(w, "XID: %d ", e.XID)
	fmt.Fprintln(w)
}

//----------------------------------QueryEvent----------------------------------------------//
// https://dev.mysql.com/doc/internals/en/query-event.html
type QueryEvent struct {
	SlaveProxyID  uint32
	ExecutionTime uint32
	ErrorCode     uint16
	StatusVars    []byte
	Schema        []byte
	Query         []byte
}

func (e *QueryEvent) Read(data []byte) error {
	pos := 0

	util.DecodeUint32(data, pos, &e.SlaveProxyID)
	pos += 4

	util.DecodeUint32(data, pos, &e.ExecutionTime)
	pos += 4

	schemaLength := uint8(data[pos])
	pos++

	util.DecodeUint16(data, pos, &e.ErrorCode)
	pos += 2

	statusVarsLength := binary.LittleEndian.Uint16(data[pos:])
	pos += 2

	e.StatusVars = data[pos: pos+int(statusVarsLength)]
	pos += int(statusVarsLength)

	e.Schema = data[pos: pos+int(schemaLength)]
	pos += int(schemaLength)

	//skip 0x00
	pos++

	e.Query = data[pos:]
	return nil
}

func (e *QueryEvent) Write(w io.Writer) {
	fmt.Fprintf(w, "Slave proxy ID: %d; ", e.SlaveProxyID)
	fmt.Fprintf(w, "Execution time: %d; ", e.ExecutionTime)
	fmt.Fprintf(w, "Error code: %d; ", e.ErrorCode)
	//fmt.Fprintf(w, "Status vars: ; %s", hex.Dump(e.StatusVars))
	fmt.Fprintf(w, "Schema: %s; ", e.Schema)
	fmt.Fprintf(w, "Query: %s; ", e.Query)
	fmt.Fprintln(w)
}

//----------------------------------TableMapEvent----------------------------------------------//

// https://dev.mysql.com/doc/internals/en/table-map-event.html
type TableMapEvent struct {
	tableIDSize int
	TableID     uint64
	Flags       uint16
	Schema      []byte
	Table       []byte
	ColumnCount uint64
	ColumnType  []byte
	ColumnMeta  []uint16


	//len = (ColumnCount + 7) / 8   len=(column_count + 8) / 7 TODO ???
	NullBitmap []byte
}

func (e *TableMapEvent) Read(data []byte) error {
	pos := 0
	e.TableID = util.FixedLengthInt(data[0:e.tableIDSize])
	pos += e.tableIDSize

	util.DecodeUint16(data, pos, &e.Flags)
	pos += 2

	schemaLength := data[pos]
	pos++

	e.Schema = data[pos: pos+int(schemaLength)]
	pos += int(schemaLength)

	//skip 0x00
	pos++

	tableLength := data[pos]
	pos++

	e.Table = data[pos: pos+int(tableLength)]
	pos += int(tableLength)

	//skip 0x00
	pos++

	var n int
	e.ColumnCount, _, n = util.LengthEncodedInt(data[pos:])
	pos += n

	e.ColumnType = data[pos:pos+int(e.ColumnCount)]
	pos += int(e.ColumnCount)

	var err error
	var metaData []byte
	if metaData, _, n, err = util.LengthEncodedString(data[pos:]); err != nil {
		return err
	}

	if err = e.decodeMeta(metaData); err != nil {
		return err
	}

	pos += n

	if len(data[pos:]) != bitmapByteSize(int(e.ColumnCount)) {
		log.Error("-----left data size not equals to NULL-bitmask length")
		return io.EOF
	}

	e.NullBitmap = data[pos:]

	return nil
}

func bitmapByteSize(columnCount int) int {
	return int(columnCount+7) / 8
}

// see mysql sql/log_event.h  metaInfo 长度
/*
	0 byte
	MYSQL_TYPE_DECIMAL
	MYSQL_TYPE_TINY
	MYSQL_TYPE_SHORT
	MYSQL_TYPE_LONG
	MYSQL_TYPE_NULL
	MYSQL_TYPE_TIMESTAMP
	MYSQL_TYPE_LONGLONG
	MYSQL_TYPE_INT24
	MYSQL_TYPE_DATE
	MYSQL_TYPE_TIME
	MYSQL_TYPE_DATETIME
	MYSQL_TYPE_YEAR

	1 byte
	MYSQL_TYPE_FLOAT
	MYSQL_TYPE_DOUBLE
	MYSQL_TYPE_BLOB
	MYSQL_TYPE_GEOMETRY

	//maybe
	MYSQL_TYPE_TIME2
	MYSQL_TYPE_DATETIME2
	MYSQL_TYPE_TIMESTAMP2

	2 byte
	MYSQL_TYPE_VARCHAR
	MYSQL_TYPE_BIT
	MYSQL_TYPE_NEWDECIMAL
	MYSQL_TYPE_VAR_STRING
	MYSQL_TYPE_STRING

	This enumeration value is only used internally and cannot exist in a binlog.
	MYSQL_TYPE_NEWDATE
	MYSQL_TYPE_ENUM
	MYSQL_TYPE_SET
	MYSQL_TYPE_TINY_BLOB
	MYSQL_TYPE_MEDIUM_BLOB
	MYSQL_TYPE_LONG_BLOB
*/
// ColumnType 每列的类型，每一列用一个byte表示
func (e *TableMapEvent) decodeMeta(data []byte) error {
	pos := 0
	e.ColumnMeta = make([]uint16, e.ColumnCount)
	for i, t := range e.ColumnType {
		switch t {
		case comm.MYSQL_TYPE_STRING:
			var x = uint16(data[pos]) << 8 //real type
			x += uint16(data[pos+1])       //pack or field length
			e.ColumnMeta[i] = x
			pos += 2
		case comm.MYSQL_TYPE_NEWDECIMAL:
			var x = uint16(data[pos]) << 8 //precision
			x += uint16(data[pos+1])       //decimals
			e.ColumnMeta[i] = x
			pos += 2
		case comm.MYSQL_TYPE_VAR_STRING,
			comm.MYSQL_TYPE_VARCHAR,
			comm.MYSQL_TYPE_BIT:
			util.DecodeUint16(data, pos, &e.ColumnMeta[i])
			pos += 2
		case comm.MYSQL_TYPE_BLOB,
			comm.MYSQL_TYPE_DOUBLE,
			comm.MYSQL_TYPE_FLOAT,
			comm.MYSQL_TYPE_GEOMETRY,
			comm.MYSQL_TYPE_JSON:
			e.ColumnMeta[i] = uint16(data[pos])
			pos++
		case comm.MYSQL_TYPE_TIME2,
			comm.MYSQL_TYPE_DATETIME2,
			comm.MYSQL_TYPE_TIMESTAMP2:
			e.ColumnMeta[i] = uint16(data[pos])
			pos++
		case comm.MYSQL_TYPE_NEWDATE,
			comm.MYSQL_TYPE_ENUM,
			comm.MYSQL_TYPE_SET,
			comm.MYSQL_TYPE_TINY_BLOB,
			comm.MYSQL_TYPE_MEDIUM_BLOB,
			comm.MYSQL_TYPE_LONG_BLOB:
			return fmt.Errorf("unsupport type in binlog %d", t)
		default:
			e.ColumnMeta[i] = 0
		}
	}

	return nil
}

func (e *TableMapEvent) Write(w io.Writer) {
	fmt.Fprintf(w, "TableID: %d; ", e.TableID)
	fmt.Fprintf(w, "TableID size: %d; ", e.tableIDSize)
	fmt.Fprintf(w, "Flags: %d; ", e.Flags)
	fmt.Fprintf(w, "Schema: %s; ", e.Schema)
	fmt.Fprintf(w, "Table: %s; ", e.Table)
	fmt.Fprintf(w, "Column count: %d; ", e.ColumnCount)
	fmt.Fprintf(w, "Column type: ; %s", hex.Dump(e.ColumnType))
	fmt.Fprintf(w, "NULL bitmap: ; %s", hex.Dump(e.NullBitmap))
	fmt.Fprintln(w)
}

//----------------------------------RowsEvent----------------------------------------------//

// RowsEventStmtEndFlag is set in the end of the statement.
const RowsEventStmtEndFlag = 0x01

type errMissingTableMapEvent error

//https://dev.mysql.com/doc/internals/en/rows-event.html
// 只支持Version2，即5.7及以上
type RowsEvent struct {
	//0, 1, 2
	Version int

	tableIDSize int
	tables      map[uint64]*TableMapEvent
	needBitmap2 bool

	Table *TableMapEvent

	TableID uint64

	Flags uint16

	//if version == 2
	ExtraData []byte

	//lenenc_int
	ColumnCount uint64
	//len = (ColumnCount + 7) / 8
	ColumnBitmap1 []byte

	//if UPDATE_ROWS_EVENTv1 or v2
	//len = (ColumnCount + 7) / 8
	ColumnBitmap2 []byte

	//rows: invalid: int64, float64, bool, []byte, string
	Rows [][]interface{}

	parseTime  bool
	useDecimal bool
}

func (e *RowsEvent) Read(data []byte) error {
	pos := 0
	e.TableID = util.FixedLengthInt(data[0:e.tableIDSize])
	pos += e.tableIDSize

	util.DecodeUint16(data, pos, &e.Flags)
	pos += 2

	if e.Version == 2 {
		dataLen := binary.LittleEndian.Uint16(data[pos:])
		pos += 2

		e.ExtraData = data[pos: pos+int(dataLen-2)]
		pos += int(dataLen - 2)
	}

	var n int
	e.ColumnCount, _, n = util.LengthEncodedInt(data[pos:])
	pos += n

	bitCount := bitmapByteSize(int(e.ColumnCount))
	e.ColumnBitmap1 = data[pos: pos+bitCount]
	pos += bitCount

	if e.needBitmap2 {
		e.ColumnBitmap2 = data[pos: pos+bitCount]
		pos += bitCount
	}

	var ok bool
	e.Table, ok = e.tables[e.TableID]
	if !ok {
		if len(e.tables) > 0 {
			return fmt.Errorf("invalid table id %d, no corresponding table map event", e.TableID)
		} else {
			return errMissingTableMapEvent(
				fmt.Errorf("invalid table id %d, no corresponding table map event", e.TableID))
		}
	}

	var err error

	// ... repeat rows until event-end
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("parse rows event panic %v, data %q, parsed rows %#v, table map %#v; %s",
				r, data, e, e.Table, util.Pstack())
		}
	}()

	for pos < len(data) {
		if n, err = e.decodeRows(data[pos:], e.Table, e.ColumnBitmap1); err != nil {
			return err
		}
		pos += n

		if e.needBitmap2 {
			if n, err = e.decodeRows(data[pos:], e.Table, e.ColumnBitmap2); err != nil {
				return err
			}
			pos += n
		}
	}

	return nil
}

// TODO
func isBitSet(bitmap []byte, i int) bool {
	return bitmap[i>>3]&(1<<(uint(i)&7)) > 0
}

func (e *RowsEvent) decodeRows(data []byte, table *TableMapEvent, bitmap []byte) (int, error) {
	row := make([]interface{}, e.ColumnCount)
	pos := 0

	count := 0
	for i := 0; i < int(e.ColumnCount); i++ {
		if isBitSet(bitmap, i) {
			count++
		}
	}
	count = (count + 7) / 8 // null-bitmap, length (bits set in 'columns-present-bitmap1'+7)/8

	nullBitmap := data[pos: pos+count]
	pos += count

	nullbitIndex := 0

	var n int
	var err error
	// 读取一行
	for i := 0; i < int(e.ColumnCount); i++ {
		if !isBitSet(bitmap, i) { // not present
			continue
		}

		isNull := (uint32(nullBitmap[nullbitIndex/8]) >> uint32(nullbitIndex%8)) & 0x01
		nullbitIndex++

		if isNull > 0 {
			row[i] = nil
			continue
		}

		row[i], n, err = e.decodeValue(data[pos:], table.ColumnType[i], table.ColumnMeta[i])

		if err != nil {
			return 0, err
		}
		pos += n
	}

	e.Rows = append(e.Rows, row)
	return pos, nil
}

func (e *RowsEvent) parseFracTime(t interface{}) interface{} {
	v, ok := t.(fracTime)
	if !ok {
		return t
	}

	if !e.parseTime {
		// Don't parse time, return string directly
		return v.String()
	}

	// return Golang time directly
	return v.Time
}

// see mysql sql/log_event.cc log_event_print_value 根据类型和元数据解析值
func (e *RowsEvent) decodeValue(data []byte, tp byte, meta uint16) (v interface{}, n int, err error) {
	var length int = 0

	if tp == comm.MYSQL_TYPE_STRING {
		if meta >= 256 {
			b0 := uint8(meta >> 8)
			b1 := uint8(meta & 0xFF)

			if b0&0x30 != 0x30 {
				length = int(uint16(b1) | (uint16((b0&0x30)^0x30) << 4))
				tp = byte(b0 | 0x30)
			} else {
				length = int(meta & 0xFF)
				tp = b0
			}
		} else {
			length = int(meta)
		}
	}

	switch tp {
	case comm.MYSQL_TYPE_NULL:
		return nil, 0, nil
	case comm.MYSQL_TYPE_LONG:
		n = 4
		v = util.ParseBinaryInt32(data)
	case comm.MYSQL_TYPE_TINY:
		n = 1
		v = util.ParseBinaryInt8(data)
	case comm.MYSQL_TYPE_SHORT:
		n = 2
		v = util.ParseBinaryInt16(data)
	case comm.MYSQL_TYPE_INT24:
		n = 3
		v = util.ParseBinaryInt24(data)
	case comm.MYSQL_TYPE_LONGLONG:
		n = 8
		v = util.ParseBinaryInt64(data)
	case comm.MYSQL_TYPE_NEWDECIMAL:
		prec := uint8(meta >> 8)
		scale := uint8(meta & 0xFF)
		v, n, err = decodeDecimal(data, int(prec), int(scale), e.useDecimal)
	case comm.MYSQL_TYPE_FLOAT:
		n = 4
		v = util.ParseBinaryFloat32(data)
	case comm.MYSQL_TYPE_DOUBLE:
		n = 8
		v = util.ParseBinaryFloat64(data)
	case comm.MYSQL_TYPE_BIT:
		nbits := ((meta >> 8) * 8) + (meta & 0xFF)
		n = int(nbits+7) / 8 // 各种转化操作都是以byte为单位， 1byte=8bits
		//use int64 for bit
		v, err = decodeBit(data, int(nbits), int(n))
	case comm.MYSQL_TYPE_TIMESTAMP:
		n = 4
		t := binary.LittleEndian.Uint32(data)
		v = e.parseFracTime(fracTime{time.Unix(int64(t), 0), 0})
	case comm.MYSQL_TYPE_TIMESTAMP2:
		v, n, err = decodeTimestamp2(data, meta)
		v = e.parseFracTime(v)
	case comm.MYSQL_TYPE_DATETIME:
		n = 8
		i64 := binary.LittleEndian.Uint64(data)
		d := i64 / 1000000
		t := i64 % 1000000
		v = e.parseFracTime(fracTime{time.Date(int(d/10000), // year
			time.Month((d%10000)/100),                       // month
			int(d%100),                                      //day
			int(t/10000),                                    //hour
			int((t%10000)/100),                              // minute
			int(t%100),                                      //second
			0,
			time.UTC), 0})
	case comm.MYSQL_TYPE_DATETIME2:
		v, n, err = decodeDatetime2(data, meta)
		v = e.parseFracTime(v)
	case comm.MYSQL_TYPE_TIME:
		n = 3
		i32 := uint32(util.FixedLengthInt(data[0:3]))
		if i32 == 0 {
			v = "00:00:00"
		} else {
			sign := ""
			if i32 < 0 {
				sign = "-"
			}
			v = fmt.Sprintf("%s%02d:%02d:%02d", sign, i32/10000, (i32%10000)/100, i32%100)
		}
	case comm.MYSQL_TYPE_TIME2:
		v, n, err = decodeTime2(data, meta)
	case comm.MYSQL_TYPE_DATE:
		n = 3
		i32 := uint32(util.FixedLengthInt(data[0:3]))
		if i32 == 0 {
			v = "0000-00-00"
		} else {
			v = fmt.Sprintf("%04d-%02d-%02d", i32/(16*32), i32/32%16, i32%32)
		}
	case comm.MYSQL_TYPE_YEAR:
		n = 1
		v = int(data[0]) + 1900
	case comm.MYSQL_TYPE_ENUM:
		l := meta & 0xFF
		switch l {
		case 1:
			v = int64(data[0])
			n = 1
		case 2:
			v = int64(binary.BigEndian.Uint16(data))
			n = 2
		default:
			err = fmt.Errorf("Unknown ENUM packlen=%d", l)
		}
	case comm.MYSQL_TYPE_SET:
		n = int(meta & 0xFF)
		nbits := n * 8
		v, err = decodeBit(data, nbits, n)
	case comm.MYSQL_TYPE_BLOB:
		v, n, err = decodeBlob(data, meta)
	case comm.MYSQL_TYPE_VARCHAR,
		comm.MYSQL_TYPE_VAR_STRING:
		length = int(meta)
		v, n = decodeString(data, length)
	case comm.MYSQL_TYPE_STRING:
		v, n = decodeString(data, length)
	case comm.MYSQL_TYPE_GEOMETRY:
		// MySQL saves Geometry as Blob in binlog
		// Seem that the binary format is SRID (4 bytes) + WKB, outer can use
		// MySQL GeoFromWKB or others to create the geometry data.
		// Refer https://dev.mysql.com/doc/refman/5.7/en/gis-wkb-functions.html
		v, n, err = decodeBlob(data, meta)
	default: // TODO 支持 MYSQL_TYPE_JSON
		err = fmt.Errorf("unsupport type %d in binlog and don't know how to handle", tp)
	}
	return
}

func decodeString(data []byte, length int) (v string, n int) {
	if length < 256 {
		length = int(data[0])

		n = int(length) + 1
		v = string(data[1:n])
	} else {
		length = int(binary.LittleEndian.Uint16(data[0:]))
		n = length + 2
		v = string(data[2:n])
	}

	return
}

const digitsPerInteger int = 9

var compressedBytes = []int{0, 1, 1, 2, 2, 3, 3, 4, 4, 4}

func decodeDecimalDecompressValue(compIndx int, data []byte, mask uint8) (size int, value uint32) {
	size = compressedBytes[compIndx]
	databuff := make([]byte, size)
	for i := 0; i < size; i++ {
		databuff[i] = data[i] ^ mask
	}
	value = uint32(util.BFixedLengthInt(databuff))
	return
}

// https://dev.mysql.com/doc/refman/5.7/en/precision-math-decimal-characteristics.html
func decodeDecimal(data []byte, precision int, decimals int, useDecimal bool) (interface{}, int, error) {
	// decimal整数和小数部分不同长度存储空间大小不一样，如下图，比如Decimal(20,6)，整数部分长14， 小数部分长为6，
	// 14 = 9*1+5, 长度9需要4bytes，长度4需要3bytes，因此整数部分需要7bytes，同理小数部分需要3bytes
	// Leftover Digits	  Number of Bytes
	//     0						0
	//    1–2					1
	//    3–4					2
	//    5–6					3
	//	  7–9					4
	integral := precision - decimals // 整数部分长度
	uncompIntegral := int(integral / digitsPerInteger)
	uncompFractional := int(decimals / digitsPerInteger)
	compIntegral := integral - (uncompIntegral * digitsPerInteger)
	compFractional := decimals - (uncompFractional * digitsPerInteger)

	binSize := uncompIntegral*4 + compressedBytes[compIntegral] +
		uncompFractional*4 + compressedBytes[compFractional]

	buf := make([]byte, binSize)
	copy(buf, data[:binSize])

	//must copy the data for later change
	data = buf

	// Support negative
	// The sign is encoded in the high bit of the the byte
	// But this bit can also be used in the value
	value := uint32(data[0])
	var res bytes.Buffer
	var mask uint32 = 0
	if value & 0x80 == 0 { // 最高位为符号位， 最高位为0表示负数 TODO???
		mask = uint32((1 << 32) - 1)
		res.WriteString("-")
	}

	//clear sign
	data[0] ^= 0x80

	pos, value := decodeDecimalDecompressValue(compIntegral, data, uint8(mask))
	res.WriteString(fmt.Sprintf("%d", value))

	for i := 0; i < uncompIntegral; i++ {
		value = binary.BigEndian.Uint32(data[pos:]) ^ mask
		pos += 4
		res.WriteString(fmt.Sprintf("%09d", value)) // TODO 09d ???
	}

	res.WriteString(".")

	for i := 0; i < uncompFractional; i++ {
		value = binary.BigEndian.Uint32(data[pos:]) ^ mask
		pos += 4
		res.WriteString(fmt.Sprintf("%09d", value))
	}

	if size, value := decodeDecimalDecompressValue(compFractional, data[pos:], uint8(mask)); size > 0 {
		res.WriteString(fmt.Sprintf("%0*d", compFractional, value))
		pos += size
	}

	log.Infof("--decode decimal: %s", string(res.Bytes()))

	if useDecimal {
		f, err := decimal.NewFromString(string(res.Bytes()))
		return f, pos, err
	}

	f, err := strconv.ParseFloat(string(res.Bytes()), 64)
	return f, pos, err
}

// https://dev.mysql.com/doc/refman/5.7/en/bit-type.html
func decodeBit(data []byte, nbits int, length int) (value int64, err error) {
	log.Infof("--decode bits, with nbits:%d, length:%d", nbits, length)
	if nbits > 1 {
		switch length {
		case 1:
			value = int64(data[0])
		case 2:
			value = int64(binary.BigEndian.Uint16(data))
		case 3:
			value = int64(util.BFixedLengthInt(data[0:3]))
		case 4:
			value = int64(binary.BigEndian.Uint32(data))
		case 5:
			value = int64(util.BFixedLengthInt(data[0:5]))
		case 6:
			value = int64(util.BFixedLengthInt(data[0:6]))
		case 7:
			value = int64(util.BFixedLengthInt(data[0:7]))
		case 8:
			value = int64(binary.BigEndian.Uint64(data))
		default:
			err = fmt.Errorf("invalid bit length %d", length)
		}
	} else {
		if length != 1 {
			err = fmt.Errorf("invalid bit length %d", length)
		} else {
			value = int64(data[0])
		}
	}
	return
}

func decodeTimestamp2(data []byte, dec uint16) (interface{}, int, error) {
	//get timestamp binary length
	n := int(4 + (dec+1)/2)
	sec := int64(binary.BigEndian.Uint32(data[0:4]))
	usec := int64(0)
	switch dec {
	case 1, 2:
		usec = int64(data[4]) * 10000
	case 3, 4:
		usec = int64(binary.BigEndian.Uint16(data[4:])) * 100
	case 5, 6:
		usec = int64(util.BFixedLengthInt(data[4:7]))
	}

	if sec == 0 {
		return formatZeroTime(int(usec), int(dec)), n, nil
	}

	return fracTime{time.Unix(sec, usec*1000), int(dec)}, n, nil
}

const DATETIMEF_INT_OFS int64 = 0x8000000000

// TODO
func decodeDatetime2(data []byte, dec uint16) (interface{}, int, error) {
	//get datetime binary length
	n := int(5 + (dec+1)/2)

	intPart := int64(util.BFixedLengthInt(data[0:5])) - DATETIMEF_INT_OFS
	var frac int64 = 0

	switch dec {
	case 1, 2:
		frac = int64(data[5]) * 10000
	case 3, 4:
		frac = int64(binary.BigEndian.Uint16(data[5:7])) * 100
	case 5, 6:
		frac = int64(util.BFixedLengthInt(data[5:8]))
	}

	if intPart == 0 {
		return formatZeroTime(int(frac), int(dec)), n, nil
	}

	tmp := intPart<<24 + frac
	//handle sign???
	if tmp < 0 {
		tmp = -tmp
	}

	// var secPart int64 = tmp % (1 << 24)
	ymdhms := tmp >> 24

	ymd := ymdhms >> 17
	ym := ymd >> 5
	hms := ymdhms % (1 << 17)

	day := int(ymd % (1 << 5))
	month := int(ym % 13)
	year := int(ym / 13)

	second := int(hms % (1 << 6))
	minute := int((hms >> 6) % (1 << 6))
	hour := int((hms >> 12))

	return fracTime{time.Date(year, time.Month(month), day, hour, minute, second, int(frac*1000), time.UTC), int(dec)}, n, nil
}

const TIMEF_OFS int64 = 0x800000000000
const TIMEF_INT_OFS int64 = 0x800000

// TODO
func decodeTime2(data []byte, dec uint16) (string, int, error) {
	//time  binary length
	n := int(3 + (dec+1)/2)

	tmp := int64(0)
	intPart := int64(0)
	frac := int64(0)
	switch dec {
	case 1:
	case 2:
		intPart = int64(util.BFixedLengthInt(data[0:3])) - TIMEF_INT_OFS
		frac = int64(data[3])
		if intPart < 0 && frac > 0 {
			/*
			   Negative values are stored with reverse fractional part order,
			   for binary sort compatibility.

			     Disk value  intpart frac   Time value   Memory value
			     800000.00    0      0      00:00:00.00  0000000000.000000
			     7FFFFF.FF   -1      255   -00:00:00.01  FFFFFFFFFF.FFD8F0
			     7FFFFF.9D   -1      99    -00:00:00.99  FFFFFFFFFF.F0E4D0
			     7FFFFF.00   -1      0     -00:00:01.00  FFFFFFFFFF.000000
			     7FFFFE.FF   -1      255   -00:00:01.01  FFFFFFFFFE.FFD8F0
			     7FFFFE.F6   -2      246   -00:00:01.10  FFFFFFFFFE.FE7960

			     Formula to convert fractional part from disk format
			     (now stored in "frac" variable) to absolute value: "0x100 - frac".
			     To reconstruct in-memory value, we shift
			     to the next integer value and then substruct fractional part.
			*/
			intPart++     /* Shift to the next integer value */
			frac -= 0x100 /* -(0x100 - frac) */
		}
		tmp = intPart<<24 + frac*10000
	case 3:
	case 4:
		intPart = int64(util.BFixedLengthInt(data[0:3])) - TIMEF_INT_OFS
		frac = int64(binary.BigEndian.Uint16(data[3:5]))
		if intPart < 0 && frac > 0 {
			/*
			   Fix reverse fractional part order: "0x10000 - frac".
			   See comments for FSP=1 and FSP=2 above.
			*/
			intPart++       /* Shift to the next integer value */
			frac -= 0x10000 /* -(0x10000-frac) */
		}
		tmp = intPart<<24 + frac*100

	case 5:
	case 6:
		tmp = int64(util.BFixedLengthInt(data[0:6])) - TIMEF_OFS
	default:
		intPart = int64(util.BFixedLengthInt(data[0:3])) - TIMEF_INT_OFS
		tmp = intPart << 24
	}

	if intPart == 0 {
		return "00:00:00", n, nil
	}

	hms := int64(0)
	sign := ""
	if tmp < 0 {
		tmp = -tmp
		sign = "-"
	}

	hms = tmp >> 24

	hour := (hms >> 12) % (1 << 10) /* 10 bits starting at 12th */
	minute := (hms >> 6) % (1 << 6) /* 6 bits starting at 6th   */
	second := hms % (1 << 6)        /* 6 bits starting at 0th   */
	secPart := tmp % (1 << 24)

	if secPart != 0 {
		return fmt.Sprintf("%s%02d:%02d:%02d.%06d", sign, hour, minute, second, secPart), n, nil
	}

	return fmt.Sprintf("%s%02d:%02d:%02d", sign, hour, minute, second), n, nil
}

// 不同类型（即长度不同）的blob
func decodeBlob(data []byte, meta uint16) (v []byte, n int, err error) {
	var length int
	switch meta {
	case 1:
		length = int(data[0])
		v = data[1: 1+length]
		n = length + 1
	case 2:
		length = int(binary.LittleEndian.Uint16(data))
		v = data[2: 2+length]
		n = length + 2
	case 3:
		length = int(util.FixedLengthInt(data[0:3]))
		v = data[3: 3+length]
		n = length + 3
	case 4:
		length = int(binary.LittleEndian.Uint32(data))
		v = data[4: 4+length]
		n = length + 4
	default:
		err = fmt.Errorf("invalid blob packlen = %d", meta)
	}

	return
}

func (e *RowsEvent) Write(w io.Writer) {
	fmt.Fprintf(w, "TableID: %d; ", e.TableID)
	fmt.Fprintf(w, "Flags: %d; ", e.Flags)
	fmt.Fprintf(w, "Column count: %d; ", e.ColumnCount)

	fmt.Fprintf(w, "Values:; ")
	for _, rows := range e.Rows {
		fmt.Fprintf(w, "--:")
		for j, d := range rows {
			if _, ok := d.([]byte); ok {
				fmt.Fprintf(w, "%d:%q; ", j, d)
			} else {
				fmt.Fprintf(w, "%d:%#v; ", j, d)
			}
		}
	}
	fmt.Fprintln(w)
}

//----------------------------------RowsQueryEvent----------------------------------------------//

type RowsQueryEvent struct {
	Query []byte
}

func (e *RowsQueryEvent) Read(data []byte) error {
	//ignore length byte 1
	e.Query = data[1:]
	return nil
}

func (e *RowsQueryEvent) Write(w io.Writer) {
	fmt.Fprintf(w, "Query: %s; ", e.Query)
	fmt.Fprintln(w)
}
