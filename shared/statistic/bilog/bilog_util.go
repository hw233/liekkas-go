package bilog

import "go.uber.org/zap/zapcore"

func boolToInt8(b bool) int8 {
	if b {
		return 1
	}

	return 0
}

func int32SliceEncoder(arr []int32) zapcore.ArrayMarshalerFunc {
	return zapcore.ArrayMarshalerFunc(func(encoder zapcore.ArrayEncoder) error {
		for _, value := range arr {
			encoder.AppendInt32(value)
		}

		return nil
	})
}
