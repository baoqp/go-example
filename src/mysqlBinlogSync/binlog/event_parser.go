package binlog

import (
	"fmt"
)

type EventError struct {
	Header *EventHeader

	//Error message
	Err string

	//Event data
	Data []byte
}

func (e *EventError) Error() string {
	return e.Err
}

type EventParser struct {
	format         *FormatDescriptionEvent
	tables         map[uint64]*TableMapEvent
	parseTime      bool
	stopProcessing uint32 // used to start/stop processing
	useDecimal     bool
}

func NewEventParser() *EventParser {
	p := new(EventParser)
	p.tables = make(map[uint64]*TableMapEvent)
	return p
}

func (parser *EventParser) Parse(data []byte) (*BinlogEvent, error) {

	rawData := data

	header, err := parser.ParseHeader(data)

	if err != nil {
		return nil, err
	}

	data = data[EventHeaderSize:]
	eventLen := int(header.EventSize) - EventHeaderSize

	if len(data) != eventLen {
		return nil, fmt.Errorf("invalid data size %d in event %s, less event length %d",
			len(data), header.EventType, eventLen)
	}

	e, err := parser.parseEvent(data, header)
	if err != nil {
		return nil, err
	}

	return &BinlogEvent{rawData, header, e}, nil

}

func (parser *EventParser) Reset() {
	parser.format = nil
}

func (parser *EventParser) ParseHeader(data []byte) (*EventHeader, error) {
	header := &EventHeader{}
	err := header.Read(data)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (parser *EventParser) parseEvent(data []byte, header *EventHeader) (Event, error) {
	var e Event

	if header.EventType == FORMAT_DESCRIPTION_EVENT {
		parser.format = &FormatDescriptionEvent{}
		e = parser.format
	} else {
		// TODO ???
		if parser.format != nil && parser.format.ChecksumAlgorithm == BINLOG_CHECKSUM_ALG_CRC32 {
			data = data[0: len(data)-4]
		}

		switch header.EventType {
		case ROTATE_EVENT:
			e = &RotateEvent{}
		case QUERY_EVENT:
			e = &QueryEvent{}
		case TABLE_MAP_EVENT:
			tme := &TableMapEvent{}
			if parser.format.EventTypeHeaderLengths[TABLE_MAP_EVENT-1] == 6 {
				tme.tableIDSize = 4
			} else {
				tme.tableIDSize = 6
			}
			e = tme
		case WRITE_ROWS_EVENTv0,
			UPDATE_ROWS_EVENTv0,
			DELETE_ROWS_EVENTv0,
			WRITE_ROWS_EVENTv1,
			DELETE_ROWS_EVENTv1,
			UPDATE_ROWS_EVENTv1,
			WRITE_ROWS_EVENTv2,
			UPDATE_ROWS_EVENTv2,
			DELETE_ROWS_EVENTv2:
			e = parser.newRowsEvent(header)
		case ROWS_QUERY_EVENT:
			e = &RowsQueryEvent{}
		default:
			e = &GenericEvent{}
		}

	}

	if err := e.Read(data); err != nil {
		return nil, &EventError{header, err.Error(), data}
	}

	if te, ok := e.(*TableMapEvent); ok {
		parser.tables[te.TableID] = te
	}

	// end of statement
	if re, ok := e.(*RowsEvent); ok {
		if (re.Flags & RowsEventStmtEndFlag) > 0 {
			parser.tables = make(map[uint64]*TableMapEvent)
		}
	}

	return e, nil
}

func (parser *EventParser) newRowsEvent(h *EventHeader) *RowsEvent {
	e := &RowsEvent{}
	if parser.format.EventTypeHeaderLengths[h.EventType-1] == 6 {
		e.tableIDSize = 4
	} else {
		e.tableIDSize = 6
	}

	e.needBitmap2 = false
	e.tables = parser.tables
	e.parseTime = parser.parseTime
	e.useDecimal = parser.useDecimal

	switch h.EventType {
	case WRITE_ROWS_EVENTv0:
		e.Version = 0
	case UPDATE_ROWS_EVENTv0:
		e.Version = 0
	case DELETE_ROWS_EVENTv0:
		e.Version = 0
	case WRITE_ROWS_EVENTv1:
		e.Version = 1
	case DELETE_ROWS_EVENTv1:
		e.Version = 1
	case UPDATE_ROWS_EVENTv1:
		e.Version = 1
		e.needBitmap2 = true
	case WRITE_ROWS_EVENTv2:
		e.Version = 2
	case UPDATE_ROWS_EVENTv2:
		e.Version = 2
		e.needBitmap2 = true
	case DELETE_ROWS_EVENTv2:
		e.Version = 2
	}

	return e
}
